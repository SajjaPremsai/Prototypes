class CharacterTokenizer:

    def __init__(self):
        chars = list("abcdefghijklmnopqrstuvwxyz")
        vocab = ["<PAD>", "<EOS>"] + chars
        self.str_to_int = {
            ch: idx
            for idx, ch in enumerate(vocab)
        }
        self.int_to_str = {
            idx: ch
            for ch, idx in self.str_to_int.items()
        }
        self.eos_token = "<EOS>"
        self.eos_id = self.str_to_int[self.eos_token]

    def encode(self, text):
        text = text.lower()
        ids = [
            self.str_to_int[ch]
            for ch in text
            if ch in self.str_to_int
        ]
        ids.append(self.eos_id)
        return ids

    def decode(self, ids):
        chars = []
        for token_id in ids:
            token = self.int_to_str[token_id]
            if token == "<EOS>":
                break
            if token == "<PAD>":
                continue
            chars.append(token)
        return "".join(chars)