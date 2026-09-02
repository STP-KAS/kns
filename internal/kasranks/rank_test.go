package kasranks

import "testing"

func TestFromSompi(t *testing.T) {
	if FromSompi(0).Title != "Plankton" {
		t.Fatal("zero")
	}
	if FromSompi(27 * 100_000_000).Title != "Shrimp" {
		t.Fatal("27 kas")
	}
	if FromSompi(2_000_000 * 100_000_000).Depth != 1 {
		t.Fatal("whale")
	}
}
