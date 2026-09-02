// Package registrar is the uniqueness design. Read this before writing
// "consensus enforces the name".
package registrar

import "kns/internal/protocol"

// KIP-20 genesis (kaspanotes / KIP-20):
//
//	covenant_id = BLAKE2b("CovenantID",
//	    outpoint.tx_id || le32(index) || authorized outputs…)
//
// That id is unique because the outpoint is unique. It does not contain the
// human label. Anyone can genesis another covenant and write "alice.kas" in
// its state. Nodes will accept both.
//
// A script cannot scan the UTXO set for "is alice taken?". That is why KCC-721
// puts ticker uniqueness on the indexer, and why inscription KNS does too.

type Path string

const (
	IndexerFCFS Path = "indexer-fcfs" // live KNS
	RootUTXO    Path = "root-registrar-utxo"
	ParentIssue Path = "parent-issued-subnames"
	BasedZK     Path = "based-zk-name-set"
)

type Route struct {
	Path        Path   `json:"path"`
	Status      protocol.Evidence `json:"status"`
	Uniqueness  string `json:"uniqueness"`
	Parallelism string `json:"parallelism"`
	Cost        string `json:"cost"`
}

func Routes() []Route {
	return []Route{
		{
			Path:        IndexerFCFS,
			Status:      protocol.Indexer,
			Uniqueness:  "First valid kns envelope in canonical order. Reorgs can theoretically flip a race; in practice KNS has run this since Jan 2025.",
			Parallelism: "Full. Anyone can reveal at once; indexer serializes the name.",
			Cost:        "You trust KNS's indexer (or run one). Consensus will not reject a second alice.kas inscription; the indexer will ignore it.",
		},
		{
			Path:        RootUTXO,
			Status:      protocol.Roadmap,
			Uniqueness:  "One singleton .kas registrar UTXO. register(label) is only valid if the successor state adds that label (or a sparse commitment to it) and pays the fee. Only this lineage may spawn Name UTXOs.",
			Parallelism: "Registrations queue on one UTXO. At 10 BPS that is a product constraint, not a physics one. Lookups stay parallel: each issued Name is its own UTXO.",
			Cost:        "Needs a deployed, reviewed Silverscript registrar and a migration for 50k inscriptions (snapshot + claim, or live indexer attestation). Not optional if you want consensus uniqueness for the global namespace.",
		},
		{
			Path:        ParentIssue,
			Status:      protocol.Roadmap,
			Uniqueness:  "pay.shop.kas is unique because only shop.kas can spawn the child. No global scan.",
			Parallelism: "Every parent is its own UTXO track. This is the Toccata-native scale path (partition state).",
			Cost:        "Requires the parent to already exist (inscription today, Name UTXO later). Do not pretend inscription KNS has this.",
		},
		{
			Path:        BasedZK,
			Status:      protocol.Research,
			Uniqueness:  "Off-chain name set, user ops on a KIP-21 lane, settlement covenant verifies a proof that the name was free.",
			Parallelism: "High. This is how you get ENS-scale without a single L1 registrar UTXO.",
			Cost:        "vProgs / proving stack are not a product. Shipping this copy as live is a lie.",
		},
	}
}
