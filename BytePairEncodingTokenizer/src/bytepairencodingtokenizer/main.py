


from bytepairencodingtokenizer import BPETokenizer


tokenizer = BPETokenizer(vocab_size=1000)
tokenizer.train('input.txt')