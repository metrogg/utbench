import xml.etree.ElementTree as ET
import csv

def task_func(xml_content, output_csv_path):
    try:
        root = ET.fromstring(xml_content)
        data = [[elem.tag, elem.text] for elem in root.iter()]
        with open(output_csv_path, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f)
            writer.writerows(data)
    except ET.ParseError as e:
        raise ET.ParseError(f'Error parsing XML: {e}') from e
    except IOError as e:
        raise IOError(f'Error writing CSV file: {e}') from e
