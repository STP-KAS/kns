# KNS — Web 4.0 names on Kaspa

**project delusional** · [@StppStp](https://x.com/StppStp)

Part of [STP-KAS/project-delusional](https://github.com/STP-KAS/project-delusional). Sisters: [gramlane](https://github.com/STP-KAS/gramlane) (`:8081`), [kaspa-till](https://github.com/STP-KAS/kaspa-till) (`:8082`).

A `.kas` name for humans and agents. Live data is the official KNS indexer. Covenant uniqueness is **not** consensus. Web 4.0 here means readable, discoverable, callable, payable — not a Kaspa hard fork.

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
| Silverscript v1-rc1 artifacts | `/silverc` |

## Honest limits

- Uniqueness today: KNS indexer FCFS.
- Silverscript: official `v1-rc1` (`@OriNewman`, kaspanet). RC ≠ tagged `v1`.
- KaChat encryption, Kassword ciphertext, and KASRANKS NFTs stay in those apps.
- Agent cards are synthesized. There is no Kaspa ERC-8004 registry.
- Work Credits are prepaid KIP-21 grams. Not a stablecoin. No L2.
- AgenC/Tetsuo is a Solana marketplace. Not wired.

## Trusted tools used

- `C:\Users\Remco\silverscript` — clone of `github.com/kaspanet/silverscript` @ `v1-rc1` (`c7d17a1`)
- `C:\Users\Remco\tools\silverc\silverc.exe` — GitHub release zip, SHA256 `fbf75851e8d1c97e1982e72cb26e8b8f6417fa5a6ed99d58693d6314890619c3`
