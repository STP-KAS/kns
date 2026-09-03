# KNS post-Toccata protocol

This is the naming layer for Kaspa, modeled on ENS but native to UTXO + covenants.

## What exists today (inscription KNS)

Official product: [app.knsdomains.org](https://app.knsdomains.org), docs, [@knsdomain](https://x.com/knsdomain).

- Envelope protocol id: `kns`
- Domain create: `{"op":"create","p":"domain","v":"<label>"}`
- Transfer: `{"op":"transfer","p":"domain","id":"<inscriptionId>","to":"<kaspa addr>"}`
- Commit-reveal P2SH, same family as KRC-20
- Reveal output 0 pays:
  - mainnet `kaspa:qyp4nvaq3pdq7609z09fvdgwtc9c7rg07fuw5zgeee7xpr085de59eseqfcmynn`
- Price by visual length: 1–2 = 4200, 3 = 2100, 4 = 525, 5+ = 35 KAS. Text = 1 KAS
- No renewal. First-come, first-served. Max 255 chars / 520-byte inscription envelope
- Profile fields as text inscriptions: avatar, website, banner, x, github, telegram, discord, email, redirectUrl
- Gateway: `https://<label>.kas.limo`
- Indexer: `https://api.knsdomains.org/mainnet`

```
GET  /api/v1/{domain}/owner
GET  /api/v1/assets
GET  /api/v1/asset/{assetId}/detail
POST /api/v1/domains/check   { domainNames, address }
GET  /api/v1/domain/{assetId}/profile
GET  /api/v1/primary-name/{owner}
POST /api/v1/domain/primary-name
```

50,000+ inscriptions (KNS, Jul 2026).

## Silverscript v1-rc1 (Ori Newman, 30 Aug 2026)

Official: https://github.com/kaspanet/silverscript/releases/tag/v1-rc1  
Commit `c7d17a1`, verified GitHub signature. Windows `silverc.exe` SHA256 `fbf75851e8d1c97e1982e72cb26e8b8f6417fa5a6ed99d58693d6314890619c3`.

This repo compiles:

- `contracts/v1/KasName.sil` → `KasName.json` (template_hash `e7f981d9…32b79f`)
- `contracts/v1/KaChatPayTimeout.sil` → `KaChatPayTimeout.json`

Language is `entry` + `validateOutputState`, not the old `#[covenant.singleton]` sketches. RC is not the same as a `v1` tag. README still says experimental / prefer testnet-10 until stable v1.

**House rule:** do not `readInputState` a foreign covenant on v1-rc1. Framing of foreign state is not pinned ([silverscript#234](https://github.com/kaspanet/silverscript/pull/234), closed unmerged). Same 42-byte `TokenState`: amount 1, vault reads 264. Own-UTXO `validateOutputState` only. See `conventions/no-foreign-state.md` and `/234`.

## Web 4.0 (what this actually is)

Web 1 read, Web 2 write, Web 3 own, Web 4 **delegate**. Agents resolve a name, call it, pay it. Cloudflare 2026: readable, discoverable, callable, payable. Crypto shorthand: Web3 + agent (MCP, A2A, ERC-8004, x402).

On Kaspa, today:

- **Readable** — HTML site + `Accept: application/json`
- **Discoverable** — `/agent/{name}.json` uses the ERC-8004 *registration file shape*. There is no Kaspa ERC-8004 contract.
- **Callable** — `POST /mcp` (`resolve_kas`, `agent_card`, `pay_kas`, `kachat_contact`, `kas_rank`, `quote_work`)
- **Payable** — `GET /api/v1/call/{name}` → HTTP 402 with a `kaspa:` payTo **and** `kaspa-work-credit` grams. Not Coinbase x402/USDC. Header `X-Kaspa-Payment` is not verified on-chain.

## Work Credits (fee voucher, not a stable)

dApps sequenced on Kaspa need predictable costs. A Toccata covenant cannot mint USD. It can conserve **grams**.

- Unit: `1 credit = 1 KIP-21 gram` of a named lane
- Policy quote: 100 sompi/gram → 1e6 grams = 1 KAS
- Script: `contracts/v1/WorkCredit.sil` (mint / consume / transfer)
- Quote: `GET /api/v1/credits/quote?grams=`
- Map: `C:\Users\Remco\Documents\kaspa\superapp\WORK-CREDITS.md`

This is the gas-token / fee-juice pattern, not USDC. Dollar apps use Igra Hyperlane USDC/USDT.

## Protocol eras (Toccata / vProgs) — separate from Web 4.0

There is no Kaspa 4.0 hard fork. Named forks: **Crescendo** (live), **Toccata** (live), **DAGKnight** (KIP-2 Proposed), a **2027 100 BPS** target with no spec. vProgs: no public testnet.

KNS naming eras (identity model, not “Web 4.0”):

| Era | Layer | Status |
| --- | --- | --- |
| 1 | Inscription + indexer | Live |
| 2 | Toccata registrar / Name UTXO | Roadmap (consensus ready, app not deployed) |
| 3 | Based ZK name set, lane `KNS1` | Research + local `/sim` |
| 4 | vProgs composition + DAGKnight ordering | Research socket |

`kas://name.kas`, `did:kas:name.kas`, and lane `KNS1` are conventions of this repo.

## KNS 4.0 (superseded product label)

There is no Kaspa 4.0 hard fork. Named forks: **Crescendo** (live), **Toccata** (live), **DAGKnight** (KIP-2 Proposed), a **2027 100 BPS** target with no spec. vProgs: no public testnet.

KNS 4.0 means the naming stack is specified across four eras so identity does not have to be redesigned at each fork:

| Era | Layer | Status |
| --- | --- | --- |
| 1 | Inscription + indexer | Live |
| 2 | Toccata registrar / Name UTXO | Roadmap (consensus ready, app not deployed) |
| 3 | Based ZK name set, lane `KNS1` | Research + local `/sim` |
| 4 | vProgs composition + DAGKnight ordering | Research socket |

`kas://name.kas`, `did:kas:name.kas`, and lane `KNS1` are conventions of this repo.

## Uniqueness (read this first)

KIP-20 genesis:

`covenant_id = BLAKE2b("CovenantID", outpoint || authorized outputs…)`

That id is unique because the outpoint is unique. It does **not** encode `alice`. Anyone can genesis another covenant and write `alice.kas` in state. Nodes accept both. A script cannot ask the UTXO set “is alice taken?”. KCC-721 says the same for tickers: uniqueness is an **indexer** rule.

So “covenant domain ⇒ consensus uniqueness” is only true if:

1. **Root registrar UTXO** — one singleton issues every global `.kas` Name child (serialized register, parallel updates), or
2. **Parent-issued subnames** — `pay.shop.kas` unique because only `shop.kas` can spawn it, or
3. **Based ZK name set** — off-chain set, KIP-21 lane, proof of unused (research).

Hashing the label and calling it a covenant id is false. This repo used to do that. It does not anymore.

Live uniqueness remains KNS indexer FCFS.

## What Toccata changes

Toccata activated 30 June 2026 at DAA `474_165_565` (rusty-kaspa v2.0.1). KIPs 16, 17, 20, 21 Active.

Kaspa is still UTXO. It is not an EVM. A covenant UTXO must create a valid successor. Covenant IDs (KIP-20) keep lineage when the P2SH hash changes. Silverscript `singleton` keeps one live UTXO.

KNS said this in public (24 Aug 2026, 2 Sep 2026):

> With a Covenant Domain, the name is the UTXO. Its rules live in the script, and consensus enforces them on every spend.
> With inscription-based .kas, the record is on-chain and an indexer derives the domain state. With a Covenant Domain, ownership and uniqueness are enforced during transaction validation.

## Layering

| Layer | Uniqueness | Records | Status |
| --- | --- | --- | --- |
| Inscription KNS | Indexer, first reveal | Profile text inscriptions | Live mainnet |
| Covenant domain | Consensus singleton + covenant id | State in UTXO preimage | Spec + overlay in this repo; tooling catching up |
| Based ZK registrar | Off-chain state, lane settlement | Huge namespaces | Roadmap (vProgs) |

## Record keys (ENS-class)

`kas`, `pay`, `igra`, `btc`, `eth`, `contenthash`, `website`, `redirectUrl`, `ipfs`, `arweave`, `kfs`, `vault`, `vaultCommit`, `avatarUrl`, `banner`, `bio`, `x`, `github`, `telegram`, `discord`, `email`

Resolution order for the web gateway: redirectUrl → website → ipfs → kfs → arweave → contenthash.

## Elevation

1. Own `alice.kas` as a KNS inscription.
2. Publish a genesis covenant UTXO whose `label` matches and whose `owner` is the inscription owner.
3. Indexers treat the covenant as canonical for records; the inscription remains the historical NFT-like asset until transferred into the covenant (or burned by convention).

## Vaults

`kaspa-data-vault`: seed-only AES-GCM, 5-minute window, commitment on Kaspa. Name record stores package id + commitment. `NameVault.sil` can lock a KAS bond to the same window.

## Honest limits (from kaspaexplained.com)

- Toccata consensus is live. Wallet/explorer/SDK support is still catching up.
- Silverscript is the authoring path; audit/release status still needs attention.
- Argent and vProgs are not production app infrastructure.
- This app talks to the live KNS indexer for inscriptions and keeps covenant/vault state locally until a mainnet registrar UTXO exists.
