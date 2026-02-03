package main

import "time"

const (
<<<<<<< HEAD
	maxEntries            = 30
	maxLargeFiles         = 30
	barWidth              = 24
	minLargeFileSize      = 100 << 20 // 100 MB
	defaultViewport       = 12        // Default viewport when terminal height is unknown
	overviewCacheTTL      = 7 * 24 * time.Hour // 7 days
	overviewCacheFile     = "overview_sizes.json"
	duTimeout             = 60 * time.Second // Increased for large directories
	mdlsTimeout           = 5 * time.Second
	maxConcurrentOverview = 3                // Scan up to 3 overview dirs concurrently
	batchUpdateSize       = 100              // Batch atomic updates every N items
	cacheModTimeGrace     = 30 * time.Minute // Ignore minor directory mtime bumps

	// Worker pool configuration
	minWorkers         = 8                // Minimum workers for better I/O throughput
	maxWorkers         = 64               // Maximum workers to avoid excessive goroutines
	cpuMultiplier      = 2                // Worker multiplier per CPU core for I/O-bound operations
	maxDirWorkers      = 16               // Maximum concurrent subdirectory scans
	openCommandTimeout = 10 * time.Second // Timeout for open/reveal commands
)

var foldDirs = map[string]bool{
	// Version control
=======
	maxEntries             = 30
	maxLargeFiles          = 20
	barWidth               = 24
	spotlightMinFileSize   = 100 << 20
	largeFileWarmupMinSize = 1 << 20
	defaultViewport        = 12
	overviewCacheTTL       = 7 * 24 * time.Hour
	overviewCacheFile      = "overview_sizes.json"
	duTimeout              = 30 * time.Second
	mdlsTimeout            = 5 * time.Second
	maxConcurrentOverview  = 8
	batchUpdateSize        = 100
	cacheModTimeGrace      = 30 * time.Minute

	// Worker pool limits.
	minWorkers         = 16
	maxWorkers         = 64
	cpuMultiplier      = 4
	maxDirWorkers      = 32
	openCommandTimeout = 10 * time.Second
)

var foldDirs = map[string]bool{
	// VCS.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	".git": true,
	".svn": true,
	".hg":  true,

<<<<<<< HEAD
	// JavaScript/Node
	"node_modules":                  true,
	".npm":                          true,
	"_npx":                          true, // ~/.npm/_npx global cache
	"_cacache":                      true, // ~/.npm/_cacache
=======
	// JavaScript/Node.
	"node_modules":                  true,
	".npm":                          true,
	"_npx":                          true,
	"_cacache":                      true,
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"_logs":                         true,
	"_locks":                        true,
	"_quick":                        true,
	"_libvips":                      true,
	"_prebuilds":                    true,
	"_update-notifier-last-checked": true,
	".yarn":                         true,
	".pnpm-store":                   true,
	".next":                         true,
	".nuxt":                         true,
	"bower_components":              true,
	".vite":                         true,
	".turbo":                        true,
	".parcel-cache":                 true,
	".nx":                           true,
	".rush":                         true,
	"tnpm":                          true,
	".tnpm":                         true,
	".bun":                          true,
	".deno":                         true,

<<<<<<< HEAD
	// Python
=======
	// Python.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"__pycache__":   true,
	".pytest_cache": true,
	".mypy_cache":   true,
	".ruff_cache":   true,
	"venv":          true,
	".venv":         true,
	"virtualenv":    true,
	".tox":          true,
	"site-packages": true,
	".eggs":         true,
	"*.egg-info":    true,
	".pyenv":        true,
	".poetry":       true,
	".pip":          true,
	".pipx":         true,

<<<<<<< HEAD
	// Ruby/Go/PHP (vendor), Java/Kotlin/Scala/Rust (target)
=======
	// Ruby/Go/PHP (vendor), Java/Kotlin/Scala/Rust (target).
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"vendor":        true,
	".bundle":       true,
	"gems":          true,
	".rbenv":        true,
	"target":        true,
	".gradle":       true,
	".m2":           true,
	".ivy2":         true,
	"out":           true,
	"pkg":           true,
	"composer.phar": true,
	".composer":     true,
	".cargo":        true,

<<<<<<< HEAD
	// Build outputs
=======
	// Build outputs.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"build":     true,
	"dist":      true,
	".output":   true,
	"coverage":  true,
	".coverage": true,

<<<<<<< HEAD
	// IDE
=======
	// IDE.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	".idea":   true,
	".vscode": true,
	".vs":     true,
	".fleet":  true,

<<<<<<< HEAD
	// Cache directories
=======
	// Cache directories.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	".cache":                  true,
	"__MACOSX":                true,
	".DS_Store":               true,
	".Trash":                  true,
	"Caches":                  true,
	".Spotlight-V100":         true,
	".fseventsd":              true,
	".DocumentRevisions-V100": true,
	".TemporaryItems":         true,
	"$RECYCLE.BIN":            true,
	".temp":                   true,
	".tmp":                    true,
	"_temp":                   true,
	"_tmp":                    true,
	".Homebrew":               true,
	".rustup":                 true,
	".sdkman":                 true,
	".nvm":                    true,

<<<<<<< HEAD
	// macOS specific
	"Application Scripts":     true,
	"Saved Application State": true,

	// iCloud
	"Mobile Documents": true,

	// Docker & Containers
	".docker":     true,
	".containerd": true,

	// Mobile development
=======
	// macOS.
	"Application Scripts":     true,
	"Saved Application State": true,

	// iCloud.
	"Mobile Documents": true,

	// Containers.
	".docker":     true,
	".containerd": true,

	// Mobile development.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"Pods":        true,
	"DerivedData": true,
	".build":      true,
	"xcuserdata":  true,
	"Carthage":    true,
<<<<<<< HEAD

	// Web frameworks
=======
	".dart_tool":  true,

	// Web frameworks.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	".angular":    true,
	".svelte-kit": true,
	".astro":      true,
	".solid":      true,

<<<<<<< HEAD
	// Databases
=======
	// Databases.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	".mysql":    true,
	".postgres": true,
	"mongodb":   true,

<<<<<<< HEAD
	// Other
=======
	// Other.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	".terraform": true,
	".vagrant":   true,
	"tmp":        true,
	"temp":       true,
}

var skipSystemDirs = map[string]bool{
	"dev":                     true,
	"tmp":                     true,
	"private":                 true,
	"cores":                   true,
	"net":                     true,
	"home":                    true,
	"System":                  true,
	"sbin":                    true,
	"bin":                     true,
	"etc":                     true,
	"var":                     true,
<<<<<<< HEAD
=======
	"opt":                     false,
	"usr":                     false,
	"Volumes":                 true,
	"Network":                 true,
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	".vol":                    true,
	".Spotlight-V100":         true,
	".fseventsd":              true,
	".DocumentRevisions-V100": true,
	".TemporaryItems":         true,
<<<<<<< HEAD
=======
	".MobileBackups":          true,
}

var defaultSkipDirs = map[string]bool{
	"nfs":         true,
	"PHD":         true,
	"Permissions": true,
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

var skipExtensions = map[string]bool{
	".go":     true,
	".js":     true,
	".ts":     true,
	".tsx":    true,
	".jsx":    true,
	".json":   true,
	".md":     true,
	".txt":    true,
	".yml":    true,
	".yaml":   true,
	".xml":    true,
	".html":   true,
	".css":    true,
	".scss":   true,
	".sass":   true,
	".less":   true,
	".py":     true,
	".rb":     true,
	".java":   true,
	".kt":     true,
	".rs":     true,
	".swift":  true,
	".m":      true,
	".mm":     true,
	".c":      true,
	".cpp":    true,
	".h":      true,
	".hpp":    true,
	".cs":     true,
	".sql":    true,
	".db":     true,
	".lock":   true,
	".gradle": true,
	".mjs":    true,
	".cjs":    true,
	".coffee": true,
	".dart":   true,
	".svelte": true,
	".vue":    true,
	".nim":    true,
	".hx":     true,
}

var spinnerFrames = []string{"|", "/", "-", "\\", "|", "/", "-", "\\"}

const (
<<<<<<< HEAD
	colorPurple = "\033[0;35m"
	colorGray   = "\033[0;90m"
	colorRed    = "\033[0;31m"
	colorYellow = "\033[1;33m"
	colorGreen  = "\033[0;32m"
	colorCyan   = "\033[0;36m"
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
=======
	colorPurple     = "\033[0;35m"
	colorPurpleBold = "\033[1;35m"
	colorGray       = "\033[0;90m"
	colorRed        = "\033[0;31m"
	colorYellow     = "\033[0;33m"
	colorGreen      = "\033[0;32m"
	colorBlue       = "\033[0;34m"
	colorCyan       = "\033[0;36m"
	colorReset      = "\033[0m"
	colorBold       = "\033[1m"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
)
