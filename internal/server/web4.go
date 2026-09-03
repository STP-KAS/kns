package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"kns/internal/era"
	"kns/internal/kachat"
	"kns/internal/kaspa"
	"kns/internal/kasranks"
	"kns/internal/web4"
	"kns/internal/workcredit"
)

func origin(r *http.Request) string {
	proto := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		proto = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}
	return proto + "://" + host
}

func (s *Server) web4Page(w http.ResponseWriter, r *http.Request) {
	s.render(w, "web4.html", page{Title: "Web4.0 · KNS", Active: "web4", Eras: era.All()})
}

func (s *Server) agentJSON(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		p := strings.TrimPrefix(r.URL.Path, "/agent/")
		p = strings.TrimPrefix(p, "/api/v1/agent/")
		if i := strings.Index(p, "/"); i >= 0 {
			name = p[:i]
		} else {
			name = strings.TrimSuffix(p, ".json")
		}
	}
	if name == "" {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "name required"})
		return
	}
	res, _, err := s.Res.Lookup(name)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	if strings.Contains(r.URL.Path, "agent-card") {
		writeJSON(w, 200, web4.A2AFrom(res, origin(r)))
		return
	}
	writeJSON(w, 200, web4.CardFrom(res, origin(r)))
}

func (s *Server) wellKnownAgent(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "kns.kas"
	}
	res, _, err := s.Res.Lookup(name)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, web4.CardFrom(res, origin(r)))
}

func (s *Server) call402(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/v1/call/")
	res, _, err := s.Res.Lookup(name)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	if tx := r.Header.Get("X-Kaspa-Payment"); tx != "" {
		writeJSON(w, 200, map[string]any{
			"ok": true, "paid": true, "txid": tx, "name": res.Name,
			"note": "Receipt accepted at the HTTP layer only. This does not verify the Kaspa tx.",
			"agent": web4.CardFrom(res, origin(r)),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPaymentRequired)
	_ = json.NewEncoder(w).Encode(web4.Challenge402(res, r.URL.Path))
}

func (s *Server) mcp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, 200, map[string]any{
			"name": "kns-web4",
			"version": "4.0",
			"protocol": "mcp-jsonrpc",
			"tools": mcpTools(),
		})
		return
	}
	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	var req struct {
		JSONRPC string `json:"jsonrpc"`
		ID      any    `json:"id"`
		Method  string `json:"method"`
		Params  struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"params"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, 400, map[string]any{"jsonrpc": "2.0", "error": map[string]any{"message": err.Error()}})
		return
	}
	switch req.Method {
	case "initialize":
		writeJSON(w, 200, map[string]any{
			"jsonrpc": "2.0", "id": req.ID,
			"result": map[string]any{
				"protocolVersion": "2025-06-18",
				"serverInfo":      map[string]any{"name": "kns-web4", "version": "4.0"},
				"capabilities":    map[string]any{"tools": map[string]any{}},
			},
		})
	case "tools/list":
		writeJSON(w, 200, map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": map[string]any{"tools": mcpTools()}})
	case "tools/call":
		out, err := s.mcpCall(req.Params.Name, req.Params.Arguments, origin(r))
		if err != nil {
			writeJSON(w, 200, map[string]any{"jsonrpc": "2.0", "id": req.ID, "error": map[string]any{"message": err.Error()}})
			return
		}
		raw, _ := json.Marshal(out)
		writeJSON(w, 200, map[string]any{
			"jsonrpc": "2.0", "id": req.ID,
			"result": map[string]any{"content": []map[string]any{{"type": "text", "text": string(raw)}}},
		})
	default:
		writeJSON(w, 200, map[string]any{"jsonrpc": "2.0", "id": req.ID, "error": map[string]any{"message": "unknown method"}})
	}
}

func mcpTools() []map[string]any {
	return []map[string]any{
		{"name": "resolve_kas", "description": "Resolve a .kas name or kaspa: address via the KNS indexer", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"q": map[string]any{"type": "string"}}, "required": []string{"q"}}},
		{"name": "agent_card", "description": "ERC-8004-shaped agent registration file for a .kas name", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}}, "required": []string{"name"}}},
		{"name": "pay_kas", "description": "Kaspa payment URI for a name (not Coinbase x402)", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}}, "required": []string{"name"}}},
		{"name": "kachat_contact", "description": "Resolve a .kas name to a KaChat contact address. Does not encrypt.", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}}, "required": []string{"name"}}},
		{"name": "kas_rank", "description": "Ocean-depth rank from live KAS balance (not the KASRANKS NFT)", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}}, "required": []string{"name"}}},
		{"name": "quote_work", "description": "Invoice sequenced work in grams (Work Credits). Not a USD stable. Optional lane.", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"grams": map[string]any{"type": "string"}, "lane": map[string]any{"type": "string"}}}},
	}
}

func (s *Server) mcpCall(tool string, args map[string]any, orig string) (any, error) {
	get := func(k string) string {
		if args == nil {
			return ""
		}
		if v, ok := args[k]; ok && v != nil {
			return strings.TrimSpace(fmt.Sprint(v))
		}
		return ""
	}
	if tool == "quote_work" {
		grams := uint64(1_000_000)
		if g := get("grams"); g != "" {
			if n, err := strconv.ParseUint(g, 10, 64); err == nil && n > 0 {
				grams = n
			}
		}
		inv, err := workcredit.Invoice(grams, get("lane"))
		if err != nil {
			return nil, err
		}
		return inv, nil
	}
	q := get("q")
	if q == "" {
		q = get("name")
	}
	res, _, err := s.Res.Lookup(q)
	if err != nil {
		return nil, err
	}
	switch tool {
	case "agent_card":
		return web4.CardFrom(res, orig), nil
	case "pay_kas":
		return map[string]any{"payUri": res.PayURI, "owner": res.Owner, "name": res.Name}, nil
	case "kachat_contact":
		return map[string]any{"contact": kachat.ContactFrom(res.Name, res.Owner), "envelopes": kachat.Envelopes(res.Name), "app": kachat.Linktree}, nil
	case "kas_rank":
		if res.Owner == "" {
			return nil, fmt.Errorf("no owner")
		}
		bal, err := kaspa.FetchBalance(res.Owner)
		if err != nil {
			return nil, err
		}
		r := kasranks.FromSompi(bal.Sompi)
		return r, nil
	default:
		return res, nil
	}
}
