package kaposts

import (
	"sync"
	"time"
)

// Board is a local KaPosts-shaped feed. KaChat 4.0 puts likes/comments/follows
// on Kaspa tied to KNS identities. This board is in-process only.

type Post struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Owner  string    `json:"owner"`
	Text   string    `json:"text"`
	At     time.Time `json:"at"`
	Local  bool      `json:"local"`
}

type Board struct {
	mu    sync.Mutex
	posts []Post
	next  int
}

func New() *Board { return &Board{next: 1} }

func (b *Board) List() []Post {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Post, len(b.posts))
	copy(out, b.posts)
	return out
}

func (b *Board) Add(name, owner, text string) Post {
	b.mu.Lock()
	defer b.mu.Unlock()
	p := Post{ID: b.next, Name: name, Owner: owner, Text: text, At: time.Now().UTC(), Local: true}
	b.next++
	b.posts = append([]Post{p}, b.posts...)
	if len(b.posts) > 100 {
		b.posts = b.posts[:100]
	}
	return p
}
