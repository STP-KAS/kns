package kcc

import "testing"

func TestDrafts(t *testing.T) {
	if len(Drafts()) != 4 {
		t.Fatal(len(Drafts()))
	}
	if Links()["repo"] == "" {
		t.Fatal("repo")
	}
}
