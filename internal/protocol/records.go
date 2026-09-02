package protocol

import "strings"

// RecordKey is the post-Toccata resolver key space.
// Inscription KNS already stores a subset as profile text inscriptions.
// Covenant domains store the full map in UTXO state.
const (
	KeyKAS         = "kas"
	KeyIgra        = "igra" // Kaspa EVM L2 (Igra)
	KeyBTC         = "btc"
	KeyETH         = "eth"
	KeyContent     = "contenthash"
	KeyWebsite     = "website"
	KeyRedirect    = "redirectUrl"
	KeyAvatar      = "avatarUrl"
	KeyBanner      = "banner"
	KeyBio         = "bio"
	KeyX           = "x"
	KeyGitHub      = "github"
	KeyTelegram    = "telegram"
	KeyDiscord     = "discord"
	KeyEmail       = "email"
	KeyVault       = "vault"
	KeyVaultCommit = "vaultCommit"
	KeyKFS         = "kfs" // Kaspa File Storage v3.4.5+
	KeyIPFS        = "ipfs"
	KeyArweave     = "arweave"
	KeyPay         = "pay"
	KeyPrimary     = "primary"
)

type Records struct {
	KAS         string `json:"kas,omitempty"`
	Igra        string `json:"igra,omitempty"`
	BTC         string `json:"btc,omitempty"`
	ETH         string `json:"eth,omitempty"`
	ContentHash string `json:"contenthash,omitempty"`
	Website     string `json:"website,omitempty"`
	Redirect    string `json:"redirectUrl,omitempty"`
	Avatar      string `json:"avatarUrl,omitempty"`
	Banner      string `json:"banner,omitempty"`
	Bio         string `json:"bio,omitempty"`
	X           string `json:"x,omitempty"`
	GitHub      string `json:"github,omitempty"`
	Telegram    string `json:"telegram,omitempty"`
	Discord     string `json:"discord,omitempty"`
	Email       string `json:"email,omitempty"`
	Vault       string `json:"vault,omitempty"`       // vault covenant id or package id
	VaultCommit string `json:"vaultCommit,omitempty"` // 32-byte hex commitment
	KFS         string `json:"kfs,omitempty"`         // kfs:<txid or file id>
	IPFS        string `json:"ipfs,omitempty"`
	Arweave     string `json:"arweave,omitempty"`
	Pay         string `json:"pay,omitempty"`
}

func (r Records) WebsiteTarget() string {
	for _, v := range []string{r.Redirect, r.Website, r.IPFS, r.KFS, r.Arweave, r.ContentHash} {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func (r Records) HasVault() bool {
	return r.Vault != "" || r.VaultCommit != ""
}

func (r Records) HasSite() bool { return r.WebsiteTarget() != "" }

func ContentKind(target string) string {
	t := strings.ToLower(strings.TrimSpace(target))
	switch {
	case strings.HasPrefix(t, "ipfs://"), strings.HasPrefix(t, "ipns://"):
		return "ipfs"
	case strings.HasPrefix(t, "ar://"):
		return "arweave"
	case strings.HasPrefix(t, "kfs:"):
		return "kfs"
	case strings.HasPrefix(t, "bzz://"):
		return "swarm"
	case strings.HasPrefix(t, "http://"), strings.HasPrefix(t, "https://"):
		return "https"
	default:
		return "unknown"
	}
}

func IPFSGatewayURL(target string) string {
	t := strings.TrimSpace(target)
	t = strings.TrimPrefix(t, "ipfs://")
	t = strings.TrimPrefix(t, "ipns://")
	if t == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(target), "ipns://") {
		return "https://ipfs.io/ipns/" + t
	}
	return "https://ipfs.io/ipfs/" + t
}
