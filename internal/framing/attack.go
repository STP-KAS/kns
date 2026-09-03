// Package framing is the concrete #234 example: same 42-byte region, two framings, two amounts.
// Not a compiler. Numbers from kaspanet/silverscript#234 (supertypo, closed unmerged).
package framing

const (
	PR     = "https://github.com/kaspanet/silverscript/pull/234"
	Status = "closed-unmerged"
	Author = "supertypo"
	When   = "2026-08-29 opened; 2026-08-30 closed pending discussion"
)

// Layout is one encoding of TokenState { int amount; byte[32] owner }.
type Layout struct {
	Name    string `json:"name"`
	Hex     string `json:"hexSketch"`
	Total   int    `json:"totalBytes"`
	Note    string `json:"note"`
	Amount  string `json:"vaultReadsAmountAs"`
}

func Canonical() Layout {
	return Layout{
		Name:   "encoder (what the vault assumes)",
		Hex:    "08 <8-byte amount>  20 <32-byte owner>",
		Total:  42,
		Note:   "1-byte push headers. amount payload at bytes 1..9. owner at 10..42.",
		Amount: "the real int",
	}
}

func Attack() Layout {
	return Layout{
		Name:   "hostile instance of the SAME template",
		Hex:    "4c 08 <8-byte amount>  1f <31-byte owner>",
		Total:  42,
		Note:   "OP_PUSHDATA1 header eats one extra byte; owner push is one byte short. Total still 42. Template hash and P2SH still match.",
		Amount: "08 concatenated with 7 bytes the attacker chose — e.g. vault sees ~2^59, UTXO is worth 1",
	}
}

func Report() map[string]any {
	return map[string]any{
		"pr":     PR,
		"status": Status,
		"author": Author,
		"when":   When,
		"tokenState": "struct { int amount; byte[32] owner }",
		"canonical":  Canonical(),
		"attack":     Attack(),
		"whatPins":   "template hash + total length. What does not pin: push headers inside the region.",
		"fixInPR":    "require each field header == data_prefix(payload_len) at the encoder offset; for readInputState also P2SH(window)==input scriptPubKey",
		"ourRule":    "this stack never calls readInputState. Own UTXO validateOutputState only.",
		"safeHere":   []string{"KasName.sil", "WorkCredit.sil", "KaChatPayTimeout.sil", "KasInvoice.sil"},
		"blocked":    "any till/vault/minter that reads a KCC-20 or WorkCredit owned by someone else",
	}
}
