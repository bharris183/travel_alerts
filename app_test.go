package main

import (
	"testing"
)

func TestGetThreatLevel(t *testing.T) {
	tCategory := []string {"Level 2: Exercise Increased Caution", "foo"}
	got := getThreatLevel(tCategory)
	want := "2"
	if got != want {
		t.Errorf("getThreatLevel Failure. Wanted %s. Got %s\n", want, got)
	}
}

func TestGetThreatLevelBadInput(t *testing.T) {
	tCategory := []string {"foo", "foo"}
	got := getThreatLevel(tCategory)
	want := ""
	if got != want {
		t.Errorf("getThreatLevel Failure. Wanted %s. Got %s\n", want, got)
	}
}