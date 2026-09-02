package kcc

// Kaspa Calls for Conventions — kaspanet/kccs. Drafts, not accepted standards.
// @kccforum is the discussion surface; the documents live in the repo.

type Draft struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	URL    string `json:"url"`
	About  string `json:"about"`
}

func Drafts() []Draft {
	base := "https://github.com/kaspanet/kccs/blob/main/"
	return []Draft{
		{ID: "KCC-0", Title: "Purpose and Guidelines", Status: "Draft", URL: base + "kcc-0000.md", About: "A KCC is a convention, not consensus. KIPs change the node."},
		{ID: "KCC-1", Title: "Covenant definition, layout, ABI", Status: "Draft", URL: base + "kcc-0001.md", About: "Shared language for state, entrypoints, templates."},
		{ID: "KCC-2", Title: "Authority schemes", Status: "Draft", URL: base + "kcc-0002.md", About: "Who may authorize a spend: key, script, or another covenant."},
		{ID: "KCC-20", Title: "Fungible token covenant", Status: "Draft", URL: base + "kcc-0020.md", About: "Token shape wallets/DEX can recognize. Merged as Draft, not Accepted."},
	}
}

func Links() map[string]string {
	return map[string]string{
		"repo":   "https://github.com/kaspanet/kccs",
		"forum":  "https://x.com/kccforum",
		"smiths": "https://kas-smiths.org/",
	}
}
