from collections import defaultdict
from operator import itemgetter
from itertools import groupby

def task_func(news_articles):
    if any((not sorted(dic.keys()) == ['category', 'id', 'title', 'title_url'] for dic in news_articles)):
        raise ValueError("input dictionaries must contain the following keys: 'category', 'id', 'title', 'title_url'")
    news_articles.sort(key=itemgetter('category', 'title'))
    grouped_articles = defaultdict(list)
    for category, group in groupby(news_articles, key=itemgetter('category')):
        grouped_articles[category] = list(group)
    return grouped_articles
