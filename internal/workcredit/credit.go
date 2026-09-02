// Package workcredit is the honest substitute for a Kaspa-native stablecoin
// when the thing dApps actually need is predictable *fees*, not a fake dollar.
//
// 1 credit = 1 gram of named-lane mass (KIP-21). Min-relay policy is 100 sompi/gram.
// Covenants can lock and conserve those grams. They cannot mint USD purchasing power.
package workcredit

import (
	"fmt"
	"math"
)

const (
	// SompiPerGramPolicy is the current min-relay policy, not a consensus constant.
	SompiPerGramPolicy uint64 = 100
	SompiPerKAS        uint64 = 100_000_000
	Unit               string = "gram"
)

// Quote is a dual-unit invoice: work in grams, settlement in KAS at policy rate.
// USD is never filled. There is no oracle in this process.
type Quote struct {
	Grams          uint64  `json:"grams"`
	Credits        uint64  `json:"credits"`
	Unit           string  `json:"unit"`
	SompiPerGram   uint64  `json:"sompiPerGram"`
	Sompi          uint64  `json:"sompi"`
	KAS            float64 `json:"kas"`
	Lane           string  `json:"lane,omitempty"`
	USD            string  `json:"usd"`
	Scheme         string  `json:"scheme"`
	NotAStablecoin string  `json:"notAStablecoin"`
	Note           string  `json:"note"`
}

func QuoteGrams(grams uint64) (Quote, error) {
	return QuoteLane(grams, "")
}

func QuoteLane(grams uint64, lane string) (Quote, error) {
	if grams == 0 {
		return Quote{}, fmt.Errorf("grams must be > 0")
	}
	if grams > math.MaxUint64/SompiPerGramPolicy {
		return Quote{}, fmt.Errorf("grams overflow sompi quote")
	}
	sompi := grams * SompiPerGramPolicy
	return Quote{
		Grams:        grams,
		Credits:      grams,
		Unit:         Unit,
		SompiPerGram: SompiPerGramPolicy,
		Sompi:        sompi,
		KAS:          float64(sompi) / float64(SompiPerKAS),
		Lane:         lane,
		USD:          "not quoted — no oracle, no reserve",
		Scheme:       "kaspa-work-credit",
		NotAStablecoin: "WorkCredit is a prepaid voucher for sequenced grams on Kaspa L1. It is not USDC, not DAI, not UST. No L2. A Kaspa L1 dollar is a different product and is not live.",
		Note: "1 credit = 1 KIP-21 gram of a named lane. Covenant tracks remaining grams. Issuer (sequencer) is the counterparty. Miners are still paid in KAS.",
	}, nil
}

// DualInvoice is how a sequenced L1 dApp should bill: work in grams, fallback in KAS.
// No L2. No dollar rail in this object.
type DualInvoice struct {
	Work     Quote  `json:"work"`
	Fallback string `json:"kasFallback"`
	Dollar   string `json:"dollarRail"`
}

func Invoice(grams uint64, lane string) (DualInvoice, error) {
	q, err := QuoteLane(grams, lane)
	if err != nil {
		return DualInvoice{}, err
	}
	return DualInvoice{
		Work:     q,
		Fallback: "pay KAS at policy sompi/gram if the user has no WorkCredit UTXO",
		Dollar:   "none on L1 today. This invoice does not use an L2. See the Kaspa Till dApp for a reserved L1 stable slot.",
	}, nil
}
