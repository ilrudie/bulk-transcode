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

func (e *Exec) Generate(args Args) []string {
	cmd := []string{"ffmpeg"}
	for _, option := range args.InputOptions {
		cmd = append(cmd, strings.Fields(option)...)
	}
	cmd = append(cmd, "-i", e.Input)
	for _, option := range args.OutputOptions {
		cmd = append(cmd, strings.Fields(option)...)
	}
	cmd = append(cmd, e.Output)
	return cmd
}

func (e *Exec) Run(args Args, logDir string) error {
	parts := e.Generate(args)
	cmd := exec.Command(parts[0], parts[1:]...)
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
		content := "=== STDOUT ===\n" + stdout.String() + "\n=== STDERR ===\n" + stderr.String()
		if writeErr := os.WriteFile(logPath, []byte(content), 0644); writeErr != nil {
			return writeErr
		}
	}
	return err
}
