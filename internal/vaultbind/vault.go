package vaultbind

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Binding attaches a kaspa-data-vault package (or Toccata DataVault covenant)
// to a .kas name. File bytes stay off-chain; only the commitment is public.
type Binding struct {
	Name        string    `json:"name"`
	PackageID   string    `json:"packageId"`
	Commitment  string    `json:"commitment"`
	UseCase     string    `json:"useCase"`
	Label       string    `json:"label"`
	WindowSecs  int       `json:"windowSecs"`
	CovenantID  string    `json:"covenantId,omitempty"`
	Status      string    `json:"status"` // sealed | sent | expired | demo
	Updated     time.Time `json:"updated"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	by   map[string][]Binding
}

func Open(path string) (*Store, error) {
	s := &Store{path: path, by: map[string][]Binding{}}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(filepath.Dir(path), 0o755)
			return s, nil
		}
		return nil, err
	}
	var list []Binding
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	for _, b := range list {
		s.by[b.Name] = append(s.by[b.Name], b)
	}
	return s, nil
}

func (s *Store) ForName(name string) []Binding {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.by[name]
	out := make([]Binding, len(src))
	copy(out, src)
	return out
}

func (s *Store) All() []Binding {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Binding
	for _, list := range s.by {
		out = append(out, list...)
	}
	return out
}

func (s *Store) Bind(b Binding) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if b.WindowSecs == 0 {
		b.WindowSecs = 300
	}
	b.Updated = time.Now().UTC()
	s.by[b.Name] = append(s.by[b.Name], b)
	return s.flush()
}

func (s *Store) flush() error {
	var list []Binding
	for _, bs := range s.by {
		list = append(list, bs...)
	}
	raw, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, raw, 0o644)
}

