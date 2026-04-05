package ffmpeg

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Exec struct {
	Input, Output string
}

func New(i, o string) *Exec {
	return &Exec{
		Input:  i,
		Output: o,
	}
}

// Generate returns the command as a slice of tokens suitable for display.
// Tokens preserve their original shell quoting (e.g. "'format=nv12,hwupload'")
// so that joining with spaces produces a copy-pasteable shell command.
// File paths are quoted with shellQuote so spaces in filenames are handled.
func (e *Exec) Generate(args Args) []string {
	cmd := []string{"ffmpeg"}
	for _, option := range args.InputOptions {
		cmd = append(cmd, strings.Fields(option)...)
	}
	cmd = append(cmd, "-i", shellQuote(e.Input))
	for _, option := range args.OutputOptions {
		cmd = append(cmd, strings.Fields(option)...)
	}
	cmd = append(cmd, shellQuote(e.Output))
	return cmd
}

// shellQuote wraps s in single quotes if it contains any character a POSIX
// shell would interpret. Embedded single quotes are escaped as '\”.
func shellQuote(s string) string {
	for _, r := range s {
		safe := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == '/' || r == ':' || r == '='
		if !safe {
			return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
		}
	}
	return s
}

func (e *Exec) Run(args Args, logDir string) error {
	// Build exec args using shellSplit so shell quoting is stripped before
	// passing to exec.Command, which does not go through a shell.
	execArgs := []string{}
	for _, option := range args.InputOptions {
		execArgs = append(execArgs, shellSplit(option)...)
	}
	execArgs = append(execArgs, "-i", e.Input)
	for _, option := range args.OutputOptions {
		execArgs = append(execArgs, shellSplit(option)...)
	}
	execArgs = append(execArgs, e.Output)

	cmd := exec.Command("ffmpeg", execArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if logDir != "" {
		if mkErr := os.MkdirAll(logDir, 0755); mkErr != nil {
			return mkErr
		}
		base := filepath.Base(e.Output)
		ext := filepath.Ext(base)
		logName := strings.TrimSuffix(base, ext) + ".log"
		logPath := filepath.Join(logDir, logName)
		displayCmd := strings.Join(e.Generate(args), " ")
		content := "=== COMMAND ===\n" + displayCmd + "\n\n=== STDOUT ===\n" + stdout.String() + "\n=== STDERR ===\n" + stderr.String()
		if writeErr := os.WriteFile(logPath, []byte(content), 0644); writeErr != nil {
			return writeErr
		}
	}
	return err
}

// shellSplit splits s into tokens the same way a POSIX shell would, handling
// single- and double-quoted spans as single tokens (stripping the quotes).
// This ensures that config values like "-vf 'format=nv12,hwupload'" produce
// ["-vf", "format=nv12,hwupload"] rather than ["-vf", "'format=nv12,hwupload'"].
func shellSplit(s string) []string {
	var tokens []string
	var cur strings.Builder
	inSingle := false
	inDouble := false

	for _, r := range s {
		switch {
		case r == '\'' && !inDouble:
			inSingle = !inSingle
		case r == '"' && !inSingle:
			inDouble = !inDouble
		case (r == ' ' || r == '\t') && !inSingle && !inDouble:
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}
