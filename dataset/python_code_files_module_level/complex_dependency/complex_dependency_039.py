import copy
import math
import os
import re
from dataclasses import dataclass
from functools import partial
from typing import TYPE_CHECKING, Optional, Union
import pyarrow as pa
import pyarrow.parquet as pq
from tqdm.contrib.concurrent import thread_map
from .download.download_config import DownloadConfig
from .naming import _split_re, filenames_for_dataset_split
from .table import InMemoryTable, MemoryMappedTable, Table, concat_tables
from .utils import logging
from .utils import tqdm as hf_tqdm
if TYPE_CHECKING:
    from .info import DatasetInfo
    from .splits import Split, SplitInfo
logger = logging.get_logger(__name__)
HF_GCP_BASE_URL = 'https://storage.googleapis.com/huggingface-nlp/cache/datasets'
_SUB_SPEC_RE = re.compile(f'\n^\n (?P<split>{_split_re[1:-1]})\n (\\[\n    ((?P<from>-?[\\d_]+)\n     (?P<from_pct>%)?)?\n    :\n    ((?P<to>-?[\\d_]+)\n     (?P<to_pct>%)?)?\n \\])?(\\((?P<rounding>[^\\)]*)\\))?\n$\n', re.X)
_ADDITION_SEP_RE = re.compile('\\s*\\+\\s*')

class DatasetNotOnHfGcsError(ConnectionError):
    pass

class MissingFilesOnHfGcsError(ConnectionError):
    pass

@dataclass(frozen=True)
class FileInstructions:
    num_examples: int
    file_instructions: list[dict]

def make_file_instructions(name: str, split_infos: list['SplitInfo'], instruction: Union[str, 'ReadInstruction'], filetype_suffix: Optional[str]=None, prefix_path: Optional[str]=None) -> FileInstructions:
    if not isinstance(name, str):
        raise TypeError(f"Expected str 'name', but got: {type(name).__name__}")
    elif not name:
        raise ValueError("Expected non-empty str 'name'")
    name2len = {info.name: info.num_examples for info in split_infos}
    name2shard_lengths = {info.name: info.shard_lengths for info in split_infos}
    name2filenames = {info.name: filenames_for_dataset_split(path=prefix_path, dataset_name=name, split=info.name, filetype_suffix=filetype_suffix, shard_lengths=name2shard_lengths[info.name]) for info in split_infos}
    if not isinstance(instruction, ReadInstruction):
        instruction = ReadInstruction.from_spec(instruction)
    absolute_instructions = instruction.to_absolute(name2len)
    file_instructions = []
    num_examples = 0
    for abs_instr in absolute_instructions:
        split_length = name2len[abs_instr.splitname]
        filenames = name2filenames[abs_instr.splitname]
        shard_lengths = name2shard_lengths[abs_instr.splitname]
        from_ = 0 if abs_instr.from_ is None else abs_instr.from_
        to = split_length if abs_instr.to is None else abs_instr.to
        if shard_lengths is None:
            for filename in filenames:
                take = to - from_
                if take == 0:
                    continue
                num_examples += take
                file_instructions.append({'filename': filename, 'skip': from_, 'take': take})
        else:
            index_start = 0
            index_end = 0
            for filename, shard_length in zip(filenames, shard_lengths):
                index_end += shard_length
                if from_ < index_end and to > index_start:
                    skip = from_ - index_start if from_ > index_start else 0
                    take = to - index_start - skip if to < index_end else -1
                    if take == 0:
                        continue
                    file_instructions.append({'filename': filename, 'skip': skip, 'take': take})
                    num_examples += shard_length - skip if take == -1 else take
                index_start += shard_length
    return FileInstructions(num_examples=num_examples, file_instructions=file_instructions)

class BaseReader:

    def __init__(self, path: str, info: Optional['DatasetInfo']):
        self._path: str = path
        self._info: Optional['DatasetInfo'] = info
        self._filetype_suffix: Optional[str] = None

    def _get_table_from_filename(self, filename_skip_take, in_memory=False) -> Table:
        raise NotImplementedError

    def _read_files(self, files, in_memory=False) -> Table:
        if len(files) == 0 or not all((isinstance(f, dict) for f in files)):
            raise ValueError('please provide valid file informations')
        files = copy.deepcopy(files)
        for f in files:
            f['filename'] = os.path.join(self._path, f['filename'])
        pa_tables = thread_map(partial(self._get_table_from_filename, in_memory=in_memory), files, tqdm_class=hf_tqdm, desc='Loading dataset shards', disable=len(files) <= 16 or None)
        pa_tables = [t for t in pa_tables if len(t) > 0]
        if not pa_tables and (self._info is None or self._info.features is None):
            raise ValueError('Tried to read an empty table. Please specify at least info.features to create an empty table with the right type.')
        pa_tables = pa_tables or [InMemoryTable.from_batches([], schema=pa.schema(self._info.features.type))]
        pa_table = concat_tables(pa_tables) if len(pa_tables) != 1 else pa_tables[0]
        return pa_table

    def get_file_instructions(self, name, instruction, split_infos):
        file_instructions = make_file_instructions(name, split_infos, instruction, filetype_suffix=self._filetype_suffix, prefix_path=self._path)
        files = file_instructions.file_instructions
        return files

    def read(self, name, instructions, split_infos, in_memory=False):
        files = self.get_file_instructions(name, instructions, split_infos)
        if not files:
            msg = f'Instruction "{instructions}" corresponds to no data!'
            raise ValueError(msg)
        return self.read_files(files=files, original_instructions=instructions, in_memory=in_memory)

    def read_files(self, files: list[dict], original_instructions: Union[None, 'ReadInstruction', 'Split']=None, in_memory=False):
        pa_table = self._read_files(files, in_memory=in_memory)
        if original_instructions is not None:
            from .splits import Split
            split = Split(str(original_instructions))
        else:
            split = None
        dataset_kwargs = {'arrow_table': pa_table, 'info': self._info, 'split': split}
        return dataset_kwargs

class ArrowReader(BaseReader):

    def __init__(self, path: str, info: Optional['DatasetInfo']):
        super().__init__(path, info)
        self._filetype_suffix = 'arrow'

    def _get_table_from_filename(self, filename_skip_take, in_memory=False) -> Table:
        filename, skip, take = (filename_skip_take['filename'], filename_skip_take['skip'] if 'skip' in filename_skip_take else None, filename_skip_take['take'] if 'take' in filename_skip_take else None)
        table = ArrowReader.read_table(filename, in_memory=in_memory)
        if take == -1:
            take = len(table) - skip
        if skip is not None and take is not None and (not (skip == 0 and take == len(table))):
            table = table.slice(skip, take)
        return table

    @staticmethod
    def read_table(filename, in_memory=False) -> Table:
        table_cls = InMemoryTable if in_memory else MemoryMappedTable
        return table_cls.from_file(filename)

class ParquetReader(BaseReader):

    def __init__(self, path: str, info: Optional['DatasetInfo']):
        super().__init__(path, info)
        self._filetype_suffix = 'parquet'

    def _get_table_from_filename(self, filename_skip_take, **kwargs):
        filename, skip, take = (filename_skip_take['filename'], filename_skip_take['skip'] if 'skip' in filename_skip_take else None, filename_skip_take['take'] if 'take' in filename_skip_take else None)
        pa_table = pq.read_table(filename, memory_map=True)
        if skip is not None and take is not None and (not (skip == 0 and take == len(pa_table))):
            pa_table = pa_table.slice(skip, take)
        return pa_table

@dataclass(frozen=True)
class _AbsoluteInstruction:
    splitname: str
    from_: int
    to: int

@dataclass(frozen=True)
class _RelativeInstruction:
    splitname: str
    from_: Optional[int] = None
    to: Optional[int] = None
    unit: Optional[str] = None
    rounding: Optional[str] = None

    def __post_init__(self):
        if self.unit is not None and self.unit not in ['%', 'abs']:
            raise ValueError('unit must be either % or abs')
        if self.rounding is not None and self.rounding not in ['closest', 'pct1_dropremainder']:
            raise ValueError('rounding must be either closest or pct1_dropremainder')
        if self.unit != '%' and self.rounding is not None:
            raise ValueError('It is forbidden to specify rounding if not using percent slicing.')
        if self.unit == '%' and self.from_ is not None and (abs(self.from_) > 100):
            raise ValueError('Percent slice boundaries must be > -100 and < 100.')
        if self.unit == '%' and self.to is not None and (abs(self.to) > 100):
            raise ValueError('Percent slice boundaries must be > -100 and < 100.')
        self.__dict__['rounding'] = 'closest' if self.rounding is None and self.unit == '%' else self.rounding

def _str_to_read_instruction(spec):
    res = _SUB_SPEC_RE.match(spec)
    if not res:
        raise ValueError(f'Unrecognized instruction format: {spec}')
    unit = '%' if res.group('from_pct') or res.group('to_pct') else 'abs'
    return ReadInstruction(split_name=res.group('split'), rounding=res.group('rounding'), from_=int(res.group('from')) if res.group('from') else None, to=int(res.group('to')) if res.group('to') else None, unit=unit)

def _pct_to_abs_pct1(boundary, num_examples):
    if num_examples < 100:
        msg = 'Using "pct1_dropremainder" rounding on a split with less than 100 elements is forbidden: it always results in an empty dataset.'
        raise ValueError(msg)
    return boundary * math.trunc(num_examples / 100.0)

def _pct_to_abs_closest(boundary, num_examples):
    return int(round(boundary * num_examples / 100.0))

def _rel_to_abs_instr(rel_instr, name2len):
    pct_to_abs = _pct_to_abs_closest if rel_instr.rounding == 'closest' else _pct_to_abs_pct1
    split = rel_instr.splitname
    if split not in name2len:
        raise ValueError(f'Unknown split "{split}". Should be one of {list(name2len)}.')
    num_examples = name2len[split]
    from_ = rel_instr.from_
    to = rel_instr.to
    if rel_instr.unit == '%':
        from_ = 0 if from_ is None else pct_to_abs(from_, num_examples)
        to = num_examples if to is None else pct_to_abs(to, num_examples)
    else:
        from_ = 0 if from_ is None else from_
        to = num_examples if to is None else to
    if from_ < 0:
        from_ = max(num_examples + from_, 0)
    if to < 0:
        to = max(num_examples + to, 0)
    from_ = min(from_, num_examples)
    to = min(to, num_examples)
    return _AbsoluteInstruction(split, from_, to)

class ReadInstruction:

    def _init(self, relative_instructions):
        self._relative_instructions = relative_instructions

    @classmethod
    def _read_instruction_from_relative_instructions(cls, relative_instructions):
        result = cls.__new__(cls)
        result._init(relative_instructions)
        return result

    def __init__(self, split_name, rounding=None, from_=None, to=None, unit=None):
        self._init([_RelativeInstruction(split_name, from_, to, unit, rounding)])

    @classmethod
    def from_spec(cls, spec):
        spec = str(spec)
        subs = _ADDITION_SEP_RE.split(spec)
        if not subs:
            raise ValueError(f'No instructions could be built out of {spec}')
        instruction = _str_to_read_instruction(subs[0])
        return sum((_str_to_read_instruction(sub) for sub in subs[1:]), instruction)

    def to_spec(self):
        rel_instr_specs = []
        for rel_instr in self._relative_instructions:
            rel_instr_spec = rel_instr.splitname
            if rel_instr.from_ is not None or rel_instr.to is not None:
                from_ = rel_instr.from_
                to = rel_instr.to
                unit = rel_instr.unit
                rounding = rel_instr.rounding
                unit = unit if unit == '%' else ''
                from_ = str(from_) + unit if from_ is not None else ''
                to = str(to) + unit if to is not None else ''
                slice_str = f'[{from_}:{to}]'
                rounding_str = f'({rounding})' if unit == '%' and rounding is not None and (rounding != 'closest') else ''
                rel_instr_spec += slice_str + rounding_str
            rel_instr_specs.append(rel_instr_spec)
        return '+'.join(rel_instr_specs)

    def __add__(self, other):
        if not isinstance(other, ReadInstruction):
            msg = 'ReadInstruction can only be added to another ReadInstruction obj.'
            raise TypeError(msg)
        self_ris = self._relative_instructions
        other_ris = other._relative_instructions
        if self_ris[0].unit != 'abs' and other_ris[0].unit != 'abs' and (self._relative_instructions[0].rounding != other_ris[0].rounding):
            raise ValueError('It is forbidden to sum ReadInstruction instances with different rounding values.')
        return self._read_instruction_from_relative_instructions(self_ris + other_ris)

    def __str__(self):
        return self.to_spec()

    def __repr__(self):
        return f'ReadInstruction({self._relative_instructions})'

    def to_absolute(self, name2len):
        return [_rel_to_abs_instr(rel_instr, name2len) for rel_instr in self._relative_instructions]
