package goquint

import (
	"fmt"
	"regexp"
	"testing"
)

var proquintPattern = regexp.MustCompile(
	`^[bdfghjklmnpqrstvz][aiou][bdfghjklmnpqrstvz][aiou][bdfghjklmnpqrstvz]-[bdfghjklmnpqrstvz][aiou][bdfghjklmnpqrstvz][aiou][bdfghjklmnpqrstvz]$`,
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		number   uint32
		expected string
	}{
		{"zero", 0, "babab-babab"},
		{"one", 1, "babab-babad"},
		{"max uint32", 0xFFFFFFFF, "vuvuv-vuvuv"},
		{"127", 127, "babab-baduv"},
		{"256", 256, "babab-bahab"},
		{"65535", 65535, "babab-vuvuv"},
		{"mid range", 0x7F7F7F7F, "lusuv-lusuv"},
		{"powers of two 16", 1 << 16, "babad-babab"},
		{"powers of two 24", 1 << 24, "bahab-babab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.number)
			if got != tt.expected {
				t.Errorf("Encode(%d) = %v, want %v", tt.number, got, tt.expected)
			}
		})
	}
}

func TestEncodeFormat(t *testing.T) {
	// Every encoded value must match the proquint pattern
	values := []uint32{0, 1, 42, 255, 1000, 65535, 1<<20, 0xDEADBEEF, 0xFFFFFFFF}
	for _, n := range values {
		pq := Encode(n)
		if !proquintPattern.MatchString(pq) {
			t.Errorf("Encode(%d) = %v, doesn't match proquint pattern", n, pq)
		}
		if len(pq) != 11 {
			t.Errorf("Encode(%d) = %v, length %d, want 11", n, pq, len(pq))
		}
	}
}

func TestEncodeDeterminism(t *testing.T) {
	for _, n := range []uint32{0, 42, 0xDEADBEEF, 0xFFFFFFFF} {
		a := Encode(n)
		b := Encode(n)
		if a != b {
			t.Errorf("Encode(%d) not deterministic: %v != %v", n, a, b)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		proquint string
		expected uint32
		wantErr  bool
	}{
		{"zero", "babab-babab", 0, false},
		{"one", "babab-babad", 1, false},
		{"max value", "vuvuv-vuvuv", 0xFFFFFFFF, false},
		{"known value", "babab-baduv", 127, false},
		{"without hyphens", "bababbabab", 0, false},
		{"multiple hyphens", "b-a-b-a-b-b-a-b-a-b", 0, false},
		{"invalid length short", "babab", 0, true},
		{"invalid length long", "babab-babab-babab", 0, true},
		{"invalid characters", "xxxxx-xxxxx", 0, true},
		{"empty string", "", 0, true},
		{"vowel in consonant slot", "aabab-babab", 0, true},
		{"consonant in vowel slot", "bbbab-babab", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, err := Decode(tt.proquint)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode(%q) error = %v, wantErr %v", tt.proquint, err, tt.wantErr)
				return
			}
			if !tt.wantErr && num != tt.expected {
				t.Errorf("Decode(%q) = %v, want %v", tt.proquint, num, tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	testNumbers := []uint32{
		0, 1, 2, 127, 255, 256, 1000, 65535, 65536,
		1000000, 0x7FFFFFFF, 0x80000000, 0xDEADBEEF, 0xFFFFFFFF,
	}

	for _, num := range testNumbers {
		t.Run(fmt.Sprintf("%d", num), func(t *testing.T) {
			proquint := Encode(num)
			result, err := Decode(proquint)
			if err != nil {
				t.Fatalf("Decode(%q) error = %v for input %d", proquint, err, num)
			}
			if result != num {
				t.Errorf("Round trip failed: %d -> %q -> %d", num, proquint, result)
			}
		})
	}
}

func TestRoundTripExhaustive16Bit(t *testing.T) {
	// Exhaustively test all 16-bit values in the low half
	for i := 0; i < 0x10000; i++ {
		num := uint32(i)
		pq := Encode(num)
		result, err := Decode(pq)
		if err != nil {
			t.Fatalf("Decode(%q) error = %v for input %d", pq, err, num)
		}
		if result != num {
			t.Fatalf("Round trip failed: %d -> %q -> %d", num, pq, result)
		}
	}
}

func TestRandom(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		pq := Random()
		if !proquintPattern.MatchString(pq) {
			t.Errorf("Random() = %v, doesn't match proquint pattern", pq)
		}
		if seen[pq] {
			t.Errorf("Random() generated duplicate: %v", pq)
		}
		seen[pq] = true
	}
}

func TestRandomUniqueness(t *testing.T) {
	// Generate 1000 random proquints and verify no collisions
	seen := make(map[string]bool, 1000)
	for i := 0; i < 1000; i++ {
		pq := Random()
		if seen[pq] {
			t.Errorf("Random() collision after %d generations: %v", i, pq)
		}
		seen[pq] = true
	}
}

func TestRandomRoundTrips(t *testing.T) {
	for i := 0; i < 100; i++ {
		pq := Random()
		num, err := Decode(pq)
		if err != nil {
			t.Fatalf("Decode(Random()) error = %v for %q", err, pq)
		}
		pq2 := Encode(num)
		if pq != pq2 {
			t.Errorf("Random round trip failed: %q -> %d -> %q", pq, num, pq2)
		}
	}
}

func TestEncodeHex(t *testing.T) {
	tests := []struct {
		name     string
		hexStr   string
		expected string
	}{
		{
			"sha256 of hello world",
			"b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
			"qojas-fitun",
		},
		{
			"sha256 of empty json",
			"44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
			"hibig-kutog",
		},
		{
			"exactly 8 hex chars",
			"00000000",
			"babab-babab",
		},
		{
			"max 8 hex chars",
			"ffffffff",
			"vuvuv-vuvuv",
		},
		{
			"deadbeef",
			"deadbeef",
			"supos-quqov",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeHex(tt.hexStr)
			if got != tt.expected {
				t.Errorf("EncodeHex(%q) = %v, want %v", tt.hexStr, got, tt.expected)
			}
		})
	}
}

func TestEncodeHexShortFallback(t *testing.T) {
	// Short hex strings should return a valid random proquint
	shortInputs := []string{"", "a", "abc", "1234567"}
	for _, input := range shortInputs {
		pq := EncodeHex(input)
		if !proquintPattern.MatchString(pq) {
			t.Errorf("EncodeHex(%q) = %v, doesn't match proquint pattern", input, pq)
		}
	}
}

func TestEncodeHexIgnoresTrailingChars(t *testing.T) {
	// Only the first 8 hex chars matter
	a := EncodeHex("deadbeef0000")
	b := EncodeHex("deadbeefffff")
	if a != b {
		t.Errorf("EncodeHex should only use first 8 chars: %v != %v", a, b)
	}
}

// Benchmarks

func BenchmarkEncode(b *testing.B) {
	for b.Loop() {
		Encode(0xDEADBEEF)
	}
}

func BenchmarkDecode(b *testing.B) {
	for b.Loop() {
		Decode("tukos-quluv")
	}
}

func BenchmarkRandom(b *testing.B) {
	for b.Loop() {
		Random()
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	for b.Loop() {
		pq := Encode(0xDEADBEEF)
		Decode(pq)
	}
}

// Examples

func ExampleEncode() {
	fmt.Println(Encode(0))
	fmt.Println(Encode(0xDEADBEEF))
	// Output:
	// babab-babab
	// supos-quqov
}

func ExampleDecode() {
	n, _ := Decode("supos-quqov")
	fmt.Printf("%08x\n", n)
	// Output:
	// deadbeef
}

func ExampleRandom() {
	pq := Random()
	fmt.Println(proquintPattern.MatchString(pq))
	// Output:
	// true
}

func ExampleEncodeHex() {
	fmt.Println(EncodeHex("deadbeef1234"))
	// Output:
	// supos-quqov
}
