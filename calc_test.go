package main

import "testing"

func TestCastVarsOK(t *testing.T) {
	wantX, wantY := -10, 10
	if gotX, gotY, gotErr := castVars("-10", "10"); gotX != wantX && gotY != wantY && gotErr != nil {
		t.Errorf("castVars() = %v, %v, %v, want %v, %v, %v", gotX, gotY, gotErr, wantX, wantY, nil)
	}
}
func TestCastVarsErrorOnSpace(t *testing.T) {
	if gotX, gotY, gotErr := castVars("-10 ", " 10 "); gotErr == nil {
		t.Errorf(`castVars("-10 ", " 10 ")`+"= %v, %v, %v, want error", gotX, gotY, gotErr)
	}
}
func TestCastVarsErrorOnFloat(t *testing.T) {
	if gotX, gotY, gotErr := castVars("-10.2", "99.2"); gotErr == nil {
		t.Errorf(`castVars("-10.2", "99.2")`+"= %v, %v, %v, want error", gotX, gotY, gotErr)
	}
}
