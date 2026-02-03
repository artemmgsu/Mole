<<<<<<< HEAD
=======
// Package main provides the mo status command for real-time system monitoring.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
package main

import (
	"fmt"
	"os"
<<<<<<< HEAD
=======
	"path/filepath"
	"strings"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const refreshInterval = time.Second

var (
	Version   = "dev"
	BuildTime = ""
)

type tickMsg struct{}
type animTickMsg struct{}

type metricsMsg struct {
	data MetricsSnapshot
	err  error
}

type model struct {
	collector   *Collector
	width       int
	height      int
	metrics     MetricsSnapshot
	errMessage  string
	ready       bool
	lastUpdated time.Time
	collecting  bool
	animFrame   int
<<<<<<< HEAD
=======
	catHidden   bool // true = hidden, false = visible
}

// getConfigPath returns the path to the status preferences file.
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "mole", "status_prefs")
}

// loadCatHidden loads the cat hidden preference from config file.
func loadCatHidden() bool {
	path := getConfigPath()
	if path == "" {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == "cat_hidden=true"
}

// saveCatHidden saves the cat hidden preference to config file.
func saveCatHidden(hidden bool) {
	path := getConfigPath()
	if path == "" {
		return
	}
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	value := "cat_hidden=false"
	if hidden {
		value = "cat_hidden=true"
	}
	_ = os.WriteFile(path, []byte(value+"\n"), 0644)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

func newModel() model {
	return model{
		collector: NewCollector(),
<<<<<<< HEAD
=======
		catHidden: loadCatHidden(),
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tickAfter(0), animTick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
<<<<<<< HEAD
=======
		case "k":
			// Toggle cat visibility and persist preference
			m.catHidden = !m.catHidden
			saveCatHidden(m.catHidden)
			return m, nil
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		if m.collecting {
			return m, nil
		}
		m.collecting = true
		return m, m.collectCmd()
	case metricsMsg:
		if msg.err != nil {
			m.errMessage = msg.err.Error()
		} else {
			m.errMessage = ""
		}
		m.metrics = msg.data
		m.lastUpdated = msg.data.CollectedAt
		m.collecting = false
<<<<<<< HEAD
		// Mark ready after first successful data collection
=======
		// Mark ready after first successful data collection.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		if !m.ready {
			m.ready = true
		}
		return m, tickAfter(refreshInterval)
	case animTickMsg:
		m.animFrame++
		return m, animTickWithSpeed(m.metrics.CPU.Usage)
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

<<<<<<< HEAD
	header := renderHeader(m.metrics, m.errMessage, m.animFrame, m.width)
=======
	header := renderHeader(m.metrics, m.errMessage, m.animFrame, m.width, m.catHidden)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	cardWidth := 0
	if m.width > 80 {
		cardWidth = maxInt(24, m.width/2-4)
	}
	cards := buildCards(m.metrics, cardWidth)

	if m.width <= 80 {
		var rendered []string
<<<<<<< HEAD
		for _, c := range cards {
			rendered = append(rendered, renderCard(c, cardWidth, 0))
		}
		return header + "\n" + lipgloss.JoinVertical(lipgloss.Left, rendered...)
	}

	return header + "\n" + renderTwoColumns(cards, m.width)
=======
		for i, c := range cards {
			if i > 0 {
				rendered = append(rendered, "")
			}
			rendered = append(rendered, renderCard(c, cardWidth, 0))
		}
		result := header + "\n" + lipgloss.JoinVertical(lipgloss.Left, rendered...)
		// Add extra newline if cat is hidden for better spacing
		if m.catHidden {
			result = header + "\n\n" + lipgloss.JoinVertical(lipgloss.Left, rendered...)
		}
		return result
	}

	twoCol := renderTwoColumns(cards, m.width)
	// Add extra newline if cat is hidden for better spacing
	if m.catHidden {
		return header + "\n\n" + twoCol
	}
	return header + "\n" + twoCol
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
}

func (m model) collectCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.collector.Collect()
		return metricsMsg{data: data, err: err}
	}
}

func tickAfter(delay time.Duration) tea.Cmd {
	return tea.Tick(delay, func(time.Time) tea.Msg { return tickMsg{} })
}

func animTick() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(time.Time) tea.Msg { return animTickMsg{} })
}

func animTickWithSpeed(cpuUsage float64) tea.Cmd {
<<<<<<< HEAD
	// Higher CPU = faster animation (50ms to 300ms)
	interval := 300 - int(cpuUsage*2.5)
	if interval < 50 {
		interval = 50
	}
=======
	// Higher CPU = faster animation.
	interval := max(300-int(cpuUsage*2.5), 50)
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	return tea.Tick(time.Duration(interval)*time.Millisecond, func(time.Time) tea.Msg { return animTickMsg{} })
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
<<<<<<< HEAD
	if err := p.Start(); err != nil {
=======
	if _, err := p.Run(); err != nil {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		fmt.Fprintf(os.Stderr, "system status error: %v\n", err)
		os.Exit(1)
	}
}
