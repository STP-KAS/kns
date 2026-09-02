package contracts

import "testing"

func TestArtifactsPresent(t *testing.T) {
	for _, n := range []string{"v1/KasName.json", "v1/KaChatPayTimeout.json", "v1/WorkCredit.json"} {
		b, err := Artifacts.ReadFile(n)
		if err != nil || len(b) < 100 {
			t.Fatalf("%s: %v len=%d", n, err, len(b))
		}
	}
}
