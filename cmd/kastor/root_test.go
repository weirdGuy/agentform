package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut string
	}{
		{
			name:    "version subcommand prints version",
			args:    []string{"version"},
			wantOut: "kastor version ",
		},
		{
			name:    "no args prints help",
			args:    []string{},
			wantOut: "Usage:",
		},
		{
			name:    "unknown subcommand errors",
			args:    []string{"frobnicate"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCmd()
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantOut != "" && !strings.Contains(out.String(), tt.wantOut) {
				t.Errorf("output = %q, want it to contain %q", out.String(), tt.wantOut)
			}
		})
	}
}

// The exact version and commit depend on how the binary was built (ldflags,
// go build info, or neither), so only the output shape is asserted.
func TestVersionOutputFormat(t *testing.T) {
	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"version"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	want := regexp.MustCompile(`^kastor version \S+ \(commit \S+\)\n$`)
	if !want.MatchString(out.String()) {
		t.Errorf("version output = %q, want it to match %q", out.String(), want)
	}
}

func TestRootCommandUse(t *testing.T) {
	if got := newRootCmd().Use; got != "kastor" {
		t.Errorf("root command Use = %q, want %q", got, "kastor")
	}
}
