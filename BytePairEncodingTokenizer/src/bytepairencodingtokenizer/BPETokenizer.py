import json
import os
import re
from typing import Dict, List, Tuple

class BPETokenizer:
    GPT2_SPLIT_PATTERN = re.compile(
        r"'s|'t|'re|'ve|'m|'ll|'d| ?\w+| ?[^\s\w]+",
        flags=re.UNICODE,
    )

    def __init__(self, vocab_size: int):
        self.vocab_size = vocab_size
        self.merges: Dict[Tuple[int, int], int] = {}
        self.merge_order: List[Tuple[int, int]] = []
        self.id_to_bytes: Dict[int, bytes] = {i: bytes([i]) for i in range(256)}
        self.bytes_to_id: Dict[bytes, int] = {bytes([i]): i for i in range(256)}
        self.next_id = 256

    @classmethod
    def train(cls, file_path: str, vocab_size: int) -> "BPETokenizer":
        tokenizer = cls(vocab_size)
        return tokenizer._train_from_file(file_path)

    @classmethod
    def from_file(cls, file_path: str, vocab_size: int) -> "BPETokenizer":
        return cls.train(file_path, vocab_size)

    @classmethod
    def load(cls, bpe_path: str, vocab_path: str) -> "BPETokenizer":
        tokenizer = cls(vocab_size=0)
        tokenizer.merges = {}
        tokenizer.merge_order = []
        tokenizer.id_to_bytes = {}
        tokenizer.bytes_to_id = {}

        with open(vocab_path, 'r', encoding='utf-8') as f:
            token_to_id = json.load(f)

        for token, token_id in token_to_id.items():
            raw = token.encode('latin-1')
            tokenizer.id_to_bytes[token_id] = raw
            tokenizer.bytes_to_id[raw] = token_id

        tokenizer.next_id = max(tokenizer.id_to_bytes) + 1
        tokenizer.vocab_size = len(tokenizer.id_to_bytes)

        with open(bpe_path, 'r', encoding='utf-8') as f:
            for line in f:
                stripped = line.strip()
                if not stripped:
                    continue
                a, b = stripped.split()
                pair = (int(a), int(b))
                tokenizer.merges[pair] = len(tokenizer.merges) + 256
                tokenizer.merge_order.append(pair)

        return tokenizer

    @staticmethod
    def _get_stats(token_ids: List[List[int]]) -> Dict[Tuple[int, int], int]:
        stats: Dict[Tuple[int, int], int] = {}
        for ids in token_ids:
            for i in range(len(ids) - 1):
                pair = (ids[i], ids[i + 1])
                stats[pair] = stats.get(pair, 0) + 1
        return stats

    @staticmethod
    def _top_pair(stats: Dict[Tuple[int, int], int]) -> Tuple[int, int]:
        return max(stats, key=stats.get)

    def _merge(self, ids: List[int], pair: Tuple[int, int]) -> List[int]:
        new_ids: List[int] = []
        i = 0
        while i < len(ids):
            if i < len(ids) - 1 and (ids[i], ids[i + 1]) == pair:
                new_ids.append(self.merges[pair])
                i += 2
            else:
                new_ids.append(ids[i])
                i += 1
        return new_ids

    @classmethod
    def _gpt2_split(cls, text: str) -> List[str]:
        return cls.GPT2_SPLIT_PATTERN.findall(text)

    def _read_file(self, file_path: str) -> List[List[int]]:
        with open(file_path, 'rb') as f:
            raw = f.read()
        text = raw.decode('utf-8', errors='replace')
        return [self._bytes_to_ids(token.encode('utf-8')) for token in self._gpt2_split(text)]

    def _bytes_to_ids(self, data: bytes) -> List[int]:
        return [self.bytes_to_id[bytes([b])] for b in data]

    def _train_from_file(self, file_path: str) -> "BPETokenizer":
        token_ids = self._read_file(file_path)
        while len(self.merges) < self.vocab_size - 256:
            stats = self._get_stats(token_ids)
            if not stats:
                break
            pair = self._top_pair(stats)
            self.merges[pair] = self.next_id
            self.merge_order.append(pair)
            token_ids = [self._merge(ids, pair) for ids in token_ids]
            self.id_to_bytes[self.next_id] = self.id_to_bytes[pair[0]] + self.id_to_bytes[pair[1]]
            self.bytes_to_id[self.id_to_bytes[self.next_id]] = self.next_id
            self.next_id += 1
        return self

    def save(self, output_prefix: str) -> None:
        bpe_path = f"{output_prefix}.bpe"
        vocab_path = f"{output_prefix}.vocab"

        os.makedirs(os.path.dirname(output_prefix) or '.', exist_ok=True)

        with open(bpe_path, 'w', encoding='utf-8') as f:
            for pair in self.merge_order:
                f.write(f"{pair[0]} {pair[1]}\n")

        token_to_id = {
            self.id_to_bytes[token_id].decode('latin-1'): token_id
            for token_id in sorted(self.id_to_bytes)
        }
        with open(vocab_path, 'w', encoding='utf-8') as f:
            json.dump(token_to_id, f, ensure_ascii=False, indent=2)

    def encode(self, text: str) -> List[int]:
        encoded: List[int] = []
        for token in self._gpt2_split(text):
            token_ids = self._bytes_to_ids(token.encode('utf-8'))
            for pair in self.merge_order:
                token_ids = self._merge(token_ids, pair)
            encoded.extend(token_ids)
        return encoded

    def decode(self, ids: List[int]) -> str:
        bytes_seq = b''.join(self.id_to_bytes[id] for id in ids)
        return bytes_seq.decode('utf-8')