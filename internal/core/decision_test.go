package core

import (
	"testing"

	"github.com/anubhavg-icpl/talon/internal/config"
)

func TestNormalizeDecision(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      Decision
		want    string
		wantErr bool
	}{
		{name: "approve lower", in: Decision{Type: "approve"}, want: "approve"},
		{name: "approve mixed", in: Decision{Type: "Approve"}, want: "approve"},
		{name: "reject upper", in: Decision{Type: "REJECT"}, want: "reject"},
		{name: "edit ok", in: Decision{Type: "Edit", EditedArgs: map[string]any{"ports": "21"}}, want: "edit"},
		{name: "edit missing args", in: Decision{Type: "edit"}, wantErr: true},
		{name: "unknown", in: Decision{Type: "maybe"}, wantErr: true},
		{name: "empty", in: Decision{Type: ""}, wantErr: true},
		{name: "whitespace", in: Decision{Type: "  approve  "}, want: "approve"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := NormalizeDecision(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Type != tc.want {
				t.Fatalf("Type=%q want %q", got.Type, tc.want)
			}
		})
	}
}

func TestSessionKeyPrefersSessionID(t *testing.T) {
	t.Parallel()
	a := RunInput{SessionID: "run-1", TargetIP: "1.2.3.4"}
	b := RunInput{SessionID: "run-2", TargetIP: "1.2.3.4"}
	if sessionKey(a) == sessionKey(b) {
		t.Fatal("distinct SessionIDs must not collide")
	}
	if sessionKey(a) != "run-1" {
		t.Fatalf("got %q", sessionKey(a))
	}
}

func TestSessionKeyFallbackComposite(t *testing.T) {
	t.Parallel()
	a := RunInput{
		TargetIP: "1.2.3.4", CVEID: "CVE-1",
		Context: config.Context{LHOST: "10.0.0.1", LPORT: 4444},
	}
	b := RunInput{
		TargetIP: "1.2.3.4", CVEID: "CVE-1",
		Context: config.Context{LHOST: "10.0.0.1", LPORT: 4444},
	}
	c := RunInput{
		TargetIP: "1.2.3.4", CVEID: "CVE-2",
		Context: config.Context{LHOST: "10.0.0.1", LPORT: 4444},
	}
	if sessionKey(a) != sessionKey(b) {
		t.Fatal("identical inputs without SessionID should share key")
	}
	if sessionKey(a) == sessionKey(c) {
		t.Fatal("different CVE should change key")
	}
}
