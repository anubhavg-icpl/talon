package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseSimpleYAML(t *testing.T) {
	raw := `
# comment
core_url: "http://example:9000"
arsenal_url: http://example:8888/health
timeout: 45s
output: json
project_dir: /opt/talon
`
	cfg := parseSimpleYAML(raw)
	if cfg.CoreURL != "http://example:9000" {
		t.Fatalf("core_url=%q", cfg.CoreURL)
	}
	if cfg.Timeout != "45s" {
		t.Fatalf("timeout=%q", cfg.Timeout)
	}
	if cfg.Output != "json" {
		t.Fatalf("output=%q", cfg.Output)
	}
	if cfg.ProjectDir != "/opt/talon" {
		t.Fatalf("project_dir=%q", cfg.ProjectDir)
	}
}

func TestLoadFileConfigMissing(t *testing.T) {
	cfg, err := LoadFileConfig(filepath.Join(t.TempDir(), "nope.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CoreURL != "" {
		t.Fatalf("expected empty config, got %+v", cfg)
	}
}

func TestLoadFileConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("core_url: http://cfg:8000\noutput: raw\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFileConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CoreURL != "http://cfg:8000" || cfg.Output != "raw" {
		t.Fatalf("%+v", cfg)
	}
}

func TestResolveConfigPrecedence(t *testing.T) {
	t.Setenv("TALON_CORE_URL", "http://env:8000")
	file := FileConfig{CoreURL: "http://file:8000", Timeout: "10s", Output: "json"}
	r := ResolveConfig(file, "http://flag:8000", "table", 30*time.Second, "")
	if r.CoreURL != "http://flag:8000" {
		t.Fatalf("flag should win: %q", r.CoreURL)
	}
	// Without flag, env wins over file
	r2 := ResolveConfig(file, "", "table", 30*time.Second, "")
	if r2.CoreURL != "http://env:8000" {
		t.Fatalf("env should win: %q", r2.CoreURL)
	}
}
