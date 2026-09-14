# KNS — Web4.0 names on Kaspa

**For the KNS team:** https://stp-kas.github.io/kns-spec/ — [STP-KAS/kns-spec](https://github.com/STP-KAS/kns-spec). This repo is the Web4 demo.

**project delusional** · [@StppStp](https://x.com/StppStp)

Part of [STP-KAS/project-delusional](https://github.com/STP-KAS/project-delusional). Sisters: [gramlane](https://github.com/STP-KAS/gramlane) (`:8081`), [kaspa-till](https://github.com/STP-KAS/kaspa-till) (`:8082`).

A `.kas` name for humans and agents. Live data is the official KNS indexer. Covenant uniqueness is **not** consensus. Web4.0 here means readable, discoverable, callable, payable.

This site never DMs you. We never ask for a seed.

## Run

```powershell
cd C:\Users\Remco\kns
go test ./...
go build -o kns.exe ./cmd/kns
.\kns.exe
```

http://localhost:8080

All three sites:

```powershell
powershell -File C:\Users\Remco\Documents\kaspa\start-local.ps1
```

If the browser says “refused to connect”, the `.exe` is not running. Run the script. Leave the minimized windows open.

## Map

| What | Where |
| --- | --- |
| **Idea** (what this URL is) | `/idea` |
| **Why** (plants, agents, boring bills — not tokens) | `/why` |
| Resolve a name or `kaspa:` address | `/app` |
| Profile + pay URI | `/name/kns.kas` |
| Generated site / JSON agent view | `/site/kns.kas` |
| Agent card (ERC-8004 *shape*) | `/agent/kns.kas.json` |
| MCP | `/mcp` |
| Kaspa HTTP 402 | `/api/v1/call/kns.kas` |
| Wallets (Kasware/Kastle inject; catalog the rest) | `/wallets` |
| Work Credits (grams, not USD) | `/credits` |
| Safety (never DMs, never seeds) | `/safety` |
| Feedback (saved on this PC) | `/feedback` |
| Claims checker | `/honest` |
| KaChat contact (not E2E) | `/kachat?q=kns.kas` |
| Kassword pointer | `/kassword` |
| Ocean rank from live balance | `/ranks?q=kns.kas` |
| KCC drafts | `/kcc` |
| Silverscript v1.0.0 artifacts | `/silverc` |
| #234 framing attack (42 bytes: amount 1 → vault 264) | `/234` |

## Honest limits

- Uniqueness today: KNS indexer FCFS.
- Silverscript: official **`v1.0.0`** (`@OriNewman` / someone235, 9 Sep 2026, `3ed9733`).
- KaChat encryption, Kassword ciphertext, and KASRANKS NFTs stay in those apps.
- Agent cards are synthesized. There is no Kaspa ERC-8004 registry.
- Work Credits are prepaid KIP-21 grams. Not a stablecoin. No L2.
- AgenC/Tetsuo is a Solana marketplace. Not wired.
- **No foreign `readInputState`** on v1.0.0 ([silverscript#234](https://github.com/kaspanet/silverscript/pull/234) closed unmerged). Own UTXO only. `conventions/no-foreign-state.md`.
- **`go test` red X:** `TestContractSourcesDoNotReadForeignState` greps the **comment** in `KasName.sil` (`// Do not readInputState…`). The function is not called. House rule still stands. Fix the linter, not the covenant. [tn10-hard-test](https://github.com/STP-KAS/tn10-hard-test).

## Trusted tools used

- `C:\Users\Remco\silverscript` — clone of `github.com/kaspanet/silverscript` @ **`v1.0.0`** (`3ed9733`)
- `C:\Users\Remco\tools\silverc-v1\bin\silverc.exe` — GitHub release zip SHA256 `3e0d660c15a9e7ac90f3960da24d348b076b1891481bfe758db18accc8a102e1`
