package protocol

import "testing"

func TestParseName(t *testing.T) {
	n, err := ParseName("KNS.kas")
	if err != nil {
		t.Fatal(err)
	}
	if n.String() != "kns.kas" || n.ApexLabel() != "kns" {
		t.Fatalf("got %+v", n)
	}
	if _, err := ParseName("foo.com"); err == nil {
		t.Fatal("expected tld error")
	}
}

func TestSubname(t *testing.T) {
	n, err := ParseName("pay.shop.kas")
	if err != nil {
		t.Fatal(err)
	}
	if !n.IsSubname() || n.Leaf() != "pay" || n.Parent() != "shop.kas" || n.Apex() != "shop.kas" {
		t.Fatalf("got %+v parent=%s apex=%s", n, n.Parent(), n.Apex())
	}
	if _, err := PlanRegister("pay.shop.kas", Mainnet); err == nil {
		t.Fatal("subname must not inscribe")
	}
}

func TestPrice(t *testing.T) {
	if PriceKAS("a") != 4200 || PriceKAS("ab") != 4200 {
		t.Fatal("1-2 char")
	}
	if PriceKAS("abc") != 2100 || PriceKAS("abcd") != 525 {
		t.Fatal("3-4 char")
	}
	if PriceKAS("alice") != 35 {
		t.Fatal("5+")
	}
}

func TestCreatePayload(t *testing.T) {
	b, err := CreateDomain("example")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"op":"create","p":"domain","v":"example"}` {
		t.Fatalf("payload %s", b)
	}
}

func TestClub(t *testing.T) {
	n, _ := ParseName("7.kas")
	if n.Club() != "99" {
		t.Fatal(n.Club())
	}
	n, _ = ParseName("01.kas")
	if n.Club() != "" {
		t.Fatal("leading zero numeric excluded")
	}
}

func TestPayURI(t *testing.T) {
	u := PayURI("kaspa:qqqq", "kns.kas")
	if u != "kaspa:qqqq?label=kns.kas" {
		t.Fatal(u)
	}
}

func TestLooksLikeAddress(t *testing.T) {
	if !LooksLikeAddress("kaspa:qzt9yuqceqvt2vk9dz7ddzayaa5flnenkymec59xvzm55ln3k72vgecxjhnjp") {
		t.Fatal("addr")
	}
	if LooksLikeAddress("kns.kas") {
		t.Fatal("name")
	}
}
