import random
import re

def task_func(target_words, n_sentences, vocabulary):
    if n_sentences < 0:
        raise ValueError('n_sentences cannot be negative.')
    if not vocabulary:
        raise ValueError('Vocabulary cannot be empty.')
    sentences = []
    for _ in range(n_sentences):
        sentence = ' '.join(random.choices(vocabulary, k=10))
        for word in target_words:
            pattern = re.compile(re.escape(word), re.IGNORECASE)
            sentence = pattern.sub(word.replace(' ', '_'), sentence)
        sentences.append(sentence.lower())
    return sentences
