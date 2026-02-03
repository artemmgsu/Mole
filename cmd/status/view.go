package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
<<<<<<< HEAD
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#5FD7FF")).Bold(true)
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C6C6C"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD75F"))
	dangerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F")).Bold(true)
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#87D787"))
	lineStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4A4A4A"))
=======
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#C79FD7")).Bold(true)
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#737373"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD75F"))
	dangerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F")).Bold(true)
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#A5D6A7"))
	lineStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#404040"))

	primaryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9"))
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
)

const (
	colWidth    = 38
<<<<<<< HEAD
	iconCPU     = "⚙"
	iconMemory  = "▦"
	iconGPU     = "▣"
	iconDisk    = "▤"
	iconNetwork = "⇅"
	iconBattery = "▮"
	iconSensors = "♨"
	iconProcs   = "▶"
)

// Mole body frames (legs animate)
=======
	iconCPU     = "◉"
	iconMemory  = "◫"
	iconGPU     = "◧"
	iconDisk    = "▥"
	iconNetwork = "⇅"
	iconBattery = "◪"
	iconSensors = "◈"
	iconProcs   = "❊"
)

// Mole body frames (facing right).
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
var moleBody = [][]string{
	{
		`     /\_/\`,
		` ___/ o o \`,
		`/___   =-= /`,
		`\____)-m-m)`,
	},
	{
		`     /\_/\`,
		` ___/ o o \`,
		`/___   =-= /`,
		`\____)mm__)`,
	},
	{
		`     /\_/\`,
		` ___/ · · \`,
		`/___   =-= /`,
		`\___)-m__m)`,
	},
	{
		`     /\_/\`,
		` ___/ o o \`,
		`/___   =-= /`,
		`\____)-mm-)`,
	},
}

<<<<<<< HEAD
// Generate frames with horizontal movement
func getMoleFrame(animFrame int, termWidth int) string {
	bodyIdx := animFrame % len(moleBody)
	body := moleBody[bodyIdx]

	// Calculate mole width (approximate)
	moleWidth := 15
	// Move across terminal width
	maxPos := termWidth - moleWidth
	if maxPos < 0 {
		maxPos = 0
	}

	// Move position: 0 -> maxPos -> 0
=======
// Mirror mole body frames (facing left).
var moleBodyMirror = [][]string{
	{
		`    /\_/\`,
		`   / o o \___`,
		`  \ =-=   ___\`,
		`  (m-m-(____/`,
	},
	{
		`    /\_/\`,
		`   / o o \___`,
		`  \ =-=   ___\`,
		`  (__mm(____/`,
	},
	{
		`    /\_/\`,
		`   / · · \___`,
		`  \ =-=   ___\`,
		`  (m__m-(___/`,
	},
	{
		`    /\_/\`,
		`   / o o \___`,
		`  \ =-=   ___\`,
		`  (-mm-(____/`,
	},
}

// getMoleFrame renders the animated mole.
func getMoleFrame(animFrame int, termWidth int) string {
	moleWidth := 15
	maxPos := max(termWidth-moleWidth, 0)

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	cycleLength := maxPos * 2
	if cycleLength == 0 {
		cycleLength = 1
	}
	pos := animFrame % cycleLength
<<<<<<< HEAD
	if pos > maxPos {
		pos = cycleLength - pos
	}

	padding := strings.Repeat(" ", pos)
	var lines []string
	for _, line := range body {
		lines = append(lines, padding+line)
	}
=======
	movingLeft := pos > maxPos
	if movingLeft {
		pos = cycleLength - pos
	}

	// Use mirror frames when moving left
	var frames [][]string
	if movingLeft {
		frames = moleBodyMirror
	} else {
		frames = moleBody
	}

	bodyIdx := animFrame % len(frames)
	body := frames[bodyIdx]

	padding := strings.Repeat(" ", pos)
	var lines []string

	for _, line := range body {
		lines = append(lines, padding+line)
	}

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	return strings.Join(lines, "\n")
}

type cardData struct {
	icon  string
	title string
	lines []string
}

<<<<<<< HEAD
func renderHeader(m MetricsSnapshot, errMsg string, animFrame int, termWidth int) string {
	// Title
	title := titleStyle.Render("Mole Status")

	// Health Score with color and label
	scoreStyle := getScoreStyle(m.HealthScore)
	scoreText := subtleStyle.Render("Health ") + scoreStyle.Render(fmt.Sprintf("● %d", m.HealthScore))

	// Hardware info
	infoParts := []string{}
	if m.Hardware.Model != "" {
		infoParts = append(infoParts, m.Hardware.Model)
	}
	if m.Hardware.CPUModel != "" {
		infoParts = append(infoParts, m.Hardware.CPUModel)
	}
	if m.Hardware.TotalRAM != "" {
		infoParts = append(infoParts, m.Hardware.TotalRAM)
	}
	if m.Hardware.DiskSize != "" {
		infoParts = append(infoParts, m.Hardware.DiskSize)
=======
func renderHeader(m MetricsSnapshot, errMsg string, animFrame int, termWidth int, catHidden bool) string {
	title := titleStyle.Render("Status")

	scoreStyle := getScoreStyle(m.HealthScore)
	scoreText := subtleStyle.Render("Health ") + scoreStyle.Render(fmt.Sprintf("● %d", m.HealthScore))

	// Hardware info for a single line.
	infoParts := []string{}
	if m.Hardware.Model != "" {
		infoParts = append(infoParts, primaryStyle.Render(m.Hardware.Model))
	}
	if m.Hardware.CPUModel != "" {
		cpuInfo := m.Hardware.CPUModel
		// Append GPU core count when available.
		if len(m.GPU) > 0 && m.GPU[0].CoreCount > 0 {
			cpuInfo += fmt.Sprintf(", %dGPU", m.GPU[0].CoreCount)
		}
		infoParts = append(infoParts, cpuInfo)
	}
	var specs []string
	if m.Hardware.TotalRAM != "" {
		specs = append(specs, m.Hardware.TotalRAM)
	}
	if m.Hardware.DiskSize != "" {
		specs = append(specs, m.Hardware.DiskSize)
	}
	if len(specs) > 0 {
		infoParts = append(infoParts, strings.Join(specs, "/"))
	}
	if m.Hardware.RefreshRate != "" {
		infoParts = append(infoParts, m.Hardware.RefreshRate)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	}
	if m.Hardware.OSVersion != "" {
		infoParts = append(infoParts, m.Hardware.OSVersion)
	}
<<<<<<< HEAD

	headerLine := title + "  " + scoreText + "  " + subtleStyle.Render(strings.Join(infoParts, " · "))

	// Running mole animation
	mole := getMoleFrame(animFrame, termWidth)

	if errMsg != "" {
		return lipgloss.JoinVertical(lipgloss.Left, headerLine, "", mole, dangerStyle.Render(errMsg), "")
	}
	return headerLine + "\n\n" + mole
}

func getScoreStyle(score int) lipgloss.Style {
	if score >= 90 {
		// Excellent - Green
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#87D787")).Bold(true)
	} else if score >= 75 {
		// Good - Light Green
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#AFD787")).Bold(true)
	} else if score >= 60 {
		// Fair - Yellow
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD75F")).Bold(true)
	} else if score >= 40 {
		// Poor - Orange
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAF5F")).Bold(true)
	} else {
		// Critical - Red
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F")).Bold(true)
	}
}

func buildCards(m MetricsSnapshot, _ int) []cardData {
	// Row 1: CPU + Memory
	// Row 2: Disk + Power
	// Row 3: Top Processes + Network
	cards := []cardData{
		renderCPUCard(m.CPU),
		renderMemoryCard(m.Memory),
		renderDiskCard(m.Disks, m.DiskIO),
		renderBatteryCard(m.Batteries, m.Thermal),
		renderProcessCard(m.TopProcesses),
		renderNetworkCard(m.Network, m.Proxy),
	}
	// Only show GPU card if there are GPUs with usage data
	if len(m.GPU) > 0 && m.GPU[0].Usage >= 0 {
		cards = append(cards, renderGPUCard(m.GPU))
	}
	// Only show sensors if we have valid temperature readings
	if hasSensorData(m.Sensors) {
		cards = append(cards, renderSensorsCard(m.Sensors))
	}
	return cards
}

func hasSensorData(sensors []SensorReading) bool {
	for _, s := range sensors {
		if s.Note == "" && s.Value > 0 {
			return true
		}
	}
	return false
}

func renderCPUCard(cpu CPUStatus) cardData {
	var lines []string
	lines = append(lines, fmt.Sprintf("Total  %s  %5.1f%%", progressBar(cpu.Usage), cpu.Usage))
	lines = append(lines, subtleStyle.Render(fmt.Sprintf("%.2f / %.2f / %.2f  (%d cores)", cpu.Load1, cpu.Load5, cpu.Load15, cpu.LogicalCPU)))

	if cpu.PerCoreEstimated {
		lines = append(lines, subtleStyle.Render("Per-core data unavailable (using averaged load)"))
	} else if len(cpu.PerCore) > 0 {
		// Show top 3 busiest cores
=======
	if m.Uptime != "" {
		infoParts = append(infoParts, subtleStyle.Render("up "+m.Uptime))
	}

	headerLine := title + "  " + scoreText + "  " + strings.Join(infoParts, " · ")

	// Show cat unless hidden
	var mole string
	if !catHidden {
		mole = getMoleFrame(animFrame, termWidth)
	}

	if errMsg != "" {
		if mole == "" {
			return lipgloss.JoinVertical(lipgloss.Left, headerLine, "", dangerStyle.Render("ERROR: "+errMsg), "")
		}
		return lipgloss.JoinVertical(lipgloss.Left, headerLine, "", mole, dangerStyle.Render("ERROR: "+errMsg), "")
	}
	if mole == "" {
		return headerLine
	}
	return headerLine + "\n" + mole
}

func getScoreStyle(score int) lipgloss.Style {
	switch {
	case score >= 90:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#87FF87")).Bold(true)
	case score >= 75:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#87D787")).Bold(true)
	case score >= 60:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD75F")).Bold(true)
	case score >= 40:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAF5F")).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	}
}

func renderCPUCard(cpu CPUStatus, thermal ThermalStatus) cardData {
	var lines []string

	// Line 1: Usage + Temp (Format: 15% @ 30.4°C)
	usageBar := progressBar(cpu.Usage)

	headerText := fmt.Sprintf("%5.1f%%", cpu.Usage)
	if thermal.CPUTemp > 0 {
		headerText += fmt.Sprintf(" @ %s°C", colorizeTemp(thermal.CPUTemp))
	}

	lines = append(lines, fmt.Sprintf("Total  %s  %s", usageBar, headerText))

	if cpu.PerCoreEstimated {
		lines = append(lines, subtleStyle.Render("Per-core data unavailable, using averaged load"))
	} else if len(cpu.PerCore) > 0 {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		type coreUsage struct {
			idx int
			val float64
		}
		var cores []coreUsage
		for i, v := range cpu.PerCore {
			cores = append(cores, coreUsage{i, v})
		}
		sort.Slice(cores, func(i, j int) bool { return cores[i].val > cores[j].val })

<<<<<<< HEAD
		maxCores := 3
		if len(cores) < maxCores {
			maxCores = len(cores)
		}
=======
		maxCores := min(len(cores), 3)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		for i := 0; i < maxCores; i++ {
			c := cores[i]
			lines = append(lines, fmt.Sprintf("Core%-2d %s  %5.1f%%", c.idx+1, progressBar(c.val), c.val))
		}
	}

<<<<<<< HEAD
	return cardData{icon: iconCPU, title: "CPU", lines: lines}
}

func renderGPUCard(gpus []GPUStatus) cardData {
	var lines []string
	if len(gpus) == 0 {
		lines = append(lines, subtleStyle.Render("No GPU detected"))
	} else {
		for _, g := range gpus {
			name := shorten(g.Name, 12)
			if g.Usage >= 0 {
				lines = append(lines, fmt.Sprintf("%-12s  %s  %5.1f%%", name, progressBar(g.Usage), g.Usage))
			} else {
				lines = append(lines, name)
			}
		}
	}
	return cardData{icon: iconGPU, title: "GPU", lines: lines}
}

func renderMemoryCard(mem MemoryStatus) cardData {
	var lines []string
	lines = append(lines, fmt.Sprintf("Used   %s  %5.1f%%", progressBar(mem.UsedPercent), mem.UsedPercent))
	lines = append(lines, subtleStyle.Render(fmt.Sprintf("%s / %s total", humanBytes(mem.Used), humanBytes(mem.Total))))
	available := mem.Total - mem.Used
	freePercent := 100 - mem.UsedPercent
	lines = append(lines, fmt.Sprintf("Free   %s  %5.1f%%", progressBar(freePercent), freePercent))
	lines = append(lines, subtleStyle.Render(fmt.Sprintf("%s available", humanBytes(available))))
	if mem.SwapTotal > 0 || mem.SwapUsed > 0 {
=======
	// Load line at the end
	if cpu.PCoreCount > 0 && cpu.ECoreCount > 0 {
		lines = append(lines, fmt.Sprintf("Load   %.2f / %.2f / %.2f, %dP+%dE",
			cpu.Load1, cpu.Load5, cpu.Load15, cpu.PCoreCount, cpu.ECoreCount))
	} else {
		lines = append(lines, fmt.Sprintf("Load   %.2f / %.2f / %.2f, %d cores",
			cpu.Load1, cpu.Load5, cpu.Load15, cpu.LogicalCPU))
	}

	return cardData{icon: iconCPU, title: "CPU", lines: lines}
}

func renderMemoryCard(mem MemoryStatus) cardData {
	// Check if swap is being used (or at least allocated).
	hasSwap := mem.SwapTotal > 0 || mem.SwapUsed > 0

	var lines []string
	// Line 1: Used
	lines = append(lines, fmt.Sprintf("Used   %s  %5.1f%%", progressBar(mem.UsedPercent), mem.UsedPercent))

	// Line 2: Free
	freePercent := 100 - mem.UsedPercent
	lines = append(lines, fmt.Sprintf("Free   %s  %5.1f%%", progressBar(freePercent), freePercent))

	if hasSwap {
		// Layout with Swap:
		// 3. Swap (progress bar + text)
		// 4. Total
		// 5. Avail
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		var swapPercent float64
		if mem.SwapTotal > 0 {
			swapPercent = (float64(mem.SwapUsed) / float64(mem.SwapTotal)) * 100.0
		}
<<<<<<< HEAD
		swapText := subtleStyle.Render(fmt.Sprintf("%s / %s swap", humanBytes(mem.SwapUsed), humanBytes(mem.SwapTotal)))
		lines = append(lines, fmt.Sprintf("Swap   %s  %5.1f%%  %s", progressBar(swapPercent), swapPercent, swapText))
	} else {
		lines = append(lines, fmt.Sprintf("Swap   %s", subtleStyle.Render("not in use")))
	}
	// Memory pressure
	if mem.Pressure != "" {
		pressureStyle := okStyle
		pressureText := "Status " + mem.Pressure
		if mem.Pressure == "warn" {
			pressureStyle = warnStyle
		} else if mem.Pressure == "critical" {
=======
		swapText := fmt.Sprintf("%s/%s", humanBytesCompact(mem.SwapUsed), humanBytesCompact(mem.SwapTotal))
		lines = append(lines, fmt.Sprintf("Swap   %s  %5.1f%% %s", progressBar(swapPercent), swapPercent, swapText))

		lines = append(lines, fmt.Sprintf("Total  %s / %s", humanBytes(mem.Used), humanBytes(mem.Total)))
		lines = append(lines, fmt.Sprintf("Avail  %s", humanBytes(mem.Total-mem.Used))) // Simplified avail logic for consistency
	} else {
		// Layout without Swap:
		// 3. Total
		// 4. Cached (if > 0)
		// 5. Avail
		lines = append(lines, fmt.Sprintf("Total  %s / %s", humanBytes(mem.Used), humanBytes(mem.Total)))

		if mem.Cached > 0 {
			lines = append(lines, fmt.Sprintf("Cached %s", humanBytes(mem.Cached)))
		}
		// Calculate available if not provided directly, or use Total-Used as proxy if needed,
		// but typically available is more nuanced. Using what we have.
		// Re-calculating available based on logic if needed, but mem.Total - mem.Used is often "Avail"
		// in simple terms for this view or we could use the passed definition.
		// Original code calculated: available := mem.Total - mem.Used
		available := mem.Total - mem.Used
		lines = append(lines, fmt.Sprintf("Avail  %s", humanBytes(available)))
	}
	// Memory pressure status.
	if mem.Pressure != "" {
		pressureStyle := okStyle
		pressureText := "Status " + mem.Pressure
		switch mem.Pressure {
		case "warn":
			pressureStyle = warnStyle
		case "critical":
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
			pressureStyle = dangerStyle
		}
		lines = append(lines, pressureStyle.Render(pressureText))
	}
	return cardData{icon: iconMemory, title: "Memory", lines: lines}
}

func renderDiskCard(disks []DiskStatus, io DiskIOStatus) cardData {
	var lines []string
	if len(disks) == 0 {
		lines = append(lines, subtleStyle.Render("Collecting..."))
	} else {
		internal, external := splitDisks(disks)
		addGroup := func(prefix string, list []DiskStatus) {
			if len(list) == 0 {
				return
			}
			for i, d := range list {
				label := diskLabel(prefix, i, len(list))
				lines = append(lines, formatDiskLine(label, d))
			}
		}
		addGroup("INTR", internal)
		addGroup("EXTR", external)
		if len(lines) == 0 {
			lines = append(lines, subtleStyle.Render("No disks detected"))
		}
	}
	readBar := ioBar(io.ReadRate)
	writeBar := ioBar(io.WriteRate)
	lines = append(lines, fmt.Sprintf("Read   %s  %.1f MB/s", readBar, io.ReadRate))
	lines = append(lines, fmt.Sprintf("Write  %s  %.1f MB/s", writeBar, io.WriteRate))
	return cardData{icon: iconDisk, title: "Disk", lines: lines}
}

func splitDisks(disks []DiskStatus) (internal, external []DiskStatus) {
	for _, d := range disks {
		if d.External {
			external = append(external, d)
		} else {
			internal = append(internal, d)
		}
	}
	return internal, external
}

func diskLabel(prefix string, index int, total int) string {
	if total <= 1 {
		return prefix
	}
	return fmt.Sprintf("%s%d", prefix, index+1)
}

func formatDiskLine(label string, d DiskStatus) string {
	if label == "" {
		label = "DISK"
	}
	bar := progressBar(d.UsedPercent)
	used := humanBytesShort(d.Used)
	total := humanBytesShort(d.Total)
<<<<<<< HEAD
	return fmt.Sprintf("%-6s %s  %5.1f%% (%s/%s)", label, bar, d.UsedPercent, used, total)
}

func ioBar(rate float64) string {
	// Scale: 0-50 MB/s maps to 0-5 blocks
	filled := int(rate / 10.0)
	if filled > 5 {
		filled = 5
	}
=======
	return fmt.Sprintf("%-6s %s  %5.1f%%, %s/%s", label, bar, d.UsedPercent, used, total)
}

func ioBar(rate float64) string {
	filled := min(int(rate/10.0), 5)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("▮", filled) + strings.Repeat("▯", 5-filled)
	if rate > 80 {
		return dangerStyle.Render(bar)
	}
	if rate > 30 {
		return warnStyle.Render(bar)
	}
	return okStyle.Render(bar)
}

func renderProcessCard(procs []ProcessInfo) cardData {
	var lines []string
	maxProcs := 3
	for i, p := range procs {
		if i >= maxProcs {
			break
		}
		name := shorten(p.Name, 12)
		cpuBar := miniBar(p.CPU)
		lines = append(lines, fmt.Sprintf("%-12s  %s  %5.1f%%", name, cpuBar, p.CPU))
	}
	if len(lines) == 0 {
		lines = append(lines, subtleStyle.Render("No data"))
	}
	return cardData{icon: iconProcs, title: "Processes", lines: lines}
}

<<<<<<< HEAD
func miniBar(percent float64) string {
	filled := int(percent / 20) // 5 chars max for 100%
	if filled > 5 {
		filled = 5
	}
=======
func buildCards(m MetricsSnapshot, width int) []cardData {
	cards := []cardData{
		renderCPUCard(m.CPU, m.Thermal),
		renderMemoryCard(m.Memory),
		renderDiskCard(m.Disks, m.DiskIO),
		renderBatteryCard(m.Batteries, m.Thermal),
		renderProcessCard(m.TopProcesses),
		renderNetworkCard(m.Network, m.NetworkHistory, m.Proxy, width),
	}
	// Sensors card disabled - redundant with CPU temp
	// if hasSensorData(m.Sensors) {
	// 	cards = append(cards, renderSensorsCard(m.Sensors))
	// }
	return cards
}

func miniBar(percent float64) string {
	filled := min(int(percent/20), 5)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if filled < 0 {
		filled = 0
	}
	return colorizePercent(percent, strings.Repeat("▮", filled)+strings.Repeat("▯", 5-filled))
}

<<<<<<< HEAD
func renderNetworkCard(netStats []NetworkStatus, proxy ProxyStatus) cardData {
=======
func renderNetworkCard(netStats []NetworkStatus, history NetworkHistory, proxy ProxyStatus, cardWidth int) cardData {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	var lines []string
	var totalRx, totalTx float64
	var primaryIP string

	for _, n := range netStats {
		totalRx += n.RxRateMBs
		totalTx += n.TxRateMBs
		if primaryIP == "" && n.IP != "" && n.Name == "en0" {
			primaryIP = n.IP
		}
	}

	if len(netStats) == 0 {
		lines = []string{subtleStyle.Render("Collecting...")}
	} else {
<<<<<<< HEAD
		rxBar := netBar(totalRx)
		txBar := netBar(totalTx)
		lines = append(lines, fmt.Sprintf("Down   %s  %s", rxBar, formatRate(totalRx)))
		lines = append(lines, fmt.Sprintf("Up     %s  %s", txBar, formatRate(totalTx)))
		// Proxy + IP
		info := ""
		if proxy.Enabled {
			info = okStyle.Render("Proxy: " + proxy.Type)
		}
		if primaryIP != "" {
			if info != "" {
				info += " · "
			}
			info += primaryIP
		}
		if info != "" {
			lines = append(lines, subtleStyle.Render(info))
=======
		// Calculate dynamic width
		// Layout: "Down   " (7) + graph + "  " (2) + rate (approx 10-12)
		// Safe margin: 22 chars.
		// We target 16 chars to match progressBar implementation for visual consistency.
		graphWidth := cardWidth - 22
		if graphWidth < 5 {
			graphWidth = 5
		}
		if graphWidth > 16 {
			graphWidth = 16 // Match progressBar fixed width
		}

		// sparkline graphs
		rxSparkline := sparkline(history.RxHistory, totalRx, graphWidth)
		txSparkline := sparkline(history.TxHistory, totalTx, graphWidth)
		lines = append(lines, fmt.Sprintf("Down   %s  %s", rxSparkline, formatRate(totalRx)))
		lines = append(lines, fmt.Sprintf("Up     %s  %s", txSparkline, formatRate(totalTx)))
		// Show proxy and IP on one line.
		var infoParts []string
		if proxy.Enabled {
			infoParts = append(infoParts, "Proxy "+proxy.Type)
		}
		if primaryIP != "" {
			infoParts = append(infoParts, primaryIP)
		}
		if len(infoParts) > 0 {
			lines = append(lines, strings.Join(infoParts, " · "))
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		}
	}
	return cardData{icon: iconNetwork, title: "Network", lines: lines}
}

<<<<<<< HEAD
func netBar(rate float64) string {
	// Scale: 0-10 MB/s maps to 0-5 blocks
	filled := int(rate / 2.0)
	if filled > 5 {
		filled = 5
	}
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("▮", filled) + strings.Repeat("▯", 5-filled)
	if rate > 8 {
		return dangerStyle.Render(bar)
	}
	if rate > 3 {
		return warnStyle.Render(bar)
	}
	return okStyle.Render(bar)
=======
// 8 levels: ▁▂▃▄▅▆▇█
func sparkline(history []float64, current float64, width int) string {
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	data := make([]float64, 0, width)
	if len(history) > 0 {
		// Take the most recent points.
		start := 0
		if len(history) > width {
			start = len(history) - width
		}
		data = append(data, history[start:]...)
	}
	// padding with zeros at the start
	for len(data) < width {
		data = append([]float64{0}, data...)
	}
	if len(data) > width {
		data = data[len(data)-width:]
	}

	maxVal := 0.1
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}

	var builder strings.Builder
	for _, v := range data {
		level := int((v / maxVal) * float64(len(blocks)-1))
		if level < 0 {
			level = 0
		}
		if level >= len(blocks) {
			level = len(blocks) - 1
		}
		builder.WriteRune(blocks[level])
	}

	result := builder.String()
	if current > 8 {
		return dangerStyle.Render(result)
	}
	if current > 3 {
		return warnStyle.Render(result)
	}
	return okStyle.Render(result)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

func renderBatteryCard(batts []BatteryStatus, thermal ThermalStatus) cardData {
	var lines []string
	if len(batts) == 0 {
		lines = append(lines, subtleStyle.Render("No battery"))
	} else {
		b := batts[0]
<<<<<<< HEAD
		// Line 1: label + percentage + bar
		lines = append(lines, fmt.Sprintf("Level  %3.0f%%  %s", b.Percent, progressBar(b.Percent)))

		// Line 2: status
		statusIcon := ""
		statusStyle := subtleStyle
		statusLower := strings.ToLower(b.Status)
=======
		statusLower := strings.ToLower(b.Status)
		percentText := fmt.Sprintf("%5.1f%%", b.Percent)
		if b.Percent < 20 && statusLower != "charging" && statusLower != "charged" {
			percentText = dangerStyle.Render(percentText)
		}
		lines = append(lines, fmt.Sprintf("Level  %s  %s", batteryProgressBar(b.Percent), percentText))

		// Add capacity line if available.
		if b.Capacity > 0 {
			capacityText := fmt.Sprintf("%5d%%", b.Capacity)
			if b.Capacity < 70 {
				capacityText = dangerStyle.Render(capacityText)
			} else if b.Capacity < 85 {
				capacityText = warnStyle.Render(capacityText)
			}
			lines = append(lines, fmt.Sprintf("Health %s  %s", batteryProgressBar(float64(b.Capacity)), capacityText))
		}

		statusIcon := ""
		statusStyle := subtleStyle
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if statusLower == "charging" || statusLower == "charged" {
			statusIcon = " ⚡"
			statusStyle = okStyle
		} else if b.Percent < 20 {
			statusStyle = dangerStyle
		}
<<<<<<< HEAD
		// Capitalize first letter
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		statusText := b.Status
		if len(statusText) > 0 {
			statusText = strings.ToUpper(statusText[:1]) + strings.ToLower(statusText[1:])
		}
		if b.TimeLeft != "" {
			statusText += " · " + b.TimeLeft
		}
<<<<<<< HEAD
		lines = append(lines, statusStyle.Render(statusText+statusIcon))

		// Line 3: Health + cycles
=======
		// Add power info.
		if statusLower == "charging" || statusLower == "charged" {
			if thermal.SystemPower > 0 {
				statusText += fmt.Sprintf(" · %.0fW", thermal.SystemPower)
			} else if thermal.AdapterPower > 0 {
				statusText += fmt.Sprintf(" · %.0fW Adapter", thermal.AdapterPower)
			}
		} else if thermal.BatteryPower > 0 {
			// Only show battery power when discharging (positive value)
			statusText += fmt.Sprintf(" · %.0fW", thermal.BatteryPower)
		}
		lines = append(lines, statusStyle.Render(statusText+statusIcon))

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		healthParts := []string{}
		if b.Health != "" {
			healthParts = append(healthParts, b.Health)
		}
		if b.CycleCount > 0 {
			healthParts = append(healthParts, fmt.Sprintf("%d cycles", b.CycleCount))
		}
<<<<<<< HEAD
		if len(healthParts) > 0 {
			lines = append(lines, subtleStyle.Render(strings.Join(healthParts, " · ")))
		}

		// Line 4: Temp + Fan combined
		var thermalParts []string
		if thermal.CPUTemp > 0 {
			tempStyle := okStyle
			if thermal.CPUTemp > 80 {
				tempStyle = dangerStyle
			} else if thermal.CPUTemp > 60 {
				tempStyle = warnStyle
			}
			thermalParts = append(thermalParts, tempStyle.Render(fmt.Sprintf("%.0f°C", thermal.CPUTemp)))
		}
		if thermal.FanSpeed > 0 {
			thermalParts = append(thermalParts, fmt.Sprintf("%d RPM", thermal.FanSpeed))
		}
		if len(thermalParts) > 0 {
			lines = append(lines, strings.Join(thermalParts, " · "))
		}
	}
	return cardData{icon: iconBattery, title: "Power", lines: lines}
}

func renderSensorsCard(sensors []SensorReading) cardData {
	var lines []string
	for _, s := range sensors {
		if s.Note != "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%-12s %s", shorten(s.Label, 12), colorizeTemp(s.Value)+s.Unit))
	}
	if len(lines) == 0 {
		lines = append(lines, subtleStyle.Render("No sensors"))
	}
	return cardData{icon: iconSensors, title: "Sensors", lines: lines}
=======

		if thermal.CPUTemp > 0 {
			tempText := colorizeTemp(thermal.CPUTemp) + "°C" // Reuse common color logic
			healthParts = append(healthParts, tempText)
		}

		if thermal.FanSpeed > 0 {
			healthParts = append(healthParts, fmt.Sprintf("%d RPM", thermal.FanSpeed))
		}

		if len(healthParts) > 0 {
			lines = append(lines, strings.Join(healthParts, " · "))
		}
	}

	return cardData{icon: iconBattery, title: "Power", lines: lines}
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

func renderCard(data cardData, width int, height int) string {
	titleText := data.icon + " " + data.title
<<<<<<< HEAD
	lineLen := width - lipgloss.Width(titleText) - 1
	if lineLen < 4 {
		lineLen = 4
	}
	header := titleStyle.Render(titleText) + " " + lineStyle.Render(strings.Repeat("─", lineLen))
	content := header + "\n" + strings.Join(data.lines, "\n") + "\n"

	// Pad to target height
=======
	lineLen := max(width-lipgloss.Width(titleText)-2, 4)
	header := titleStyle.Render(titleText) + "  " + lineStyle.Render(strings.Repeat("╌", lineLen))
	content := header + "\n" + strings.Join(data.lines, "\n")

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func progressBar(percent float64) string {
<<<<<<< HEAD
	total := 18
=======
	total := 16
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := int(percent / 100 * float64(total))
<<<<<<< HEAD
	if filled > total {
		filled = total
	}

	var builder strings.Builder
	for i := 0; i < total; i++ {
=======

	var builder strings.Builder
	for i := range total {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if i < filled {
			builder.WriteString("█")
		} else {
			builder.WriteString("░")
		}
	}
	return colorizePercent(percent, builder.String())
}

<<<<<<< HEAD
func colorizePercent(percent float64, s string) string {
	switch {
	case percent >= 90:
		return dangerStyle.Render(s)
	case percent >= 70:
=======
func batteryProgressBar(percent float64) string {
	total := 16
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := int(percent / 100 * float64(total))

	var builder strings.Builder
	for i := range total {
		if i < filled {
			builder.WriteString("█")
		} else {
			builder.WriteString("░")
		}
	}
	return colorizeBattery(percent, builder.String())
}

func colorizePercent(percent float64, s string) string {
	switch {
	case percent >= 85:
		return dangerStyle.Render(s)
	case percent >= 60:
		return warnStyle.Render(s)
	default:
		return okStyle.Render(s)
	}
}

func colorizeBattery(percent float64, s string) string {
	switch {
	case percent < 20:
		return dangerStyle.Render(s)
	case percent < 50:
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		return warnStyle.Render(s)
	default:
		return okStyle.Render(s)
	}
}

func colorizeTemp(t float64) string {
	switch {
<<<<<<< HEAD
	case t >= 85:
		return dangerStyle.Render(fmt.Sprintf("%.1f", t))
	case t >= 70:
		return warnStyle.Render(fmt.Sprintf("%.1f", t))
	default:
		return subtleStyle.Render(fmt.Sprintf("%.1f", t))
=======
	case t >= 76:
		return dangerStyle.Render(fmt.Sprintf("%.1f", t))
	case t >= 56:
		return warnStyle.Render(fmt.Sprintf("%.1f", t))
	default:
		return okStyle.Render(fmt.Sprintf("%.1f", t))
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	}
}

func formatRate(mb float64) string {
	if mb < 0.01 {
		return "0 MB/s"
	}
	if mb < 1 {
		return fmt.Sprintf("%.2f MB/s", mb)
	}
	if mb < 10 {
		return fmt.Sprintf("%.1f MB/s", mb)
	}
	return fmt.Sprintf("%.0f MB/s", mb)
}

func humanBytes(v uint64) string {
	switch {
	case v > 1<<40:
		return fmt.Sprintf("%.1f TB", float64(v)/(1<<40))
	case v > 1<<30:
		return fmt.Sprintf("%.1f GB", float64(v)/(1<<30))
	case v > 1<<20:
		return fmt.Sprintf("%.1f MB", float64(v)/(1<<20))
	case v > 1<<10:
		return fmt.Sprintf("%.1f KB", float64(v)/(1<<10))
	default:
		return strconv.FormatUint(v, 10) + " B"
	}
}

func humanBytesShort(v uint64) string {
	switch {
	case v >= 1<<40:
		return fmt.Sprintf("%.0fT", float64(v)/(1<<40))
	case v >= 1<<30:
		return fmt.Sprintf("%.0fG", float64(v)/(1<<30))
	case v >= 1<<20:
		return fmt.Sprintf("%.0fM", float64(v)/(1<<20))
	case v >= 1<<10:
		return fmt.Sprintf("%.0fK", float64(v)/(1<<10))
	default:
		return strconv.FormatUint(v, 10)
	}
}

<<<<<<< HEAD
func shorten(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
=======
func humanBytesCompact(v uint64) string {
	switch {
	case v >= 1<<40:
		return fmt.Sprintf("%.1fT", float64(v)/(1<<40))
	case v >= 1<<30:
		return fmt.Sprintf("%.1fG", float64(v)/(1<<30))
	case v >= 1<<20:
		return fmt.Sprintf("%.1fM", float64(v)/(1<<20))
	case v >= 1<<10:
		return fmt.Sprintf("%.1fK", float64(v)/(1<<10))
	default:
		return strconv.FormatUint(v, 10)
	}
}

func shorten(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

func renderTwoColumns(cards []cardData, width int) string {
	if len(cards) == 0 {
		return ""
	}
	cw := colWidth
	if width > 0 && width/2-2 > cw {
		cw = width/2 - 2
	}
	var rows []string
	for i := 0; i < len(cards); i += 2 {
		left := renderCard(cards[i], cw, 0)
		right := ""
		if i+1 < len(cards) {
			right = renderCard(cards[i+1], cw, 0)
		}
		targetHeight := maxInt(lipgloss.Height(left), lipgloss.Height(right))
		left = renderCard(cards[i], cw, targetHeight)
		if right != "" {
			right = renderCard(cards[i+1], cw, targetHeight)
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right))
		} else {
			rows = append(rows, left)
		}
	}
<<<<<<< HEAD
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
=======

	var spacedRows []string
	for i, r := range rows {
		if i > 0 {
			spacedRows = append(spacedRows, "")
		}
		spacedRows = append(spacedRows, r)
	}
	return lipgloss.JoinVertical(lipgloss.Left, spacedRows...)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
