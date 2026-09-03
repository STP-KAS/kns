package protocol

// Evidence is the kaspaexplained-style label. Do not mix these.
type Evidence string

const (
	Live     Evidence = "live"      // mainnet, you can re-derive it
	Indexer  Evidence = "indexer"   // true if you trust the KNS indexer
	Local    Evidence = "local"     // this process only
	Roadmap  Evidence = "roadmap"   // designed, not deployed
	Research Evidence = "research"  // paper/KIP, not a product
	Wrong    Evidence = "wrong"     // sounds right, contradicts the spec
)

type Claim struct {
	Claim   string   `json:"claim"`
	Status  Evidence `json:"status"`
	Meaning string   `json:"meaning"`
	Source  string   `json:"source"`
}

func Claims() []Claim {
	return []Claim{
		{"You can send KAS to kns.kas in supporting wallets", Indexer, "The KNS indexer maps the first valid domain inscription to an owner address. Wallets that call that API can pay a name. Consensus does not know the name.", "api.knsdomains.org/mainnet + KNS envelope spec"},
		{"A .kas name is unique on L1", Wrong, "Uniqueness is an indexer FCFS rule. A covenant cannot see the rest of the UTXO set. KIP-20 ids are hashed from an outpoint, not from a label. KCC-721 says ticker uniqueness the same way.", "KIP-20 genesis formula; KCC-721: ticker uniqueness is an indexer validity rule"},
		{"sha256(name) can be the covenant_id", Wrong, "Genesis covenant_id = BLAKE2b(CovenantID, outpoint || authorized outputs). You cannot pick the id to encode a name.", "KIP-20"},
		{"Toccata is live", Live, "Activated 30 Jun 2026 at DAA 474,165,565. rusty-kaspa v2.0.1. KIPs 16, 17, 20, 21 Active.", "docs.kaspa.org/toccata, kaspaexplained.com/status"},
		{"name.kas.limo is a website", Indexer, "KNS stores redirectUrl as a text inscription and a gateway forwards. That is a redirect, not content addressed on Kaspa. IPFS/KFS still have to be set and hosted.", "KNS redirect-url docs"},
		{"Covenant domain = the name is the UTXO", Roadmap, "KNS said this in Aug–Sep 2026. True only if a registrar singleton issues child name UTXOs (or a based-ZK name set). A random covenant with the label in its state is not globally unique.", "@knsdomain 24 Aug / 2 Sep 2026; KIP-20"},
		{"This app registers covenant names on mainnet", Wrong, "It resolves live inscriptions. It does not submit Toccata txs. Silverscript here is a spec, not a deployment.", "this repo"},
		{"Vault bytes are on-chain", Wrong, "kaspa-data-vault keeps ciphertext off-chain. Only a commitment and a 5-minute window are anchored. Binding a vault to a name is a record convention, not live on kns.kas.", "kaspa-data-vault README"},
		{"Subnames (pay.shop.kas) exist in KNS inscriptions", Wrong, "Inscription create takes a single label. Subnames are a parent-issued covenant design: uniqueness is local to the parent, which is the part that actually scales on a UTXO DAG.", "KNS create spec"},
		{"Kaspa has ENS-class mature naming", Roadmap, "ENS has a hierarchical on-chain registry, Universal Resolver, DNSSEC import, 37M names, 575+ integrations. KNS has 50k+ inscriptions and an indexer. That is a start, not the same object.", "ens.domains; KNS 50k post Jul 2026"},
		{"vProgs / based ZK registrar can hold the name set", Research, "KIP-21 lanes exist. vProgs repo is early development. Do not ship product copy that says the name set is already a proof-verified off-chain state.", "kaspaexplained.com/status"},
		{"KNS 4.0 means Kaspa 4.0 / DAGKnight is live", Wrong, "There is no Kaspa 4.0 fork. Web 4.0 means agents as users of the name. DAGKnight is KIP-2 Proposed.", "kaspaexplained.com/status; Cloudflare Agentic Internet 2026"},
		{"A local name-set simulator is a vProg", Wrong, "The /sim set is an in-process map with a toy hash root. No RISC Zero proof, no lane, no settlement covenant.", "this repo"},
		{"Every .kas name is a Web 4.0 agent", Indexer, "This app synthesizes an ERC-8004-shaped agent card and a Kaspa 402 challenge from indexer records. There is no on-chain agent registry on Kaspa. MCP/A2A endpoints are this process.", "this repo; EIP-8004 registration file shape"},
		{"x402 payments work on Kaspa like Base USDC", Wrong, "Coinbase x402 is an HTTP 402 rail for USDC on EVM. /api/v1/call/{name} returns 402 with a kaspa: payTo. It does not verify settlement.", "x402 spec vs this handler"},
		{"This app is KaChat", Wrong, "KaChat encrypts ciph_msg payloads in the client. We resolve a .kas name to a contact and show envelope shapes. We cannot seal a message without the key.", "kachat.org"},
		{"Kassword passwords are in this database", Wrong, "Kassword has no server. Ciphertext is browser + optional DAG backup. We link the product and a vaultCommit convention.", "kassword.com"},
		{"Ocean rank is the KASRANKS NFT", Wrong, "/ranks uses api.kaspa.org balance. kasranks.com NFTs are a separate 782-supply collection.", "kasranks.com"},
		{"KCCs are live protocol", Wrong, "KCC-0/1/2/20 are Draft. They are not KIPs. Wallets may converge; nodes do not have to.", "github.com/kaspanet/kccs"},
		{"Silverscript v1 is the mainnet language release", Indexer, "v1-rc1 is an official kaspanet pre-release (30 Aug 2026, @OriNewman). Functionally aimed at v1. Mainnet v1 was gated on a week of feedback. Treat RC as RC until v1 is tagged.", "github.com/kaspanet/silverscript/releases/tag/v1-rc1"},
		{"WorkCredits are a USD stablecoin", Wrong, "1 credit = 1 KIP-21 gram of a named lane. The covenant is a prepaid work ledger. It has no oracle and no reserve. Igra bridged USDC/USDT is the dollar rail.", "KIP-21; Igra Hyperlane warp routes Apr 2026; Documents/kaspa/superapp/WORK-CREDITS.md"},
		{"Covenants can mint USD purchasing power", Wrong, "A script can conserve grams and require signatures. It cannot see Circle's reserve or print dollars. Algorithmic kas-USD is UST with extra steps.", "this repo; failed UST/Luna"},
		{"Kaspa L1 fees are already a KAS-denominated work unit", Live, "Min-relay policy 100 sompi/gram. 1e6 grams = 1 KAS at that policy. USD(fee) still floats with KAS/USD.", "Kaspa fee policy; KIP-21 lanes 50/block, 1e9 gas/lane"},
		{"Igra has capital-backed stables", Live, "Hyperlane warp routes on Igra: USDC, iKAS, cbBTC, wstETH. Kaskad-style apps also list USDT. That is bridged reserves, not an L1 covenant.", "@Igra_Labs 9 Apr 2026; Hyperlane"},
		{"WorkCredit.sil is a mainnet voucher you can spend today", Local, "silverc v1-rc1 compiled the artifact. No genesis UTXO is submitted by this app. Quotes at /credits are local invoices.", "contracts/v1/WorkCredit.sil"},
		{"This stack integrates AgenC / Tetsuo", Wrong, "AgenC is a live Solana agent marketplace (agenc.ag, @tetsuoai). Looked at; not wired. Settlement is SOL, not Kaspa L1. Telegram dump in Documents is chat export, not a runtime.", "agenc.ag; github.com/tetsuo-ai"},
		{"Silverscript v1-rc1 foreign state reads are safe", Wrong, "readInputState offsets are compile-time. Non-minimal push encodings can slide fields while keeping template hash and length. silverscript#234 proposed a framing pin; it was closed unmerged on RC day. This repo only transitions its own UTXOs.", "https://github.com/kaspanet/silverscript/pull/234 ; conventions/no-foreign-state.md"},
	}
}
