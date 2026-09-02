# Draft convention: Work Credits (fee vouchers, not stables)

Status: idea for a future KCC. Not submitted to kaspanet/kccs.

A sequenced dApp on Kaspa SHOULD bill in **grams** (KIP-21 mass units), not in USD, unless it actually holds a capital-backed dollar (Igra USDC/USDT).

## Unit

- `1 credit = 1 gram` of a named lane
- Settlement fallback: `sompi = grams × sompiPerGram` (policy today: 100)
- HTTP 402 `Accepts` MAY include scheme `kaspa-work-credit` next to `kaspa`

## Covenant

`WorkCredit` UTXO state: `issuer`, `holder`, `credits`, `lane`.

- `mint` — issuer signature (sale of prepaid work; KAS paid off-script)
- `consume` — holder + issuer signatures, `credits` decreases
- `transfer` — holder signature, whole voucher moves

The covenant is the ledger of remaining work. It is not a reserve. It cannot mint USD.

## What this is not

- Not USDC, DAI, or any $1 token
- Not an algorithmic stable
- Not a claim that miners accept grams as fees (miners still take KAS)
- Not a replacement for Igra when the dApp truly needs dollars
