package resolver

import (
	"fmt"
	"strings"

	"kns/internal/knsapi"
	"kns/internal/protocol"
)

type Layer string

const (
	LayerNone        Layer = "none"
	LayerInscription Layer = "inscription"
)

type Result struct {
	Name        string           `json:"name"`
	Available   bool             `json:"available"`
	Reserved    bool             `json:"reserved"`
	PriceKAS    int              `json:"priceKas"`
	Club        string           `json:"club,omitempty"`
	Layer       Layer            `json:"layer"`
	Evidence    protocol.Evidence `json:"evidence"`
	Owner       string           `json:"owner,omitempty"`
	AssetID     string           `json:"assetId,omitempty"`
	TxID        string           `json:"txId,omitempty"`
	Created     string           `json:"created,omitempty"`
	Verified    bool             `json:"verified"`
	Primary     bool             `json:"primary"`
	Subname     bool             `json:"subname,omitempty"`
	Parent      string           `json:"parent,omitempty"`
	Records     protocol.Records `json:"records"`
	PayURI      string           `json:"payUri,omitempty"`
	KasURI      string           `json:"kasUri,omitempty"`
	DID         string           `json:"did,omitempty"`
	GatewayPath string           `json:"gatewayPath"`
	LimoURL     string           `json:"limoUrl"`
	ExplorerURL string           `json:"explorerUrl,omitempty"`
	SiteKind    string           `json:"siteKind"` // profile | redirect | ipfs | none
	Warning     string           `json:"warning,omitempty"`
}

type Resolver struct {
	API *knsapi.Client
}

func (r *Resolver) Lookup(raw string) (*Result, *knsapi.Primary, error) {
	raw = strings.TrimSpace(raw)
	if protocol.LooksLikeAddress(raw) {
		p, err := r.Reverse(raw)
		if err != nil {
			return nil, nil, err
		}
		if p == nil || p.Domain == nil || p.Domain.FullName == "" {
			return nil, p, fmt.Errorf("no primary .kas name for that address")
		}
		res, err := r.Resolve(p.Domain.FullName)
		return res, p, err
	}
	res, err := r.Resolve(raw)
	return res, nil, err
}

func (r *Resolver) Reverse(addr string) (*knsapi.Primary, error) {
	if r.API == nil {
		return nil, fmt.Errorf("no indexer")
	}
	return r.API.Primary(strings.TrimSpace(addr))
}

func (r *Resolver) Resolve(raw string) (*Result, error) {
	n, err := protocol.ParseName(raw)
	if err != nil {
		return nil, err
	}
	full := n.String()
	out := &Result{
		Name:        full,
		PriceKAS:    protocol.PriceKAS(n.ApexLabel()),
		Club:        n.Club(),
		Layer:       LayerNone,
		Evidence:    protocol.Live,
		Available:   true,
		Subname:     n.IsSubname(),
		Parent:      n.Parent(),
		KasURI:      protocol.KasURI(full, ""),
		DID:         protocol.DID(full),
		GatewayPath: "/site/" + full,
		LimoURL:     "https://" + strings.TrimSuffix(full, ".kas") + ".kas.limo",
		SiteKind:    "none",
	}
	if n.IsSubname() {
		out.Evidence = protocol.Roadmap
		out.Warning = "Subnames are not an inscription KNS feature. This lookup checks the parent apex on the indexer and sketches the parent-issued design."
		parent, err := r.Resolve(n.Parent())
		if err != nil {
			return out, nil
		}
		out.Available = false
		out.Reserved = parent.Reserved
		if parent.Owner != "" {
			out.Warning = "Parent " + n.Parent() + " is taken on the indexer (owner below). A Name UTXO for that parent could issue " + full + ". Inscription KNS cannot. That issuance is not live."
			out.Owner = parent.Owner
		} else {
			out.Warning = "Parent " + n.Parent() + " is free on the indexer. Inscribe the parent first. Subnames are not a KNS inscription op."
		}
		return out, nil
	}

	if r.API == nil {
		return out, nil
	}
	if checks, err := r.API.Check([]string{full}, ""); err == nil && len(checks) > 0 {
		out.Available = checks[0].Available
		out.Reserved = checks[0].IsReservedDomain
		out.Evidence = protocol.Indexer
	}
	own, err := r.API.Owner(full)
	if err != nil || own == nil {
		return out, nil
	}
	out.Available = false
	out.Owner = own.Owner
	out.AssetID = own.AssetID
	out.Layer = LayerInscription
	out.Evidence = protocol.Indexer
	out.Records.KAS = own.Owner
	out.Records.Pay = own.Owner
	out.PayURI = protocol.PayURI(own.Owner, full)
	if ast, err := r.API.AssetByName(full); err == nil && ast != nil {
		out.Created = ast.CreationBlockTime
		out.Verified = ast.IsVerifiedDomain
		out.TxID = ast.TransactionID
		if ast.AssetID != "" {
			out.AssetID = ast.AssetID
		}
	}
	if out.AssetID != "" {
		if p, err := r.API.Profile(out.AssetID); err == nil && p != nil {
			applyProfile(&out.Records, p.Profile)
		}
	}
	if prim, err := r.API.Primary(own.Owner); err == nil && prim != nil && prim.Domain != nil {
		out.Primary = strings.EqualFold(prim.Domain.FullName, full)
	}
	if out.TxID != "" {
		out.ExplorerURL = "https://explorer.kaspa.org/txs/" + out.TxID
	}
	out.SiteKind = siteKind(out)
	return out, nil
}

func siteKind(out *Result) string {
	t := out.Records.WebsiteTarget()
	if t == "" {
		if out.Owner != "" {
			return "profile"
		}
		return "none"
	}
	switch protocol.ContentKind(t) {
	case "ipfs":
		return "ipfs"
	case "https":
		return "redirect"
	default:
		return "profile"
	}
}

func applyProfile(rec *protocol.Records, m map[string]any) {
	if m == nil {
		return
	}
	get := func(k string) string {
		v, ok := m[k]
		if !ok || v == nil {
			return ""
		}
		s := strings.TrimSpace(fmt.Sprint(v))
		if s == "" || s == "<nil>" {
			return ""
		}
		return s
	}
	rec.Website = get("website")
	rec.Redirect = get("redirectUrl")
	rec.Avatar = get("avatarUrl")
	rec.Banner = get("banner")
	rec.Bio = get("bio")
	rec.X = get("x")
	rec.GitHub = get("github")
	rec.Telegram = get("telegram")
	rec.Discord = get("discord")
	rec.Email = get("email")
}
