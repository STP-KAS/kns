package registrar

import (
	"testing"

	"kns/internal/protocol"
)

func TestRoutesAreHonest(t *testing.T) {
	rs := Routes()
	if len(rs) != 4 {
		t.Fatal(len(rs))
	}
	if rs[0].Status != protocol.Indexer || rs[1].Status != protocol.Roadmap {
		t.Fatal("live path is indexer; root registrar is roadmap")
	}
	if rs[3].Status != protocol.Research {
		t.Fatal("based zk is research")
	}
}
