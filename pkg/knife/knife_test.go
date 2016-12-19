package knife

import (
	"testing"
)

func TestNew(t *testing.T) {
	k := New()

	if k.Username == "" {
		t.Error("Expected a username here")
	}

	if k.HomeDir == "" {
		t.Error("Expected a homedir here")
	}
}
