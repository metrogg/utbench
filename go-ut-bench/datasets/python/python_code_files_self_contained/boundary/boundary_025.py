import bisect
import statistics

def task_func(df, column, value):
    if column not in df.columns:
        raise ValueError(f"Column '{column}' does not exist in DataFrame")
    if not isinstance(value, (int, float)):
        raise ValueError('Value must be a number')
    data = df[column].values
    avg = statistics.mean(data)
    greater_avg = data[data > avg]
    data.sort()
    bpoint = bisect.bisect_right(data, value)
    num_greater_value = len(data) - bpoint
    ax = df.hist(column=column, bins=10)[0][0]
    return (greater_avg, num_greater_value, ax)
