import torch
import torch.nn as nn
import torch.nn.functional as F

class MultiHeadAttention(nn.Module):
    def __init__(self, emb_size, num_heads):
        super(MultiHeadAttention, self).__init__()
        self.num_heads = num_heads
        self.emb_size = emb_size
        self.head_dim = emb_size // num_heads

        assert self.head_dim * num_heads == self.emb_size, "emb_size must be divisible by num_heads"

        self.keys = nn.Linear(emb_size, emb_size)
        self.queries = nn.Linear(emb_size, emb_size)
        self.values = nn.Linear(emb_size, emb_size)
        self.fc_out = nn.Linear(emb_size, emb_size)

    def forward(self, key, query, value, mask=None):
        N = query.shape[0]
        key_len, query_len, value_len = key.shape[1], query.shape[1], value.shape[1]

        # Split the embedding into self.num_heads different pieces
        keys = self.keys(key).view(N, key_len, self.num_heads, self.head_dim)
        queries = self.queries(query).view(N, query_len, self.num_heads, self.head_dim)
        values = self.values(value).view(N, value_len, self.num_heads, self.head_dim)

        keys = keys.transpose(1, 2)  # (N, heads, key_len, head_dim)
        queries = queries.transpose(1, 2)  # (N, heads, query_len, head_dim)
        values = values.transpose(1, 2)  # (N, heads, value_len, head_dim)

        # Scaled dot-product attention
        energy = torch.einsum("nqhd,nkhd->nhqk", [queries, keys])
        if mask is not None:
            energy = energy.masked_fill(mask == 0, float("-1e20"))

        attention = torch.softmax(energy / (self.emb_size ** (1 / 2)), dim=3)
        out = torch.einsum("nhql,nlhd->nqhd", [attention, values]).reshape(N, query_len, self.emb_size)

        return self.fc_out(out)
class TransformerEncoderLayer(nn.Module):
    def __init__(self, emb_size, num_heads, dropout, forward_expansion):
        super(TransformerEncoderLayer, self).__init__()
        self.attention = MultiHeadAttention(emb_size, num_heads)
        self.norm1 = nn.LayerNorm(emb_size)
        self.norm2 = nn.LayerNorm(emb_size)

        self.feed_forward = nn.Sequential(
            nn.Linear(emb_size, forward_expansion * emb_size),
            nn.ReLU(),
            nn.Linear(forward_expansion * emb_size, emb_size),
        )

        self.dropout = nn.Dropout(dropout)

    def forward(self, x, mask):
        attention = self.attention(x, x, x, mask)
        x = self.norm1(attention + x)
        x = self.dropout(x)
        forward = self.feed_forward(x)
        out = self.norm2(forward + x)
        out = self.dropout(out)
        return out

class TransformerDecoderLayer(nn.Module):
    def __init__(self, emb_size, num_heads, forward_expansion, dropout):
        super(TransformerDecoderLayer, self).__init__()
        self.attention = MultiHeadAttention(emb_size, num_heads)
        self.norm1 = nn.LayerNorm(emb_size)
        self.norm2 = nn.LayerNorm(emb_size)
        self.norm3 = nn.LayerNorm(emb_size)

        self.feed_forward = nn.Sequential(
            nn.Linear(emb_size, forward_expansion * emb_size),
            nn.ReLU(),
            nn.Linear(forward_expansion * emb_size, emb_size),
        )

        self.dropout = nn.Dropout(dropout)

    def forward(self, x, value, key, src_mask, trg_mask):
        attention = self.attention(x, x, x, trg_mask)
        x = self.dropout(self.norm1(attention + x))
        attention = self.attention(x, key, value, src_mask)
        x = self.dropout(self.norm2(attention + x))
        forward = self.feed_forward(x)
        out = self.dropout(self.norm3(forward + x))
        return out
