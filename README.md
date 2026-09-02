# KNS — Web 4.0 names on Kaspa

# project delusional · [ @StppStp ](https://x.com/StppStp)

A `.kas` name for humans and agents. Live data is the official KNS indexer. Covenant uniqueness is not consensus. Web 4.0 here means readable, discoverable, callable, payable — not a Kaspa hard fork.

```powershell
cd C:\Users\Remco\kns
go test ./...
go run ./cmd/kns
```

http://localhost:8080

## Map

| What | Where |
| --- | --- |
| Resolve a name or `kaspa:` address | `/app` |
| Profile + pay URI | `/name/kns.kas` |
| Generated site / JSON agent view | `/site/kns.kas` |
| Agent card (ERC-8004 *shape*) | `/agent/kns.kas.json` |
| MCP | `/mcp` |
| Kaspa HTTP 402 | `/api/v1/call/kns.kas` |
| KaChat contact (not E2E) | `/kachat?q=kns.kas` |
| Kassword pointer | `/kassword` |
| Ocean rank from live balance | `/ranks?q=kns.kas` |
| KCC drafts | `/kcc` |
| Silverscript v1-rc1 artifacts | `/silverc` |
| Work Credits (grams, not USD) | `/credits` |
| Claims checker | `/honest` |
| Superapp map + guide | `Documents\kaspa\superapp` |
| Kaspa wallets (inject Kasware/Kastle, catalog the rest) | `/wallets` |

## Honest limits

- Uniqueness today: KNS indexer FCFS.
- Silverscript: official `v1-rc1` (`@OriNewman`, kaspanet). RC ≠ tagged `v1`.
- KaChat encryption, Kassword ciphertext, and KASRANKS NFTs stay in those apps.
- This process synthesizes agent cards. There is no Kaspa ERC-8004 registry.
- Work Credits are prepaid KIP-21 grams. They are not a stablecoin. Igra USDC/USDT is the dollar rail.

## Trusted tools used

- `C:\Users\Remco\silverscript` — clone of `github.com/kaspanet/silverscript` @ `v1-rc1` (`c7d17a1`)
- `C:\Users\Remco\tools\silverc\silverc.exe` — GitHub release zip, SHA256 `fbf75851e8d1c97e1982e72cb26e8b8f6417fa5a6ed99d58693d6314890619c3`
