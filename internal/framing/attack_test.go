package framing

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCanonicalIs42AndVaultReadsOne(t *testing.T) {
	b := CanonicalEncode(Amount1, DemoOwner())
	if len(b) != TokenStateLen {
		t.Fatalf("len %d", len(b))
	}
	if b[0] != 0x08 || b[9] != 0x20 {
		t.Fatalf("headers %02x %02x", b[0], b[9])
	}
	n, ok := VaultAmount(b)
	if !ok || n != 1 {
		t.Fatalf("vault %d ok=%v", n, ok)
	}
	if err := PinHeaders(b); err != nil {
		t.Fatal(err)
	}
}

func TestAttackSameLengthVaultReads264(t *testing.T) {
	b := AttackEncode(Amount1, DemoOwner())
	if len(b) != TokenStateLen {
		t.Fatalf("len %d", len(b))
	}
	c := CanonicalEncode(Amount1, DemoOwner())
	if len(b) != len(c) {
		t.Fatal("attack must preserve length")
	}
	if b[0] != 0x4c || b[1] != 0x08 {
		t.Fatalf("attack header %02x %02x", b[0], b[1])
	}
	n, ok := VaultAmount(b)
	if !ok || n != VaultSeesForAmount1 {
		t.Fatalf("vault %d want %d", n, VaultSeesForAmount1)
	}
	// Explicit: window is 08 01 00 00 00 00 00 00
	want := binary.LittleEndian.Uint64([]byte{0x08, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if n != want {
		t.Fatalf("window %d want %d", n, want)
	}
	if err := PinHeaders(b); err == nil {
		t.Fatal("pin must reject the reframe")
	}
}

func TestHugeReframeIsTwoToTheFiftyNinePlusEight(t *testing.T) {
	b := AttackEncode(HugeReal, DemoOwner())
	n, ok := VaultAmount(b)
	if !ok || n != HugeVault {
		t.Fatalf("vault %d want %d", n, HugeVault)
	}
	c := CanonicalEncode(HugeReal, DemoOwner())
	got, _ := VaultAmount(c)
	if got != HugeReal {
		t.Fatalf("canonical huge %d", got)
	}
}

func TestDataPrefixWidths(t *testing.T) {
	if got := DataPrefix(8); string(got) != "\x08" {
		t.Fatalf("8: %x", got)
	}
	if got := DataPrefix(32); string(got) != "\x20" {
		t.Fatalf("32: %x", got)
	}
	if got := DataPrefix(76); string(got) != "\x4c\x4c" {
		t.Fatalf("76: %x", got)
	}
	if got := DataPrefix(256); string(got) != "\x4d\x00\x01" {
		t.Fatalf("256: %x", got)
	}
	// Pinning only the first byte of a 76-byte field would miss 4c 4b.
	wide := DataPrefix(76)
	fake := []byte{0x4c, 0x4b}
	if wide[0] != fake[0] {
		t.Fatal("first bytes should match — that is why a 1-byte pin is not enough")
	}
	if bytesEqual(wide, fake) {
		t.Fatal("full header must differ")
	}
}

func TestDecodeHexAttack(t *testing.T) {
	p, err := DecodeHex(hexOf(AttackEncode(Amount1, DemoOwner())))
	if err != nil {
		t.Fatal(err)
	}
	if p.VaultAmount != VaultSeesForAmount1 {
		t.Fatal(p.VaultAmount)
	}
	if p.PinOK {
		t.Fatal("pin")
	}
}

func TestContractSourcesDoNotReadForeignState(t *testing.T) {
	root := filepath.Join("..", "..", "contracts")
	var seenForbidden bool
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".sil") {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(b)
		if strings.Contains(info.Name(), "FORBIDDEN") {
			seenForbidden = true
			if !strings.Contains(s, "readInputState") {
				t.Errorf("%s should illustrate the forbidden call", path)
			}
			return nil
		}
		if strings.Contains(s, "readInputState") {
			t.Errorf("%s uses readInputState — forbidden on v1-rc1 (#234)", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !seenForbidden {
		t.Fatal("missing FORBIDDEN-*.sil illustration")
	}
}

func hexOf(b []byte) string {
	const hexdigits = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, x := range b {
		out[i*2] = hexdigits[x>>4]
		out[i*2+1] = hexdigits[x&0x0f]
	}
	return string(out)
}
