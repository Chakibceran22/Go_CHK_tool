package steps

import (
	"strings"

	"chk/styles"
)

func Footer(parts ...string) string {
	var items []string
	for _, p := range parts {
		split := strings.SplitN(p, " ", 2)
		if len(split) == 2 {
			items = append(items, styles.Key.Render(split[0])+" "+styles.Desc.Render(split[1]))
		}
	}
	return "  " + strings.Join(items, styles.Dim.Render("  ·  ")) + "\n"
}

func Header() string {
	return "\n" + styles.Title.Render("  chk") + styles.Subtitle.Render("  project scaffolder") + "\n"
}
