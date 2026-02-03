package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func displayPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return path
	}
	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1)
	}
	return path
}

<<<<<<< HEAD
// truncateMiddle truncates string in the middle, keeping head and tail.
=======
// truncateMiddle trims the middle, keeping head and tail.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
func truncateMiddle(s string, maxWidth int) string {
	runes := []rune(s)
	currentWidth := displayWidth(s)

	if currentWidth <= maxWidth {
		return s
	}

<<<<<<< HEAD
	// Reserve 3 width for "..."
	if maxWidth < 10 {
		// Simple truncation for very small width
=======
	if maxWidth < 10 {
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		width := 0
		for i, r := range runes {
			width += runeWidth(r)
			if width > maxWidth {
				return string(runes[:i])
			}
		}
		return s
	}

<<<<<<< HEAD
	// Keep more of the tail (filename usually more important)
	targetHeadWidth := (maxWidth - 3) / 3
	targetTailWidth := maxWidth - 3 - targetHeadWidth

	// Find head cutoff point based on display width
=======
	targetHeadWidth := (maxWidth - 3) / 3
	targetTailWidth := maxWidth - 3 - targetHeadWidth

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	headWidth := 0
	headIdx := 0
	for i, r := range runes {
		w := runeWidth(r)
		if headWidth+w > targetHeadWidth {
			break
		}
		headWidth += w
		headIdx = i + 1
	}

<<<<<<< HEAD
	// Find tail cutoff point
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	tailWidth := 0
	tailIdx := len(runes)
	for i := len(runes) - 1; i >= 0; i-- {
		w := runeWidth(runes[i])
		if tailWidth+w > targetTailWidth {
			break
		}
		tailWidth += w
		tailIdx = i
	}

	return string(runes[:headIdx]) + "..." + string(runes[tailIdx:])
}

func formatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}

func humanizeBytes(size int64) string {
	if size < 0 {
		return "0 B"
	}
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	value := float64(size) / float64(div)
	return fmt.Sprintf("%.1f %cB", value, "KMGTPE"[exp])
}

<<<<<<< HEAD
func coloredProgressBar(value, max int64, percent float64) string {
	if max <= 0 {
		return colorGray + strings.Repeat("░", barWidth) + colorReset
	}

	filled := int((value * int64(barWidth)) / max)
	if filled > barWidth {
		filled = barWidth
	}

	// Choose color based on percentage
=======
func coloredProgressBar(value, maxValue int64, percent float64) string {
	if maxValue <= 0 {
		return colorGray + strings.Repeat("░", barWidth) + colorReset
	}

	filled := min(int((value*int64(barWidth))/maxValue), barWidth)

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	var barColor string
	if percent >= 50 {
		barColor = colorRed
	} else if percent >= 20 {
		barColor = colorYellow
	} else if percent >= 5 {
<<<<<<< HEAD
		barColor = colorCyan
=======
		barColor = colorBlue
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
	} else {
		barColor = colorGreen
	}

<<<<<<< HEAD
	bar := barColor
	for i := 0; i < barWidth; i++ {
		if i < filled {
			if i < filled-1 {
				bar += "█"
			} else {
				remainder := (value * int64(barWidth)) % max
				if remainder > max/2 {
					bar += "█"
				} else if remainder > max/4 {
					bar += "▓"
				} else {
					bar += "▒"
				}
			}
		} else {
			bar += colorGray + "░" + barColor
		}
	}
	return bar + colorReset
}

// Calculate display width considering CJK characters.
func runeWidth(r rune) int {
	if r >= 0x4E00 && r <= 0x9FFF ||
		r >= 0x3400 && r <= 0x4DBF ||
		r >= 0xAC00 && r <= 0xD7AF ||
		r >= 0xFF00 && r <= 0xFFEF {
=======
	var bar strings.Builder
	bar.WriteString(barColor)
	for i := range barWidth {
		if i < filled {
			if i < filled-1 {
				bar.WriteString("█")
			} else {
				remainder := (value * int64(barWidth)) % maxValue
				if remainder > maxValue/2 {
					bar.WriteString("█")
				} else if remainder > maxValue/4 {
					bar.WriteString("▓")
				} else {
					bar.WriteString("▒")
				}
			}
		} else {
			bar.WriteString(colorGray + "░" + barColor)
		}
	}
	return bar.String() + colorReset
}

// runeWidth returns display width for wide characters and emoji.
func runeWidth(r rune) int {
	if r >= 0x4E00 && r <= 0x9FFF || // CJK Unified Ideographs
		r >= 0x3400 && r <= 0x4DBF || // CJK Extension A
		r >= 0x20000 && r <= 0x2A6DF || // CJK Extension B
		r >= 0x2A700 && r <= 0x2B73F || // CJK Extension C
		r >= 0x2B740 && r <= 0x2B81F || // CJK Extension D
		r >= 0x2B820 && r <= 0x2CEAF || // CJK Extension E
		r >= 0x3040 && r <= 0x30FF || // Hiragana and Katakana
		r >= 0x31F0 && r <= 0x31FF || // Katakana Phonetic Extensions
		r >= 0xAC00 && r <= 0xD7AF || // Hangul Syllables
		r >= 0xFF00 && r <= 0xFFEF || // Fullwidth Forms
		r >= 0x1F300 && r <= 0x1F6FF || // Miscellaneous Symbols and Pictographs (includes Transport)
		r >= 0x1F900 && r <= 0x1F9FF || // Supplemental Symbols and Pictographs
		r >= 0x2600 && r <= 0x26FF || // Miscellaneous Symbols
		r >= 0x2700 && r <= 0x27BF || // Dingbats
		r >= 0xFE10 && r <= 0xFE1F || // Vertical Forms
		r >= 0x1F000 && r <= 0x1F02F { // Mahjong Tiles
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		return 2
	}
	return 1
}

func displayWidth(s string) int {
	width := 0
	for _, r := range s {
		width += runeWidth(r)
	}
	return width
}

<<<<<<< HEAD
func trimName(name string) string {
	const (
		maxWidth      = 28
=======
// calculateNameWidth computes name column width from terminal width.
func calculateNameWidth(termWidth int) int {
	const fixedWidth = 61
	available := termWidth - fixedWidth

	if available < 24 {
		return 24
	}
	if available > 60 {
		return 60
	}
	return available
}

func trimNameWithWidth(name string, maxWidth int) string {
	const (
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
		ellipsis      = "..."
		ellipsisWidth = 3
	)

	runes := []rune(name)
	widths := make([]int, len(runes))
	for i, r := range runes {
		widths[i] = runeWidth(r)
	}

	currentWidth := 0
	for i, w := range widths {
		if currentWidth+w > maxWidth {
			subWidth := currentWidth
			j := i
			for j > 0 && subWidth+ellipsisWidth > maxWidth {
				j--
				subWidth -= widths[j]
			}
			if j == 0 {
				return ellipsis
			}
			return string(runes[:j]) + ellipsis
		}
		currentWidth += w
	}

	return name
}

func padName(name string, targetWidth int) string {
	currentWidth := displayWidth(name)
	if currentWidth >= targetWidth {
		return name
	}
	return name + strings.Repeat(" ", targetWidth-currentWidth)
}

<<<<<<< HEAD
// formatUnusedTime formats the time since last access in a compact way.
=======
// formatUnusedTime formats time since last access.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
func formatUnusedTime(lastAccess time.Time) string {
	if lastAccess.IsZero() {
		return ""
	}

	duration := time.Since(lastAccess)
	days := int(duration.Hours() / 24)

	if days < 90 {
		return ""
	}

	months := days / 30
	years := days / 365

	if years >= 2 {
		return fmt.Sprintf(">%dyr", years)
	} else if years >= 1 {
		return ">1yr"
	} else if months >= 3 {
		return fmt.Sprintf(">%dmo", months)
	}

	return ""
}
