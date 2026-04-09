import io
import re
from dataclasses import dataclass
from math import ceil
from os.path import getsize
from typing import TYPE_CHECKING, BinaryIO, Iterable, Optional, TypedDict
from urllib.parse import unquote
from huggingface_hub import constants
from .utils import build_hf_headers, fix_hf_endpoint_in_url, hf_raise_for_status, http_backoff, logging, validate_hf_hub_args
from .utils._lfs import SliceFileObj
from .utils.sha import sha256, sha_fileobj
if TYPE_CHECKING:
    from ._commit_api import CommitOperationAdd
logger = logging.get_logger(__name__)
OID_REGEX = re.compile('^[0-9a-f]{40}$')
LFS_MULTIPART_UPLOAD_COMMAND = 'lfs-multipart-upload'
LFS_HEADERS = {'Accept': 'application/vnd.git-lfs+json', 'Content-Type': 'application/vnd.git-lfs+json'}

@dataclass
class UploadInfo:
    sha256: bytes
    size: int
    sample: bytes

    @classmethod
    def from_path(cls, path: str):
        size = getsize(path)
        with io.open(path, 'rb') as file:
            sample = file.peek(512)[:512]
            sha = sha_fileobj(file)
        return cls(size=size, sha256=sha, sample=sample)

    @classmethod
    def from_bytes(cls, data: bytes):
        sha = sha256(data).digest()
        return cls(size=len(data), sample=data[:512], sha256=sha)

    @classmethod
    def from_fileobj(cls, fileobj: BinaryIO):
        sample = fileobj.read(512)
        fileobj.seek(0, io.SEEK_SET)
        sha = sha_fileobj(fileobj)
        size = fileobj.tell()
        fileobj.seek(0, io.SEEK_SET)
        return cls(size=size, sha256=sha, sample=sample)

@validate_hf_hub_args
def post_lfs_batch_info(upload_infos: Iterable[UploadInfo], token: Optional[str], repo_type: str, repo_id: str, revision: Optional[str]=None, endpoint: Optional[str]=None, headers: Optional[dict[str, str]]=None, transfers: Optional[list[str]]=None) -> tuple[list[dict], list[dict], Optional[str]]:
    endpoint = endpoint if endpoint is not None else constants.ENDPOINT
    url_prefix = ''
    if repo_type in constants.REPO_TYPES_URL_PREFIXES:
        url_prefix = constants.REPO_TYPES_URL_PREFIXES[repo_type]
    batch_url = f'{endpoint}/{url_prefix}{repo_id}.git/info/lfs/objects/batch'
    payload: dict = {'operation': 'upload', 'transfers': transfers if transfers is not None else ['basic', 'multipart'], 'objects': [{'oid': upload.sha256.hex(), 'size': upload.size} for upload in upload_infos], 'hash_algo': 'sha256'}
    if revision is not None:
        payload['ref'] = {'name': unquote(revision)}
    headers = {**LFS_HEADERS, **build_hf_headers(token=token), **(headers or {})}
    resp = http_backoff('POST', batch_url, headers=headers, json=payload)
    hf_raise_for_status(resp)
    batch_info = resp.json()
    objects = batch_info.get('objects', None)
    if not isinstance(objects, list):
        raise ValueError('Malformed response from server')
    chosen_transfer = batch_info.get('transfer')
    chosen_transfer = chosen_transfer if isinstance(chosen_transfer, str) else None
    return ([_validate_batch_actions(obj) for obj in objects if 'error' not in obj], [_validate_batch_error(obj) for obj in objects if 'error' in obj], chosen_transfer)

class PayloadPartT(TypedDict):
    partNumber: int
    etag: str

class CompletionPayloadT(TypedDict):
    oid: str
    parts: list[PayloadPartT]

def lfs_upload(operation: 'CommitOperationAdd', lfs_batch_action: dict, token: Optional[str]=None, headers: Optional[dict[str, str]]=None, endpoint: Optional[str]=None) -> None:
    _validate_batch_actions(lfs_batch_action)
    actions = lfs_batch_action.get('actions')
    if actions is None:
        logger.debug(f'Content of file {operation.path_in_repo} is already present upstream - skipping upload')
        return
    upload_action = lfs_batch_action['actions']['upload']
    _validate_lfs_action(upload_action)
    verify_action = lfs_batch_action['actions'].get('verify')
    if verify_action is not None:
        _validate_lfs_action(verify_action)
    header = upload_action.get('header', {})
    chunk_size = header.get('chunk_size')
    upload_url = fix_hf_endpoint_in_url(upload_action['href'], endpoint=endpoint)
    if chunk_size is not None:
        try:
            chunk_size = int(chunk_size)
        except (ValueError, TypeError):
            raise ValueError(f"Malformed response from LFS batch endpoint: `chunk_size` should be an integer. Got '{chunk_size}'.")
        _upload_multi_part(operation=operation, header=header, chunk_size=chunk_size, upload_url=upload_url)
    else:
        _upload_single_part(operation=operation, upload_url=upload_url)
    if verify_action is not None:
        _validate_lfs_action(verify_action)
        verify_url = fix_hf_endpoint_in_url(verify_action['href'], endpoint)
        verify_resp = http_backoff('POST', verify_url, headers=build_hf_headers(token=token, headers=headers), json={'oid': operation.upload_info.sha256.hex(), 'size': operation.upload_info.size})
        hf_raise_for_status(verify_resp)
    logger.debug(f'{operation.path_in_repo}: Upload successful')

def _validate_lfs_action(lfs_action: dict):
    if not (isinstance(lfs_action.get('href'), str) and (lfs_action.get('header') is None or isinstance(lfs_action.get('header'), dict))):
        raise ValueError('lfs_action is improperly formatted')
    return lfs_action

def _validate_batch_actions(lfs_batch_actions: dict):
    if not (isinstance(lfs_batch_actions.get('oid'), str) and isinstance(lfs_batch_actions.get('size'), int)):
        raise ValueError('lfs_batch_actions is improperly formatted')
    upload_action = lfs_batch_actions.get('actions', {}).get('upload')
    verify_action = lfs_batch_actions.get('actions', {}).get('verify')
    if upload_action is not None:
        _validate_lfs_action(upload_action)
    if verify_action is not None:
        _validate_lfs_action(verify_action)
    return lfs_batch_actions

def _validate_batch_error(lfs_batch_error: dict):
    if not (isinstance(lfs_batch_error.get('oid'), str) and isinstance(lfs_batch_error.get('size'), int)):
        raise ValueError('lfs_batch_error is improperly formatted')
    error_info = lfs_batch_error.get('error')
    if not (isinstance(error_info, dict) and isinstance(error_info.get('message'), str) and isinstance(error_info.get('code'), int)):
        raise ValueError('lfs_batch_error is improperly formatted')
    return lfs_batch_error

def _upload_single_part(operation: 'CommitOperationAdd', upload_url: str) -> None:
    with operation.as_file(with_tqdm=True) as fileobj:
        response = http_backoff('PUT', upload_url, data=fileobj)
        hf_raise_for_status(response)

def _upload_multi_part(operation: 'CommitOperationAdd', header: dict, chunk_size: int, upload_url: str) -> None:
    sorted_parts_urls = _get_sorted_parts_urls(header=header, upload_info=operation.upload_info, chunk_size=chunk_size)
    response_headers = _upload_parts_iteratively(operation=operation, sorted_parts_urls=sorted_parts_urls, chunk_size=chunk_size)
    completion_res = http_backoff('POST', upload_url, json=_get_completion_payload(response_headers, operation.upload_info.sha256.hex()), headers=LFS_HEADERS)
    hf_raise_for_status(completion_res)

def _get_sorted_parts_urls(header: dict, upload_info: UploadInfo, chunk_size: int) -> list[str]:
    sorted_part_upload_urls = [upload_url for _, upload_url in sorted([(int(part_num, 10), upload_url) for part_num, upload_url in header.items() if part_num.isdigit() and len(part_num) > 0], key=lambda t: t[0])]
    num_parts = len(sorted_part_upload_urls)
    if num_parts != ceil(upload_info.size / chunk_size):
        raise ValueError('Invalid server response to upload large LFS file')
    return sorted_part_upload_urls

def _get_completion_payload(response_headers: list[dict], oid: str) -> CompletionPayloadT:
    parts: list[PayloadPartT] = []
    for part_number, header in enumerate(response_headers):
        etag = header.get('etag')
        if etag is None or etag == '':
            raise ValueError(f'Invalid etag (`{etag}`) returned for part {part_number + 1}')
        parts.append({'partNumber': part_number + 1, 'etag': etag})
    return {'oid': oid, 'parts': parts}

def _upload_parts_iteratively(operation: 'CommitOperationAdd', sorted_parts_urls: list[str], chunk_size: int) -> list[dict]:
    headers = []
    with operation.as_file(with_tqdm=True) as fileobj:
        for part_idx, part_upload_url in enumerate(sorted_parts_urls):
            with SliceFileObj(fileobj, seek_from=chunk_size * part_idx, read_limit=chunk_size) as fileobj_slice:
                part_upload_res = http_backoff('PUT', part_upload_url, data=fileobj_slice)
                hf_raise_for_status(part_upload_res)
                headers.append(part_upload_res.headers)
    return headers
