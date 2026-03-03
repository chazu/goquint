# goquint

Go library for encoding and decoding [proquints](https://arxiv.org/html/0901.4016) — pronounceable representations of 32-bit and 64-bit integers.

A proquint like `lusab-babad` encodes a `uint32` as two 5-character quintuplets of alternating consonants and vowels (CVCVC), separated by a hyphen. Four quintuplets encode a `uint64`.

## Install

```
go get github.com/chazu/goquint
```

## Usage

```go
import "github.com/chazu/goquint"

// Encode a uint32
goquint.Encode(0)          // "babab-babab"
goquint.Encode(0xDEADBEEF) // "supos-quqov"

// Decode back to uint32
n, err := goquint.Decode("supos-quqov") // 0xDEADBEEF

// Generate a random proquint (crypto/rand)
pq := goquint.Random() // e.g. "kisof-baduh"

// Encode from a hex string (uses first 32 bits)
goquint.EncodeHex("b94d27b9934d3e08...") // "qojas-fitun"

// 64-bit support
goquint.Encode64(0xFFFFFFFFFFFFFFFF) // "vuvuv-vuvuv-vuvuv-vuvuv"
n64, err := goquint.Decode64("vuvuv-vuvuv-vuvuv-vuvuv")

pq64 := goquint.Random64()                  // e.g. "kisof-baduh-lumab-disov"
goquint.EncodeHex64("b94d27b9934d3e08...") // uses first 64 bits
```

## API

| Function | Description |
|---|---|
| `Encode(n uint32) string` | Convert a uint32 to a proquint |
| `Decode(s string) (uint32, error)` | Convert a proquint to a uint32 |
| `Random() string` | Generate a cryptographically random proquint |
| `EncodeHex(hex string) string` | Encode the first 32 bits of a hex string |
| `Encode64(n uint64) string` | Convert a uint64 to a four-quintuplet proquint |
| `Decode64(s string) (uint64, error)` | Convert a four-quintuplet proquint to a uint64 |
| `Random64() string` | Generate a cryptographically random 64-bit proquint |
| `EncodeHex64(hex string) string` | Encode the first 64 bits of a hex string |

## Encoding

Each 16-bit half is encoded as 5 characters: `CVCVC`

- **C** — consonant (4 bits, 16 values): `b d f g h j k l m n p q r s t v z`
- **V** — vowel (2 bits, 4 values): `a i o u`

Two quintuplets joined by `-` encode a 32-bit value. Four quintuplets encode a 64-bit value.
