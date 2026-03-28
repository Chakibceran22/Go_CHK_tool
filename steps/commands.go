package steps

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"chk/data"
	"chk/scaffold"

	tea "github.com/charmbracelet/bubbletea"
)

// Messages
type CreateDoneMsg struct{ Err error }
type DetectMsg struct{ IsReact bool }
type InstallDoneMsg struct{ Err error }
type ScaffoldDoneMsg struct{}

func RunCreate(kind data.ProjectKind, name string) tea.Cmd {
	var cmd *exec.Cmd
	if kind == data.KindVite {
		cmd = exec.Command("npm", "create", "vite@latest", name)
	} else {
		cmd = exec.Command("npx", "create-next-app@latest", name)
	}
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return CreateDoneMsg{Err: err}
	})
}

func DetectReact(projectDir string) tea.Cmd {
	return func() tea.Msg {
		pkgPath := filepath.Join(projectDir, "package.json")
		raw, err := os.ReadFile(pkgPath)
		if err != nil {
			return DetectMsg{IsReact: false}
		}
		var pkg map[string]interface{}
		if err := json.Unmarshal(raw, &pkg); err != nil {
			return DetectMsg{IsReact: false}
		}
		deps, _ := pkg["dependencies"].(map[string]interface{})
		_, hasReact := deps["react"]
		return DetectMsg{IsReact: hasReact}
	}
}

func InstallPackages(projectDir string, extras []data.ExtraPkg) tea.Cmd {
	return func() tea.Msg {
		var pkgs []string
		for _, e := range extras {
			if !e.Checked {
				continue
			}
			if strings.Contains(e.Name, "TanStack") {
				pkgs = append(pkgs, "@tanstack/react-query", "@tanstack/react-query-devtools")
			}
			if strings.Contains(e.Name, "Zustand") {
				pkgs = append(pkgs, "zustand")
			}
		}
		if len(pkgs) == 0 {
			return InstallDoneMsg{}
		}
		args := append([]string{"install"}, pkgs...)
		cmd := exec.Command("npm", args...)
		cmd.Dir = projectDir
		return InstallDoneMsg{Err: cmd.Run()}
	}
}

func ScaffoldProject(projectDir string, extras []data.ExtraPkg) tea.Cmd {
	return func() tea.Msg {
		srcDir := filepath.Join(projectDir, "src")

		wantTanstack := false
		wantZustand := false
		for _, e := range extras {
			if !e.Checked {
				continue
			}
			if strings.Contains(e.Name, "TanStack") {
				wantTanstack = true
			}
			if strings.Contains(e.Name, "Zustand") {
				wantZustand = true
			}
		}

		if wantTanstack {
			scaffold.WriteQueryProvider(srcDir)
			scaffold.WrapMainWithProvider(srcDir)
		}
		if wantZustand {
			scaffold.WriteZustandStore(srcDir)
		}

		return ScaffoldDoneMsg{}
	}
}
