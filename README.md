# chk — Project Scaffolder CLI

A terminal UI tool built with Go + [Bubble Tea](https://github.com/charmbracelet/bubbletea) that scaffolds frontend projects with your preferred packages pre-configured.

## What it does

1. **Asks project type** — Vite or Next.js (more coming: Electron, Capacitor, etc.)
2. **Asks project name** — text input with default `my-app`
3. **Runs the create command** — hands your terminal to `npm create vite@latest` or `npx create-next-app@latest` so you interact with their prompts directly
4. **Detects React** — reads the generated `package.json` to check if `react` is in dependencies
5. **Offers extras** — if React detected, shows a checkbox list:
   - TanStack Query + DevTools (async state management)
   - Zustand (lightweight state management)
6. **Installs packages** — runs `npm install` with selected packages
7. **Scaffolds files** — writes boilerplate into your project:
   - `src/providers/QueryProvider.tsx` — QueryClient setup with DevTools
   - `src/store/useAppStore.ts` — Zustand store with example state
   - Updates `src/main.tsx` — wraps `<App />` with `<QueryProvider>`

## Project structure

```
chk-cli/
├── main.go              # App model, update loop, view rendering
├── go.mod               # Go module definition
├── styles/
│   └── styles.go        # Catppuccin Frappé color palette + lipgloss styles
├── data/
│   └── projects.go      # Project type definitions (Vite, Next.js) + extras list
├── steps/
│   ├── commands.go      # Tea commands: create, detect react, install, scaffold
│   └── view.go          # Shared view helpers (header, footer)
├── scaffold/
│   ├── templates.go     # File templates (QueryProvider.tsx, useAppStore.ts)
│   └── scaffold.go      # File writing logic + main.tsx wrapping
```

## How it works internally

### main.go
The app is a state machine with these steps:
- `stepSelectType` → list selection with j/k + enter
- `stepEnterName` → bubbletea `textinput` component
- `stepRunCreate` → `tea.ExecProcess` suspends the TUI and gives the terminal to npm/npx
- `stepDetect` → reads `package.json` in a goroutine, sends a message back
- `stepSelectExtras` → checkbox list with space to toggle, enter to confirm
- `stepInstalling` → spinner runs while `npm install` happens in a goroutine
- `stepScaffold` → writes template files, wraps main.tsx
- `stepDone` → summary screen

### styles/styles.go
All colors use the **Catppuccin Frappé** palette. Every text style (Title, Prompt, Active, Dim, etc.) is defined once here and imported everywhere.

### data/projects.go
Defines `ProjectChoice` structs (name, description, icon, kind) and `ExtraPkg` structs (name, description, default checked state). To add a new project type, add an entry to `ProjectChoices`. To add a new extra package, add to `DefaultExtras()`.

### steps/commands.go
Each async operation is a `tea.Cmd` (a function that returns a `tea.Msg`):
- `RunCreate()` — uses `tea.ExecProcess` to hand terminal control to npm
- `DetectReact()` — parses package.json for react dependency
- `InstallPackages()` — runs `npm install` with the selected package names
- `ScaffoldProject()` — calls functions from the scaffold package

### scaffold/
- `templates.go` — raw string constants for the generated files
- `scaffold.go` — `WriteQueryProvider()`, `WriteZustandStore()`, `WrapMainWithProvider()` which reads main.tsx, adds the import, and wraps `<App />` with `<QueryProvider>`

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework (Elm architecture)
- [bubbles](https://github.com/charmbracelet/bubbles) — pre-built components (textinput, spinner)
- [lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling/colors

## Build & run

```bash
# Build
go build -o chk .

# Run
./chk

# Install globally
cp chk ~/.local/bin/chk
```

## Keybindings

| Key       | Action              |
|-----------|---------------------|
| `j` / `↓` | Move down           |
| `k` / `↑` | Move up             |
| `enter`   | Confirm / select    |
| `space`   | Toggle checkbox     |
| `esc`     | Go back / skip      |
| `q`       | Quit                |

## Adding new project types

1. Add a new `KindXxx` constant in `data/projects.go`
2. Add a `ProjectChoice` entry to `ProjectChoices`
3. Add the create command in `steps/commands.go` → `RunCreate()`
4. Add any new extras to `DefaultExtras()` if needed
5. Add scaffold templates in `scaffold/templates.go` and writing logic in `scaffold/scaffold.go`
