package workcredit

import "testing"

func TestQuoteMillionGramsIsOneKAS(t *testing.T) {
	q, err := QuoteGrams(1_000_000)
	if err != nil {
		t.Fatal(err)
	}
	if q.Credits != 1_000_000 || q.Sompi != 100_000_000 || q.KAS != 1 {
		t.Fatalf("%+v", q)
	}
	if q.USD == "" || q.Scheme != "kaspa-work-credit" {
		t.Fatalf("missing honesty fields: %+v", q)
	}
}

func TestQuoteRejectsZero(t *testing.T) {
	if _, err := QuoteGrams(0); err == nil {
		t.Fatal("expected error")
	}
}

func TestInvoiceNamesTheDollarRail(t *testing.T) {
	inv, err := Invoice(50_000, "KNS1")
	if err != nil {
		t.Fatal(err)
	}
	if inv.Work.Lane != "KNS1" || inv.Work.Sompi != 5_000_000 {
		t.Fatalf("%+v", inv)
	}
	if inv.Dollar == "" {
		t.Fatal("expected explicit no-dollar note")
	}
}
