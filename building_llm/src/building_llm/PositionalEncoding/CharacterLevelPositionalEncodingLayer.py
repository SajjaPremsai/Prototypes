import torch
import torch.nn as nn

class CharacterLevelPositionalEncodingLayer(nn.Module):

    def __init__(
        self,
        max_sequence_length,
        embedding_dim
    ):
        super().__init__()

        self.position_embedding = nn.Embedding(
            max_sequence_length,
            embedding_dim
        )

    def forward(self, sequence_length):

        positions = torch.arange(
            sequence_length,
            dtype=torch.long
        )

        return self.position_embedding(
            positions
        )