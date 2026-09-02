package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"kns/contracts"
	"kns/internal/feedback"
	"kns/internal/kachat"
	"kns/internal/wallets"
	"kns/internal/kaspa"
	"kns/internal/kasranks"
	"kns/internal/kcc"
	"kns/internal/workcredit"
)

func (s *Server) safetyPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "safety.html", page{Title: "Safety · KNS", Active: "safety"})
}

func (s *Server) feedbackPage(w http.ResponseWriter, r *http.Request) {
	p := page{Title: "Feedback · KNS", Active: "feedback"}
	if r.Method == http.MethodPost {
		n, err := feedback.Save("kns", r.FormValue("text"), r.FormValue("contact"))
		if err != nil {
			p.Error = err.Error()
		} else {
			p.Query = n.ID
		}
	}
	s.render(w, "feedback.html", p)
}

func (s *Server) apiFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]any{"ok": false, "error": "POST"})
		return
	}
	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<16))
	var req struct{ Text, Contact string }
	_ = json.Unmarshal(body, &req)
	if req.Text == "" {
		req.Text = r.FormValue("text")
		req.Contact = r.FormValue("contact")
	}
	n, err := feedback.Save("kns", req.Text, req.Contact)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "id": n.ID, "dir": feedback.Dir(), "note": "This site never DMs you."})
}

func (s *Server) walletsPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "wallets.html", page{
		Title:   "Wallets · KNS",
		Active:  "wallets",
		Wallets: wallets.All(),
	})
}

func (s *Server) kachatPage(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	p := page{
		Title:      "KaChat · KNS",
		Active:     "kachat",
		Query:      q,
		KaChatDocs: kachat.Docs,
		Envelopes:  kachat.Envelopes(q),
	}
	if q != "" {
		res, _, err := s.Res.Lookup(q)
		if err != nil {
			p.Error = err.Error()
		} else {
			p.Result = res
			c := kachat.ContactFrom(res.Name, res.Owner)
			p.Contact = &c
			p.Envelopes = kachat.Envelopes(res.Name)
		}
	}
	s.render(w, "kachat.html", p)
}

func (s *Server) kasswordPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "kassword.html", page{Title: "Kassword · KNS", Active: "kassword"})
}

func (s *Server) ranksPage(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	p := page{Title: "Ranks · KNS", Active: "ranks", Query: q}
	if q != "" {
		res, _, err := s.Res.Lookup(q)
		if err != nil {
			p.Error = err.Error()
		} else {
			p.Result = res
			if res.Owner != "" {
				if bal, err := kaspa.FetchBalance(res.Owner); err != nil {
					p.Error = err.Error()
				} else {
					rk := kasranks.FromSompi(bal.Sompi)
					p.Rank = &rk
				}
			}
		}
	}
	s.render(w, "ranks.html", p)
}

func (s *Server) silvercPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "silverc.html", page{Title: "Silverscript v1-rc1 · KNS", Active: "silverc"})
}

func (s *Server) apiArtifact(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/v1/artifact/")
	name = strings.TrimSuffix(name, ".json")
	allowed := map[string]string{
		"KasName":          "v1/KasName.json",
		"KaChatPayTimeout": "v1/KaChatPayTimeout.json",
		"WorkCredit":       "v1/WorkCredit.json",
	}
	rel, ok := allowed[name]
	if !ok {
		writeJSON(w, 404, map[string]any{"ok": false, "error": "unknown artifact"})
		return
	}
	raw, err := contracts.Artifacts.ReadFile(rel)
	if err != nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(raw)
}

func parseGrams(r *http.Request) uint64 {
	raw := strings.TrimSpace(r.URL.Query().Get("grams"))
	if raw == "" {
		raw = strings.TrimSpace(r.URL.Query().Get("credits"))
	}
	if raw == "" {
		return 1_000_000
	}
	n, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || n == 0 {
		return 1_000_000
	}
	return n
}

func (s *Server) creditsPage(w http.ResponseWriter, r *http.Request) {
	grams := parseGrams(r)
	lane := strings.TrimSpace(r.URL.Query().Get("lane"))
	q, err := workcredit.QuoteLane(grams, lane)
	p := page{Title: "Work Credits · KNS", Active: "credits", Query: strconv.FormatUint(grams, 10)}
	if err != nil {
		p.Error = err.Error()
	} else {
		p.Quote = &q
	}
	s.render(w, "credits.html", p)
}

func (s *Server) apiCreditsQuote(w http.ResponseWriter, r *http.Request) {
	grams := parseGrams(r)
	lane := strings.TrimSpace(r.URL.Query().Get("lane"))
	inv, err := workcredit.Invoice(grams, lane)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{
		"ok":     true,
		"local":  true,
		"deploy": false,
		"data":   inv,
		"map":    "C:\\Users\\Remco\\Documents\\kaspa\\superapp\\WORK-CREDITS.md",
	})
}

func (s *Server) kccPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "kcc.html", page{
		Title:    "KCC · KNS",
		Active:   "kcc",
		Drafts:   kcc.Drafts(),
		KCCLinks: kcc.Links(),
	})
}

func (s *Server) kapostsPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "kaposts.html", page{
		Title:  "KaPosts · local",
		Active: "kachat",
		Posts:  s.Posts.List(),
	})
}

func (s *Server) apiRank(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		q = r.URL.Query().Get("name")
	}
	res, _, err := s.Res.Lookup(q)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	if res.Owner == "" {
		writeJSON(w, 404, map[string]any{"ok": false, "error": "no owner"})
		return
	}
	bal, err := kaspa.FetchBalance(res.Owner)
	if err != nil {
		writeJSON(w, 502, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	rk := kasranks.FromSompi(bal.Sompi)
	writeJSON(w, 200, map[string]any{"ok": true, "name": res.Name, "owner": res.Owner, "rank": rk})
}

func (s *Server) apiKaposts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, 200, map[string]any{"ok": true, "local": true, "data": s.Posts.List(), "note": "Not KaChat's on-chain KaPosts indexer."})
		return
	}
	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<16))
	var req struct{ Name, Owner, Text string }
	_ = json.Unmarshal(body, &req)
	if req.Text == "" {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "text required"})
		return
	}
	p := s.Posts.Add(req.Name, req.Owner, req.Text)
	writeJSON(w, 200, map[string]any{"ok": true, "local": true, "data": p})
}
