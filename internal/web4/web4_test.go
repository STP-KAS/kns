package web4

import (
	"testing"

	"kns/internal/resolver"
)

func TestCardFrom(t *testing.T) {
	c := CardFrom(&resolver.Result{Name: "kns.kas", Owner: "kaspa:q", PayURI: "kaspa:q?label=kns.kas"}, "http://localhost:8080")
	if c.Name != "kns.kas" || len(c.Services) < 4 || !c.X402Support {
		t.Fatalf("%+v", c)
	}
	found := false
	for _, s := range c.Services {
		if s.Name == "work-credits" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing work-credits service")
	}
}

func TestChallenge402HasGramScheme(t *testing.T) {
	ch := Challenge402(&resolver.Result{Name: "kns.kas", Owner: "kaspa:q"}, "/api/v1/call/kns.kas")
	if len(ch.Accepts) != 2 {
		t.Fatalf("accepts=%d", len(ch.Accepts))
	}
	if ch.Accepts[0]["scheme"] != "kaspa" || ch.Accepts[1]["scheme"] != "kaspa-work-credit" {
		t.Fatalf("%+v", ch.Accepts)
	}
}
