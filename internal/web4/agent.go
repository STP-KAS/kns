// Package web4 builds agent discovery documents from a KNS name.
// Web 4.0 here = readable, discoverable, callable, payable (Cloudflare 2026).
// Crypto-native shorthand: Web3 + agents (x402, ERC-8004, MCP, A2A).
// Kaspa has no ERC-8004 registry and no USDC x402 rail. These files are
// the same *shape* so agents can consume them; payments are kaspa: URIs.
package web4

import (
	"strings"

	"kns/internal/protocol"
	"kns/internal/resolver"
)

type Service struct {
	Name     string   `json:"name"`
	Endpoint string   `json:"endpoint"`
	Version  string   `json:"version,omitempty"`
	Skills   []string `json:"skills,omitempty"`
}

type Card struct {
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Image        string    `json:"image,omitempty"`
	Services     []Service `json:"services"`
	X402Support  bool      `json:"x402Support"`
	X402Note     string    `json:"x402Note"`
	Active       bool      `json:"active"`
	Owner        string    `json:"owner,omitempty"`
	PayURI       string    `json:"payUri,omitempty"`
	DID          string    `json:"did,omitempty"`
	KasURI       string    `json:"kasUri,omitempty"`
	Evidence     string    `json:"evidence,omitempty"`
	Registrations []any    `json:"registrations"`
	SupportedTrust []string `json:"supportedTrust"`
}

type A2ACard struct {
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	URL          string         `json:"url"`
	Version      string         `json:"version"`
	Capabilities map[string]any `json:"capabilities"`
	Skills       []map[string]any `json:"skills"`
	DefaultInputModes []string `json:"defaultInputModes"`
	DefaultOutputModes []string `json:"defaultOutputModes"`
}

type Pay402 struct {
	X402Version int              `json:"x402Version"`
	Error       string           `json:"error"`
	Accepts     []map[string]any `json:"accepts"`
	Note        string           `json:"note"`
}

func CardFrom(res *resolver.Result, origin string) *Card {
	if res == nil {
		return nil
	}
	origin = strings.TrimRight(origin, "/")
	desc := res.Records.Bio
	if desc == "" {
		desc = res.Name + " is a Kaspa Name Service identity. Agents can resolve, pay, and fetch an agent card. Uniqueness is indexer-backed, not consensus."
	}
	c := &Card{
		Type:        "https://eips.ethereum.org/EIPS/eip-8004#registration-v1",
		Name:        res.Name,
		Description: desc,
		Image:       res.Records.Avatar,
		X402Support: res.PayURI != "",
		X402Note:    "Not Coinbase x402/USDC. HTTP 402 here asks for a Kaspa payment to the name's owner.",
		Active:      res.Owner != "",
		Owner:       res.Owner,
		PayURI:      res.PayURI,
		DID:         res.DID,
		KasURI:      res.KasURI,
		Evidence:    string(res.Evidence),
		Registrations: []any{},
		SupportedTrust: []string{},
	}
	c.Services = []Service{
		{Name: "web", Endpoint: origin + "/site/" + res.Name, Version: "1"},
		{Name: "KNS", Endpoint: res.Name, Version: "inscription"},
		{Name: "DID", Endpoint: protocol.DID(res.Name), Version: "v1"},
		{Name: "MCP", Endpoint: origin + "/mcp", Version: "2025-06-18"},
		{Name: "A2A", Endpoint: origin + "/agent/" + res.Name + "/.well-known/agent-card.json", Version: "0.3.0"},
		{Name: "pay", Endpoint: res.PayURI, Version: "kaspa"},
		{Name: "KaChat", Endpoint: origin + "/kachat?q=" + res.Name, Version: "ciph_msg:1", Skills: []string{"e2e", "pay-memo"}},
		{Name: "Kassword", Endpoint: "https://kassword.com/", Version: "toccata-vault"},
		{Name: "KasRanks", Endpoint: origin + "/ranks?q=" + res.Name, Version: "balance-depth"},
		{Name: "KCC", Endpoint: "https://github.com/kaspanet/kccs", Version: "draft"},
		{Name: "work-credits", Endpoint: origin + "/credits", Version: "gram-voucher", Skills: []string{"quote", "not-usd"}},
	}
	if res.Records.Website != "" {
		c.Services = append(c.Services, Service{Name: "website", Endpoint: res.Records.Website})
	}
	if res.Records.Redirect != "" {
		c.Services = append(c.Services, Service{Name: "redirect", Endpoint: res.Records.Redirect})
	}
	if res.Records.Email != "" {
		c.Services = append(c.Services, Service{Name: "email", Endpoint: res.Records.Email})
	}
	if res.Records.HasVault() {
		c.Services = append(c.Services, Service{Name: "vault", Endpoint: res.Records.VaultCommit, Skills: []string{"commitment-only"}})
	}
	return c
}

func A2AFrom(res *resolver.Result, origin string) *A2ACard {
	if res == nil {
		return nil
	}
	origin = strings.TrimRight(origin, "/")
	return &A2ACard{
		Name:        res.Name,
		Description: "Kaspa name. Resolve, pay, fetch records.",
		URL:         origin + "/site/" + res.Name,
		Version:     "0.3.0",
		Capabilities: map[string]any{
			"streaming": false,
			"pushNotifications": false,
		},
		DefaultInputModes:  []string{"text", "application/json"},
		DefaultOutputModes: []string{"application/json"},
		Skills: []map[string]any{
			{"id": "resolve", "name": "resolve", "description": "Return owner, pay URI, profile", "tags": []string{"kns", "kaspa"}},
			{"id": "pay", "name": "pay", "description": "Kaspa payment URI for this name", "tags": []string{"kaspa"}},
		},
	}
}

func Challenge402(res *resolver.Result, resource string) *Pay402 {
	payTo := ""
	if res != nil {
		payTo = res.Owner
	}
	return &Pay402{
		X402Version: 1,
		Error:       "Payment Required",
		Note:        "Kaspa-native 402. Not Coinbase x402/USDC. Pay KAS, or burn WorkCredit grams with the sequencer. Retry with X-Kaspa-Payment: <txid> (unchecked).",
		Accepts: []map[string]any{
			{
				"scheme":            "kaspa",
				"network":           "kaspa-mainnet",
				"payTo":             payTo,
				"maxAmountRequired": "100000000",
				"asset":             "KAS",
				"resource":          resource,
				"payUri":            protocol.PayURI(payTo, resName(res)),
			},
			{
				"scheme":            "kaspa-work-credit",
				"network":           "kaspa-mainnet",
				"asset":             "GRAM",
				"maxAmountRequired": "1000000",
				"resource":          resource,
				"note":              "1 credit = 1 KIP-21 gram. Not USD. Covenant voucher; miners still take KAS.",
			},
		},
	}
}

func resName(res *resolver.Result) string {
	if res == nil {
		return ""
	}
	return res.Name
}

func WantsJSON(rHeader string) bool {
	h := strings.ToLower(rHeader)
	return strings.Contains(h, "application/json") && !strings.Contains(h, "text/html")
}
