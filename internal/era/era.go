package era

import "kns/internal/protocol"

// KNS 4.0 is a product version, not a Kaspa hard fork.
// Kaspa named forks: Crescendo (live), Toccata (live), DAGKnight (KIP-2, not mainnet),
// then a 2027 100 BPS target. vProgs have no testnet. Do not call DAGKnight "4.0".

type ID int

const (
	Inscription ID = 1
	Covenant    ID = 2
	BasedZK     ID = 3
	VProgs      ID = 4
)

type Era struct {
	ID       ID                 `json:"id"`
	Name     string             `json:"name"`
	Kaspa    string             `json:"kaspaLayer"`
	Status   protocol.Evidence  `json:"status"`
	What     string             `json:"what"`
	Uniq     string             `json:"uniqueness"`
	Ready    string             `json:"readyMeans"`
}

func All() []Era {
	return []Era{
		{
			ID: Inscription, Name: "1 · Inscription",
			Kaspa:  "L1 payloads, KNS indexer",
			Status: protocol.Indexer,
			What:   "Commit-reveal envelope kns. 50k+ names. Pay, primary, profile, kas.limo redirect.",
			Uniq:   "Indexer first-come, first-served.",
			Ready:  "This app resolves it now.",
		},
		{
			ID: Covenant, Name: "2 · Toccata covenant",
			Kaspa:  "KIP-17/20 singleton UTXOs (consensus live)",
			Status: protocol.Roadmap,
			What:   "A Name is a UTXO issued by one .kas registrar. Records and transfer in script. Parent-issued subnames in parallel.",
			Uniq:   "Only the registrar lineage may mint global labels. KIP-20 id is still an outpoint hash, not hash(name).",
			Ready:  "Silverscript spec in-repo. Not deployed. Tooling/audit still catching up.",
		},
		{
			ID: BasedZK, Name: "3 · Based name set",
			Kaspa:  "KIP-16 + KIP-21 lanes (consensus live; apps are not)",
			Status: protocol.Research,
			What:   "Off-chain SMT of names. User ops in a reserved lane. Settlement covenant verifies a proof that the name was free.",
			Uniq:   "The name set lives in the app state commitment, not in a single L1 UTXO. This is how you scale past a serialized registrar.",
			Ready:  "Lane + ZK precompile exist. No production name-set prover. Local simulator at /sim is not L1.",
		},
		{
			ID: VProgs, Name: "4 · Composed vProgs",
			Kaspa:  "vProgs + DAGKnight (neither is mainnet)",
			Status: protocol.Research,
			What:   "Names, vaults, and payments as guests that compose atomically. DAGKnight (KIP-2) would tighten ordering; 100 BPS is a 2027 target with no spec.",
			Ready:  "vProgs: zero releases, no public testnet. DAGKnight: Proposed. This era is a socket, not a switch.",
		},
	}
}
