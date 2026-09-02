package nameset

import "testing"

func TestRegisterAndSpawn(t *testing.T) {
	s := New()
	if _, err := s.Register("shop.kas", "kaspa:aaa"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Register("shop.kas", "kaspa:bbb"); err == nil {
		t.Fatal("double register")
	}
	child, err := s.Spawn("shop.kas", "kaspa:aaa", "pay", "kaspa:ccc")
	if err != nil {
		t.Fatal(err)
	}
	if child.Name != "pay.shop.kas" || child.Parent != "shop.kas" {
		t.Fatalf("%+v", child)
	}
	if _, err := s.Spawn("shop.kas", "kaspa:zzz", "x", "kaspa:ccc"); err == nil {
		t.Fatal("non-owner spawn")
	}
}
