package main

import (
	"bytes"
<<<<<<< HEAD
=======
	"container/heap"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sync/singleflight"
)

var scanGroup singleflight.Group

<<<<<<< HEAD
func scanPathConcurrent(root string, filesScanned, dirsScanned, bytesScanned *int64, currentPath *string) (scanResult, error) {
=======
// trySend attempts to send an item to a channel with a timeout.
// Returns true if the item was sent, false if the timeout was reached.
func trySend[T any](ch chan<- T, item T, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case ch <- item:
		if !timer.Stop() {
			<-timer.C
		}
		return true
	case <-timer.C:
		return false
	}
}

func scanPathConcurrent(root string, filesScanned, dirsScanned, bytesScanned *int64, currentPath *atomic.Value) (scanResult, error) {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	children, err := os.ReadDir(root)
	if err != nil {
		return scanResult{}, err
	}

	var total int64
<<<<<<< HEAD
	entries := make([]dirEntry, 0, len(children))
	largeFiles := make([]fileEntry, 0, maxLargeFiles*2)

	// Use worker pool for concurrent directory scanning
	// For I/O-bound operations, use more workers than CPU count
	numWorkers := runtime.NumCPU() * cpuMultiplier
	if numWorkers < minWorkers {
		numWorkers = minWorkers
	}
=======

	// Keep Top N heaps.
	entriesHeap := &entryHeap{}
	heap.Init(entriesHeap)

	largeFilesHeap := &largeFileHeap{}
	heap.Init(largeFilesHeap)
	largeFileMinSize := int64(largeFileWarmupMinSize)

	// Worker pool sized for I/O-bound scanning.
	numWorkers := max(runtime.NumCPU()*cpuMultiplier, minWorkers)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if numWorkers > maxWorkers {
		numWorkers = maxWorkers
	}
	if numWorkers > len(children) {
		numWorkers = len(children)
	}
	if numWorkers < 1 {
		numWorkers = 1
	}
	sem := make(chan struct{}, numWorkers)
<<<<<<< HEAD
	var wg sync.WaitGroup

	// Use channels to collect results without lock contention
	entryChan := make(chan dirEntry, len(children))
	largeFileChan := make(chan fileEntry, maxLargeFiles*2)

	// Start goroutines to collect from channels
=======
	duSem := make(chan struct{}, min(4, runtime.NumCPU()))        // limits concurrent du processes
	duQueueSem := make(chan struct{}, min(4, runtime.NumCPU())*2) // limits how many goroutines may be waiting to run du
	var wg sync.WaitGroup

	// Collect results via channels.
	// Cap buffer size to prevent memory spikes with huge directories.
	entryBufSize := len(children)
	if entryBufSize > 4096 {
		entryBufSize = 4096
	}
	if entryBufSize < 1 {
		entryBufSize = 1
	}
	entryChan := make(chan dirEntry, entryBufSize)
	largeFileChan := make(chan fileEntry, maxLargeFiles*2)

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	var collectorWg sync.WaitGroup
	collectorWg.Add(2)
	go func() {
		defer collectorWg.Done()
		for entry := range entryChan {
<<<<<<< HEAD
			entries = append(entries, entry)
=======
			if entriesHeap.Len() < maxEntries {
				heap.Push(entriesHeap, entry)
			} else if entry.Size > (*entriesHeap)[0].Size {
				heap.Pop(entriesHeap)
				heap.Push(entriesHeap, entry)
			}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		}
	}()
	go func() {
		defer collectorWg.Done()
		for file := range largeFileChan {
<<<<<<< HEAD
			largeFiles = append(largeFiles, file)
=======
			if largeFilesHeap.Len() < maxLargeFiles {
				heap.Push(largeFilesHeap, file)
				if largeFilesHeap.Len() == maxLargeFiles {
					atomic.StoreInt64(&largeFileMinSize, (*largeFilesHeap)[0].Size)
				}
			} else if file.Size > (*largeFilesHeap)[0].Size {
				heap.Pop(largeFilesHeap)
				heap.Push(largeFilesHeap, file)
				atomic.StoreInt64(&largeFileMinSize, (*largeFilesHeap)[0].Size)
			}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		}
	}()

	isRootDir := root == "/"
<<<<<<< HEAD
=======
	home := os.Getenv("HOME")
	isHomeDir := home != "" && root == home
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b

	for _, child := range children {
		fullPath := filepath.Join(root, child.Name())

<<<<<<< HEAD
		// Skip symlinks to avoid following them into unexpected locations
		// Use Type() instead of IsDir() to check without following symlinks
		if child.Type()&fs.ModeSymlink != 0 {
			// For symlinks, get their target info but mark them specially
=======
		// Skip symlinks to avoid following unexpected targets.
		if child.Type()&fs.ModeSymlink != 0 {
			targetInfo, err := os.Stat(fullPath)
			isDir := false
			if err == nil && targetInfo.IsDir() {
				isDir = true
			}

			// Count link size only to avoid double-counting targets.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			info, err := child.Info()
			if err != nil {
				continue
			}
			size := getActualFileSize(fullPath, info)
			atomic.AddInt64(&total, size)

<<<<<<< HEAD
			entryChan <- dirEntry{
				Name:       child.Name() + " →",  // Add arrow to indicate symlink
				Path:       fullPath,
				Size:       size,
				IsDir:      false,  // Don't allow navigation into symlinks
				LastAccess: getLastAccessTimeFromInfo(info),
			}
			continue
		}

		if child.IsDir() {
			// In root directory, skip system directories completely
=======
			trySend(entryChan, dirEntry{
				Name:       child.Name() + " →",
				Path:       fullPath,
				Size:       size,
				IsDir:      isDir,
				LastAccess: getLastAccessTimeFromInfo(info),
			}, 100*time.Millisecond)
			continue

		}

		if child.IsDir() {
			if defaultSkipDirs[child.Name()] {
				continue
			}

			// Skip system dirs at root.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			if isRootDir && skipSystemDirs[child.Name()] {
				continue
			}

<<<<<<< HEAD
			// For folded directories, calculate size quickly without expanding
			if shouldFoldDirWithPath(child.Name(), fullPath) {
				wg.Add(1)
				go func(name, path string) {
					defer wg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()

					// Try du command first for folded dirs (much faster)
					size, err := getDirectorySizeFromDu(path)
					if err != nil || size <= 0 {
						// Fallback to walk if du fails
=======
			// ~/Library is scanned separately; reuse cache when possible.
			if isHomeDir && child.Name() == "Library" {
				sem <- struct{}{}
				wg.Add(1)
				go func(name, path string) {
					defer wg.Done()
					defer func() { <-sem }()

					var size int64
					if cached, err := loadStoredOverviewSize(path); err == nil && cached > 0 {
						size = cached
					} else if cached, err := loadCacheFromDisk(path); err == nil {
						size = cached.TotalSize
					} else {
						size = calculateDirSizeConcurrent(path, largeFileChan, &largeFileMinSize, duSem, duQueueSem, filesScanned, dirsScanned, bytesScanned, currentPath)
					}
					atomic.AddInt64(&total, size)
					atomic.AddInt64(dirsScanned, 1)

					trySend(entryChan, dirEntry{
						Name:       name,
						Path:       path,
						Size:       size,
						IsDir:      true,
						LastAccess: time.Time{},
					}, 100*time.Millisecond)
				}(child.Name(), fullPath)
				continue
			}

			// Folded dirs: fast size without expanding.
			if shouldFoldDirWithPath(child.Name(), fullPath) {
				duQueueSem <- struct{}{}
				wg.Add(1)
				go func(name, path string) {
					defer wg.Done()
					defer func() { <-duQueueSem }()

					size, err := func() (int64, error) {
						duSem <- struct{}{}
						defer func() { <-duSem }()
						return getDirectorySizeFromDu(path)
					}()
					if err != nil || size <= 0 {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
						size = calculateDirSizeFast(path, filesScanned, dirsScanned, bytesScanned, currentPath)
					}
					atomic.AddInt64(&total, size)
					atomic.AddInt64(dirsScanned, 1)

<<<<<<< HEAD
					entryChan <- dirEntry{
=======
					trySend(entryChan, dirEntry{
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
						Name:       name,
						Path:       path,
						Size:       size,
						IsDir:      true,
<<<<<<< HEAD
						LastAccess: time.Time{}, // Lazy load when displayed
					}
=======
						LastAccess: time.Time{},
					}, 100*time.Millisecond)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				}(child.Name(), fullPath)
				continue
			}

<<<<<<< HEAD
			// Normal directory: full scan with detail
			wg.Add(1)
			go func(name, path string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				size := calculateDirSizeConcurrent(path, largeFileChan, filesScanned, dirsScanned, bytesScanned, currentPath)
				atomic.AddInt64(&total, size)
				atomic.AddInt64(dirsScanned, 1)

				entryChan <- dirEntry{
=======
			sem <- struct{}{}
			wg.Add(1)
			go func(name, path string) {
				defer wg.Done()
				defer func() { <-sem }()

				size := calculateDirSizeConcurrent(path, largeFileChan, &largeFileMinSize, duSem, duQueueSem, filesScanned, dirsScanned, bytesScanned, currentPath)
				atomic.AddInt64(&total, size)
				atomic.AddInt64(dirsScanned, 1)

				trySend(entryChan, dirEntry{
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
					Name:       name,
					Path:       path,
					Size:       size,
					IsDir:      true,
<<<<<<< HEAD
					LastAccess: time.Time{}, // Lazy load when displayed
				}
=======
					LastAccess: time.Time{},
				}, 100*time.Millisecond)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			}(child.Name(), fullPath)
			continue
		}

		info, err := child.Info()
		if err != nil {
			continue
		}
<<<<<<< HEAD
		// Get actual disk usage for sparse files and cloud files
=======
		// Actual disk usage for sparse/cloud files.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		size := getActualFileSize(fullPath, info)
		atomic.AddInt64(&total, size)
		atomic.AddInt64(filesScanned, 1)
		atomic.AddInt64(bytesScanned, size)

<<<<<<< HEAD
		entryChan <- dirEntry{
=======
		trySend(entryChan, dirEntry{
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			Name:       child.Name(),
			Path:       fullPath,
			Size:       size,
			IsDir:      false,
			LastAccess: getLastAccessTimeFromInfo(info),
<<<<<<< HEAD
		}
		// Only track large files that are not code/text files
		if !shouldSkipFileForLargeTracking(fullPath) && size >= minLargeFileSize {
			largeFileChan <- fileEntry{Name: child.Name(), Path: fullPath, Size: size}
=======
		}, 100*time.Millisecond)

		// Track large files only.
		if !shouldSkipFileForLargeTracking(fullPath) {
			minSize := atomic.LoadInt64(&largeFileMinSize)
			if size >= minSize {
				trySend(largeFileChan, fileEntry{Name: child.Name(), Path: fullPath, Size: size}, 100*time.Millisecond)
			}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		}
	}

	wg.Wait()

<<<<<<< HEAD
	// Close channels and wait for collectors to finish
=======
	// Close channels and wait for collectors.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	close(entryChan)
	close(largeFileChan)
	collectorWg.Wait()

<<<<<<< HEAD
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})
	if len(entries) > maxEntries {
		entries = entries[:maxEntries]
	}

	// Try to use Spotlight for faster large file discovery
	if spotlightFiles := findLargeFilesWithSpotlight(root, minLargeFileSize); len(spotlightFiles) > 0 {
		largeFiles = spotlightFiles
	} else {
		// Sort and trim large files collected from scanning
		sort.Slice(largeFiles, func(i, j int) bool {
			return largeFiles[i].Size > largeFiles[j].Size
		})
		if len(largeFiles) > maxLargeFiles {
			largeFiles = largeFiles[:maxLargeFiles]
		}
=======
	// Convert heaps to sorted slices (descending).
	entries := make([]dirEntry, entriesHeap.Len())
	for i := len(entries) - 1; i >= 0; i-- {
		entries[i] = heap.Pop(entriesHeap).(dirEntry)
	}

	largeFiles := make([]fileEntry, largeFilesHeap.Len())
	for i := len(largeFiles) - 1; i >= 0; i-- {
		largeFiles[i] = heap.Pop(largeFilesHeap).(fileEntry)
	}

	// Use Spotlight for large files when it expands the list.
	if spotlightFiles := findLargeFilesWithSpotlight(root, spotlightMinFileSize); len(spotlightFiles) > len(largeFiles) {
		largeFiles = spotlightFiles
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	}

	return scanResult{
		Entries:    entries,
		LargeFiles: largeFiles,
		TotalSize:  total,
<<<<<<< HEAD
=======
		TotalFiles: atomic.LoadInt64(filesScanned),
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	}, nil
}

func shouldFoldDirWithPath(name, path string) bool {
<<<<<<< HEAD
	// Check basic fold list first
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if foldDirs[name] {
		return true
	}

<<<<<<< HEAD
	// Special case: npm cache directories - fold all subdirectories
	// This includes: .npm/_quick/*, .npm/_cacache/*, .npm/a-z/*, .tnpm/*
	if strings.Contains(path, "/.npm/") || strings.Contains(path, "/.tnpm/") {
		// Get the parent directory name
		parent := filepath.Base(filepath.Dir(path))
		// If parent is a cache folder (_quick, _cacache, etc) or npm dir itself, fold it
		if parent == ".npm" || parent == ".tnpm" || strings.HasPrefix(parent, "_") {
			return true
		}
		// Also fold single-letter subdirectories (npm cache structure like .npm/a/, .npm/b/)
=======
	// Handle npm cache structure.
	if strings.Contains(path, "/.npm/") || strings.Contains(path, "/.tnpm/") {
		parent := filepath.Base(filepath.Dir(path))
		if parent == ".npm" || parent == ".tnpm" || strings.HasPrefix(parent, "_") {
			return true
		}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if len(name) == 1 {
			return true
		}
	}

	return false
}

func shouldSkipFileForLargeTracking(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return skipExtensions[ext]
}

<<<<<<< HEAD
// calculateDirSizeFast performs fast directory size calculation without detailed tracking or large file detection.
// Updates progress counters in batches to reduce atomic operation overhead.
func calculateDirSizeFast(root string, filesScanned, dirsScanned, bytesScanned *int64, currentPath *string) int64 {
	var total int64
	var localFiles, localDirs int64
	var batchBytes int64

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		// Check for timeout
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err != nil {
			return nil
		}
		if d.IsDir() {
			localDirs++
			// Batch update every N dirs to reduce atomic operations
			if localDirs%batchUpdateSize == 0 {
				atomic.AddInt64(dirsScanned, batchUpdateSize)
				localDirs = 0
			}
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		// Get actual disk usage for sparse files and cloud files
		size := getActualFileSize(path, info)
		total += size
		batchBytes += size
		localFiles++
		if currentPath != nil {
			*currentPath = path
		}
		// Batch update every N files to reduce atomic operations
		if localFiles%batchUpdateSize == 0 {
			atomic.AddInt64(filesScanned, batchUpdateSize)
			atomic.AddInt64(bytesScanned, batchBytes)
			localFiles = 0
			batchBytes = 0
		}
		return nil
	}

	_ = filepath.WalkDir(root, walkFunc)

	// Final update for remaining counts
	if localFiles > 0 {
		atomic.AddInt64(filesScanned, localFiles)
	}
	if localDirs > 0 {
		atomic.AddInt64(dirsScanned, localDirs)
	}
	if batchBytes > 0 {
		atomic.AddInt64(bytesScanned, batchBytes)
	}
=======
// calculateDirSizeFast performs concurrent dir sizing using os.ReadDir.
func calculateDirSizeFast(root string, filesScanned, dirsScanned, bytesScanned *int64, currentPath *atomic.Value) int64 {
	var total int64
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	concurrency := min(runtime.NumCPU()*4, 64)
	sem := make(chan struct{}, concurrency)

	var walk func(string)
	walk = func(dirPath string) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if currentPath != nil && atomic.LoadInt64(filesScanned)%int64(batchUpdateSize) == 0 {
			currentPath.Store(dirPath)
		}

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return
		}

		var localBytes, localFiles int64

		for _, entry := range entries {
			if entry.IsDir() {
				subDir := filepath.Join(dirPath, entry.Name())
				sem <- struct{}{}
				wg.Add(1)
				go func(p string) {
					defer wg.Done()
					defer func() { <-sem }()
					walk(p)
				}(subDir)
				atomic.AddInt64(dirsScanned, 1)
			} else {
				info, err := entry.Info()
				if err == nil {
					size := getActualFileSize(filepath.Join(dirPath, entry.Name()), info)
					localBytes += size
					localFiles++
				}
			}
		}

		if localBytes > 0 {
			atomic.AddInt64(&total, localBytes)
			atomic.AddInt64(bytesScanned, localBytes)
		}
		if localFiles > 0 {
			atomic.AddInt64(filesScanned, localFiles)
		}
	}

	walk(root)
	wg.Wait()
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b

	return total
}

<<<<<<< HEAD
// Use Spotlight (mdfind) to quickly find large files in a directory
func findLargeFilesWithSpotlight(root string, minSize int64) []fileEntry {
	// mdfind query: files >= minSize in the specified directory
=======
// Use Spotlight (mdfind) to quickly find large files.
func findLargeFilesWithSpotlight(root string, minSize int64) []fileEntry {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	query := fmt.Sprintf("kMDItemFSSize >= %d", minSize)

	ctx, cancel := context.WithTimeout(context.Background(), mdlsTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "mdfind", "-onlyin", root, query)
	output, err := cmd.Output()
	if err != nil {
<<<<<<< HEAD
		// Fallback: mdfind not available or failed
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []fileEntry

	for _, line := range lines {
=======
		return nil
	}

	var files []fileEntry

	for line := range strings.Lines(strings.TrimSpace(string(output))) {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if line == "" {
			continue
		}

<<<<<<< HEAD
		// Filter out code files first (cheapest check, no I/O)
=======
		// Filter code files first (cheap).
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if shouldSkipFileForLargeTracking(line) {
			continue
		}

<<<<<<< HEAD
		// Filter out files in folded directories (cheap string check)
=======
		// Filter folded directories (cheap string check).
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if isInFoldedDir(line) {
			continue
		}

<<<<<<< HEAD
		// Use Lstat instead of Stat (faster, doesn't follow symlinks)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		info, err := os.Lstat(line)
		if err != nil {
			continue
		}

<<<<<<< HEAD
		// Skip if it's a directory or symlink
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			continue
		}

<<<<<<< HEAD
		// Get actual disk usage for sparse files and cloud files
=======
		// Actual disk usage for sparse/cloud files.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		actualSize := getActualFileSize(line, info)
		files = append(files, fileEntry{
			Name: filepath.Base(line),
			Path: line,
			Size: actualSize,
		})
	}

<<<<<<< HEAD
	// Sort by size (descending)
=======
	// Sort by size (descending).
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

<<<<<<< HEAD
	// Return top N
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if len(files) > maxLargeFiles {
		files = files[:maxLargeFiles]
	}

	return files
}

<<<<<<< HEAD
// isInFoldedDir checks if a path is inside a folded directory (optimized)
func isInFoldedDir(path string) bool {
	// Split path into components for faster checking
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
=======
// isInFoldedDir checks if a path is inside a folded directory.
func isInFoldedDir(path string) bool {
	parts := strings.SplitSeq(path, string(os.PathSeparator))
	for part := range parts {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if foldDirs[part] {
			return true
		}
	}
	return false
}

<<<<<<< HEAD
func calculateDirSizeConcurrent(root string, largeFileChan chan<- fileEntry, filesScanned, dirsScanned, bytesScanned *int64, currentPath *string) int64 {
	// Read immediate children
=======
func calculateDirSizeConcurrent(root string, largeFileChan chan<- fileEntry, largeFileMinSize *int64, duSem, duQueueSem chan struct{}, filesScanned, dirsScanned, bytesScanned *int64, currentPath *atomic.Value) int64 {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	children, err := os.ReadDir(root)
	if err != nil {
		return 0
	}

	var total int64
	var wg sync.WaitGroup

<<<<<<< HEAD
	// Limit concurrent subdirectory scans to avoid too many goroutines
	maxConcurrent := runtime.NumCPU() * 2
	if maxConcurrent > maxDirWorkers {
		maxConcurrent = maxDirWorkers
	}
=======
	// Limit concurrent subdirectory scans.
	maxConcurrent := min(runtime.NumCPU()*2, maxDirWorkers)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	sem := make(chan struct{}, maxConcurrent)

	for _, child := range children {
		fullPath := filepath.Join(root, child.Name())

<<<<<<< HEAD
		// Skip symlinks to avoid following them into unexpected locations
		if child.Type()&fs.ModeSymlink != 0 {
			// For symlinks, just count their size without following
=======
		if child.Type()&fs.ModeSymlink != 0 {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			info, err := child.Info()
			if err != nil {
				continue
			}
			size := getActualFileSize(fullPath, info)
			total += size
			atomic.AddInt64(filesScanned, 1)
			atomic.AddInt64(bytesScanned, size)
			continue
		}

		if child.IsDir() {
<<<<<<< HEAD
			// Check if this is a folded directory
			if shouldFoldDirWithPath(child.Name(), fullPath) {
				// Use du for folded directories (much faster)
				wg.Add(1)
				go func(path string) {
					defer wg.Done()
					size, err := getDirectorySizeFromDu(path)
					if err == nil && size > 0 {
						atomic.AddInt64(&total, size)
						atomic.AddInt64(bytesScanned, size)
						atomic.AddInt64(dirsScanned, 1)
					}
=======
			if shouldFoldDirWithPath(child.Name(), fullPath) {
				duQueueSem <- struct{}{}
				wg.Add(1)
				go func(path string) {
					defer wg.Done()
					defer func() { <-duQueueSem }()

					size, err := func() (int64, error) {
						duSem <- struct{}{}
						defer func() { <-duSem }()
						return getDirectorySizeFromDu(path)
					}()
					if err != nil || size <= 0 {
						size = calculateDirSizeFast(path, filesScanned, dirsScanned, bytesScanned, currentPath)
					} else {
						atomic.AddInt64(bytesScanned, size)
					}
					atomic.AddInt64(&total, size)
					atomic.AddInt64(dirsScanned, 1)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				}(fullPath)
				continue
			}

<<<<<<< HEAD
			// Recursively scan subdirectory in parallel
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				size := calculateDirSizeConcurrent(path, largeFileChan, filesScanned, dirsScanned, bytesScanned, currentPath)
=======
			sem <- struct{}{}
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				defer func() { <-sem }()

				size := calculateDirSizeConcurrent(path, largeFileChan, largeFileMinSize, duSem, duQueueSem, filesScanned, dirsScanned, bytesScanned, currentPath)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
				atomic.AddInt64(&total, size)
				atomic.AddInt64(dirsScanned, 1)
			}(fullPath)
			continue
		}

<<<<<<< HEAD
		// Handle files
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		info, err := child.Info()
		if err != nil {
			continue
		}

		size := getActualFileSize(fullPath, info)
		total += size
		atomic.AddInt64(filesScanned, 1)
		atomic.AddInt64(bytesScanned, size)

<<<<<<< HEAD
		// Track large files
		if !shouldSkipFileForLargeTracking(fullPath) && size >= minLargeFileSize {
			largeFileChan <- fileEntry{Name: child.Name(), Path: fullPath, Size: size}
		}

		// Update current path
		if currentPath != nil {
			*currentPath = fullPath
=======
		if !shouldSkipFileForLargeTracking(fullPath) && largeFileMinSize != nil {
			minSize := atomic.LoadInt64(largeFileMinSize)
			if size >= minSize {
				trySend(largeFileChan, fileEntry{Name: child.Name(), Path: fullPath, Size: size}, 100*time.Millisecond)
			}
		}

		// Update current path occasionally to prevent UI jitter.
		if currentPath != nil && atomic.LoadInt64(filesScanned)%int64(batchUpdateSize) == 0 {
			currentPath.Store(fullPath)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		}
	}

	wg.Wait()
	return total
}

// measureOverviewSize calculates the size of a directory using multiple strategies.
<<<<<<< HEAD
=======
// When scanning Home, it excludes ~/Library to avoid duplicate counting.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
func measureOverviewSize(path string) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("empty path")
	}

	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		return 0, fmt.Errorf("path must be absolute: %s", path)
	}

	if _, err := os.Stat(path); err != nil {
		return 0, fmt.Errorf("cannot access path: %v", err)
	}

<<<<<<< HEAD
	if cached, err := loadStoredOverviewSize(path); err == nil && cached > 0 {
		return cached, nil
	}

	if duSize, err := getDirectorySizeFromDu(path); err == nil && duSize > 0 {
=======
	// Determine if we should exclude ~/Library (when scanning Home)
	home := os.Getenv("HOME")
	excludePath := ""
	if home != "" && path == home {
		excludePath = filepath.Join(home, "Library")
	}

	if duSize, err := getDirectorySizeFromDuWithExclude(path, excludePath); err == nil && duSize > 0 {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		_ = storeOverviewSize(path, duSize)
		return duSize, nil
	}

<<<<<<< HEAD
	if logicalSize, err := getDirectoryLogicalSize(path); err == nil && logicalSize > 0 {
=======
	if logicalSize, err := getDirectoryLogicalSizeWithExclude(path, excludePath); err == nil && logicalSize > 0 {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		_ = storeOverviewSize(path, logicalSize)
		return logicalSize, nil
	}

	if cached, err := loadCacheFromDisk(path); err == nil {
		_ = storeOverviewSize(path, cached.TotalSize)
		return cached.TotalSize, nil
	}

	return 0, fmt.Errorf("unable to measure directory size with fast methods")
}

func getDirectorySizeFromDu(path string) (int64, error) {
<<<<<<< HEAD
	ctx, cancel := context.WithTimeout(context.Background(), duTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "du", "-sk", path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return 0, fmt.Errorf("du timeout after %v", duTimeout)
		}
		if stderr.Len() > 0 {
			return 0, fmt.Errorf("du failed: %v (%s)", err, stderr.String())
		}
		return 0, fmt.Errorf("du failed: %v", err)
	}
	fields := strings.Fields(stdout.String())
	if len(fields) == 0 {
		return 0, fmt.Errorf("du output empty")
	}
	kb, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse du output: %v", err)
	}
	if kb <= 0 {
		return 0, fmt.Errorf("du size invalid: %d", kb)
	}
	return kb * 1024, nil
}

func getDirectoryLogicalSize(path string) (int64, error) {
=======
	return getDirectorySizeFromDuWithExclude(path, "")
}

func getDirectorySizeFromDuWithExclude(path string, excludePath string) (int64, error) {
	runDuSize := func(target string) (int64, error) {
		if _, err := os.Stat(target); err != nil {
			return 0, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), duTimeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, "du", "-skP", target)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return 0, fmt.Errorf("du timeout after %v", duTimeout)
			}
			if stderr.Len() > 0 {
				return 0, fmt.Errorf("du failed: %v, %s", err, stderr.String())
			}
			return 0, fmt.Errorf("du failed: %v", err)
		}
		fields := strings.Fields(stdout.String())
		if len(fields) == 0 {
			return 0, fmt.Errorf("du output empty")
		}
		kb, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse du output: %v", err)
		}
		if kb <= 0 {
			return 0, fmt.Errorf("du size invalid: %d", kb)
		}
		return kb * 1024, nil
	}

	// When excluding a path (e.g., ~/Library), subtract only that exact directory instead of ignoring every "Library"
	if excludePath != "" {
		totalSize, err := runDuSize(path)
		if err != nil {
			return 0, err
		}
		excludeSize, err := runDuSize(excludePath)
		if err != nil {
			if !os.IsNotExist(err) {
				return 0, err
			}
			excludeSize = 0
		}
		if excludeSize > totalSize {
			excludeSize = 0
		}
		return totalSize - excludeSize, nil
	}

	return runDuSize(path)
}

func getDirectoryLogicalSizeWithExclude(path string, excludePath string) (int64, error) {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	var total int64
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return filepath.SkipDir
			}
			return nil
		}
<<<<<<< HEAD
=======
		// Skip excluded path
		if excludePath != "" && p == excludePath {
			return filepath.SkipDir
		}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		total += getActualFileSize(p, info)
		return nil
	})
	if err != nil && err != filepath.SkipDir {
		return 0, err
	}
	return total, nil
}

func getActualFileSize(_ string, info fs.FileInfo) int64 {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return info.Size()
	}

	actualSize := stat.Blocks * 512
	if actualSize < info.Size() {
		return actualSize
	}
	return info.Size()
}

func getLastAccessTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return getLastAccessTimeFromInfo(info)
}

func getLastAccessTimeFromInfo(info fs.FileInfo) time.Time {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}
	}
	return time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
}
