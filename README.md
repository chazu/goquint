# goquint

Go library for encoding and decoding [proquints](https://arxiv.org/html/0901.4016) — pronounceable representations of 32-bit integers.

A proquint like `lusab-babad` encodes a `uint32` as two 5-character quintuplets of alternating consonants and vowels (CVCVC), separated by a hyphen.

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
```

## API

| Function | Description |
|---|---|
| `Encode(n uint32) string` | Convert a uint32 to a proquint |
| `Decode(s string) (uint32, error)` | Convert a proquint to a uint32 |
| `Random() string` | Generate a cryptographically random proquint |
| `EncodeHex(hex string) string` | Encode the first 32 bits of a hex string |

## Encoding

Each 16-bit half is encoded as 5 characters: `CVCVC`

- **C** — consonant (4 bits, 16 values): `b d f g h j k l m n p q r s t v z`
- **V** — vowel (2 bits, 4 values): `a i o u`

Two quintuplets joined by `-` encode a full 32-bit value.
