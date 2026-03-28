package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"chk/data"
	"chk/steps"
	"chk/styles"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	stepSelectType step = iota
	stepEnterName
	stepRunCreate
	stepDetect
	stepSelectExtras
	stepInstalling
	stepScaffold
	stepDone
)

type model struct {
	step        step
	kind        data.ProjectKind
	projectName string
	projectDir  string
	isReact     bool
	cursor      int
	extras      []data.ExtraPkg
	textInput   textinput.Model
	spinner     spinner.Model
	statusMsg   string
	errMsg      string
	width       int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "my-app"
	ti.CharLimit = 64
	ti.Width = 30
	ti.PromptStyle = styles.Prompt
	ti.TextStyle = lipgloss.NewStyle().Foreground(styles.ColText)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(styles.ColMauve)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.ColMauve)

	return model{
		step:      stepSelectType,
		textInput: ti,
		spinner:   s,
		extras:    data.DefaultExtras(),
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

// ── Update ───────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width

	case spinner.TickMsg:
		if m.step == stepInstalling || m.step == stepScaffold {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case steps.CreateDoneMsg:
		if msg.Err != nil {
			m.errMsg = fmt.Sprintf("Create failed: %v", msg.Err)
			m.step = stepDone
			return m, nil
		}
		m.step = stepDetect
		return m, steps.DetectReact(m.projectDir)

	case steps.DetectMsg:
		m.isReact = msg.IsReact
		if m.isReact {
			m.step = stepSelectExtras
			m.cursor = 0
		} else {
			m.step = stepDone
			m.statusMsg = "Project created"
		}
		return m, nil

	case steps.InstallDoneMsg:
		if msg.Err != nil {
			m.errMsg = fmt.Sprintf("Install failed: %v", msg.Err)
			m.step = stepDone
			return m, nil
		}
		m.step = stepScaffold
		m.statusMsg = "Scaffolding files..."
		return m, tea.Batch(m.spinner.Tick, steps.ScaffoldProject(m.projectDir, m.extras))

	case steps.ScaffoldDoneMsg:
		m.step = stepDone
		m.statusMsg = "Project created and configured"
		return m, nil
	}

	switch m.step {
	case stepSelectType:
		return m.updateSelectType(msg)
	case stepEnterName:
		return m.updateEnterName(msg)
	case stepSelectExtras:
		return m.updateSelectExtras(msg)
	case stepDone:
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "q" || key.String() == "ctrl+c") {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) updateSelectType(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "j", "down":
		if m.cursor < len(data.ProjectChoices)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "enter":
		m.kind = data.ProjectChoices[m.cursor].Kind
		m.step = stepEnterName
		m.textInput.Focus()
		m.cursor = 0
		return m, textinput.Blink
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m model) updateEnterName(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if ok {
		switch key.String() {
		case "enter":
			name := strings.TrimSpace(m.textInput.Value())
			if name == "" {
				name = "my-app"
			}
			m.projectName = name
			m.projectDir, _ = filepath.Abs(name)
			m.step = stepRunCreate
			m.textInput.Blur()
			return m, steps.RunCreate(m.kind, m.projectName)
		case "esc":
			m.step = stepSelectType
			m.textInput.Blur()
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) updateSelectExtras(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "j", "down":
		if m.cursor < len(m.extras)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case " ", "x":
		m.extras[m.cursor].Checked = !m.extras[m.cursor].Checked
	case "enter":
		m.step = stepInstalling
		m.statusMsg = "Installing packages..."
		return m, tea.Batch(m.spinner.Tick, steps.InstallPackages(m.projectDir, m.extras))
	case "esc":
		m.step = stepDone
		m.statusMsg = "Project created"
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

// ── View ─────────────────────────────────────────────────────────────

func (m model) View() string {
	var s strings.Builder
	w := m.width
	if w == 0 || w > 60 {
		w = 60
	}

	s.WriteString(steps.Header())
	s.WriteString(styles.Dim.Render(strings.Repeat("─", w)) + "\n\n")

	switch m.step {
	case stepSelectType:
		s.WriteString(styles.Prompt.Render("  What type of project?") + "\n\n")
		for i, c := range data.ProjectChoices {
			cursor := "  "
			style := styles.Item
			if i == m.cursor {
				cursor = styles.Accent.Render("❯ ")
				style = styles.Active
			}
			s.WriteString("  " + cursor + style.Render(c.Icon+" "+c.Name))
			s.WriteString(styles.Dim.Render("  "+c.Desc) + "\n")
		}
		s.WriteString("\n" + steps.Footer("↑/↓ navigate", "enter select", "q quit"))

	case stepEnterName:
		kindName := data.ProjectChoices[m.kind].Name
		s.WriteString(styles.Prompt.Render("  Project name") + styles.Dim.Render(" ("+kindName+")") + "\n\n")
		s.WriteString("  " + m.textInput.View() + "\n")
		s.WriteString("\n" + steps.Footer("enter confirm", "esc back"))

	case stepRunCreate:
		s.WriteString(styles.Dim.Render("  Running create command...") + "\n")

	case stepSelectExtras:
		s.WriteString(styles.Prompt.Render("  React detected!") + " " + styles.Dim.Render("Select packages to install:") + "\n\n")
		for i, e := range m.extras {
			cursor := "  "
			style := styles.Item
			if i == m.cursor {
				cursor = styles.Accent.Render("❯ ")
				style = styles.Active
			}
			check := styles.CheckOff.Render("[ ]")
			if e.Checked {
				check = styles.CheckOn.Render("[✓]")
			}
			s.WriteString("  " + cursor + check + " " + style.Render(e.Name))
			s.WriteString(styles.Dim.Render("  "+e.Desc) + "\n")
		}
		s.WriteString("\n" + steps.Footer("↑/↓ navigate", "space toggle", "enter install", "esc skip"))

	case stepInstalling:
		s.WriteString("  " + m.spinner.View() + " " + styles.Dim.Render(m.statusMsg) + "\n")

	case stepScaffold:
		s.WriteString("  " + m.spinner.View() + " " + styles.Dim.Render("Scaffolding files...") + "\n")

	case stepDone:
		if m.errMsg != "" {
			s.WriteString("  " + styles.Error.Render("✗ "+m.errMsg) + "\n")
		} else {
			s.WriteString("  " + styles.Success.Render("✓ "+m.statusMsg) + "\n\n")
			s.WriteString(styles.Dim.Render("  Project: ") + lipgloss.NewStyle().Foreground(styles.ColText).Render(m.projectName) + "\n")
			s.WriteString(styles.Dim.Render("  Path:    ") + lipgloss.NewStyle().Foreground(styles.ColText).Render(m.projectDir) + "\n")

			var installed []string
			for _, e := range m.extras {
				if e.Checked {
					installed = append(installed, e.Name)
				}
			}
			if len(installed) > 0 {
				s.WriteString(styles.Dim.Render("  Extras:  ") + lipgloss.NewStyle().Foreground(styles.ColTeal).Render(strings.Join(installed, ", ")) + "\n")
			}
			s.WriteString("\n" + styles.Dim.Render("  Get started:") + "\n")
			s.WriteString(styles.Accent.Render("    cd "+m.projectName) + "\n")
			s.WriteString(styles.Accent.Render("    npm run dev") + "\n")
		}
		s.WriteString("\n" + steps.Footer("q quit"))
	}

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
