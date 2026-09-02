package knsapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"kns/internal/protocol"
)

const DefaultBase = "https://api.knsdomains.org/mainnet"

type Client struct {
	Base   string
	HTTP   *http.Client
	Net    protocol.Network
}

func New(base string, net protocol.Network) *Client {
	if base == "" {
		base = DefaultBase
		if net == protocol.TN10 {
			base = "https://api.knsdomains.org/tn10"
		}
	}
	return &Client{
		Base: strings.TrimRight(base, "/"),
		Net:  net,
		HTTP: &http.Client{Timeout: 18 * time.Second},
	}
}

type Envelope[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type Owner struct {
	ID      string `json:"id"`
	AssetID string `json:"assetId"`
	Asset   string `json:"asset"`
	Owner   string `json:"owner"`
}

type Asset struct {
	ID                string `json:"id"`
	AssetID           string `json:"assetId"`
	MimeType          string `json:"mimeType"`
	Asset             string `json:"asset"`
	Owner             string `json:"owner"`
	CreationBlockTime string `json:"creationBlockTime"`
	IsDomain          bool   `json:"isDomain"`
	IsVerifiedDomain  bool   `json:"isVerifiedDomain"`
	Status            string `json:"status"`
	TransactionID     string `json:"transactionId"`
}

type AssetsPage struct {
	Assets []Asset `json:"assets"`
}

type CheckItem struct {
	Domain           string `json:"domain"`
	Available        bool   `json:"available"`
	IsReservedDomain bool   `json:"isReservedDomain"`
}

type Profile struct {
	AssetID string         `json:"assetId"`
	Owner   string         `json:"owner"`
	Name    string         `json:"name"`
	TLD     string         `json:"tld"`
	Profile map[string]any `json:"profile"`
}

type Primary struct {
	OwnerAddress  string `json:"ownerAddress"`
	InscriptionID string `json:"inscriptionId"`
	Domain        *struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		TLD        string `json:"tld"`
		FullName   string `json:"fullName"`
		IsVerified bool   `json:"isVerified"`
		Status     string `json:"status"`
	} `json:"domain"`
}

func (c *Client) get(path string, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.Base+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kns-platform/1.0")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return fmt.Errorf("kns api %s: %s %s", path, res.Status, truncate(body, 240))
	}
	return json.Unmarshal(body, out)
}

func (c *Client) post(path string, in, out any) error {
	raw, err := json.Marshal(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.Base+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kns-platform/1.0")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return fmt.Errorf("kns api %s: %s %s", path, res.Status, truncate(body, 240))
	}
	return json.Unmarshal(body, out)
}

func (c *Client) Owner(domain string) (*Owner, error) {
	var env Envelope[Owner]
	path := "/api/v1/" + domain + "/owner"
	if err := c.get(path, &env); err != nil {
		return nil, err
	}
	if !env.Success || env.Data.Owner == "" {
		return nil, fmt.Errorf("not found")
	}
	return &env.Data, nil
}

func (c *Client) Assets(query url.Values) ([]Asset, error) {
	var env Envelope[AssetsPage]
	if err := c.get("/api/v1/assets?"+query.Encode(), &env); err != nil {
		return nil, err
	}
	return env.Data.Assets, nil
}

func (c *Client) AssetByName(name string) (*Asset, error) {
	q := url.Values{}
	q.Set("asset", name)
	q.Set("type", "domain")
	q.Set("pageSize", "20")
	assets, err := c.Assets(q)
	if err != nil {
		return nil, err
	}
	for i := range assets {
		if strings.EqualFold(assets[i].Asset, name) {
			return &assets[i], nil
		}
	}
	if len(assets) > 0 {
		return &assets[0], nil
	}
	return nil, fmt.Errorf("not found")
}

func (c *Client) Check(names []string, address string) ([]CheckItem, error) {
	if address == "" {
		address = protocol.CheckDummyAddress
	}
	var env Envelope[struct {
		Domains []CheckItem `json:"domains"`
	}]
	err := c.post("/api/v1/domains/check", map[string]any{
		"domainNames": names,
		"address":     address,
	}, &env)
	if err != nil {
		return nil, err
	}
	return env.Data.Domains, nil
}

func (c *Client) Profile(assetID string) (*Profile, error) {
	var env Envelope[Profile]
	keys := "redirectUrl,bio,avatarUrl,website,x,github,telegram,discord,email,banner"
	path := "/api/v1/domain/" + url.PathEscape(assetID) + "/profile?keys=" + url.QueryEscape(keys)
	if err := c.get(path, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

func (c *Client) Primary(owner string) (*Primary, error) {
	var env Envelope[Primary]
	path := "/api/v1/primary-name/" + owner
	if err := c.get(path, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

func truncate(b []byte, n int) string {
	s := strings.TrimSpace(string(b))
	if len(s) > n {
		return s[:n]
	}
	return s
}
