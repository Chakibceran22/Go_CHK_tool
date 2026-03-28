package data

type ProjectKind int

const (
	KindVite ProjectKind = iota
	KindNext
)

type ProjectChoice struct {
	Name string
	Desc string
	Icon string
	Kind ProjectKind
}

var ProjectChoices = []ProjectChoice{
	{"Vite", "Fast build tool for modern web apps", "⚡", KindVite},
	{"Next.js", "React framework with SSR & routing", "▲", KindNext},
}

type ExtraPkg struct {
	Name    string
	Desc    string
	Checked bool
}

func DefaultExtras() []ExtraPkg {
	return []ExtraPkg{
		{"TanStack Query + DevTools", "Async state management", true},
		{"Zustand", "Lightweight state management", true},
	}
}
