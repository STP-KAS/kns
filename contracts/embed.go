package contracts

import "embed"

//go:embed v1/KasName.json v1/KaChatPayTimeout.json v1/WorkCredit.json
var Artifacts embed.FS
