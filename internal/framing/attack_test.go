package framing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLayoutsSameTotal(t *testing.T) {
	if Canonical().Total != Attack().Total {
		t.Fatal("attack must preserve length")
	}
	if Canonical().Total != 42 {
		t.Fatal(Canonical().Total)
	}
}

func TestContractSourcesDoNotReadForeignState(t *testing.T) {
	root := filepath.Join("..", "..", "contracts")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".sil") {
			return err
		}
		if strings.Contains(info.Name(), "FORBIDDEN") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(b)
		if strings.Contains(s, "readInputState") {
			t.Errorf("%s uses readInputState — forbidden on v1-rc1 (#234)", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
