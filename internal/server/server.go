package server

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"kns/internal/era"
	"kns/internal/framing"
	"kns/internal/kachat"
	"kns/internal/kaposts"
	"kns/internal/kasranks"
	"kns/internal/kcc"
	"kns/internal/knsapi"
	"kns/internal/nameset"
	"kns/internal/protocol"
	"kns/internal/registrar"
	"kns/internal/resolver"
	"kns/internal/wallets"
	"kns/internal/web4"
	"kns/internal/workcredit"
	"kns/web"
)

type Server struct {
	Addr  string
	Res   *resolver.Resolver
	Set   *nameset.Set
	Posts *kaposts.Board
	T     *template.Template
}

type page struct {
	Title      string
	Active     string
	Query      string
	Result     *resolver.Result
	Primary    *knsapi.Primary
	Plan       *protocol.RegisterPlan
	Claims     []protocol.Claim
	Routes     []registrar.Route
	Eras       []era.Era
	Assets     []knsapi.Asset
	Sim        []nameset.Record
	Ops        []nameset.Op
	Root       string
	Error      string
	Rank       *kasranks.Rank
	Envelopes  []kachat.Envelope
	Contact    *kachat.Contact
	Drafts     []kcc.Draft
	KCCLinks   map[string]string
	Posts      []kaposts.Post
	KaChatDocs string
	Quote      *workcredit.Quote
	Wallets    []wallets.Wallet
	Framing    *framing.View
}

func New(addr string, res *resolver.Resolver, set *nameset.Set) (*Server, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"short": func(s string) string {
			if len(s) <= 22 {
				return s
			}
			return s[:10] + "…" + s[len(s)-8:]
		},
		"kind": protocol.ContentKind,
		"xurl": func(h string) string {
			h = strings.TrimPrefix(strings.TrimSpace(h), "@")
			if h == "" {
				return ""
			}
			return "https://x.com/" + h
		},
	}).ParseFS(web.Templates, "templates/*.html")
	if err != nil {
		return nil, err
	}
	if set == nil {
		set = nameset.New()
	}
	return &Server{Addr: addr, Res: res, Set: set, Posts: kaposts.New(), T: t}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	static, err := fs.Sub(web.Static, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, static, "robots.txt")
	})
	mux.HandleFunc("/ai.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, static, "ai.txt")
	})
	mux.HandleFunc("/", s.home)
	mux.HandleFunc("/app", s.appPage)
	mux.HandleFunc("/name/", s.namePage)
	mux.HandleFunc("/register", s.registerPage)
	mux.HandleFunc("/vaults", s.vaultsPage)
	mux.HandleFunc("/covenant", s.covenantPage)
	mux.HandleFunc("/docs", s.docsPage)
	mux.HandleFunc("/honest", s.honestPage)
	mux.HandleFunc("/gateway", s.gatewayPage)
	mux.HandleFunc("/ecosystem", s.ecosystemPage)
	mux.HandleFunc("/kachat", s.kachatPage)
	mux.HandleFunc("/kassword", s.kasswordPage)
	mux.HandleFunc("/ranks", s.ranksPage)
	mux.HandleFunc("/kcc", s.kccPage)
	mux.HandleFunc("/silverc", s.silvercPage)
	mux.HandleFunc("/wallets", s.walletsPage)
	mux.HandleFunc("/safety", s.safetyPage)
	mux.HandleFunc("/idea", s.ideaPage)
	mux.HandleFunc("/explain", s.ideaPage)
	mux.HandleFunc("/why", s.whyPage)
	mux.HandleFunc("/234", s.framingPage)
	mux.HandleFunc("/framing", s.framingPage)
	mux.HandleFunc("/api/v1/framing", func(w http.ResponseWriter, r *http.Request) {
		v := framing.Demo()
		if hx := strings.TrimSpace(r.URL.Query().Get("hex")); hx != "" {
			got, err := framing.DecodeHex(hx)
			if err != nil {
				writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
				return
			}
			v.Custom = &got
		}
		writeJSON(w, 200, map[string]any{"ok": true, "data": v, "numbers": framing.Report()["numbers"]})
	})
	mux.HandleFunc("/feedback", s.feedbackPage)
	mux.HandleFunc("/api/v1/feedback", s.apiFeedback)
	mux.HandleFunc("/api/v1/wallets", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{
			"ok":     true,
			"inject": []string{"kasware", "kastle"},
			"ledger": "https://kasvault.io",
			"source": "https://wiki.kaspa.org/wallet",
			"data":   wallets.All(),
			"note":   "Only kasware and kastle inject a Kaspa L1 provider. Ledger uses KasVault or Kastle. No L2 required.",
		})
	})
	mux.HandleFunc("/credits", s.creditsPage)
	mux.HandleFunc("/api/v1/credits/quote", s.apiCreditsQuote)
	mux.HandleFunc("/api/v1/artifact/", s.apiArtifact)
	mux.HandleFunc("/kaposts", s.kapostsPage)
	mux.HandleFunc("/api/v1/rank", s.apiRank)
	mux.HandleFunc("/api/v1/kaposts", s.apiKaposts)
	mux.HandleFunc("/web4", s.web4Page)
	mux.HandleFunc("/4", s.web4Page)
	mux.HandleFunc("/mcp", s.mcp)
	mux.HandleFunc("/.well-known/agent.json", s.wellKnownAgent)
	mux.HandleFunc("/.well-known/agent-card.json", s.wellKnownAgent)
	mux.HandleFunc("/agent/", s.agentJSON)
	mux.HandleFunc("/api/v1/agent", s.agentJSON)
	mux.HandleFunc("/api/v1/call/", s.call402)
	mux.HandleFunc("/me", s.mePage)
	mux.HandleFunc("/sim", s.simPage)
	mux.HandleFunc("/sdk", s.sdkPage)
	mux.HandleFunc("/directory", s.directoryPage)
	mux.HandleFunc("/site/", s.site)
	mux.HandleFunc("/.well-known/kaspa-name.json", s.wellKnown)
	mux.HandleFunc("/api/v1/resolve", s.apiResolve)
	mux.HandleFunc("/api/v1/batch", s.apiBatch)
	mux.HandleFunc("/api/v1/names", s.apiNames)
	mux.HandleFunc("/api/v1/directory", s.apiDirectory)
	mux.HandleFunc("/api/v4/sim", s.apiSim)
	mux.HandleFunc("/api/v4/sim/", s.apiSimName)
	mux.HandleFunc("/api/openapi.json", s.openapi)
	mux.HandleFunc("/api/v1/check", s.apiCheck)
	mux.HandleFunc("/api/v1/plan", s.apiPlan)
	mux.HandleFunc("/api/v1/suggest", s.apiSuggest)
	mux.HandleFunc("/api/v1/reverse", s.apiReverse)
	mux.HandleFunc("/api/v1/claims", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"ok": true, "data": protocol.Claims()})
	})
	mux.HandleFunc("/llms.txt", s.llms)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{
			"ok":           true,
			"product":      "kns-web4",
			"web4":         "readable-discoverable-callable-payable",
			"kns":          "inscription-indexer",
			"toccata":      "live-consensus-not-this-app",
			"x402":         "kaspa-native-402-not-usdc",
			"erc8004":      "card-shape-only-no-registry",
			"overlays":     false,
			"simRoot":      s.Set.Root(),
			"workCredits":  "gram-voucher-not-usd",
			"sompiPerGram": 100,
		})
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/static/") {
			log.Printf("%s %s", r.Method, r.URL.Path)
		}
		mux.ServeHTTP(w, r)
	})
}

func (s *Server) render(w http.ResponseWriter, name string, p page) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.T.ExecuteTemplate(w, name, p); err != nil {
		log.Println("template", name, err)
		http.Error(w, err.Error(), 500)
	}
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	s.render(w, "home.html", page{Title: "KNS — Web4.0 names on Kaspa", Active: "home", Eras: era.All()})
}

func (s *Server) lookupPage(w http.ResponseWriter, q, tmpl, title, active string) {
	p := page{Title: title, Active: active, Query: q}
	if q != "" {
		res, prim, err := s.Res.Lookup(q)
		if err != nil {
			p.Error = err.Error()
		} else {
			p.Result = res
			p.Primary = prim
		}
	}
	s.render(w, tmpl, p)
}

func (s *Server) appPage(w http.ResponseWriter, r *http.Request) {
	s.lookupPage(w, strings.TrimSpace(r.URL.Query().Get("q")), "app.html", "App · KNS", "app")
}

func (s *Server) namePage(w http.ResponseWriter, r *http.Request) {
	raw := strings.TrimPrefix(r.URL.Path, "/name/")
	s.lookupPage(w, raw, "name.html", raw+" · KNS", "app")
}

func (s *Server) registerPage(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	p := page{Title: "Register · KNS", Active: "app", Query: q}
	if q != "" {
		plan, err := protocol.PlanRegister(q, protocol.Mainnet)
		if err != nil {
			p.Error = err.Error()
		} else {
			p.Plan = &plan
		}
		if res, _, err := s.Res.Lookup(q); err == nil {
			p.Result = res
		}
	}
	s.render(w, "register.html", p)
}

func (s *Server) vaultsPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "vaults.html", page{Title: "Vaults · KNS", Active: "vaults"})
}

func (s *Server) covenantPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "covenant.html", page{
		Title:  "Uniqueness · KNS",
		Active: "covenant",
		Routes: registrar.Routes(),
	})
}

func (s *Server) docsPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "docs.html", page{Title: "Docs · KNS", Active: "docs"})
}

func (s *Server) honestPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "honest.html", page{Title: "Claims · KNS", Active: "honest", Claims: protocol.Claims()})
}

func (s *Server) gatewayPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "gateway.html", page{Title: "Sites · KNS", Active: "gateway"})
}

func (s *Server) ecosystemPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "ecosystem.html", page{Title: "Ecosystem · KNS", Active: "ecosystem"})
}

func (s *Server) site(w http.ResponseWriter, r *http.Request) {
	raw := strings.TrimPrefix(r.URL.Path, "/site/")
	res, _, err := s.Res.Lookup(raw)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if web4.WantsJSON(r.Header.Get("Accept")) || r.URL.Query().Get("format") == "json" {
		writeJSON(w, 200, web4.CardFrom(res, origin(r)))
		return
	}
	if r.URL.Query().Get("go") == "1" {
		target := res.Records.WebsiteTarget()
		if protocol.ContentKind(target) == "ipfs" {
			http.Redirect(w, r, protocol.IPFSGatewayURL(target), http.StatusFound)
			return
		}
		if protocol.ContentKind(target) == "https" {
			http.Redirect(w, r, target, http.StatusFound)
			return
		}
	}
	s.render(w, "site.html", page{Title: res.Name, Active: "gateway", Result: res})
}

func (s *Server) apiResolve(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("name")
	if q == "" {
		q = r.URL.Query().Get("q")
	}
	res, prim, err := s.Res.Lookup(q)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "data": res, "primary": prim})
}

func (s *Server) apiSuggest(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		writeJSON(w, 200, map[string]any{"ok": true, "data": nil})
		return
	}
	res, prim, err := s.Res.Lookup(q)
	if err != nil {
		writeJSON(w, 200, map[string]any{"ok": true, "error": err.Error(), "q": q})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "data": res, "primary": prim})
}

func (s *Server) apiReverse(w http.ResponseWriter, r *http.Request) {
	p, err := s.Res.Reverse(r.URL.Query().Get("address"))
	if err != nil {
		writeJSON(w, 502, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "data": p})
}

func (s *Server) apiCheck(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("name")
	n, err := protocol.ParseName(q)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	items, err := s.Res.API.Check([]string{n.String()}, r.URL.Query().Get("address"))
	if err != nil {
		writeJSON(w, 502, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "data": items, "priceKas": protocol.PriceKAS(n.ApexLabel())})
}

func (s *Server) apiPlan(w http.ResponseWriter, r *http.Request) {
	plan, err := protocol.PlanRegister(r.URL.Query().Get("name"), protocol.Mainnet)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "data": plan})
}

func (s *Server) llms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(`# KNS Web4.0
Web4.0 = agents as users of a name (read, discover, call, pay).
Live: KNS indexer resolve.
Agent card (ERC-8004 file shape, no Kaspa registry): GET /agent/{name}.json
MCP: GET|POST /mcp
Kaspa HTTP 402 (not Coinbase x402): GET /api/v1/call/{name}
Work credits (grams, not USD): GET /credits  GET /api/v1/credits/quote?grams=1000000
JSON site: GET /site/{name}?format=json
Superapp map: C:\\Users\\Remco\\Documents\\kaspa\\superapp
`))
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
