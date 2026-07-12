package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileConfig is the optional on-disk operator config
// (~/.config/talon/config.yaml by default).
//
// Precedence (highest wins): CLI flags > env vars > config file > built-in defaults.
type FileConfig struct {
	CoreURL     string `yaml:"core_url"`
	ArsenalURL  string `yaml:"arsenal_url"`
	MSF         string `yaml:"msf"`  // host:port
	AMQP        string `yaml:"amqp"` // host:port
	Timeout     string `yaml:"timeout"`
	Output      string `yaml:"output"`
	ComposeFile string `yaml:"compose_file"`
	// ProjectDir is the Talon repo / compose project root for logs.
	ProjectDir string `yaml:"project_dir"`
}

// DefaultConfigPath returns $XDG_CONFIG_HOME/talon/config.yaml or ~/.config/talon/config.yaml.
func DefaultConfigPath() string {
	if p := os.Getenv("TALON_CONFIG"); p != "" {
		return p
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "talon", "config.yaml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "talon", "config.yaml")
}

// LoadFileConfig reads path if it exists. Missing file is not an error.
// Parsing is intentionally minimal YAML (key: value lines) to avoid a new
// dependency — enough for operator config without pulling gopkg.in/yaml.
func LoadFileConfig(path string) (FileConfig, error) {
	var cfg FileConfig
	if path == "" {
		return cfg, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config %s: %w", path, err)
	}
	return parseSimpleYAML(string(raw)), nil
}

// parseSimpleYAML supports a tiny subset: top-level `key: value` lines,
// # comments, blank lines. Values may be quoted with " or '.
func parseSimpleYAML(s string) FileConfig {
	var cfg FileConfig
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// strip inline comments only when not inside quotes (best-effort)
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if i := strings.Index(val, " #"); i >= 0 && !strings.HasPrefix(val, `"`) && !strings.HasPrefix(val, `'`) {
			val = strings.TrimSpace(val[:i])
		}
		val = strings.Trim(val, `"'`)
		switch key {
		case "core_url":
			cfg.CoreURL = val
		case "arsenal_url":
			cfg.ArsenalURL = val
		case "msf":
			cfg.MSF = val
		case "amqp":
			cfg.AMQP = val
		case "timeout":
			cfg.Timeout = val
		case "output":
			cfg.Output = val
		case "compose_file":
			cfg.ComposeFile = val
		case "project_dir":
			cfg.ProjectDir = val
		}
	}
	return cfg
}

// ResolvedConfig is the fully merged operator configuration.
type ResolvedConfig struct {
	CoreURL     string
	ArsenalURL  string
	MSF         string
	AMQP        string
	Timeout     time.Duration
	Output      string
	ComposeFile string
	ProjectDir  string
	ConfigPath  string
}

// ResolveConfig merges file + env + flag overrides into a ResolvedConfig.
// emptyFlag means "not set on CLI" for string flags (use file/env/default).
func ResolveConfig(file FileConfig, flagCoreURL, flagOutput string, flagTimeout time.Duration, configPath string) ResolvedConfig {
	r := ResolvedConfig{
		CoreURL:     "http://localhost:8000",
		ArsenalURL:  "http://localhost:8888/health",
		MSF:         "localhost:5554",
		AMQP:        "localhost:5672",
		Timeout:     30 * time.Second,
		Output:      "table",
		ComposeFile: "docker-compose.yml",
		ConfigPath:  configPath,
	}

	// File layer
	if file.CoreURL != "" {
		r.CoreURL = file.CoreURL
	}
	if file.ArsenalURL != "" {
		r.ArsenalURL = file.ArsenalURL
	}
	if file.MSF != "" {
		r.MSF = file.MSF
	}
	if file.AMQP != "" {
		r.AMQP = file.AMQP
	}
	if file.Output != "" {
		r.Output = file.Output
	}
	if file.ComposeFile != "" {
		r.ComposeFile = file.ComposeFile
	}
	if file.ProjectDir != "" {
		r.ProjectDir = file.ProjectDir
	}
	if file.Timeout != "" {
		if d, err := time.ParseDuration(file.Timeout); err == nil {
			r.Timeout = d
		}
	}

	// Env layer
	if v := os.Getenv("TALON_CORE_URL"); v != "" {
		r.CoreURL = v
	}
	if v := os.Getenv("TALON_ARSENAL_URL"); v != "" {
		r.ArsenalURL = v
	}
	if v := os.Getenv("TALON_OUTPUT"); v != "" {
		r.Output = v
	}
	if v := os.Getenv("TALON_PROJECT_DIR"); v != "" {
		r.ProjectDir = v
	}
	if v := os.Getenv("TALON_COMPOSE_FILE"); v != "" {
		r.ComposeFile = v
	}
	if v := os.Getenv("TALON_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			r.Timeout = d
		}
	}
	// MSF/AMQP from stack env when present
	if host := os.Getenv("MSF_SERVER"); host != "" {
		port := os.Getenv("MSF_PORT")
		if port == "" {
			port = "5554"
		}
		r.MSF = host + ":" + port
	}
	if u := os.Getenv("AMQP_URL"); u != "" {
		if h, p, ok := amqpHostPort(u); ok {
			r.AMQP = h + ":" + p
		}
	}

	// Flag layer (only if user set them — empty string = not set for CoreURL/Output)
	if flagCoreURL != "" {
		r.CoreURL = flagCoreURL
	}
	if flagOutput != "" && flagOutput != "table" {
		// Cobra always sets default "table"; treat non-default as explicit.
		// Also allow explicit table via config; flags that match default still
		// win only when Changed — handled by caller via Flag.Changed.
		r.Output = flagOutput
	}
	if flagTimeout > 0 {
		r.Timeout = flagTimeout
	}

	return r
}
