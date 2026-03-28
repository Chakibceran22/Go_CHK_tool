package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

func WriteQueryProvider(srcDir string) error {
	dir := filepath.Join(srcDir, "providers")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "QueryProvider.tsx"), []byte(QueryProviderTmpl), 0644)
}

func WriteZustandStore(srcDir string) error {
	dir := filepath.Join(srcDir, "store")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "useAppStore.ts"), []byte(ZustandStoreTmpl), 0644)
}

func WrapMainWithProvider(srcDir string) error {
	mainPath := filepath.Join(srcDir, "main.tsx")
	data, err := os.ReadFile(mainPath)
	if err != nil {
		return nil // not fatal, file might not exist
	}

	content := string(data)
	if strings.Contains(content, "QueryProvider") {
		return nil
	}

	importLine := "import { QueryProvider } from './providers/QueryProvider'\n"

	lastImport := strings.LastIndex(content, "import ")
	if lastImport == -1 {
		content = importLine + content
	} else {
		endOfLine := strings.Index(content[lastImport:], "\n")
		if endOfLine != -1 {
			pos := lastImport + endOfLine + 1
			content = content[:pos] + importLine + content[pos:]
		}
	}

	content = strings.Replace(content, "<App />", "<QueryProvider>\n      <App />\n    </QueryProvider>", 1)
	content = strings.Replace(content, "<App/>", "<QueryProvider>\n      <App />\n    </QueryProvider>", 1)

	return os.WriteFile(mainPath, []byte(content), 0644)
}
