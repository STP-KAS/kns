# House rule: do not read foreign covenant state on v1-rc1

Status: project delusional rule. Not a KCC. Not consensus.

Silverscript v1-rc1 can decode another input's state (`readInputState`, `readInputStateWithTemplate`) at compile-time byte offsets. Kaspa's engine accepts more than one push encoding for the same bytes. A hostile UTXO can keep the same total length and the same template hash / P2SH and still **slide field reads** (fake token amounts, fake owners).

Documented in [kaspanet/silverscript#234](https://github.com/kaspanet/silverscript/pull/234) (supertypo, 29–30 Aug 2026). **Closed, not merged.** v1-rc1 does not include the guard.

## Rule

Until a merged Silverscript release pins foreign-state framing (or we pin a later tag that does):

1. Every `.sil` in this stack **only** `validateOutputState`s **its own** UTXO.
2. No vault, till, DEX, or minter that `readInputState`s an instance an attacker can genesis.
3. `KasName`, `WorkCredit`, `KaChatPayTimeout`, `KasInvoice` already obey this. Keep them that way.
4. A future kUSD / KCC-20 reader is **blocked** by this rule, not by missing UI.

Recompiling with current `silverc` v1-rc1 does not fix the hole. Do not treat a closed PR as a patch.
