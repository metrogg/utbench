import random
import re
WORD_LIST = ['sample', 'text', 'contains', 'several', 'words', 'including']

def task_func(n_sentences):
    sentences = []
    for _ in range(n_sentences):
        sentence_len = random.randint(5, 10)
        sentence = ' '.join((random.choice(WORD_LIST) for _ in range(sentence_len))) + '.'
        sentences.append(sentence)
    text = ' '.join(sentences)
    text = re.sub('[^\\w\\s.]', '', text).lower()
    text = re.sub('\\s+\\.', '.', text)
    text = re.sub('\\s+', ' ', text)
    return text.strip()
