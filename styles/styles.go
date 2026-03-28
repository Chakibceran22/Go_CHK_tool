package styles

import "github.com/charmbracelet/lipgloss"

// Catppuccin Frappé palette
var (
	ColBg       = lipgloss.Color("#303446")
	ColText     = lipgloss.Color("#c6d0f5")
	ColSubtext0 = lipgloss.Color("#a5adce")
	ColSubtext1 = lipgloss.Color("#b5bfe2")
	ColOverlay0 = lipgloss.Color("#737994")
	ColSurface1 = lipgloss.Color("#51576d")
	ColBlue     = lipgloss.Color("#8caaee")
	ColGreen    = lipgloss.Color("#a6d189")
	ColMauve    = lipgloss.Color("#ca9ee6")
	ColRed      = lipgloss.Color("#e78284")
	ColTeal     = lipgloss.Color("#81c8be")
)

var (
	Title    = lipgloss.NewStyle().Foreground(ColMauve).Bold(true)
	Subtitle = lipgloss.NewStyle().Foreground(ColOverlay0)
	Dim      = lipgloss.NewStyle().Foreground(ColOverlay0)
	Prompt   = lipgloss.NewStyle().Foreground(ColBlue).Bold(true)
	Accent   = lipgloss.NewStyle().Foreground(ColMauve)
	Item     = lipgloss.NewStyle().Foreground(ColSubtext1)
	Active   = lipgloss.NewStyle().Foreground(ColMauve).Bold(true)
	Success  = lipgloss.NewStyle().Foreground(ColGreen).Bold(true)
	Error    = lipgloss.NewStyle().Foreground(ColRed).Bold(true)
	CheckOn  = lipgloss.NewStyle().Foreground(ColGreen).Bold(true)
	CheckOff = lipgloss.NewStyle().Foreground(ColSurface1)
	Key      = lipgloss.NewStyle().Foreground(ColMauve).Bold(true)
	Desc     = lipgloss.NewStyle().Foreground(ColSubtext0)
)
