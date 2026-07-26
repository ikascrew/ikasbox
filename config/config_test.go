package config

import "testing"

// gConfig is a package-level singleton, so every test isolates itself by
// saving/restoring it via t.Cleanup to avoid state leaking across tests.
func withCleanConfig(t *testing.T) {
	t.Helper()
	prev := gConfig
	gConfig = nil
	t.Cleanup(func() { gConfig = prev })
}

func TestGetReturnsDefaultsWhenUnset(t *testing.T) {
	withCleanConfig(t)

	c := Get()
	if c.ContentPath != "/tmp" {
		t.Errorf("ContentPath: got %q, want /tmp", c.ContentPath)
	}
	if c.DatabasePath != "ikasbox.db" {
		t.Errorf("DatabasePath: got %q, want ikasbox.db", c.DatabasePath)
	}
	if c.Host != "" {
		t.Errorf("Host: got %q, want empty", c.Host)
	}
	if c.Port != 5555 {
		t.Errorf("Port: got %d, want 5555", c.Port)
	}
	if c.Debug {
		t.Errorf("Debug: got true, want false")
	}
}

func TestGetIsSingleton(t *testing.T) {
	withCleanConfig(t)

	a := Get()
	b := Get()
	if a != b {
		t.Errorf("Get() should return the same singleton instance across calls")
	}
}

func TestSetAppliesOptionsOverDefaults(t *testing.T) {
	withCleanConfig(t)

	if err := Set(LocalOnly(), Path("custom.db"), Port(1234), Extension("*.mp4,*.png")); err != nil {
		t.Fatalf("Set() error: %+v", err)
	}

	c := Get()
	if c.Host != "localhost" {
		t.Errorf("Host: got %q, want localhost", c.Host)
	}
	if c.DatabasePath != "custom.db" {
		t.Errorf("DatabasePath: got %q, want custom.db", c.DatabasePath)
	}
	if c.Port != 1234 {
		t.Errorf("Port: got %d, want 1234", c.Port)
	}
	if len(c.Extensions) != 2 || c.Extensions[0] != "*.mp4" || c.Extensions[1] != "*.png" {
		t.Errorf("Extensions: got %v", c.Extensions)
	}
	// unrelated fields should still hold their defaults
	if c.ContentPath != "/tmp" {
		t.Errorf("ContentPath: got %q, want default /tmp", c.ContentPath)
	}
}

func TestSetWithNoOptionsIsNoop(t *testing.T) {
	withCleanConfig(t)

	if err := Set(Port(9999)); err != nil {
		t.Fatalf("Set() error: %+v", err)
	}

	// Set() with zero options must leave the existing config untouched
	// (opts == nil short-circuits before resetting to defaultConfig()).
	if err := Set(); err != nil {
		t.Fatalf("Set() with no options errored: %+v", err)
	}

	if Get().Port != 9999 {
		t.Errorf("Set() with no options must not reset config: got Port=%d, want 9999", Get().Port)
	}
}

func TestSetPropagatesOptionError(t *testing.T) {
	withCleanConfig(t)

	err := Set(Argument([]string{}))
	if err == nil {
		t.Fatalf("expected error for Argument([]) with no subcommand")
	}
}

func TestArgumentOption(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		wantSub  string
		wantFunc string
		wantArgs []string
	}{
		{name: "start needs no function", args: []string{"start"}, wantSub: "start"},
		{name: "init needs no function", args: []string{"init"}, wantSub: "init"},
		{name: "group missing function errors", args: []string{"group"}, wantErr: true},
		{
			name:     "group with function and arguments",
			args:     []string{"group", "import", "1", "/path"},
			wantSub:  "group",
			wantFunc: "import",
			wantArgs: []string{"1", "/path"},
		},
		{
			name:     "content with function only",
			args:     []string{"content", "register"},
			wantSub:  "content",
			wantFunc: "register",
			wantArgs: nil,
		},
		{name: "unknown subcommand errors", args: []string{"bogus"}, wantErr: true},
		{name: "empty args errors", args: []string{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCleanConfig(t)

			err := Set(Argument(tt.args))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			c := Get()
			if c.SubCommand != tt.wantSub {
				t.Errorf("SubCommand: got %q, want %q", c.SubCommand, tt.wantSub)
			}
			if c.Function != tt.wantFunc {
				t.Errorf("Function: got %q, want %q", c.Function, tt.wantFunc)
			}
			if len(c.Arguments) != len(tt.wantArgs) {
				t.Errorf("Arguments: got %v, want %v", c.Arguments, tt.wantArgs)
			}
		})
	}
}
