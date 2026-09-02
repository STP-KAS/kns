// Package protocol implements the KNS inscription format and the
// post-Toccata covenant-domain record model.
package protocol

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	ProtocolID     = "kns"
	ProtocolDomain = "domain"
	TLD            = "kas"
	MaxLabelRunes  = 255
	MaxInscription = 520 // pre-Toccata envelope limit; Toccata raises script elements to 1MB
	TextFeeKAS     = 1
)

// Official KNS reveal-output-0 payment addresses (fee model).
const (
	MainnetFeeAddress = "kaspa:qyp4nvaq3pdq7609z09fvdgwtc9c7rg07fuw5zgeee7xpr085de59eseqfcmynn"
	TN10FeeAddress    = "kaspatest:qq9h47etjv6x8jgcla0ecnp8mgrkfxm70ch3k60es5a50ypsf4h6sak3g0lru"
)

// CheckDummyAddress is used when availability is checked without a connected wallet.
// The official check endpoint requires an address field.
const CheckDummyAddress = "kaspa:qzt9yuqceqvt2vk9dz7ddzayaa5flnenkymec59xvzm55ln3k72vgecxjhnjp"

var (
	labelRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,253}[a-z0-9])?$`)
	numRe   = regexp.MustCompile(`^[0-9]+$`)
)

type Network string

const (
	Mainnet Network = "mainnet"
	TN10    Network = "tn10"
)

func FeeAddress(n Network) string {
	if n == TN10 {
		return TN10FeeAddress
	}
	return MainnetFeeAddress
}

// Name is a normalized .kas name. Subnames are parent-scoped (pay.shop.kas).
type Name struct {
	Labels []string // left-to-right, no TLD. pay.shop.kas => ["pay","shop"]
	TLD    string
}

func (n Name) Label() string {
	if len(n.Labels) == 0 {
		return ""
	}
	return strings.Join(n.Labels, ".")
}

func (n Name) Leaf() string {
	if len(n.Labels) == 0 {
		return ""
	}
	return n.Labels[0]
}

func (n Name) ApexLabel() string {
	if len(n.Labels) == 0 {
		return ""
	}
	return n.Labels[len(n.Labels)-1]
}

func (n Name) IsSubname() bool { return len(n.Labels) > 1 }

func (n Name) Parent() string {
	if !n.IsSubname() {
		return ""
	}
	return strings.Join(n.Labels[1:], ".") + "." + n.tld()
}

func (n Name) Apex() string { return n.ApexLabel() + "." + n.tld() }

func (n Name) tld() string {
	if n.TLD == "" {
		return TLD
	}
	return n.TLD
}

func (n Name) String() string {
	if len(n.Labels) == 0 {
		return ""
	}
	return n.Label() + "." + n.tld()
}

func (n Name) IsNumeric() bool { return numRe.MatchString(n.ApexLabel()) }

func (n Name) Club() string {
	if !n.IsNumeric() {
		return ""
	}
	lab := n.ApexLabel()
	if lab == "0" {
		return "99"
	}
	if strings.HasPrefix(lab, "0") {
		return ""
	}
	switch len(lab) {
	case 1, 2:
		return "99"
	case 3:
		return "999"
	case 4:
		return "10k"
	default:
		return ""
	}
}

// ParseName accepts "alice", "alice.kas", "Alice.KAS".
func ParseName(raw string) (Name, error) {
	s := strings.TrimSpace(strings.ToLower(raw))
	s = strings.TrimPrefix(s, "@")
	s = strings.TrimSuffix(s, "/")
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	s = strings.TrimSuffix(s, ".limo")
	s = strings.TrimSuffix(s, ".kas.limo")
	if strings.Contains(s, "/") {
		s = strings.SplitN(s, "/", 2)[0]
	}
	parts := strings.Split(s, ".")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return Name{}, fmt.Errorf("empty label")
		}
		parts[i] = p
	}
	tld := TLD
	if len(parts) >= 2 {
		last := parts[len(parts)-1]
		if last == TLD {
			parts = parts[:len(parts)-1]
		} else {
			return Name{}, fmt.Errorf("unsupported tld %q (only .kas)", last)
		}
	}
	if len(parts) == 0 {
		return Name{}, fmt.Errorf("empty name")
	}
	for _, p := range parts {
		if utf8.RuneCountInString(p) > MaxLabelRunes {
			return Name{}, fmt.Errorf("name longer than %d characters", MaxLabelRunes)
		}
	}
	return Name{Labels: parts, TLD: tld}, nil
}

func LooksLikeAddress(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	return strings.HasPrefix(s, "kaspa:") || strings.HasPrefix(s, "kaspatest:")
}

// PayURI is a Kaspa BIP-21-style string wallets already understand.
// Names are not kaspa: URIs (that prefix is the address encoding).
func PayURI(addr, name string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return ""
	}
	if name == "" {
		return addr
	}
	return addr + "?label=" + strings.ReplaceAll(name, " ", "")
}

func ValidInscriptionLabel(label string) bool {
	if label == "" || utf8.RuneCountInString(label) > MaxLabelRunes {
		return false
	}
	// ENS-normalize subset: ASCII LDH plus we accept unicode for emoji domains
	// as KNS does (UTS-46 + graphemer). ASCII path:
	if isASCII(label) {
		return labelRe.MatchString(label)
	}
	for _, r := range label {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

// VisualLength approximates KNS grapheme length (emoji = 1).
func VisualLength(s string) int {
	n := 0
	for _, r := range s {
		if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || r == 0x200D || r == 0xFE0F {
			continue
		}
		n++
	}
	if n == 0 {
		return utf8.RuneCountInString(s)
	}
	return n
}

func PriceKAS(label string) int {
	n := VisualLength(label)
	switch {
	case n <= 2:
		return 4200
	case n == 3:
		return 2100
	case n == 4:
		return 525
	default:
		return 35
	}
}

// Envelope is the KNS commit-reveal JSON payload.
type Envelope struct {
	Op string `json:"op"`
	P  string `json:"p,omitempty"`
	V  string `json:"v,omitempty"`
	S  string `json:"s,omitempty"`
	ID string `json:"id,omitempty"`
	To string `json:"to,omitempty"`
}

func CreateDomain(label string) ([]byte, error) {
	if !ValidInscriptionLabel(label) {
		return nil, fmt.Errorf("invalid label %q", label)
	}
	return json.Marshal(Envelope{Op: "create", P: ProtocolDomain, V: label})
}

func TransferDomain(inscriptionID, to string) ([]byte, error) {
	if inscriptionID == "" || to == "" {
		return nil, fmt.Errorf("id and to required")
	}
	return json.Marshal(Envelope{Op: "transfer", P: ProtocolDomain, ID: inscriptionID, To: to})
}

func TextInscription(text string) ([]byte, error) {
	if text == "" {
		return nil, fmt.Errorf("empty text")
	}
	return []byte(text), nil
}

// ScriptSketch documents the KNS envelope (Kasware buildScript type=KNS).
func ScriptSketch(payload []byte) string {
	return fmt.Sprintf("<xonly_pubkey> OP_CHECKSIG OP_FALSE OP_IF <%s> <0> <%s> OP_ENDIF", ProtocolID, string(payload))
}

type RegisterPlan struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	PriceKAS    int    `json:"priceKas"`
	TextFeeKAS  int    `json:"textFeeKas,omitempty"`
	FeeAddress  string `json:"feeAddress"`
	Payload     string `json:"payload"`
	Envelope    string `json:"envelope"`
	Network     string `json:"network"`
	Layer       string `json:"layer"` // inscription | covenant
	Notes       string `json:"notes"`
}

func PlanRegister(raw string, net Network) (RegisterPlan, error) {
	n, err := ParseName(raw)
	if err != nil {
		return RegisterPlan{}, err
	}
	if n.IsSubname() {
		return RegisterPlan{}, fmt.Errorf("inscription KNS cannot register %s — subnames are issued by the parent (%s), not by the global indexer", n.String(), n.Parent())
	}
	if !ValidInscriptionLabel(n.ApexLabel()) {
		return RegisterPlan{}, fmt.Errorf("label %q is not a valid inscription name", n.ApexLabel())
	}
	payload, err := CreateDomain(n.ApexLabel())
	if err != nil {
		return RegisterPlan{}, err
	}
	return RegisterPlan{
		Name:       n.String(),
		Label:      n.ApexLabel(),
		PriceKAS:   PriceKAS(n.ApexLabel()),
		FeeAddress: FeeAddress(net),
		Payload:    string(payload),
		Envelope:   ScriptSketch(payload),
		Network:    string(net),
		Layer:      "inscription",
		Notes:      "Reveal tx output 0 must pay the KNS fee address. First-come, first-served. Indexer uniqueness, not consensus. No renewal.",
	}, nil
}
