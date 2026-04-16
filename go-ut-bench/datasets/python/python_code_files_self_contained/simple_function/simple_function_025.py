from collections import Counter
import re

def task_func(result):
    regex = re.compile('^(?:http|ftp)s?://(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\\.)+(?:[A-Z]{2,6}\\.?|[A-Z0-9-]{2,}\\.?)|localhost|\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})(?::\\d+)?(?:/?|[/?]\\S+)$', re.IGNORECASE)
    from_user_values = []
    for l_res in result:
        for j in l_res:
            if re.match(regex, j):
                from_user_values.append(l_res[j])
    counter = Counter(from_user_values)
    most_common = dict(counter.most_common(1))
    return most_common
