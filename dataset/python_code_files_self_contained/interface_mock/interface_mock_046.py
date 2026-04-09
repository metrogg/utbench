import csv
from ipaddress import IPv4Network

def task_func(ip_range, csv_path):
    with open(csv_path, 'w', newline='') as csvfile:
        fieldnames = ['IP Address']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        for ip in IPv4Network(ip_range):
            writer.writerow({'IP Address': str(ip)})
    return csv_path
