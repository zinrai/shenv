package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env map[string]string `yaml:"env"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func executeCommand(envVars string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	// Parse shell environment variables
	for _, line := range strings.Split(envVars, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract variable name and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		value := strings.Trim(parts[1], "\"")

		// Expand variables in the value
		value = os.ExpandEnv(value)

		// Set environment variable
		os.Setenv(name, value)
	}

	// Prepare command
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

func main() {
	// Parse flags
	configPath := flag.String("config", filepath.Join(os.Getenv("HOME"), ".config", "shenv", "config.yaml"), "Path to config file")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: shenv [flags] <env-name> <command> [args...]")
		os.Exit(1)
	}

	// Load config
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get environment variables for specified environment
	envName := args[0]
	envVars, ok := config.Env[envName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: environment '%s' not found in config\n", envName)
		os.Exit(1)
	}

	// Execute command with environment variables
	if err := executeCommand(envVars, args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
