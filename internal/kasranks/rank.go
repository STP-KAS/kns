package kasranks

import "kns/internal/kaspa"

// Rank is inspired by KasRanks' 11-depth ocean order (kasranks.com):
// the bag ranks you, not a waitlist. This is NOT a lookup of the KASRANKS NFT
// collection. It is a live KAS balance → depth mapping for a resolved name.

type Rank struct {
	Depth int     `json:"depth"`
	Title string  `json:"title"`
	Glyph string  `json:"glyph"`
	KAS   float64 `json:"kas"`
	Sompi uint64  `json:"sompi"`
	Note  string  `json:"note"`
}

func FromSompi(sompi uint64) Rank {
	kas := float64(sompi) / kaspa.SompiPerKAS
	r := Rank{KAS: kas, Sompi: sompi, Note: "Live api.kaspa.org balance. Not the KASRANKS NFT token id."}
	switch {
	case kas >= 1_000_000:
		r.Depth, r.Title, r.Glyph = 1, "Aquaman", "🧜"
	case kas >= 250_000:
		r.Depth, r.Title, r.Glyph = 2, "Humpback", "🐋"
	case kas >= 50_000:
		r.Depth, r.Title, r.Glyph = 3, "Killer Whale", "🐋"
	case kas >= 10_000:
		r.Depth, r.Title, r.Glyph = 4, "Shark", "🦈"
	case kas >= 2_500:
		r.Depth, r.Title, r.Glyph = 5, "Dolphin", "🐬"
	case kas >= 500:
		r.Depth, r.Title, r.Glyph = 6, "Fish", "🐟"
	case kas >= 100:
		r.Depth, r.Title, r.Glyph = 7, "Crab", "🦀"
	case kas >= 25:
		r.Depth, r.Title, r.Glyph = 8, "Shrimp", "🦐"
	case kas >= 5:
		r.Depth, r.Title, r.Glyph = 9, "Krill", "🦐"
	case kas >= 1:
		r.Depth, r.Title, r.Glyph = 10, "Algae", "🌿"
	default:
		r.Depth, r.Title, r.Glyph = 11, "Plankton", "🦠"
	}
	return r
}
