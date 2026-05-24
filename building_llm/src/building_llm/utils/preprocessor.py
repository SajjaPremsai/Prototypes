import re

def preprocess_text(text):
    preprocessed_text = re.split(r'([,.:;?_!"()\']|--|\s)',text)
    all_words = sorted(set(preprocessed_text))
    vocab = {s:i for i, s in enumerate(all_words)}
    return preprocessed_text, vocab