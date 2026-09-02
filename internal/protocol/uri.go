package protocol

import (
	"net/url"
	"strings"
)

// ProposedLane is a KIP-21 user-lane namespace for a based name set.
// User lanes are 4-byte ASCII + 16 zero bytes. Not allocated on mainnet.
const ProposedLane = "KNS1"

func KasURI(name string, action string) string {
	n, err := ParseName(name)
	if err != nil {
		return ""
	}
	u := "kas://" + n.String()
	if action != "" {
		u += "/" + strings.Trim(action, "/")
	}
	return u
}

func DID(name string) string {
	n, err := ParseName(name)
	if err != nil {
		return ""
	}
	return "did:kas:" + n.String()
}

func ParseKasURI(raw string) (name string, action string, err error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "web+")
	if strings.HasPrefix(s, "kas://") {
		s = strings.TrimPrefix(s, "kas://")
	}
	parts := strings.SplitN(s, "/", 2)
	n, err := ParseName(parts[0])
	if err != nil {
		return "", "", err
	}
	if len(parts) == 2 {
		action = parts[1]
	}
	return n.String(), action, nil
}

func WellKnown(name, owner, payURI string) map[string]any {
	return map[string]any{
		"name":    name,
		"did":     DID(name),
		"kas":     KasURI(name, ""),
		"owner":   owner,
		"pay":     payURI,
		"lane":    ProposedLane,
		"version": "4.0",
		"note":    "kas:// and did:kas are conventions of this app. They are not ICANN or W3C registered.",
	}
}

func QueryEscapeName(name string) string { return url.QueryEscape(name) }
