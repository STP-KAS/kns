# Draft convention: KNS records for Web4.0 agents

Status: idea for a future KCC. Not submitted to kaspanet/kccs.

A `.kas` name SHOULD be resolvable to:

- owner Kaspa address (live: KNS indexer)
- pay URI `kaspa:<addr>?label=<name>`
- agent registration file (ERC-8004 *shape*, no EVM registry)
- KaChat contact = owner address (`ciph_msg:1:` sealed in the client)
- optional `vaultCommit` pointing at a Kassword/DAG sealed blob
- optional rank is derived from live balance, not stored

This is a convention for wallets and agents. It is not consensus.
