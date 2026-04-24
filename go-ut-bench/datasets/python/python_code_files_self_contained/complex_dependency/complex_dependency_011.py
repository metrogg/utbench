import re
import os

def task_func(request):
    match = re.match('^GET /([\\w\\.\\-]+) HTTP/1\\.1$', request)
    if match:
        file_name = match.group(1)
        if os.path.exists(file_name):
            try:
                with open(file_name, 'rb') as file:
                    content = file.read()
                    response = f"HTTP/1.1 200 OK\r\nContent-Length: {len(content)}\r\n\r\n{content.decode('utf-8')}"
            except Exception:
                response = 'HTTP/1.1 500 INTERNAL SERVER ERROR\r\n\r\nInternal Server Error'
        else:
            response = 'HTTP/1.1 404 NOT FOUND\r\n\r\nFile Not Found'
    else:
        response = 'HTTP/1.1 400 BAD REQUEST\r\n\r\nBad Request'
    return response
