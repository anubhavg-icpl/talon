package strike

import (
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestDecodeMSFMapIntegerKeys(t *testing.T) {
	// Simulate session.list: integer session id -> info map.
	raw, err := msgpack.Marshal(map[int]any{
		7: map[string]any{
			"type":         "shell",
			"exploit_uuid": "abc123",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Default Unmarshal must fail — that is the production bug.
	var broken map[string]any
	if err := msgpack.Unmarshal(raw, &broken); err == nil {
		t.Fatal("expected default unmarshal into map[string]any to fail on int keys")
	}
	out, err := decodeMSFMap(raw)
	if err != nil {
		t.Fatalf("decodeMSFMap: %v", err)
	}
	info, ok := out["7"].(map[string]any)
	if !ok {
		t.Fatalf("expected session 7 map, got %#v", out)
	}
	if info["type"] != "shell" {
		t.Fatalf("type=%v", info["type"])
	}
	if info["exploit_uuid"] != "abc123" {
		t.Fatalf("uuid=%v", info["exploit_uuid"])
	}
}

func TestAsStringKeyedMapEmpty(t *testing.T) {
	out, err := asStringKeyedMap(map[string]any{})
	if err != nil || len(out) != 0 {
		t.Fatalf("out=%v err=%v", out, err)
	}
	out, err = asStringKeyedMap(nil)
	if err != nil || len(out) != 0 {
		t.Fatalf("nil out=%v err=%v", out, err)
	}
}
