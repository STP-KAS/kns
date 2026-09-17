// Package wallets is a catalog of official Kaspa wallets.
// This desk does not ship in-page inject. Pay path: QR / kaspa: URI / paste txid.
package wallets

type Kind string

const (
	Hardware Kind = "hardware"
	Native   Kind = "native"
	Multi    Kind = "multi"
)

type Connect string

const (
	Inject  Connect = "inject" // withdrawn — kept so callers compile; Injected() returns nil
	Open    Connect = "open"
	Install Connect = "install"
)

type Wallet struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Kind      Kind     `json:"kind"`
	Connect   Connect  `json:"connect"`
	Inject    string   `json:"inject,omitempty"`
	URL       string   `json:"url"`
	Store     string   `json:"store,omitempty"`
	Platforms []string `json:"platforms"`
	KNS       bool     `json:"kns"`
	KRC20     bool     `json:"krc20"`
	NFT       bool     `json:"nft"`
	Note      string   `json:"note"`
	Source    string   `json:"source"`
}

func All() []Wallet {
	return []Wallet{
		{ID: "tangem", Name: "Tangem", Kind: Hardware, Connect: Install, URL: "https://tangem.com", Platforms: []string{"card", "mobile"}, Note: "Hardware. No inject from this desk.", Source: "wiki.kaspa.org/wallet"},
		{ID: "ledger", Name: "Ledger", Kind: Hardware, Connect: Open, URL: "https://kasvault.io", Platforms: []string{"Nano"}, Note: "Sign in KasVault. No inject from this desk.", Source: "kasvault.io"},
		{ID: "kaspium", Name: "Kaspium", Kind: Native, Connect: Install, URL: "https://kaspium.io", Platforms: []string{"iOS", "Android"}, Note: "Mobile. No browser inject.", Source: "kaspium.io"},
		{ID: "kasware", Name: "Kasware", Kind: Native, Connect: Open, URL: "https://www.kasware.xyz", Platforms: []string{"Chrome extension", "Android APK"}, KNS: true, Note: "Official site. This desk does not ship inject. Pay with QR / kaspa: URI / paste txid.", Source: "docs.kasware.xyz"},
		{ID: "kastle", Name: "Kastle", Kind: Native, Connect: Open, URL: "https://kastle.cc", Platforms: []string{"Chrome extension", "iOS", "Android"}, KNS: true, Note: "Official site. This desk does not ship inject.", Source: "docs.kastle.cc"},
		{ID: "kaspa-ng", Name: "Kaspa NG", Kind: Native, Connect: Open, URL: "https://kaspa-ng.org", Platforms: []string{"web", "desktop"}, Note: "Own wallet UI, not an inject on this origin.", Source: "wiki.kaspa.org/wallet"},
	}
}

func Injected() []Wallet {
	return nil
}

func ByKind(k Kind) []Wallet {
	var out []Wallet
	for _, w := range All() {
		if w.Kind == k {
			out = append(out, w)
		}
	}
	return out
}
