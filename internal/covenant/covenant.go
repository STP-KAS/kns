package covenant

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"kns/internal/protocol"
)

// Domain is a local sketch of a Name UTXO. It is not a live mainnet object.
// KIP-20 genesis ids come from the outpoint, not from the label.
type Domain struct {
	Name         string            `json:"name"`
	CovenantID   string            `json:"covenantId"`
	Owner        string            `json:"owner"`
	Records      protocol.Records  `json:"records"`
	Subnames     []string          `json:"subnames,omitempty"`
	Parent       string            `json:"parent,omitempty"`
	ElevatedFrom string            `json:"elevatedFrom,omitempty"` // inscription id
	Generation   int               `json:"generation"`
	Updated      time.Time         `json:"updated"`
	Status       string            `json:"status"` // live | pending | demo
}

// IDFor is a local map key. It is not a KIP-20 covenant_id.
func IDFor(name string) string {
	sum := sha256.Sum256([]byte("kns.local.name|" + name))
	return hex.EncodeToString(sum[:])
}

type Store struct {
	mu   sync.RWMutex
	path string
	by   map[string]*Domain
}

func Open(path string) (*Store, error) {
	s := &Store{path: path, by: map[string]*Domain{}}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	var list []*Domain
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	for _, d := range list {
		if d == nil || d.Name == "" || d.Status == "demo" {
			continue
		}
		s.by[d.Name] = d
	}
	return s, nil
}

func (s *Store) Get(name string) *Domain {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if d := s.by[name]; d != nil {
		cp := *d
		return &cp
	}
	return nil
}

func (s *Store) Put(d *Domain) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.CovenantID == "" {
		d.CovenantID = IDFor(d.Name)
	}
	d.Updated = time.Now().UTC()
	s.by[d.Name] = d
	return s.flush()
}

func (s *Store) List() []*Domain {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Domain, 0, len(s.by))
	for _, d := range s.by {
		cp := *d
		out = append(out, &cp)
	}
	return out
}

func (s *Store) flush() error {
	list := make([]*Domain, 0, len(s.by))
	for _, d := range s.by {
		list = append(list, d)
	}
	raw, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, raw, 0o644)
}
