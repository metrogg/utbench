import json
import logging
from contextlib import AsyncExitStack
from datetime import timedelta
from pathlib import Path
from typing import TYPE_CHECKING, Any, AsyncIterable, Literal, Optional, TypedDict, Union, overload
from typing_extensions import NotRequired, TypeAlias, Unpack
from ...utils._runtime import get_hf_hub_version
from .._generated._async_client import AsyncInferenceClient
from .._generated.types import ChatCompletionInputMessage, ChatCompletionInputTool, ChatCompletionStreamOutput, ChatCompletionStreamOutputDeltaToolCall
from .._providers import PROVIDER_OR_POLICY_T
from .utils import format_result
if TYPE_CHECKING:
    from mcp import ClientSession
logger = logging.getLogger(__name__)
ToolName: TypeAlias = str
ServerType: TypeAlias = Literal['stdio', 'sse', 'http']

class StdioServerParameters_T(TypedDict):
    command: str
    args: NotRequired[list[str]]
    env: NotRequired[dict[str, str]]
    cwd: NotRequired[Union[str, Path, None]]

class SSEServerParameters_T(TypedDict):
    url: str
    headers: NotRequired[dict[str, Any]]
    timeout: NotRequired[float]
    sse_read_timeout: NotRequired[float]

class StreamableHTTPParameters_T(TypedDict):
    url: str
    headers: NotRequired[dict[str, Any]]
    timeout: NotRequired[timedelta]
    sse_read_timeout: NotRequired[timedelta]
    terminate_on_close: NotRequired[bool]

class MCPClient:

    def __init__(self, *, model: Optional[str]=None, provider: Optional[PROVIDER_OR_POLICY_T]=None, base_url: Optional[str]=None, api_key: Optional[str]=None):
        self.sessions: dict[ToolName, 'ClientSession'] = {}
        self.exit_stack = AsyncExitStack()
        self.available_tools: list[ChatCompletionInputTool] = []
        if model is None and base_url is None:
            raise ValueError('At least one of `model` or `base_url` should be set in `MCPClient`.')
        self.payload_model = model
        self.client = AsyncInferenceClient(model=None if base_url is not None else model, provider=provider, api_key=api_key, base_url=base_url)

    async def __aenter__(self):
        await self.client.__aenter__()
        await self.exit_stack.__aenter__()
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        await self.client.__aexit__(exc_type, exc_val, exc_tb)
        await self.cleanup()

    async def cleanup(self):
        await self.client.close()
        await self.exit_stack.aclose()

    @overload
    async def add_mcp_server(self, type: Literal['stdio'], **params: Unpack[StdioServerParameters_T]):
        ...

    @overload
    async def add_mcp_server(self, type: Literal['sse'], **params: Unpack[SSEServerParameters_T]):
        ...

    @overload
    async def add_mcp_server(self, type: Literal['http'], **params: Unpack[StreamableHTTPParameters_T]):
        ...

    async def add_mcp_server(self, type: ServerType, **params: Any):
        from mcp import ClientSession, StdioServerParameters
        from mcp import types as mcp_types
        allowed_tools = params.pop('allowed_tools', None)
        if type == 'stdio':
            from mcp.client.stdio import stdio_client
            logger.info(f"Connecting to stdio MCP server with command: {params['command']} {params.get('args', [])}")
            client_kwargs = {'command': params['command']}
            for key in ['args', 'env', 'cwd']:
                if params.get(key) is not None:
                    client_kwargs[key] = params[key]
            server_params = StdioServerParameters(**client_kwargs)
            read, write = await self.exit_stack.enter_async_context(stdio_client(server_params))
        elif type == 'sse':
            from mcp.client.sse import sse_client
            logger.info(f"Connecting to SSE MCP server at: {params['url']}")
            client_kwargs = {'url': params['url']}
            for key in ['headers', 'timeout', 'sse_read_timeout']:
                if params.get(key) is not None:
                    client_kwargs[key] = params[key]
            read, write = await self.exit_stack.enter_async_context(sse_client(**client_kwargs))
        elif type == 'http':
            from mcp.client.streamable_http import streamablehttp_client
            logger.info(f"Connecting to StreamableHTTP MCP server at: {params['url']}")
            client_kwargs = {'url': params['url']}
            for key in ['headers', 'timeout', 'sse_read_timeout', 'terminate_on_close']:
                if params.get(key) is not None:
                    client_kwargs[key] = params[key]
            read, write, _ = await self.exit_stack.enter_async_context(streamablehttp_client(**client_kwargs))
        else:
            raise ValueError(f'Unsupported server type: {type}')
        session = await self.exit_stack.enter_async_context(ClientSession(read_stream=read, write_stream=write, client_info=mcp_types.Implementation(name='huggingface_hub.MCPClient', version=get_hf_hub_version())))
        logger.debug('Initializing session...')
        await session.initialize()
        response = await session.list_tools()
        logger.debug('Connected to server with tools:', [tool.name for tool in response.tools])
        filtered_tools = response.tools
        if allowed_tools is not None:
            filtered_tools = [tool for tool in response.tools if tool.name in allowed_tools]
            logger.debug(f'Tool filtering applied. Using {len(filtered_tools)} of {len(response.tools)} available tools: {[tool.name for tool in filtered_tools]}')
        for tool in filtered_tools:
            if tool.name in self.sessions:
                logger.warning(f"Tool '{tool.name}' already defined by another server. Skipping.")
                continue
            self.sessions[tool.name] = session
            self.available_tools.append(ChatCompletionInputTool.parse_obj_as_instance({'type': 'function', 'function': {'name': tool.name, 'description': tool.description, 'parameters': tool.inputSchema}}))

    async def process_single_turn_with_tools(self, messages: list[Union[dict, ChatCompletionInputMessage]], exit_loop_tools: Optional[list[ChatCompletionInputTool]]=None, exit_if_first_chunk_no_tool: bool=False) -> AsyncIterable[Union[ChatCompletionStreamOutput, ChatCompletionInputMessage]]:
        tools = self.available_tools
        if exit_loop_tools is not None:
            tools = [*exit_loop_tools, *self.available_tools]
        response = await self.client.chat.completions.create(model=self.payload_model, messages=messages, tools=tools, tool_choice='auto', stream=True)
        message: dict[str, Any] = {'role': 'unknown', 'content': ''}
        final_tool_calls: dict[int, ChatCompletionStreamOutputDeltaToolCall] = {}
        num_of_chunks = 0
        async for chunk in response:
            num_of_chunks += 1
            delta = chunk.choices[0].delta if chunk.choices and len(chunk.choices) > 0 else None
            if not delta:
                continue
            if delta.role:
                message['role'] = delta.role
            if delta.content:
                message['content'] += delta.content
            if delta.tool_calls:
                for tool_call in delta.tool_calls:
                    idx = tool_call.index
                    if idx not in final_tool_calls:
                        final_tool_calls[idx] = tool_call
                        if final_tool_calls[idx].function.arguments is None:
                            final_tool_calls[idx].function.arguments = ''
                        continue
                    if final_tool_calls[idx].function.arguments is None:
                        final_tool_calls[idx].function.arguments = ''
                    if tool_call.function.arguments:
                        final_tool_calls[idx].function.arguments += tool_call.function.arguments
            if exit_if_first_chunk_no_tool and num_of_chunks <= 2 and (len(final_tool_calls) == 0):
                return
            yield chunk
        if message['content'] or final_tool_calls:
            if message.get('role') == 'unknown':
                message['role'] = 'assistant'
            if final_tool_calls:
                tool_calls_list: list[dict[str, Any]] = []
                for tc in final_tool_calls.values():
                    tool_calls_list.append({'id': tc.id, 'type': 'function', 'function': {'name': tc.function.name, 'arguments': tc.function.arguments or '{}'}})
                message['tool_calls'] = tool_calls_list
            messages.append(message)
        for tool_call in final_tool_calls.values():
            function_name = tool_call.function.name
            if function_name is None:
                message = ChatCompletionInputMessage.parse_obj_as_instance({'role': 'tool', 'tool_call_id': tool_call.id, 'content': 'Invalid tool call with no function name.'})
                messages.append(message)
                yield message
                continue
            try:
                function_args = json.loads(tool_call.function.arguments or '{}')
            except json.JSONDecodeError as err:
                tool_message = {'role': 'tool', 'tool_call_id': tool_call.id, 'name': function_name, 'content': f'Invalid JSON generated by the model: {err}'}
                tool_message_as_obj = ChatCompletionInputMessage.parse_obj_as_instance(tool_message)
                messages.append(tool_message_as_obj)
                yield tool_message_as_obj
                continue
            tool_message = {'role': 'tool', 'tool_call_id': tool_call.id, 'content': '', 'name': function_name}
            if exit_loop_tools and function_name in [t.function.name for t in exit_loop_tools]:
                tool_message_as_obj = ChatCompletionInputMessage.parse_obj_as_instance(tool_message)
                messages.append(tool_message_as_obj)
                yield tool_message_as_obj
                return
            session = self.sessions.get(function_name)
            if session is not None:
                try:
                    result = await session.call_tool(function_name, function_args)
                    tool_message['content'] = format_result(result)
                except Exception as err:
                    tool_message['content'] = f'Error: MCP tool call failed with error message: {err}'
            else:
                tool_message['content'] = f'Error: No session found for tool: {function_name}'
            tool_message_as_obj = ChatCompletionInputMessage.parse_obj_as_instance(tool_message)
            messages.append(tool_message_as_obj)
            yield tool_message_as_obj
