package kachat

import "testing"

func TestEnvelopes(t *testing.T) {
	es := Envelopes("kns.kas")
	if len(es) < 4 {
		t.Fatal(len(es))
	}
	if ContactFrom("kns.kas", "kaspa:q").Address != "kaspa:q" {
		t.Fatal("contact")
	}
}
