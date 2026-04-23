import csv
import collections
import operator

def task_func(csv_file_path):
    with open(csv_file_path, 'r') as f:
        reader = csv.reader(f)
        next(reader)
        sales_data = collections.defaultdict(int)
        for row in reader:
            product, quantity = (row[0], int(row[1]))
            sales_data[product] += quantity
    top_selling_product = max(sales_data.items(), key=operator.itemgetter(1))[0]
    return top_selling_product
