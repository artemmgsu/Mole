//go:build darwin

package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
<<<<<<< HEAD
=======
	"slices"
	"sort"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"strings"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type dirEntry struct {
	Name       string
	Path       string
	Size       int64
	IsDir      bool
	LastAccess time.Time
}

type fileEntry struct {
	Name string
	Path string
	Size int64
}

type scanResult struct {
	Entries    []dirEntry
	LargeFiles []fileEntry
	TotalSize  int64
<<<<<<< HEAD
=======
	TotalFiles int64
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

type cacheEntry struct {
	Entries    []dirEntry
	LargeFiles []fileEntry
	TotalSize  int64
<<<<<<< HEAD
=======
	TotalFiles int64
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	ModTime    time.Time
	ScanTime   time.Time
}

type historyEntry struct {
	Path          string
	Entries       []dirEntry
	LargeFiles    []fileEntry
	TotalSize     int64
<<<<<<< HEAD
=======
	TotalFiles    int64
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	Selected      int
	EntryOffset   int
	LargeSelected int
	LargeOffset   int
	Dirty         bool
<<<<<<< HEAD
=======
	IsOverview    bool
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

type scanResultMsg struct {
	result scanResult
	err    error
}

type overviewSizeMsg struct {
	Path  string
	Index int
	Size  int64
	Err   error
}

type tickMsg time.Time

type deleteProgressMsg struct {
	done  bool
	err   error
	count int64
	path  string
}

type model struct {
	path                 string
	history              []historyEntry
	entries              []dirEntry
	largeFiles           []fileEntry
	selected             int
	offset               int
	status               string
	totalSize            int64
	scanning             bool
	spinner              int
	filesScanned         *int64
	dirsScanned          *int64
	bytesScanned         *int64
<<<<<<< HEAD
	currentPath          *string
=======
	currentPath          *atomic.Value
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	showLargeFiles       bool
	isOverview           bool
	deleteConfirm        bool
	deleteTarget         *dirEntry
	deleting             bool
	deleteCount          *int64
	cache                map[string]historyEntry
	largeSelected        int
	largeOffset          int
	overviewSizeCache    map[string]int64
	overviewFilesScanned *int64
	overviewDirsScanned  *int64
	overviewBytesScanned *int64
	overviewCurrentPath  *string
	overviewScanning     bool
	overviewScanningSet  map[string]bool // Track which paths are currently being scanned
	width                int             // Terminal width
	height               int             // Terminal height
<<<<<<< HEAD
=======
	multiSelected        map[string]bool // Track multi-selected items by path (safer than index)
	largeMultiSelected   map[string]bool // Track multi-selected large files by path (safer than index)
	totalFiles           int64           // Total files found in current/last scan
	lastTotalFiles       int64           // Total files from previous scan (for progress bar)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

func (m model) inOverviewMode() bool {
	return m.isOverview && m.path == "/"
}

func main() {
	target := os.Getenv("MO_ANALYZE_PATH")
	if target == "" && len(os.Args) > 1 {
		target = os.Args[1]
	}

	var abs string
	var isOverview bool

	if target == "" {
<<<<<<< HEAD
		// Default to overview mode
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		isOverview = true
		abs = "/"
	} else {
		var err error
		abs, err = filepath.Abs(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot resolve %q: %v\n", target, err)
			os.Exit(1)
		}
		isOverview = false
	}

<<<<<<< HEAD
	// Prefetch overview cache in background (non-blocking)
	go prefetchOverviewCache()

	p := tea.NewProgram(newModel(abs, isOverview), tea.WithAltScreen())
	if err := p.Start(); err != nil {
=======
	// Warm overview cache in background.
	prefetchCtx, prefetchCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer prefetchCancel()
	go prefetchOverviewCache(prefetchCtx)

	p := tea.NewProgram(newModel(abs, isOverview), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		fmt.Fprintf(os.Stderr, "analyzer error: %v\n", err)
		os.Exit(1)
	}
}

func newModel(path string, isOverview bool) model {
	var filesScanned, dirsScanned, bytesScanned int64
<<<<<<< HEAD
	currentPath := ""
=======
	currentPath := &atomic.Value{}
	currentPath.Store("")
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	var overviewFilesScanned, overviewDirsScanned, overviewBytesScanned int64
	overviewCurrentPath := ""

	m := model{
		path:                 path,
		selected:             0,
		status:               "Preparing scan...",
		scanning:             !isOverview,
		filesScanned:         &filesScanned,
		dirsScanned:          &dirsScanned,
		bytesScanned:         &bytesScanned,
<<<<<<< HEAD
		currentPath:          &currentPath,
=======
		currentPath:          currentPath,
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		showLargeFiles:       false,
		isOverview:           isOverview,
		cache:                make(map[string]historyEntry),
		overviewFilesScanned: &overviewFilesScanned,
		overviewDirsScanned:  &overviewDirsScanned,
		overviewBytesScanned: &overviewBytesScanned,
		overviewCurrentPath:  &overviewCurrentPath,
		overviewSizeCache:    make(map[string]int64),
		overviewScanningSet:  make(map[string]bool),
<<<<<<< HEAD
	}

	// In overview mode, create shortcut entries
=======
		multiSelected:        make(map[string]bool),
		largeMultiSelected:   make(map[string]bool),
	}

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if isOverview {
		m.scanning = false
		m.hydrateOverviewEntries()
		m.selected = 0
		m.offset = 0
		if nextPendingOverviewIndex(m.entries) >= 0 {
			m.overviewScanning = true
			m.status = "Checking system folders..."
		} else {
			m.status = "Ready"
		}
	}

<<<<<<< HEAD
=======
	// Try to peek last total files for progress bar, even if cache is stale
	if !isOverview {
		if total, err := peekCacheTotalFiles(path); err == nil && total > 0 {
			m.lastTotalFiles = total
		}
	}

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	return m
}

func createOverviewEntries() []dirEntry {
	home := os.Getenv("HOME")
	entries := []dirEntry{}

<<<<<<< HEAD
	if home != "" {
		entries = append(entries,
			dirEntry{Name: "Home (~)", Path: home, IsDir: true, Size: -1},
			dirEntry{Name: "Library (~/Library)", Path: filepath.Join(home, "Library"), IsDir: true, Size: -1},
		)
=======
	// Separate Home and ~/Library to avoid double counting.
	if home != "" {
		entries = append(entries, dirEntry{Name: "Home", Path: home, IsDir: true, Size: -1})

		userLibrary := filepath.Join(home, "Library")
		if _, err := os.Stat(userLibrary); err == nil {
			entries = append(entries, dirEntry{Name: "App Library", Path: userLibrary, IsDir: true, Size: -1})
		}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	}

	entries = append(entries,
		dirEntry{Name: "Applications", Path: "/Applications", IsDir: true, Size: -1},
		dirEntry{Name: "System Library", Path: "/Library", IsDir: true, Size: -1},
	)

<<<<<<< HEAD
	// Add Volumes shortcut only when it contains real mounted folders (e.g., external disks)
=======
	// Include Volumes only when real mounts exist.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if hasUsefulVolumeMounts("/Volumes") {
		entries = append(entries, dirEntry{Name: "Volumes", Path: "/Volumes", IsDir: true, Size: -1})
	}

	return entries
}

func hasUsefulVolumeMounts(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		name := entry.Name()
<<<<<<< HEAD
		// Skip hidden control entries for Spotlight/TimeMachine etc.
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if strings.HasPrefix(name, ".") {
			continue
		}

		info, err := os.Lstat(filepath.Join(path, name))
		if err != nil {
			continue
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			continue // Ignore the synthetic MacintoshHD link
		}
		if info.IsDir() {
			return true
		}
	}
	return false
}

func (m *model) hydrateOverviewEntries() {
	m.entries = createOverviewEntries()
	if m.overviewSizeCache == nil {
		m.overviewSizeCache = make(map[string]int64)
	}
	for i := range m.entries {
		if size, ok := m.overviewSizeCache[m.entries[i].Path]; ok {
			m.entries[i].Size = size
			continue
		}
		if size, err := loadOverviewCachedSize(m.entries[i].Path); err == nil {
			m.entries[i].Size = size
			m.overviewSizeCache[m.entries[i].Path] = size
		}
	}
	m.totalSize = sumKnownEntrySizes(m.entries)
}

<<<<<<< HEAD
=======
func (m *model) sortOverviewEntriesBySize() {
	// Stable sort by size.
	sort.SliceStable(m.entries, func(i, j int) bool {
		return m.entries[i].Size > m.entries[j].Size
	})
}

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
func (m *model) scheduleOverviewScans() tea.Cmd {
	if !m.inOverviewMode() {
		return nil
	}

<<<<<<< HEAD
	// Find pending entries (not scanned and not currently scanning)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	var pendingIndices []int
	for i, entry := range m.entries {
		if entry.Size < 0 && !m.overviewScanningSet[entry.Path] {
			pendingIndices = append(pendingIndices, i)
			if len(pendingIndices) >= maxConcurrentOverview {
				break
			}
		}
	}

<<<<<<< HEAD
	// No more work to do
	if len(pendingIndices) == 0 {
		m.overviewScanning = false
		if !hasPendingOverviewEntries(m.entries) {
=======
	if len(pendingIndices) == 0 {
		m.overviewScanning = false
		if !hasPendingOverviewEntries(m.entries) {
			m.sortOverviewEntriesBySize()
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			m.status = "Ready"
		}
		return nil
	}

<<<<<<< HEAD
	// Mark all as scanning
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	var cmds []tea.Cmd
	for _, idx := range pendingIndices {
		entry := m.entries[idx]
		m.overviewScanningSet[entry.Path] = true
		cmd := scanOverviewPathCmd(entry.Path, idx)
		cmds = append(cmds, cmd)
	}

	m.overviewScanning = true
	remaining := 0
	for _, e := range m.entries {
		if e.Size < 0 {
			remaining++
		}
	}
	if len(pendingIndices) > 0 {
		firstEntry := m.entries[pendingIndices[0]]
		if len(pendingIndices) == 1 {
<<<<<<< HEAD
			m.status = fmt.Sprintf("Scanning %s... (%d left)", firstEntry.Name, remaining)
		} else {
			m.status = fmt.Sprintf("Scanning %d directories... (%d left)", len(pendingIndices), remaining)
=======
			m.status = fmt.Sprintf("Scanning %s..., %d left", firstEntry.Name, remaining)
		} else {
			m.status = fmt.Sprintf("Scanning %d directories..., %d left", len(pendingIndices), remaining)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		}
	}

	cmds = append(cmds, tickCmd())
	return tea.Batch(cmds...)
}

func (m *model) getScanProgress() (files, dirs, bytes int64) {
	if m.filesScanned != nil {
		files = atomic.LoadInt64(m.filesScanned)
	}
	if m.dirsScanned != nil {
		dirs = atomic.LoadInt64(m.dirsScanned)
	}
	if m.bytesScanned != nil {
		bytes = atomic.LoadInt64(m.bytesScanned)
	}
	return
}

func (m model) Init() tea.Cmd {
	if m.inOverviewMode() {
		return m.scheduleOverviewScans()
	}
	return tea.Batch(m.scanCmd(m.path), tickCmd())
}

func (m model) scanCmd(path string) tea.Cmd {
	return func() tea.Msg {
<<<<<<< HEAD
		// Try to load from persistent cache first
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if cached, err := loadCacheFromDisk(path); err == nil {
			result := scanResult{
				Entries:    cached.Entries,
				LargeFiles: cached.LargeFiles,
				TotalSize:  cached.TotalSize,
<<<<<<< HEAD
=======
				TotalFiles: 0, // Cache doesn't store file count currently, minor UI limitation
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			}
			return scanResultMsg{result: result, err: nil}
		}

<<<<<<< HEAD
		// Use singleflight to avoid duplicate scans of the same path
		// If multiple goroutines request the same path, only one scan will be performed
		v, err, _ := scanGroup.Do(path, func() (interface{}, error) {
=======
		v, err, _ := scanGroup.Do(path, func() (any, error) {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			return scanPathConcurrent(path, m.filesScanned, m.dirsScanned, m.bytesScanned, m.currentPath)
		})

		if err != nil {
			return scanResultMsg{err: err}
		}

		result := v.(scanResult)

<<<<<<< HEAD
		// Save to persistent cache asynchronously with error logging
		go func(p string, r scanResult) {
			if err := saveCacheToDisk(p, r); err != nil {
				// Log error but don't fail the scan
=======
		go func(p string, r scanResult) {
			if err := saveCacheToDisk(p, r); err != nil {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				_ = err // Cache save failure is not critical
			}
		}(path, result)

		return scanResultMsg{result: result, err: nil}
	}
}

func tickCmd() tea.Cmd {
<<<<<<< HEAD
	return tea.Tick(time.Millisecond*120, func(t time.Time) tea.Msg {
=======
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.updateKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case deleteProgressMsg:
		if msg.done {
			m.deleting = false
<<<<<<< HEAD
=======
			m.multiSelected = make(map[string]bool)
			m.largeMultiSelected = make(map[string]bool)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			if msg.err != nil {
				m.status = fmt.Sprintf("Failed to delete: %v", msg.err)
			} else {
				if msg.path != "" {
					m.removePathFromView(msg.path)
					invalidateCache(msg.path)
				}
				invalidateCache(m.path)
				m.status = fmt.Sprintf("Deleted %d items", msg.count)
<<<<<<< HEAD
				// Mark all caches as dirty
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				for i := range m.history {
					m.history[i].Dirty = true
				}
				for path := range m.cache {
					entry := m.cache[path]
					entry.Dirty = true
					m.cache[path] = entry
				}
<<<<<<< HEAD
				// Refresh the view
				m.scanning = true
				// Reset scan counters for rescan
=======
				m.scanning = true
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				atomic.StoreInt64(m.filesScanned, 0)
				atomic.StoreInt64(m.dirsScanned, 0)
				atomic.StoreInt64(m.bytesScanned, 0)
				if m.currentPath != nil {
<<<<<<< HEAD
					*m.currentPath = ""
=======
					m.currentPath.Store("")
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				}
				return m, tea.Batch(m.scanCmd(m.path), tickCmd())
			}
		}
		return m, nil
	case scanResultMsg:
		m.scanning = false
		if msg.err != nil {
			m.status = fmt.Sprintf("Scan failed: %v", msg.err)
			return m, nil
		}
<<<<<<< HEAD
		m.entries = msg.result.Entries
		m.largeFiles = msg.result.LargeFiles
		m.totalSize = msg.result.TotalSize
=======
		filteredEntries := make([]dirEntry, 0, len(msg.result.Entries))
		for _, e := range msg.result.Entries {
			if e.Size > 0 {
				filteredEntries = append(filteredEntries, e)
			}
		}
		m.entries = filteredEntries
		m.largeFiles = msg.result.LargeFiles
		m.totalSize = msg.result.TotalSize
		m.totalFiles = msg.result.TotalFiles
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		m.status = fmt.Sprintf("Scanned %s", humanizeBytes(m.totalSize))
		m.clampEntrySelection()
		m.clampLargeSelection()
		m.cache[m.path] = cacheSnapshot(m)
		if m.totalSize > 0 {
			if m.overviewSizeCache == nil {
				m.overviewSizeCache = make(map[string]int64)
			}
			m.overviewSizeCache[m.path] = m.totalSize
			go func(path string, size int64) {
				_ = storeOverviewSize(path, size)
			}(m.path, m.totalSize)
		}
		return m, nil
	case overviewSizeMsg:
<<<<<<< HEAD
		// Remove from scanning set
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		delete(m.overviewScanningSet, msg.Path)

		if msg.Err == nil {
			if m.overviewSizeCache == nil {
				m.overviewSizeCache = make(map[string]int64)
			}
			m.overviewSizeCache[msg.Path] = msg.Size
		}

		if m.inOverviewMode() {
<<<<<<< HEAD
			// Update entry with result
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			for i := range m.entries {
				if m.entries[i].Path == msg.Path {
					if msg.Err == nil {
						m.entries[i].Size = msg.Size
					} else {
						m.entries[i].Size = 0
					}
					break
				}
			}
			m.totalSize = sumKnownEntrySizes(m.entries)

<<<<<<< HEAD
			// Show error briefly if any
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			if msg.Err != nil {
				m.status = fmt.Sprintf("Unable to measure %s: %v", displayPath(msg.Path), msg.Err)
			}

<<<<<<< HEAD
			// Schedule next batch of scans
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			cmd := m.scheduleOverviewScans()
			return m, cmd
		}
		return m, nil
	case tickMsg:
<<<<<<< HEAD
		// Keep spinner running if scanning or deleting or if there are pending overview items
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		hasPending := false
		if m.inOverviewMode() {
			for _, entry := range m.entries {
				if entry.Size < 0 {
					hasPending = true
					break
				}
			}
		}
		if m.scanning || m.deleting || (m.inOverviewMode() && (m.overviewScanning || hasPending)) {
			m.spinner = (m.spinner + 1) % len(spinnerFrames)
<<<<<<< HEAD
			// Update delete progress status
			if m.deleting && m.deleteCount != nil {
				count := atomic.LoadInt64(m.deleteCount)
				if count > 0 {
					m.status = fmt.Sprintf("Deleting... %s items removed", formatNumber(count))
=======
			if m.deleting && m.deleteCount != nil {
				count := atomic.LoadInt64(m.deleteCount)
				if count > 0 {
					m.status = fmt.Sprintf("Moving to Trash... %s items", formatNumber(count))
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				}
			}
			return m, tickCmd()
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
<<<<<<< HEAD
	// Handle delete confirmation
	if m.deleteConfirm {
		if msg.String() == "delete" || msg.String() == "backspace" {
			// Confirm delete - start async deletion
			if m.deleteTarget != nil {
				m.deleteConfirm = false
				m.deleting = true
				var deleteCount int64
				m.deleteCount = &deleteCount
				targetPath := m.deleteTarget.Path
				targetName := m.deleteTarget.Name
				m.deleteTarget = nil
				m.status = fmt.Sprintf("Deleting %s...", targetName)
				return m, tea.Batch(deletePathCmd(targetPath, m.deleteCount), tickCmd())
			}
			m.deleteConfirm = false
			m.deleteTarget = nil
			return m, nil
		} else if msg.String() == "esc" || msg.String() == "q" {
			// Cancel delete with ESC or Q
=======
	// Delete confirm flow.
	if m.deleteConfirm {
		switch msg.String() {
		case "enter":
			m.deleteConfirm = false
			m.deleting = true
			var deleteCount int64
			m.deleteCount = &deleteCount

			// Collect paths (safer than indices).
			var pathsToDelete []string
			if m.showLargeFiles {
				if len(m.largeMultiSelected) > 0 {
					for path := range m.largeMultiSelected {
						pathsToDelete = append(pathsToDelete, path)
					}
				} else if m.deleteTarget != nil {
					pathsToDelete = append(pathsToDelete, m.deleteTarget.Path)
				}
			} else {
				if len(m.multiSelected) > 0 {
					for path := range m.multiSelected {
						pathsToDelete = append(pathsToDelete, path)
					}
				} else if m.deleteTarget != nil {
					pathsToDelete = append(pathsToDelete, m.deleteTarget.Path)
				}
			}

			m.deleteTarget = nil
			if len(pathsToDelete) == 0 {
				m.deleting = false
				m.status = "Nothing to delete"
				return m, nil
			}

			if len(pathsToDelete) == 1 {
				targetPath := pathsToDelete[0]
				m.status = fmt.Sprintf("Deleting %s...", filepath.Base(targetPath))
				return m, tea.Batch(deletePathCmd(targetPath, m.deleteCount), tickCmd())
			}

			m.status = fmt.Sprintf("Deleting %d items...", len(pathsToDelete))
			return m, tea.Batch(deleteMultiplePathsCmd(pathsToDelete, m.deleteCount), tickCmd())
		case "esc", "q":
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			m.status = "Cancelled"
			m.deleteConfirm = false
			m.deleteTarget = nil
			return m, nil
<<<<<<< HEAD
		} else {
			// Any other key also cancels
			m.status = "Cancelled"
			m.deleteConfirm = false
			m.deleteTarget = nil
=======
		default:
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			return m, nil
		}
	}

	switch msg.String() {
<<<<<<< HEAD
	case "q", "ctrl+c":
=======
	case "q", "ctrl+c", "Q":
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		return m, tea.Quit
	case "esc":
		if m.showLargeFiles {
			m.showLargeFiles = false
			return m, nil
		}
		return m, tea.Quit
<<<<<<< HEAD
	case "up", "k":
=======
	case "up", "k", "K":
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if m.showLargeFiles {
			if m.largeSelected > 0 {
				m.largeSelected--
				if m.largeSelected < m.largeOffset {
					m.largeOffset = m.largeSelected
				}
			}
		} else if len(m.entries) > 0 && m.selected > 0 {
			m.selected--
			if m.selected < m.offset {
				m.offset = m.selected
			}
		}
<<<<<<< HEAD
	case "down", "j":
=======
	case "down", "j", "J":
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if m.showLargeFiles {
			if m.largeSelected < len(m.largeFiles)-1 {
				m.largeSelected++
				viewport := calculateViewport(m.height, true)
				if m.largeSelected >= m.largeOffset+viewport {
					m.largeOffset = m.largeSelected - viewport + 1
				}
			}
		} else if len(m.entries) > 0 && m.selected < len(m.entries)-1 {
			m.selected++
			viewport := calculateViewport(m.height, false)
			if m.selected >= m.offset+viewport {
				m.offset = m.selected - viewport + 1
			}
		}
<<<<<<< HEAD
	case "enter", "right", "l":
=======
	case "enter", "right", "l", "L":
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if m.showLargeFiles {
			return m, nil
		}
		return m.enterSelectedDir()
<<<<<<< HEAD
	case "b", "left", "h":
=======
	case "b", "left", "h", "B", "H":
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if m.showLargeFiles {
			m.showLargeFiles = false
			return m, nil
		}
		if len(m.history) == 0 {
<<<<<<< HEAD
			// Return to overview if at top level
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			if !m.inOverviewMode() {
				return m, m.switchToOverviewMode()
			}
			return m, nil
		}
		last := m.history[len(m.history)-1]
		m.history = m.history[:len(m.history)-1]
		m.path = last.Path
		m.selected = last.Selected
		m.offset = last.EntryOffset
		m.largeSelected = last.LargeSelected
		m.largeOffset = last.LargeOffset
<<<<<<< HEAD
		m.isOverview = false
		if last.Dirty {
=======
		m.isOverview = last.IsOverview
		if last.Dirty {
			// On overview return, refresh cached entries.
			if last.IsOverview {
				m.hydrateOverviewEntries()
				m.totalSize = sumKnownEntrySizes(m.entries)
				m.status = "Ready"
				m.scanning = false
				if nextPendingOverviewIndex(m.entries) >= 0 {
					m.overviewScanning = true
					return m, m.scheduleOverviewScans()
				}
				return m, nil
			}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			m.status = "Scanning..."
			m.scanning = true
			return m, tea.Batch(m.scanCmd(m.path), tickCmd())
		}
		m.entries = last.Entries
		m.largeFiles = last.LargeFiles
		m.totalSize = last.TotalSize
		m.clampEntrySelection()
		m.clampLargeSelection()
		if len(m.entries) == 0 {
			m.selected = 0
		} else if m.selected >= len(m.entries) {
			m.selected = len(m.entries) - 1
		}
		if m.selected < 0 {
			m.selected = 0
		}
		m.status = fmt.Sprintf("Scanned %s", humanizeBytes(m.totalSize))
		m.scanning = false
		return m, nil
<<<<<<< HEAD
	case "r":
		// Invalidate cache before rescanning to ensure fresh data
		invalidateCache(m.path)
		m.status = "Refreshing..."
		m.scanning = true
		// Reset scan counters for refresh
=======
	case "r", "R":
		m.multiSelected = make(map[string]bool)
		m.largeMultiSelected = make(map[string]bool)

		if m.inOverviewMode() {
			// Explicitly invalidate cache for all overview entries to force re-scan
			for _, entry := range m.entries {
				invalidateCache(entry.Path)
			}

			m.overviewSizeCache = make(map[string]int64)
			m.overviewScanningSet = make(map[string]bool)
			m.hydrateOverviewEntries() // Reset sizes to pending

			for i := range m.entries {
				m.entries[i].Size = -1
			}
			m.totalSize = 0

			m.status = "Refreshing..."
			m.overviewScanning = true
			return m, tea.Batch(m.scheduleOverviewScans(), tickCmd())
		}

		invalidateCache(m.path)
		m.status = "Refreshing..."
		m.scanning = true
		if m.totalFiles > 0 {
			m.lastTotalFiles = m.totalFiles
		}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		atomic.StoreInt64(m.filesScanned, 0)
		atomic.StoreInt64(m.dirsScanned, 0)
		atomic.StoreInt64(m.bytesScanned, 0)
		if m.currentPath != nil {
<<<<<<< HEAD
			*m.currentPath = ""
		}
		return m, tea.Batch(m.scanCmd(m.path), tickCmd())
	case "L":
		m.showLargeFiles = !m.showLargeFiles
		if m.showLargeFiles {
			m.largeSelected = 0
			m.largeOffset = 0
		}
	case "o":
		// Open selected entry
		if m.showLargeFiles {
			if len(m.largeFiles) > 0 {
				selected := m.largeFiles[m.largeSelected]
=======
			m.currentPath.Store("")
		}
		return m, tea.Batch(m.scanCmd(m.path), tickCmd())
	case "t", "T":
		if !m.inOverviewMode() {
			m.showLargeFiles = !m.showLargeFiles
			if m.showLargeFiles {
				m.largeSelected = 0
				m.largeOffset = 0
				m.largeMultiSelected = make(map[string]bool)
			} else {
				m.multiSelected = make(map[string]bool)
			}
			m.status = fmt.Sprintf("Scanned %s", humanizeBytes(m.totalSize))
		}
	case "o", "O":
		// Open selected entries (multi-select aware).
		const maxBatchOpen = 20
		if m.showLargeFiles {
			if len(m.largeFiles) > 0 {
				if len(m.largeMultiSelected) > 0 {
					count := len(m.largeMultiSelected)
					if count > maxBatchOpen {
						m.status = fmt.Sprintf("Too many items to open, max %d, selected %d", maxBatchOpen, count)
						return m, nil
					}
					for path := range m.largeMultiSelected {
						go func(p string) {
							ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
							defer cancel()
							_ = exec.CommandContext(ctx, "open", p).Run()
						}(path)
					}
					m.status = fmt.Sprintf("Opening %d items...", count)
				} else {
					selected := m.largeFiles[m.largeSelected]
					go func(path string) {
						ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
						defer cancel()
						_ = exec.CommandContext(ctx, "open", path).Run()
					}(selected.Path)
					m.status = fmt.Sprintf("Opening %s...", selected.Name)
				}
			}
		} else if len(m.entries) > 0 {
			if len(m.multiSelected) > 0 {
				count := len(m.multiSelected)
				if count > maxBatchOpen {
					m.status = fmt.Sprintf("Too many items to open, max %d, selected %d", maxBatchOpen, count)
					return m, nil
				}
				for path := range m.multiSelected {
					go func(p string) {
						ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
						defer cancel()
						_ = exec.CommandContext(ctx, "open", p).Run()
					}(path)
				}
				m.status = fmt.Sprintf("Opening %d items...", count)
			} else {
				selected := m.entries[m.selected]
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				go func(path string) {
					ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
					defer cancel()
					_ = exec.CommandContext(ctx, "open", path).Run()
				}(selected.Path)
				m.status = fmt.Sprintf("Opening %s...", selected.Name)
			}
<<<<<<< HEAD
		} else if len(m.entries) > 0 {
			selected := m.entries[m.selected]
			go func(path string) {
				ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
				defer cancel()
				_ = exec.CommandContext(ctx, "open", path).Run()
			}(selected.Path)
			m.status = fmt.Sprintf("Opening %s...", selected.Name)
		}
	case "f", "F":
		// Reveal selected entry in Finder
		if m.showLargeFiles {
			if len(m.largeFiles) > 0 {
				selected := m.largeFiles[m.largeSelected]
=======
		}
	case "f", "F":
		// Reveal in Finder (multi-select aware).
		const maxBatchReveal = 20
		if m.showLargeFiles {
			if len(m.largeFiles) > 0 {
				if len(m.largeMultiSelected) > 0 {
					count := len(m.largeMultiSelected)
					if count > maxBatchReveal {
						m.status = fmt.Sprintf("Too many items to reveal, max %d, selected %d", maxBatchReveal, count)
						return m, nil
					}
					for path := range m.largeMultiSelected {
						go func(p string) {
							ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
							defer cancel()
							_ = exec.CommandContext(ctx, "open", "-R", p).Run()
						}(path)
					}
					m.status = fmt.Sprintf("Showing %d items in Finder...", count)
				} else {
					selected := m.largeFiles[m.largeSelected]
					go func(path string) {
						ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
						defer cancel()
						_ = exec.CommandContext(ctx, "open", "-R", path).Run()
					}(selected.Path)
					m.status = fmt.Sprintf("Showing %s in Finder...", selected.Name)
				}
			}
		} else if len(m.entries) > 0 {
			if len(m.multiSelected) > 0 {
				count := len(m.multiSelected)
				if count > maxBatchReveal {
					m.status = fmt.Sprintf("Too many items to reveal, max %d, selected %d", maxBatchReveal, count)
					return m, nil
				}
				for path := range m.multiSelected {
					go func(p string) {
						ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
						defer cancel()
						_ = exec.CommandContext(ctx, "open", "-R", p).Run()
					}(path)
				}
				m.status = fmt.Sprintf("Showing %d items in Finder...", count)
			} else {
				selected := m.entries[m.selected]
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				go func(path string) {
					ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
					defer cancel()
					_ = exec.CommandContext(ctx, "open", "-R", path).Run()
				}(selected.Path)
				m.status = fmt.Sprintf("Showing %s in Finder...", selected.Name)
			}
<<<<<<< HEAD
		} else if len(m.entries) > 0 {
			selected := m.entries[m.selected]
			go func(path string) {
				ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
				defer cancel()
				_ = exec.CommandContext(ctx, "open", "-R", path).Run()
			}(selected.Path)
			m.status = fmt.Sprintf("Showing %s in Finder...", selected.Name)
		}
	case "delete", "backspace":
		// Delete selected file or directory
		if m.showLargeFiles {
			if len(m.largeFiles) > 0 {
				selected := m.largeFiles[m.largeSelected]
				m.deleteConfirm = true
				m.deleteTarget = &dirEntry{
					Name:  selected.Name,
					Path:  selected.Path,
					Size:  selected.Size,
					IsDir: false,
				}
			}
		} else if len(m.entries) > 0 && !m.inOverviewMode() {
			selected := m.entries[m.selected]
			m.deleteConfirm = true
			m.deleteTarget = &selected
=======
		}
	case " ":
		// Toggle multi-select (paths as keys).
		if m.showLargeFiles {
			if len(m.largeFiles) > 0 && m.largeSelected < len(m.largeFiles) {
				if m.largeMultiSelected == nil {
					m.largeMultiSelected = make(map[string]bool)
				}
				selectedPath := m.largeFiles[m.largeSelected].Path
				if m.largeMultiSelected[selectedPath] {
					delete(m.largeMultiSelected, selectedPath)
				} else {
					m.largeMultiSelected[selectedPath] = true
				}
				count := len(m.largeMultiSelected)
				if count > 0 {
					var totalSize int64
					for path := range m.largeMultiSelected {
						for _, file := range m.largeFiles {
							if file.Path == path {
								totalSize += file.Size
								break
							}
						}
					}
					m.status = fmt.Sprintf("%d selected, %s", count, humanizeBytes(totalSize))
				} else {
					m.status = fmt.Sprintf("Scanned %s", humanizeBytes(m.totalSize))
				}
			}
		} else if len(m.entries) > 0 && !m.inOverviewMode() && m.selected < len(m.entries) {
			if m.multiSelected == nil {
				m.multiSelected = make(map[string]bool)
			}
			selectedPath := m.entries[m.selected].Path
			if m.multiSelected[selectedPath] {
				delete(m.multiSelected, selectedPath)
			} else {
				m.multiSelected[selectedPath] = true
			}
			count := len(m.multiSelected)
			if count > 0 {
				var totalSize int64
				for path := range m.multiSelected {
					for _, entry := range m.entries {
						if entry.Path == path {
							totalSize += entry.Size
							break
						}
					}
				}
				m.status = fmt.Sprintf("%d selected, %s", count, humanizeBytes(totalSize))
			} else {
				m.status = fmt.Sprintf("Scanned %s", humanizeBytes(m.totalSize))
			}
		}
	case "delete", "backspace":
		if m.showLargeFiles {
			if len(m.largeFiles) > 0 {
				if len(m.largeMultiSelected) > 0 {
					m.deleteConfirm = true
					for path := range m.largeMultiSelected {
						for _, file := range m.largeFiles {
							if file.Path == path {
								m.deleteTarget = &dirEntry{
									Name:  file.Name,
									Path:  file.Path,
									Size:  file.Size,
									IsDir: false,
								}
								break
							}
						}
						break // Only need first one for display
					}
				} else if m.largeSelected < len(m.largeFiles) {
					selected := m.largeFiles[m.largeSelected]
					m.deleteConfirm = true
					m.deleteTarget = &dirEntry{
						Name:  selected.Name,
						Path:  selected.Path,
						Size:  selected.Size,
						IsDir: false,
					}
				}
			}
		} else if len(m.entries) > 0 && !m.inOverviewMode() {
			if len(m.multiSelected) > 0 {
				m.deleteConfirm = true
				for path := range m.multiSelected {
					// Resolve entry by path.
					for i := range m.entries {
						if m.entries[i].Path == path {
							m.deleteTarget = &m.entries[i]
							break
						}
					}
					break // Only need first one for display
				}
			} else if m.selected < len(m.entries) {
				selected := m.entries[m.selected]
				m.deleteConfirm = true
				m.deleteTarget = &selected
			}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		}
	}
	return m, nil
}

func (m *model) switchToOverviewMode() tea.Cmd {
	m.isOverview = true
	m.path = "/"
	m.scanning = false
	m.showLargeFiles = false
	m.largeFiles = nil
	m.largeSelected = 0
	m.largeOffset = 0
	m.deleteConfirm = false
	m.deleteTarget = nil
	m.selected = 0
	m.offset = 0
	m.hydrateOverviewEntries()
	cmd := m.scheduleOverviewScans()
	if cmd == nil {
		m.status = "Ready"
		return nil
	}
<<<<<<< HEAD
	// Start tick to animate spinner while scanning
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	return tea.Batch(cmd, tickCmd())
}

func (m model) enterSelectedDir() (tea.Model, tea.Cmd) {
	if len(m.entries) == 0 {
		return m, nil
	}
	selected := m.entries[m.selected]
	if selected.IsDir {
<<<<<<< HEAD
		if !m.inOverviewMode() {
=======
		if len(m.history) == 0 || m.history[len(m.history)-1].Path != m.path {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			m.history = append(m.history, snapshotFromModel(m))
		}
		m.path = selected.Path
		m.selected = 0
		m.offset = 0
		m.status = "Scanning..."
		m.scanning = true
		m.isOverview = false
<<<<<<< HEAD

		// Reset scan counters for new scan
=======
		m.multiSelected = make(map[string]bool)
		m.largeMultiSelected = make(map[string]bool)

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		atomic.StoreInt64(m.filesScanned, 0)
		atomic.StoreInt64(m.dirsScanned, 0)
		atomic.StoreInt64(m.bytesScanned, 0)
		if m.currentPath != nil {
<<<<<<< HEAD
			*m.currentPath = ""
		}

		if cached, ok := m.cache[m.path]; ok && !cached.Dirty {
			m.entries = cloneDirEntries(cached.Entries)
			m.largeFiles = cloneFileEntries(cached.LargeFiles)
			m.totalSize = cached.TotalSize
=======
			m.currentPath.Store("")
		}

		if cached, ok := m.cache[m.path]; ok && !cached.Dirty {
			m.entries = slices.Clone(cached.Entries)
			m.largeFiles = slices.Clone(cached.LargeFiles)
			m.totalSize = cached.TotalSize
			m.totalFiles = cached.TotalFiles
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			m.selected = cached.Selected
			m.offset = cached.EntryOffset
			m.largeSelected = cached.LargeSelected
			m.largeOffset = cached.LargeOffset
			m.clampEntrySelection()
			m.clampLargeSelection()
			m.status = fmt.Sprintf("Cached view for %s", displayPath(m.path))
			m.scanning = false
			return m, nil
		}
<<<<<<< HEAD
		return m, tea.Batch(m.scanCmd(m.path), tickCmd())
	}
	m.status = fmt.Sprintf("File: %s (%s)", selected.Name, humanizeBytes(selected.Size))
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	fmt.Fprintln(&b)

	if m.inOverviewMode() {
		fmt.Fprintf(&b, "%sAnalyze Disk%s\n", colorPurple, colorReset)
		if m.overviewScanning {
			// Check if we're in initial scan (all entries are pending)
			allPending := true
			for _, entry := range m.entries {
				if entry.Size >= 0 {
					allPending = false
					break
				}
			}

			if allPending {
				// Show prominent loading screen for initial scan
				fmt.Fprintf(&b, "%s%s%s%s Analyzing disk usage, please wait...%s\n",
					colorCyan, colorBold,
					spinnerFrames[m.spinner],
					colorReset, colorReset)
				return b.String()
			} else {
				// Progressive scanning - show subtle indicator
				fmt.Fprintf(&b, "%sSelect a location to explore:%s  ", colorGray, colorReset)
				fmt.Fprintf(&b, "%s%s%s%s Scanning...\n\n", colorCyan, colorBold, spinnerFrames[m.spinner], colorReset)
			}
		} else {
			// Check if there are still pending items
			hasPending := false
			for _, entry := range m.entries {
				if entry.Size < 0 {
					hasPending = true
					break
				}
			}
			if hasPending {
				fmt.Fprintf(&b, "%sSelect a location to explore:%s  ", colorGray, colorReset)
				fmt.Fprintf(&b, "%s%s%s%s Scanning...\n\n", colorCyan, colorBold, spinnerFrames[m.spinner], colorReset)
			} else {
				fmt.Fprintf(&b, "%sSelect a location to explore:%s\n\n", colorGray, colorReset)
			}
		}
	} else {
		fmt.Fprintf(&b, "%sAnalyze Disk%s  %s%s%s", colorPurple, colorReset, colorGray, displayPath(m.path), colorReset)
		if !m.scanning {
			fmt.Fprintf(&b, "  |  Total: %s", humanizeBytes(m.totalSize))
		}
		fmt.Fprintf(&b, "\n\n")
	}

	if m.deleting {
		// Show delete progress
		count := int64(0)
		if m.deleteCount != nil {
			count = atomic.LoadInt64(m.deleteCount)
		}

		fmt.Fprintf(&b, "%s%s%s%s Deleting: %s%s items%s removed, please wait...\n",
			colorCyan, colorBold,
			spinnerFrames[m.spinner],
			colorReset,
			colorYellow, formatNumber(count), colorReset)

		return b.String()
	}

	if m.scanning {
		filesScanned, dirsScanned, bytesScanned := m.getScanProgress()

		fmt.Fprintf(&b, "%s%s%s%s Scanning: %s%s files%s, %s%s dirs%s, %s%s%s\n",
			colorCyan, colorBold,
			spinnerFrames[m.spinner],
			colorReset,
			colorYellow, formatNumber(filesScanned), colorReset,
			colorYellow, formatNumber(dirsScanned), colorReset,
			colorGreen, humanizeBytes(bytesScanned), colorReset)

		if m.currentPath != nil {
			currentPath := *m.currentPath
			if currentPath != "" {
				shortPath := displayPath(currentPath)
				shortPath = truncateMiddle(shortPath, 50)
				fmt.Fprintf(&b, "%s%s%s\n", colorGray, shortPath, colorReset)
			}
		}

		return b.String()
	}

	if m.showLargeFiles {
		if len(m.largeFiles) == 0 {
			fmt.Fprintln(&b, "  No large files found (>=100MB)")
		} else {
			viewport := calculateViewport(m.height, true)
			start := m.largeOffset
			if start < 0 {
				start = 0
			}
			end := start + viewport
			if end > len(m.largeFiles) {
				end = len(m.largeFiles)
			}
			maxLargeSize := int64(1)
			for _, file := range m.largeFiles {
				if file.Size > maxLargeSize {
					maxLargeSize = file.Size
				}
			}
			for idx := start; idx < end; idx++ {
				file := m.largeFiles[idx]
				shortPath := displayPath(file.Path)
				shortPath = truncateMiddle(shortPath, 35)
				paddedPath := padName(shortPath, 35)
				entryPrefix := "   "
				nameColor := ""
				sizeColor := colorGray
				numColor := ""
				if idx == m.largeSelected {
					entryPrefix = fmt.Sprintf(" %s%s▶%s ", colorCyan, colorBold, colorReset)
					nameColor = colorCyan
					sizeColor = colorCyan
					numColor = colorCyan
				}
				size := humanizeBytes(file.Size)
				bar := coloredProgressBar(file.Size, maxLargeSize, 0)
				fmt.Fprintf(&b, "%s%s%2d.%s %s  |  📄 %s%s%s  %s%10s%s\n",
					entryPrefix, numColor, idx+1, colorReset, bar, nameColor, paddedPath, colorReset, sizeColor, size, colorReset)
			}
		}
	} else {
		if len(m.entries) == 0 {
			fmt.Fprintln(&b, "  Empty directory")
		} else {
			if m.inOverviewMode() {
				maxSize := int64(1)
				for _, entry := range m.entries {
					if entry.Size > maxSize {
						maxSize = entry.Size
					}
				}
				totalSize := m.totalSize
				for idx, entry := range m.entries {
					icon := "📁"
					sizeVal := entry.Size
					barValue := sizeVal
					if barValue < 0 {
						barValue = 0
					}
					var percent float64
					if totalSize > 0 && sizeVal >= 0 {
						percent = float64(sizeVal) / float64(totalSize) * 100
					} else {
						percent = 0
					}
					percentStr := fmt.Sprintf("%5.1f%%", percent)
					if totalSize == 0 || sizeVal < 0 {
						percentStr = "  --  "
					}
					bar := coloredProgressBar(barValue, maxSize, percent)
					sizeText := "pending.."
					if sizeVal >= 0 {
						sizeText = humanizeBytes(sizeVal)
					}
					sizeColor := colorGray
					if sizeVal >= 0 && totalSize > 0 {
						switch {
						case percent >= 50:
							sizeColor = colorRed
						case percent >= 20:
							sizeColor = colorYellow
						case percent >= 5:
							sizeColor = colorCyan
						default:
							sizeColor = colorGray
						}
					}
					entryPrefix := "   "
					name := trimName(entry.Name)
					paddedName := padName(name, 28)
					nameSegment := fmt.Sprintf("%s %s", icon, paddedName)
					numColor := ""
					percentColor := ""
					if idx == m.selected {
						entryPrefix = fmt.Sprintf(" %s%s▶%s ", colorCyan, colorBold, colorReset)
						nameSegment = fmt.Sprintf("%s%s %s%s", colorCyan, icon, paddedName, colorReset)
						numColor = colorCyan
						percentColor = colorCyan
						sizeColor = colorCyan
					}
					displayIndex := idx + 1

					// Priority: cleanable > unused time
					var hintLabel string
					if entry.IsDir && isCleanableDir(entry.Path) {
						hintLabel = fmt.Sprintf("%s🧹%s", colorYellow, colorReset)
					} else {
						// For overview mode, get access time on-demand if not set
						lastAccess := entry.LastAccess
						if lastAccess.IsZero() && entry.Path != "" {
							lastAccess = getLastAccessTime(entry.Path)
						}
						if unusedTime := formatUnusedTime(lastAccess); unusedTime != "" {
							hintLabel = fmt.Sprintf("%s%s%s", colorGray, unusedTime, colorReset)
						}
					}

					if hintLabel == "" {
						fmt.Fprintf(&b, "%s%s%2d.%s %s %s%s%s  |  %s %s%10s%s\n",
							entryPrefix, numColor, displayIndex, colorReset, bar, percentColor, percentStr, colorReset,
							nameSegment, sizeColor, sizeText, colorReset)
					} else {
						fmt.Fprintf(&b, "%s%s%2d.%s %s %s%s%s  |  %s %s%10s%s  %s\n",
							entryPrefix, numColor, displayIndex, colorReset, bar, percentColor, percentStr, colorReset,
							nameSegment, sizeColor, sizeText, colorReset, hintLabel)
					}
				}
			} else {
				// Normal mode with sizes and progress bars
				maxSize := int64(1)
				for _, entry := range m.entries {
					if entry.Size > maxSize {
						maxSize = entry.Size
					}
				}

				viewport := calculateViewport(m.height, false)
				start := m.offset
				if start < 0 {
					start = 0
				}
				end := start + viewport
				if end > len(m.entries) {
					end = len(m.entries)
				}

				for idx := start; idx < end; idx++ {
					entry := m.entries[idx]
					icon := "📄"
					if entry.IsDir {
						icon = "📁"
					}
					size := humanizeBytes(entry.Size)
					name := trimName(entry.Name)
					paddedName := padName(name, 28)

					// Calculate percentage
					percent := float64(entry.Size) / float64(m.totalSize) * 100
					percentStr := fmt.Sprintf("%5.1f%%", percent)

					// Get colored progress bar
					bar := coloredProgressBar(entry.Size, maxSize, percent)

					// Color the size based on magnitude
					var sizeColor string
					if percent >= 50 {
						sizeColor = colorRed
					} else if percent >= 20 {
						sizeColor = colorYellow
					} else if percent >= 5 {
						sizeColor = colorCyan
					} else {
						sizeColor = colorGray
					}

					// Keep chart columns aligned even when arrow is shown
					entryPrefix := "   "
					nameSegment := fmt.Sprintf("%s %s", icon, paddedName)
					numColor := ""
					percentColor := ""
					if idx == m.selected {
						entryPrefix = fmt.Sprintf(" %s%s▶%s ", colorCyan, colorBold, colorReset)
						nameSegment = fmt.Sprintf("%s%s %s%s", colorCyan, icon, paddedName, colorReset)
						numColor = colorCyan
						percentColor = colorCyan
						sizeColor = colorCyan
					}

					displayIndex := idx + 1

					// Priority: cleanable > unused time
					var hintLabel string
					if entry.IsDir && isCleanableDir(entry.Path) {
						hintLabel = fmt.Sprintf("%s🧹%s", colorYellow, colorReset)
					} else {
						// Get access time on-demand if not set
						lastAccess := entry.LastAccess
						if lastAccess.IsZero() && entry.Path != "" {
							lastAccess = getLastAccessTime(entry.Path)
						}
						if unusedTime := formatUnusedTime(lastAccess); unusedTime != "" {
							hintLabel = fmt.Sprintf("%s%s%s", colorGray, unusedTime, colorReset)
						}
					}

					if hintLabel == "" {
						fmt.Fprintf(&b, "%s%s%2d.%s %s %s%s%s  |  %s %s%10s%s\n",
							entryPrefix, numColor, displayIndex, colorReset, bar, percentColor, percentStr, colorReset,
							nameSegment, sizeColor, size, colorReset)
					} else {
						fmt.Fprintf(&b, "%s%s%2d.%s %s %s%s%s  |  %s %s%10s%s  %s\n",
							entryPrefix, numColor, displayIndex, colorReset, bar, percentColor, percentStr, colorReset,
							nameSegment, sizeColor, size, colorReset, hintLabel)
					}
				}
			}
		}
	}

	fmt.Fprintln(&b)
	if m.inOverviewMode() {
		fmt.Fprintf(&b, "%s↑↓→  |  Enter  |  R Refresh  |  O Open  |  F Show  |  Q Quit%s\n", colorGray, colorReset)
	} else if m.showLargeFiles {
		fmt.Fprintf(&b, "%s↑↓  |  R Refresh  |  O Open  |  F Show  |  ⌫ Delete  |  L Back  |  Q Quit%s\n", colorGray, colorReset)
	} else {
		largeFileCount := len(m.largeFiles)
		if largeFileCount > 0 {
			fmt.Fprintf(&b, "%s↑↓←→  |  Enter  |  R Refresh  |  O Open  |  F Show  |  ⌫ Delete  |  L Large(%d)  |  Q Quit%s\n", colorGray, largeFileCount, colorReset)
		} else {
			fmt.Fprintf(&b, "%s↑↓←→  |  Enter  |  R Refresh  |  O Open  |  F Show  |  ⌫ Delete  |  Q Quit%s\n", colorGray, colorReset)
		}
	}
	if m.deleteConfirm && m.deleteTarget != nil {
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "%sDelete:%s %s (%s)  %sPress ⌫ again  |  ESC cancel%s\n",
			colorRed, colorReset,
			m.deleteTarget.Name, humanizeBytes(m.deleteTarget.Size),
			colorGray, colorReset)
	}
	return b.String()
}

=======
		m.lastTotalFiles = 0
		if total, err := peekCacheTotalFiles(m.path); err == nil && total > 0 {
			m.lastTotalFiles = total
		}
		return m, tea.Batch(m.scanCmd(m.path), tickCmd())
	}
	m.status = fmt.Sprintf("File: %s, %s", selected.Name, humanizeBytes(selected.Size))
	return m, nil
}

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
func (m *model) clampEntrySelection() {
	if len(m.entries) == 0 {
		m.selected = 0
		m.offset = 0
		return
	}
	if m.selected >= len(m.entries) {
		m.selected = len(m.entries) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
	viewport := calculateViewport(m.height, false)
<<<<<<< HEAD
	maxOffset := len(m.entries) - viewport
	if maxOffset < 0 {
		maxOffset = 0
	}
=======
	maxOffset := max(len(m.entries)-viewport, 0)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.selected < m.offset {
		m.offset = m.selected
	}
	if m.selected >= m.offset+viewport {
		m.offset = m.selected - viewport + 1
	}
}

func (m *model) clampLargeSelection() {
	if len(m.largeFiles) == 0 {
		m.largeSelected = 0
		m.largeOffset = 0
		return
	}
	if m.largeSelected >= len(m.largeFiles) {
		m.largeSelected = len(m.largeFiles) - 1
	}
	if m.largeSelected < 0 {
		m.largeSelected = 0
	}
	viewport := calculateViewport(m.height, true)
<<<<<<< HEAD
	maxOffset := len(m.largeFiles) - viewport
	if maxOffset < 0 {
		maxOffset = 0
	}
=======
	maxOffset := max(len(m.largeFiles)-viewport, 0)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if m.largeOffset > maxOffset {
		m.largeOffset = maxOffset
	}
	if m.largeSelected < m.largeOffset {
		m.largeOffset = m.largeSelected
	}
	if m.largeSelected >= m.largeOffset+viewport {
		m.largeOffset = m.largeSelected - viewport + 1
	}
}

func sumKnownEntrySizes(entries []dirEntry) int64 {
	var total int64
	for _, entry := range entries {
		if entry.Size > 0 {
			total += entry.Size
		}
	}
	return total
}

func nextPendingOverviewIndex(entries []dirEntry) int {
	for i, entry := range entries {
		if entry.Size < 0 {
			return i
		}
	}
	return -1
}

func hasPendingOverviewEntries(entries []dirEntry) bool {
	for _, entry := range entries {
		if entry.Size < 0 {
			return true
		}
	}
	return false
}

func (m *model) removePathFromView(path string) {
	if path == "" {
		return
	}

	var removedSize int64
	for i, entry := range m.entries {
		if entry.Path == path {
			if entry.Size > 0 {
				removedSize = entry.Size
			}
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			break
		}
	}

	for i := 0; i < len(m.largeFiles); i++ {
		if m.largeFiles[i].Path == path {
			m.largeFiles = append(m.largeFiles[:i], m.largeFiles[i+1:]...)
			break
		}
	}

	if removedSize > 0 {
		if removedSize > m.totalSize {
			m.totalSize = 0
		} else {
			m.totalSize -= removedSize
		}
		m.clampEntrySelection()
	}
	m.clampLargeSelection()
}

func scanOverviewPathCmd(path string, index int) tea.Cmd {
	return func() tea.Msg {
		size, err := measureOverviewSize(path)
		return overviewSizeMsg{
			Path:  path,
			Index: index,
			Size:  size,
			Err:   err,
		}
	}
}
<<<<<<< HEAD

// calculateViewport dynamically calculates the viewport size based on terminal height
func calculateViewport(termHeight int, isLargeFiles bool) int {
	if termHeight <= 0 {
		// Terminal height unknown, use default
		return defaultViewport
	}

	// Calculate reserved space for UI elements
	reserved := 6 // header (3-4 lines) + footer (2 lines)
	if isLargeFiles {
		reserved = 5 // Large files view has less overhead
	}

	available := termHeight - reserved

	// Ensure minimum and maximum bounds
	if available < 1 {
		return 1 // Minimum 1 line for very short terminals
	}
	if available > 30 {
		return 30 // Maximum 30 lines to avoid information overload
	}

	return available
}
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
