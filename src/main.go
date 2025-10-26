package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilrudie/bulk-transcode/src/pkg/config"
	"github.com/ilrudie/bulk-transcode/src/pkg/ffmpeg"
	"github.com/spf13/cobra"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

const sep = "----------------------------------------"

var (
	cfgFile   string
	recursive bool
	inputDir  string
	outputDir string
	mark      string
	execute   bool
	verbose   bool

	rootCmd = &cobra.Command{
		Use:   "bulk-transcode",
		Short: "bulk-transcode is a tool to scan a directory and transcode media files in bulk.",
		Run:   run,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (optional)")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "enable recursive directory scanning")
	rootCmd.PersistentFlags().StringVarP(&inputDir, "input-dir", "i", "", "input directory")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", "", "output directory")
	rootCmd.PersistentFlags().StringVarP(&mark, "mark", "m", "", "output file mark")
	rootCmd.PersistentFlags().BoolVarP(&execute, "execute", "e", false, "execute the transcoding commands")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	if verbose {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	cfg := config.DefaultConfig()
	if cfgFile == "" {
		log.Warn("No config file specified, using default configuration")
	} else {
		log.Info("Using config file", "path", cfgFile)
		loadedCfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			log.Error("Failed to load config file", "error", err)
			panic("Cannot continue without valid configuration")
		}
		cfg = loadedCfg
	}
	cfg.ArgOverrides(inputDir, outputDir, mark, execute, cmd.Flags().Changed("execute"), recursive, cmd.Flags().Changed("recursive"))
	log.Info("Configuration", "config", cfg)
	// Further processing would go here
	jobs := make([]ffmpeg.Exec, 0)

	walkErr := filepath.WalkDir(cfg.InputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Error("Error accessing path", "path", path, "error", err)
			return err
		}
		if !d.IsDir() {
			log.Info("Found file", "path", path)
			// Here you would add logic to transcode the file
			_, file := filepath.Split(path)
			parts := strings.Split(file, ".")
			if len(parts) < 2 {
				log.Warn("Skipping file with no extension", "file", file)
				return nil
			}
			ext := parts[len(parts)-1]
			name := strings.Join(parts[:len(parts)-1], ".")
			if ext != "mp4" {
				log.Warn("Skipping unsupported file type", "file", path)
				return nil
			}
			if cfg.OutputMark != "" {
				name = fmt.Sprintf("%s.%s", name, cfg.OutputMark)
			}
			outputPath := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s.%s", name, ext))
			if exists(outputPath) {
				log.Info("Output file already exists, skipping", "output", outputPath)
			} else {
				exec := ffmpeg.New(path, outputPath)
				jobs = append(jobs, *exec)
				log.Info("Prepared transcoding job", "input", path, "output", outputPath, "command", exec.Generate(cfg.CommandArguments))
			}
		} else {
			if !cfg.Recursive && path != cfg.InputDir {
				return filepath.SkipDir
			}
			if cfg.Recursive && path == cfg.OutputDir {
				log.Warn("Skipping output directory during recursive scan", "output_dir", cfg.OutputDir)
				return filepath.SkipDir
			}
		}
		return nil
	})
	if walkErr != nil {
		log.Error("Error walking the input directory", "error", walkErr)
	}
	if !cfg.Exec {
		log.Info("Execution flag not set; displaying prepared commands only")
		fmt.Println(sep)
		for _, job := range jobs {
			fmt.Println(job.Generate(cfg.CommandArguments))
		}
		fmt.Println(sep)
	}

}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
