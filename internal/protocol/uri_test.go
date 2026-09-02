package protocol

import "testing"

func TestKasURI(t *testing.T) {
	if KasURI("KNS", "pay") != "kas://kns.kas/pay" {
		t.Fatal(KasURI("KNS", "pay"))
	}
	if DID("alice.kas") != "did:kas:alice.kas" {
		t.Fatal(DID("alice.kas"))
	}
	n, a, err := ParseKasURI("kas://pay.shop.kas/site")
	if err != nil || n != "pay.shop.kas" || a != "site" {
		t.Fatalf("%s %s %v", n, a, err)
	}
}
