// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

// Package main provides the codescan CLI tool for generating OpenAPI/Swagger
// specifications from annotated Go code.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/3idey/codescan/codescan"
	"github.com/go-openapi/spec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "codescan",
	Short: "Generate OpenAPI/Swagger spec from annotated Go code",
	Long: `codescan is a tool that scans Go source code for swagger annotations
and generates an OpenAPI 2.0 (Swagger) specification document.

It parses Go packages and extracts API metadata from specially formatted
comments to build a complete API specification.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("codescan %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
	},
}

var (
	// generate command flags
	outputFile              string
	outputFormat            string
	workDir                 string
	buildTags               string
	scanModels              bool
	excludeDeps             bool
	includes                []string
	excludes                []string
	includeTags             []string
	excludeTags             []string
	inputSpec               string
	setXNullableForPointers bool
	refAliases              bool
	transparentAliases      bool
	descWithRef             bool
	compact                 bool
)

var generateCmd = &cobra.Command{
	Use:   "generate [packages...]",
	Short: "Generate a swagger spec from annotated Go code",
	Long: `Scans the specified Go packages for swagger annotations and generates
an OpenAPI 2.0 specification.

Examples:
  # Generate spec for current package
  codescan generate ./...

  # Generate spec with custom output
  codescan generate -o api.yaml --format yaml ./cmd/server

  # Generate spec with build tags
  codescan generate --tags=integration ./...`,
	Args: cobra.MinimumNArgs(1),
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(generateCmd)

	// Output flags
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file (default: stdout)")
	generateCmd.Flags().StringVar(&outputFormat, "format", "json", "output format: json or yaml")

	// Scan options
	generateCmd.Flags().StringVarP(&workDir, "work-dir", "w", "", "working directory for package resolution")
	generateCmd.Flags().StringVar(&buildTags, "tags", "", "build tags to use when scanning")
	generateCmd.Flags().BoolVar(&scanModels, "scan-models", false, "include models that are not referenced by operations")
	generateCmd.Flags().BoolVar(&excludeDeps, "exclude-deps", false, "exclude dependencies from scanning")

	// Include/Exclude filters
	generateCmd.Flags().StringSliceVar(&includes, "include", nil, "patterns to include")
	generateCmd.Flags().StringSliceVar(&excludes, "exclude", nil, "patterns to exclude")
	generateCmd.Flags().StringSliceVar(&includeTags, "include-tags", nil, "tags to include")
	generateCmd.Flags().StringSliceVar(&excludeTags, "exclude-tags", nil, "tags to exclude")

	// Input spec
	generateCmd.Flags().StringVarP(&inputSpec, "input", "i", "", "input swagger spec to merge with")

	// Schema options
	generateCmd.Flags().BoolVar(&setXNullableForPointers, "x-nullable-pointers", false, "set x-nullable for pointer types")
	generateCmd.Flags().BoolVar(&refAliases, "ref-aliases", false, "use $ref for type aliases")
	generateCmd.Flags().BoolVar(&transparentAliases, "transparent-aliases", false, "make type aliases completely transparent")
	generateCmd.Flags().BoolVar(&descWithRef, "desc-with-ref", false, "allow descriptions together with $ref")

	// Output formatting
	generateCmd.Flags().BoolVar(&compact, "compact", false, "produce compact JSON output")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	opts := &codescan.Options{
		Packages:                args,
		ScanModels:              scanModels,
		WorkDir:                 workDir,
		BuildTags:               buildTags,
		ExcludeDeps:             excludeDeps,
		Include:                 includes,
		Exclude:                 excludes,
		IncludeTags:             includeTags,
		ExcludeTags:             excludeTags,
		SetXNullableForPointers: setXNullableForPointers,
		RefAliases:              refAliases,
		TransparentAliases:      transparentAliases,
		DescWithRef:             descWithRef,
	}

	// Load input spec if provided
	if inputSpec != "" {
		spec, err := loadInputSpec(inputSpec)
		if err != nil {
			return fmt.Errorf("failed to load input spec: %w", err)
		}
		opts.InputSpec = spec
	}

	// Run the scanner
	swspec, err := codescan.Run(opts)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Marshal the output
	var output []byte
	switch strings.ToLower(outputFormat) {
	case "yaml", "yml":
		output, err = yaml.Marshal(swspec)
	case "json":
		if compact {
			output, err = json.Marshal(swspec)
		} else {
			output, err = json.MarshalIndent(swspec, "", "  ")
		}
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal spec: %w", err)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, output, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Spec written to %s\n", outputFile)
	} else {
		fmt.Println(string(output))
	}

	return nil
}

func loadInputSpec(path string) (*spec.Swagger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var swspec spec.Swagger
	// Try JSON first
	if err := json.Unmarshal(data, &swspec); err != nil {
		// Fall back to YAML
		if err := yaml.Unmarshal(data, &swspec); err != nil {
			return nil, fmt.Errorf("failed to parse as JSON or YAML: %w", err)
		}
	}
	return &swspec, nil
}
