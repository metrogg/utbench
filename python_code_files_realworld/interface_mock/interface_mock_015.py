import datetime
import hashlib
import logging
import os
import time
import urllib.parse
import warnings
from dataclasses import dataclass
from typing import TYPE_CHECKING, Literal, Optional, Union
from . import constants
from .hf_api import whoami
from .utils import experimental, get_token
logger = logging.getLogger(__name__)
if TYPE_CHECKING:
    import fastapi

@dataclass
class OAuthOrgInfo:
    sub: str
    name: str
    preferred_username: str
    picture: str
    plan: Optional[str] = None
    can_pay: Optional[bool] = None
    role_in_org: Optional[str] = None
    security_restrictions: Optional[list[Literal['ip', 'token-policy', 'mfa', 'sso']]] = None

@dataclass
class OAuthUserInfo:
    sub: str
    name: str
    preferred_username: str
    email_verified: Optional[bool]
    email: Optional[str]
    picture: str
    profile: str
    website: Optional[str]
    is_pro: bool
    can_pay: Optional[bool]
    orgs: Optional[list[OAuthOrgInfo]]

@dataclass
class OAuthInfo:
    access_token: str
    access_token_expires_at: datetime.datetime
    user_info: OAuthUserInfo
    state: Optional[str]
    scope: str

@experimental
def attach_huggingface_oauth(app: 'fastapi.FastAPI', route_prefix: str='/'):
    try:
        from starlette.middleware.sessions import SessionMiddleware
    except ImportError as e:
        raise ImportError('Cannot initialize OAuth to due a missing library. Please run `pip install huggingface_hub[oauth]` or add `huggingface_hub[oauth]` to your requirements.txt file in order to install the required dependencies.') from e
    session_secret = (constants.OAUTH_CLIENT_SECRET or '') + '-v1'
    app.add_middleware(SessionMiddleware, secret_key=hashlib.sha256(session_secret.encode()).hexdigest(), same_site='none', https_only=True)
    route_prefix = route_prefix.strip('/')
    if os.getenv('SPACE_ID') is not None:
        logger.info('OAuth is enabled in the Space. Adding OAuth routes.')
        _add_oauth_routes(app, route_prefix=route_prefix)
    else:
        logger.info('App is not running in a Space. Adding mocked OAuth routes.')
        _add_mocked_oauth_routes(app, route_prefix=route_prefix)

def parse_huggingface_oauth(request: 'fastapi.Request') -> Optional[OAuthInfo]:
    if 'oauth_info' not in request.session:
        logger.debug('No OAuth info in session.')
        return None
    logger.debug('Parsing OAuth info from session.')
    oauth_data = request.session['oauth_info']
    user_data = oauth_data.get('userinfo', {})
    orgs_data = user_data.get('orgs', [])
    orgs = [OAuthOrgInfo(sub=org.get('sub'), name=org.get('name'), preferred_username=org.get('preferred_username'), picture=org.get('picture'), plan=org.get('plan'), can_pay=org.get('canPay'), role_in_org=org.get('roleInOrg'), security_restrictions=org.get('securityRestrictions')) for org in orgs_data] if orgs_data else None
    user_info = OAuthUserInfo(sub=user_data.get('sub'), name=user_data.get('name'), preferred_username=user_data.get('preferred_username'), email_verified=user_data.get('email_verified'), email=user_data.get('email'), picture=user_data.get('picture'), profile=user_data.get('profile'), website=user_data.get('website'), is_pro=user_data.get('isPro'), can_pay=user_data.get('canPay'), orgs=orgs)
    return OAuthInfo(access_token=oauth_data.get('access_token'), access_token_expires_at=datetime.datetime.fromtimestamp(oauth_data.get('expires_at')), user_info=user_info, state=oauth_data.get('state'), scope=oauth_data.get('scope'))

def _add_oauth_routes(app: 'fastapi.FastAPI', route_prefix: str) -> None:
    try:
        import fastapi
        from authlib.integrations.base_client.errors import MismatchingStateError
        from authlib.integrations.starlette_client import OAuth
        from fastapi.responses import RedirectResponse
    except ImportError as e:
        raise ImportError('Cannot initialize OAuth to due a missing library. Please run `pip install huggingface_hub[oauth]` or add `huggingface_hub[oauth]` to your requirements.txt file.') from e
    msg = "OAuth is required but '{}' environment variable is not set. Make sure you've enabled OAuth in your Space by setting `hf_oauth: true` in the Space metadata."
    if constants.OAUTH_CLIENT_ID is None:
        raise ValueError(msg.format('OAUTH_CLIENT_ID'))
    if constants.OAUTH_CLIENT_SECRET is None:
        raise ValueError(msg.format('OAUTH_CLIENT_SECRET'))
    if constants.OAUTH_SCOPES is None:
        raise ValueError(msg.format('OAUTH_SCOPES'))
    if constants.OPENID_PROVIDER_URL is None:
        raise ValueError(msg.format('OPENID_PROVIDER_URL'))
    oauth = OAuth()
    oauth.register(name='huggingface', client_id=constants.OAUTH_CLIENT_ID, client_secret=constants.OAUTH_CLIENT_SECRET, client_kwargs={'scope': constants.OAUTH_SCOPES}, server_metadata_url=constants.OPENID_PROVIDER_URL + '/.well-known/openid-configuration')
    login_uri, callback_uri, logout_uri = _get_oauth_uris(route_prefix)

    @app.get(login_uri)
    async def oauth_login(request: fastapi.Request) -> RedirectResponse:
        redirect_uri = _generate_redirect_uri(request)
        return await oauth.huggingface.authorize_redirect(request, redirect_uri)

    @app.get(callback_uri)
    async def oauth_redirect_callback(request: fastapi.Request) -> RedirectResponse:
        try:
            oauth_info = await oauth.huggingface.authorize_access_token(request)
        except MismatchingStateError:
            nb_redirects = int(request.query_params.get('_nb_redirects', 0))
            target_url = request.query_params.get('_target_url')
            query_params: dict[str, Union[int, str]] = {'_nb_redirects': nb_redirects + 1}
            if target_url:
                query_params['_target_url'] = target_url
            redirect_uri = f'{login_uri}?{urllib.parse.urlencode(query_params)}'
            if nb_redirects > constants.OAUTH_MAX_REDIRECTS:
                host = os.environ.get('SPACE_HOST')
                if host is None:
                    raise RuntimeError('App is not running in a Space (SPACE_HOST environment variable is not set). Cannot redirect to non-iframe view.') from None
                host_url = 'https://' + host.rstrip('/')
                return RedirectResponse(host_url + redirect_uri)
            return RedirectResponse(redirect_uri)
        logger.debug('Successfully logged in with OAuth. Storing user info in session.')
        request.session['oauth_info'] = oauth_info
        return RedirectResponse(_get_redirect_target(request))

    @app.get(logout_uri)
    async def oauth_logout(request: fastapi.Request) -> RedirectResponse:
        logger.debug('Logged out with OAuth. Removing user info from session.')
        request.session.pop('oauth_info', None)
        return RedirectResponse(_get_redirect_target(request))

def _add_mocked_oauth_routes(app: 'fastapi.FastAPI', route_prefix: str='/') -> None:
    try:
        import fastapi
        from fastapi.responses import RedirectResponse
        from starlette.datastructures import URL
    except ImportError as e:
        raise ImportError('Cannot initialize OAuth to due a missing library. Please run `pip install huggingface_hub[oauth]` or add `huggingface_hub[oauth]` to your requirements.txt file.') from e
    warnings.warn('OAuth is not supported outside of a Space environment. To help you debug your app locally, the oauth endpoints are mocked to return your profile and token. To make it work, your machine must be logged in to Huggingface.')
    mocked_oauth_info = _get_mocked_oauth_info()
    login_uri, callback_uri, logout_uri = _get_oauth_uris(route_prefix)

    @app.get(login_uri)
    async def oauth_login(request: fastapi.Request) -> RedirectResponse:
        redirect_uri = _generate_redirect_uri(request)
        return RedirectResponse(callback_uri + '?' + urllib.parse.urlencode({'_target_url': redirect_uri}))

    @app.get(callback_uri)
    async def oauth_redirect_callback(request: fastapi.Request) -> RedirectResponse:
        request.session['oauth_info'] = mocked_oauth_info
        return RedirectResponse(_get_redirect_target(request))

    @app.get(logout_uri)
    async def oauth_logout(request: fastapi.Request) -> RedirectResponse:
        request.session.pop('oauth_info', None)
        logout_url = URL('/').include_query_params(**request.query_params)
        return RedirectResponse(url=logout_url, status_code=302)

def _generate_redirect_uri(request: 'fastapi.Request') -> str:
    if '_target_url' in request.query_params:
        target = request.query_params['_target_url']
    else:
        target = '/?' + urllib.parse.urlencode(request.query_params)
    redirect_uri = request.url_for('oauth_redirect_callback').include_query_params(_target_url=target)
    redirect_uri_as_str = str(redirect_uri)
    if redirect_uri.netloc.endswith('.hf.space'):
        redirect_uri_as_str = redirect_uri_as_str.replace('http://', 'https://')
    return redirect_uri_as_str

def _get_redirect_target(request: 'fastapi.Request', default_target: str='/') -> str:
    return request.query_params.get('_target_url', default_target)

def _get_mocked_oauth_info() -> dict:
    token = get_token()
    if token is None:
        raise ValueError('Your machine must be logged in to HF to debug an OAuth app locally. Please run `hf auth login` or set `HF_TOKEN` as environment variable with one of your access token. You can generate a new token in your settings page (https://huggingface.co/settings/tokens).')
    user = whoami()
    if user['type'] != 'user':
        raise ValueError('Your machine is not logged in with a personal account. Please use a personal access token. You can generate a new token in your settings page (https://huggingface.co/settings/tokens).')
    return {'access_token': token, 'token_type': 'bearer', 'expires_in': 8 * 60 * 60, 'id_token': 'FOOBAR', 'scope': 'openid profile', 'refresh_token': 'hf_oauth__refresh_token', 'expires_at': int(time.time()) + 8 * 60 * 60, 'userinfo': {'sub': '0123456789', 'name': user['fullname'], 'preferred_username': user['name'], 'profile': f"https://huggingface.co/{user['name']}", 'picture': user['avatarUrl'], 'website': '', 'aud': '00000000-0000-0000-0000-000000000000', 'auth_time': 1691672844, 'nonce': 'aaaaaaaaaaaaaaaaaaa', 'iat': 1691672844, 'exp': 1691676444, 'iss': 'https://huggingface.co'}}

def _get_oauth_uris(route_prefix: str='/') -> tuple[str, str, str]:
    route_prefix = route_prefix.strip('/')
    if route_prefix:
        route_prefix = f'/{route_prefix}'
    return (f'{route_prefix}/oauth/huggingface/login', f'{route_prefix}/oauth/huggingface/callback', f'{route_prefix}/oauth/huggingface/logout')
