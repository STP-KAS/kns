// Package nameset is a local simulator of a based-ZK name set (KNS era 3/4).
// It is not L1, not a prover, not mixed into live indexer resolve.
package nameset

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"kns/internal/protocol"
)

type Record struct {
	Name        string            `json:"name"`
	Owner       string            `json:"owner"`
	Records     protocol.Records  `json:"records"`
	Parent      string            `json:"parent,omitempty"`
	Generation  int               `json:"generation"`
	Updated     time.Time         `json:"updated"`
	Commit      string            `json:"stateCommit"`
}

type Op struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Actor  string `json:"actor"`
	Detail string `json:"detail,omitempty"`
	At     time.Time `json:"at"`
}

type Set struct {
	mu   sync.RWMutex
	by   map[string]*Record
	log  []Op
	root string
}

func New() *Set {
	s := &Set{by: map[string]*Record{}}
	s.recommit()
	return s
}

func (s *Set) Root() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.root
}

func (s *Set) Get(name string) *Record {
	n, err := protocol.ParseName(name)
	if err != nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if r := s.by[n.String()]; r != nil {
		cp := *r
		return &cp
	}
	return nil
}

func (s *Set) List() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Record, 0, len(s.by))
	for _, r := range s.by {
		out = append(out, *r)
	}
	return out
}

func (s *Set) Log() []Op {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Op, len(s.log))
	copy(out, s.log)
	return out
}

func (s *Set) Register(raw, owner string) (*Record, error) {
	n, err := protocol.ParseName(raw)
	if err != nil {
		return nil, err
	}
	if n.IsSubname() {
		return nil, fmt.Errorf("use spawn for subnames; register is apex only")
	}
	if owner == "" {
		return nil, fmt.Errorf("owner required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	key := n.String()
	if _, ok := s.by[key]; ok {
		return nil, fmt.Errorf("%s taken in local name set", key)
	}
	rec := &Record{
		Name: key, Owner: owner, Generation: 1, Updated: time.Now().UTC(),
		Records: protocol.Records{KAS: owner, Pay: owner},
	}
	s.by[key] = rec
	s.note(Op{Kind: "register", Name: key, Actor: owner})
	s.recommit()
	cp := *rec
	cp.Commit = s.root
	return &cp, nil
}

func (s *Set) Transfer(raw, actor, to string) (*Record, error) {
	return s.mutate(raw, actor, func(r *Record) error {
		if to == "" {
			return fmt.Errorf("to required")
		}
		r.Owner = to
		r.Records.KAS = to
		r.Records.Pay = to
		return nil
	}, "transfer "+to)
}

func (s *Set) SetKas(raw, actor, addr string) (*Record, error) {
	return s.mutate(raw, actor, func(r *Record) error {
		r.Records.KAS = addr
		r.Records.Pay = addr
		return nil
	}, "set_kas")
}

func (s *Set) SetContent(raw, actor, content string) (*Record, error) {
	return s.mutate(raw, actor, func(r *Record) error {
		r.Records.Website = content
		r.Records.ContentHash = content
		return nil
	}, "set_content")
}

func (s *Set) BindVault(raw, actor, commit string) (*Record, error) {
	return s.mutate(raw, actor, func(r *Record) error {
		if len(commit) < 8 {
			return fmt.Errorf("commitment too short")
		}
		r.Records.VaultCommit = commit
		r.Records.Vault = commit
		return nil
	}, "bind_vault")
}

func (s *Set) Spawn(parentRaw, actor, childLabel, childOwner string) (*Record, error) {
	p, err := protocol.ParseName(parentRaw)
	if err != nil {
		return nil, err
	}
	if childLabel == "" || childOwner == "" {
		return nil, fmt.Errorf("child label and owner required")
	}
	child, err := protocol.ParseName(childLabel + "." + p.String())
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	parent := s.by[p.String()]
	if parent == nil {
		return nil, fmt.Errorf("parent %s not in local set", p.String())
	}
	if parent.Owner != actor {
		return nil, fmt.Errorf("not parent owner")
	}
	key := child.String()
	if _, ok := s.by[key]; ok {
		return nil, fmt.Errorf("%s taken", key)
	}
	rec := &Record{
		Name: key, Owner: childOwner, Parent: p.String(), Generation: 1, Updated: time.Now().UTC(),
		Records: protocol.Records{KAS: childOwner, Pay: childOwner},
	}
	s.by[key] = rec
	s.note(Op{Kind: "spawn", Name: key, Actor: actor, Detail: childOwner})
	s.recommit()
	cp := *rec
	cp.Commit = s.root
	return &cp, nil
}

func (s *Set) mutate(raw, actor string, fn func(*Record) error, detail string) (*Record, error) {
	n, err := protocol.ParseName(raw)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.by[n.String()]
	if r == nil {
		return nil, fmt.Errorf("%s not in local set", n.String())
	}
	if r.Owner != actor {
		return nil, fmt.Errorf("not owner")
	}
	if err := fn(r); err != nil {
		return nil, err
	}
	r.Generation++
	r.Updated = time.Now().UTC()
	s.note(Op{Kind: detail, Name: n.String(), Actor: actor, Detail: detail})
	s.recommit()
	cp := *r
	cp.Commit = s.root
	return &cp, nil
}

func (s *Set) note(op Op) {
	op.At = time.Now().UTC()
	s.log = append(s.log, op)
	if len(s.log) > 200 {
		s.log = s.log[len(s.log)-200:]
	}
}

func (s *Set) recommit() {
	h := sha256.New()
	for k, r := range s.by {
		h.Write([]byte(k))
		h.Write([]byte(r.Owner))
		h.Write([]byte(r.Records.KAS))
		h.Write([]byte(r.Records.ContentHash))
		h.Write([]byte(r.Records.VaultCommit))
	}
	s.root = hex.EncodeToString(h.Sum(nil))
}
