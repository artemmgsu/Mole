package main

import (
	"path/filepath"
	"strings"
)

<<<<<<< HEAD
// isCleanableDir checks if a directory is safe to manually delete
// but NOT cleaned by mo clean (so user might want to delete it manually)
=======
// isCleanableDir marks paths safe to delete manually (not handled by mo clean).
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
func isCleanableDir(path string) bool {
	if path == "" {
		return false
	}

<<<<<<< HEAD
	// Exclude paths that mo clean will handle automatically
	// These are system caches/logs that mo clean already processes
=======
	// Exclude paths mo clean already handles.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if isHandledByMoClean(path) {
		return false
	}

	baseName := filepath.Base(path)

<<<<<<< HEAD
	// Only mark project dependencies and build outputs
	// These are safe to delete but mo clean won't touch them
=======
	// Project dependencies and build outputs are safe.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if projectDependencyDirs[baseName] {
		return true
	}

	return false
}

<<<<<<< HEAD
// isHandledByMoClean checks if this path will be cleaned by mo clean
func isHandledByMoClean(path string) bool {
	// Paths that mo clean handles (from clean.sh)
=======
// isHandledByMoClean checks if a path is cleaned by mo clean.
func isHandledByMoClean(path string) bool {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	cleanPaths := []string{
		"/Library/Caches/",
		"/Library/Logs/",
		"/Library/Saved Application State/",
		"/.Trash/",
		"/Library/DiagnosticReports/",
	}

	for _, p := range cleanPaths {
		if strings.Contains(path, p) {
			return true
		}
	}

	return false
}

<<<<<<< HEAD
// Project dependency and build directories
// These are safe to delete manually but mo clean won't touch them
var projectDependencyDirs = map[string]bool{
	// JavaScript/Node dependencies
	"node_modules":     true,
	"bower_components": true,
	".yarn":            true, // Yarn local cache
	".pnpm-store":      true, // pnpm store

	// Python dependencies and outputs
	"venv":                 true,
	".venv":                true,
	"virtualenv":           true,
	"__pycache__":          true,
	".pytest_cache":        true,
	".mypy_cache":          true,
	".ruff_cache":          true,
	".tox":                 true,
	".eggs":                true,
	"htmlcov":              true, // Coverage reports
	".ipynb_checkpoints":   true, // Jupyter checkpoints

	// Ruby dependencies
	"vendor":  true,
	".bundle": true,

	// Java/Kotlin/Scala
	".gradle": true, // Project-level Gradle cache
	"out":     true, // IntelliJ IDEA build output

	// Build outputs (can be rebuilt)
=======
// Project dependency and build directories.
var projectDependencyDirs = map[string]bool{
	// JavaScript/Node.
	"node_modules":     true,
	"bower_components": true,
	".yarn":            true,
	".pnpm-store":      true,

	// Python.
	"venv":               true,
	".venv":              true,
	"virtualenv":         true,
	"__pycache__":        true,
	".pytest_cache":      true,
	".mypy_cache":        true,
	".ruff_cache":        true,
	".tox":               true,
	".eggs":              true,
	"htmlcov":            true,
	".ipynb_checkpoints": true,

	// Ruby.
	"vendor":  true,
	".bundle": true,

	// Java/Kotlin/Scala.
	".gradle": true,
	"out":     true,

	// Build outputs.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"build":         true,
	"dist":          true,
	"target":        true,
	".next":         true,
	".nuxt":         true,
	".output":       true,
	".parcel-cache": true,
	".turbo":        true,
<<<<<<< HEAD
	".vite":         true, // Vite cache
	".nx":           true, // Nx cache
	"coverage":      true,
	".coverage":     true,
	".nyc_output":   true, // NYC coverage

	// Frontend framework outputs
	".angular":     true, // Angular CLI cache
	".svelte-kit":  true, // SvelteKit build
	".astro":       true, // Astro cache
	".docusaurus":  true, // Docusaurus build

	// iOS/macOS development
=======
	".vite":         true,
	".nx":           true,
	"coverage":      true,
	".coverage":     true,
	".nyc_output":   true,

	// Frontend framework outputs.
	".angular":    true,
	".svelte-kit": true,
	".astro":      true,
	".docusaurus": true,

	// Apple dev.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"DerivedData": true,
	"Pods":        true,
	".build":      true,
	"Carthage":    true,
<<<<<<< HEAD

	// Other tools
	".terraform": true, // Terraform plugins
=======
	".dart_tool":  true,

	// Other tools.
	".terraform": true,
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}
