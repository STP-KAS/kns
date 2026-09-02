package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"kns/internal/era"
	"kns/internal/protocol"
)

func (s *Server) eraPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "eras.html", page{Title: "KNS 4.0 eras", Active: "4", Eras: era.All()})
}

func (s *Server) mePage(w http.ResponseWriter, r *http.Request) {
	addr := strings.TrimSpace(r.URL.Query().Get("address"))
	p := page{Title: "My names · KNS", Active: "me", Query: addr}
	if addr != "" && s.Res.API != nil {
		q := url.Values{}
		q.Set("owner", addr)
		q.Set("type", "domain")
		q.Set("pageSize", "100")
		if assets, err := s.Res.API.Assets(q); err != nil {
			p.Error = err.Error()
		} else {
			p.Assets = assets
		}
		if prim, err := s.Res.Reverse(addr); err == nil {
			p.Primary = prim
		}
	}
	s.render(w, "me.html", p)
}

func (s *Server) simPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "sim.html", page{
		Title:  "Name set sim · KNS",
		Active: "sim",
		Sim:    s.Set.List(),
		Ops:    s.Set.Log(),
		Root:   s.Set.Root(),
	})
}

func (s *Server) sdkPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "sdk.html", page{Title: "SDK · KNS", Active: "sdk"})
}

func (s *Server) directoryPage(w http.ResponseWriter, r *http.Request) {
	p := page{Title: "Directory · KNS", Active: "directory"}
	club := r.URL.Query().Get("club")
	if club == "" {
		club = "99"
	}
	if s.Res.API != nil {
		q := url.Values{}
		q.Set("type", "domain")
		q.Set("collection", club)
		q.Set("pageSize", "40")
		if assets, err := s.Res.API.Assets(q); err != nil {
			p.Error = err.Error()
		} else {
			p.Assets = assets
			p.Query = club
		}
	}
	s.render(w, "directory.html", p)
}

func (s *Server) wellKnown(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "kns.kas"
	}
	res, _, err := s.Res.Lookup(name)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, protocol.WellKnown(res.Name, res.Owner, res.PayURI))
}

func (s *Server) apiNames(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("owner")
	if addr == "" {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "owner required"})
		return
	}
	q := url.Values{}
	q.Set("owner", addr)
	q.Set("type", "domain")
	q.Set("pageSize", "100")
	assets, err := s.Res.API.Assets(q)
	if err != nil {
		writeJSON(w, 502, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "data": assets})
}

func (s *Server) apiDirectory(w http.ResponseWriter, r *http.Request) {
	club := r.URL.Query().Get("club")
	if club == "" {
		club = "99"
	}
	q := url.Values{}
	q.Set("type", "domain")
	q.Set("collection", club)
	q.Set("pageSize", "40")
	assets, err := s.Res.API.Assets(q)
	if err != nil {
		writeJSON(w, 502, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "data": assets})
}

func (s *Server) apiBatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Names []string `json:"names"`
	}
	if r.Method == http.MethodGet {
		req.Names = strings.Split(r.URL.Query().Get("names"), ",")
	} else {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<16))
		_ = json.Unmarshal(body, &req)
	}
	if len(req.Names) > 25 {
		req.Names = req.Names[:25]
	}
	var out []any
	for _, n := range req.Names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		res, _, err := s.Res.Lookup(n)
		if err != nil {
			out = append(out, map[string]any{"name": n, "error": err.Error()})
			continue
		}
		out = append(out, res)
	}
	writeJSON(w, 200, map[string]any{"ok": true, "data": out})
}

func (s *Server) apiSim(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, 200, map[string]any{
			"ok": true, "local": true, "root": s.Set.Root(),
			"names": s.Set.List(), "log": s.Set.Log(),
			"warning": "in-process map, not a ZK proof",
		})
		return
	}
	var req struct {
		Op, Name, Owner, Actor, To, Addr, Content, Commit, Child, ChildOwner string
	}
	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<16))
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	var (
		rec any
		err error
	)
	switch req.Op {
	case "register":
		rec, err = s.Set.Register(req.Name, req.Owner)
	case "transfer":
		rec, err = s.Set.Transfer(req.Name, req.Actor, req.To)
	case "set_kas":
		rec, err = s.Set.SetKas(req.Name, req.Actor, req.Addr)
	case "set_content":
		rec, err = s.Set.SetContent(req.Name, req.Actor, req.Content)
	case "bind_vault":
		rec, err = s.Set.BindVault(req.Name, req.Actor, req.Commit)
	case "spawn":
		rec, err = s.Set.Spawn(req.Name, req.Actor, req.Child, req.ChildOwner)
	default:
		writeJSON(w, 400, map[string]any{"ok": false, "error": "op=register|transfer|set_kas|set_content|bind_vault|spawn"})
		return
	}
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "local": true, "root": s.Set.Root(), "data": rec})
}

func (s *Server) apiSimName(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/v4/sim/")
	rec := s.Set.Get(name)
	if rec == nil {
		writeJSON(w, 404, map[string]any{"ok": false, "error": "not in local set"})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "local": true, "data": rec})
}

func (s *Server) openapi(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"openapi": "3.0.0",
		"info": map[string]any{
			"title":   "KNS 4.0",
			"version": "4.0.0",
			"description": "Live indexer resolve plus a local based-name-set simulator. Simulator is not L1.",
		},
		"paths": map[string]any{
			"/api/v1/resolve":    map[string]any{"get": map[string]any{"summary": "Resolve name or kaspa address"}},
			"/api/v1/reverse":    map[string]any{"get": map[string]any{"summary": "Primary name for address"}},
			"/api/v1/names":      map[string]any{"get": map[string]any{"summary": "Domains owned by address"}},
			"/api/v1/batch":      map[string]any{"post": map[string]any{"summary": "Resolve up to 25 names"}},
			"/api/v1/claims":     map[string]any{"get": map[string]any{"summary": "Honest claims table"}},
			"/api/v4/sim":        map[string]any{"get": map[string]any{"summary": "Local name set"}, "post": map[string]any{"summary": "Mutate local name set"}},
			"/static/js/kns.js":  map[string]any{"get": map[string]any{"summary": "Browser SDK"}},
		},
	})
}
