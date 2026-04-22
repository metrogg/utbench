import csv
import io

def task_func(filename, from_encoding='cp1251', to_encoding='utf8', delimiter=','):
    with io.open(filename, 'r', encoding=from_encoding) as file:
        content = file.read()
    content = content.encode(from_encoding).decode(to_encoding)
    file_like = io.StringIO(content)
    reader = csv.DictReader(file_like, delimiter=delimiter)
    data = list(reader)
    output = io.StringIO()
    fieldnames = reader.fieldnames if reader.fieldnames else ['Column']
    writer = csv.DictWriter(output, fieldnames=fieldnames, delimiter=delimiter)
    writer.writeheader()
    writer.writerows(data)
    converted_csv = output.getvalue().replace('\r\n', '\n')
    return (data, converted_csv)
