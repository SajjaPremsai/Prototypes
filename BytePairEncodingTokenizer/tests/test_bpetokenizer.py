import os
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SRC_DIR = ROOT / "src"
sys.path.insert(0, str(SRC_DIR))

from bytepairencodingtokenizer.BPETokenizer import BPETokenizer

DATA_FILE = ROOT / "data" / "sample.txt"


def test_train_returns_tokenizer_instance():
    tokenizer = BPETokenizer.train(str(DATA_FILE), vocab_size=300)
    assert isinstance(tokenizer, BPETokenizer)
    assert tokenizer.merge_order, "merge_order should be populated after training"


def test_encode_decode_roundtrip_after_training():
    tokenizer = BPETokenizer.train(str(DATA_FILE), vocab_size=300)
    text = "Hello, world!"
    encoded = tokenizer.encode(text)

    assert isinstance(encoded, list)
    assert encoded, "Encoded result should not be empty"
    assert all(isinstance(token_id, int) for token_id in encoded)

    decoded = tokenizer.decode(encoded)
    assert decoded == text


def test_save_and_load_preserves_encoding(tmp_path):
    tokenizer = BPETokenizer.train(str(DATA_FILE), vocab_size=300)
    output_prefix = tmp_path / "tokenizer"
    tokenizer.save(str(output_prefix))

    bpe_path = output_prefix.with_suffix(".bpe")
    vocab_path = output_prefix.with_suffix(".vocab")

    assert bpe_path.exists()
    assert vocab_path.exists()

    loaded = BPETokenizer.load(str(bpe_path), str(vocab_path))
    assert isinstance(loaded, BPETokenizer)

    text = "Hello, world!"
    assert loaded.encode(text) == tokenizer.encode(text)
