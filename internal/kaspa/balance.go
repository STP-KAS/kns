package kaspa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const SompiPerKAS = 100_000_000

type Balance struct {
	Address string `json:"address"`
	Sompi   uint64 `json:"balance"`
	KAS     float64 `json:"kas"`
}

func FetchBalance(addr string) (*Balance, error) {
	u := "https://api.kaspa.org/addresses/" + addr + "/balance"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "kns-web4/1.0")
	c := &http.Client{Timeout: 12 * time.Second}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("kaspa api %s", res.Status)
	}
	var raw struct {
		Address string `json:"address"`
		Balance uint64 `json:"balance"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return &Balance{Address: raw.Address, Sompi: raw.Balance, KAS: float64(raw.Balance) / SompiPerKAS}, nil
}
