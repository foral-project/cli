// Testes do comando `foral init`.
package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCreatesFiles(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "test-project")

	err := runInit(projectPath, "application", "test-org", "experimental", "github")
	if err != nil {
		t.Fatalf("foral init falhou: %v", err)
	}

	// Verificar catalog-info.yaml
	catalogPath := filepath.Join(projectPath, "catalog-info.yaml")
	if _, err := os.Stat(catalogPath); os.IsNotExist(err) {
		t.Error("catalog-info.yaml não foi criado")
	}

	// Verificar workflow
	workflowPath := filepath.Join(projectPath, ".github", "workflows", "foral.yml")
	if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
		t.Error(".github/workflows/foral.yml não foi criado")
	}

	// Verificar .gitignore
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Error(".gitignore não foi criado")
	}

	// Verificar conteúdo do catalog-info.yaml
	data, err := os.ReadFile(catalogPath)
	if err != nil {
		t.Fatalf("Erro ao ler catalog-info.yaml: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "name: test-project") {
		t.Error("catalog-info.yaml não contém name: test-project")
	}
	if !strings.Contains(content, "owner: test-org") {
		t.Error("catalog-info.yaml não contém owner: test-org")
	}
	if !strings.Contains(content, "lifecycle: experimental") {
		t.Error("catalog-info.yaml não contém lifecycle: experimental")
	}
}

func TestInitWithoutCI(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "no-ci-project")

	err := runInit(projectPath, "application", "test-org", "experimental", "none")
	if err != nil {
		t.Fatalf("foral init falhou: %v", err)
	}

	workflowPath := filepath.Join(projectPath, ".github", "workflows", "foral.yml")
	if _, err := os.Stat(workflowPath); !os.IsNotExist(err) {
		t.Error("workflow NÃO deveria existir com --ci=none")
	}
}

func TestInitRejectsInvalidDNSLabel(t *testing.T) {
	dir := t.TempDir()

	cases := []struct {
		name string
		path string
	}{
		{"uppercase", filepath.Join(dir, "MyProject")},
		{"underscore", filepath.Join(dir, "my_project")},
		{"dot", filepath.Join(dir, ".hidden")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := runInit(tc.path, "application", "test-org", "experimental", "github")
			if err == nil {
				t.Errorf("Esperava erro para nome inválido '%s'", filepath.Base(tc.path))
			}
		})
	}
}

func TestInitRejectsInvalidArchetype(t *testing.T) {
	dir := t.TempDir()
	err := runInit(filepath.Join(dir, "valid-name"), "invalid-type", "test-org", "experimental", "github")
	if err == nil {
		t.Error("Esperava erro para archetype inválido")
	}
}

func TestInitRejectsInvalidLifecycle(t *testing.T) {
	dir := t.TempDir()
	err := runInit(filepath.Join(dir, "valid-name"), "application", "test-org", "invalid", "github")
	if err == nil {
		t.Error("Esperava erro para lifecycle inválido")
	}
}

func TestInitRejectsInvalidOwner(t *testing.T) {
	dir := t.TempDir()
	err := runInit(filepath.Join(dir, "valid-name"), "application", "BAD_OWNER", "experimental", "github")
	if err == nil {
		t.Error("Esperava erro para owner inválido")
	}
}
