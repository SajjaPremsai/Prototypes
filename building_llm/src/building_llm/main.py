import torch

from building_llm.EmbeddingLayers import CharacterLevelEmbeddingLayer
from building_llm.EmbeddingLayers import CharacterLevelPositionalEncodingLayer
from building_llm.tokenizers.CharacterTokenzer import CharacterTokenizer


def main():

    tokenizer = CharacterTokenizer()

    text = "cat"

    print("=" * 50)
    print("Original Text")
    print("=" * 50)
    print(text)

    token_ids = tokenizer.encode(text)

    print("\n" + "=" * 50)
    print("Token IDs")
    print("=" * 50)
    print(token_ids)

    token_ids = torch.tensor(
        token_ids,
        dtype=torch.long
    )

    vocab_size = len(
        tokenizer.str_to_int
    )

    embedding_dim = 4

    embedding_layer = CharacterLevelEmbeddingLayer(
        vocab_size=vocab_size,
        embedding_dim=embedding_dim
    )

    embeddings = embedding_layer(
        token_ids
    )

    print("\n" + "=" * 50)
    print("Embeddings")
    print("=" * 50)
    print(embeddings)
    print("Shape:", embeddings.shape)

    positional_layer = CharacterLevelPositionalEncodingLayer(
        max_sequence_length=20,
        embedding_dim=embedding_dim
    )

    positional_embeddings = positional_layer(
        len(token_ids)
    )

    print("\n" + "=" * 50)
    print("Positional Embeddings")
    print("=" * 50)
    print(positional_embeddings)
    print("Shape:", positional_embeddings.shape)

    transformer_input = (
        embeddings +
        positional_embeddings
    )

    print("\n" + "=" * 50)
    print("Embedding + Position")
    print("=" * 50)
    print(transformer_input)
    print("Shape:", transformer_input.shape)


if __name__ == "__main__":
    main()