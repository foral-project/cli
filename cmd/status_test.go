// Testes do comando `foral status`.
package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStatusPassesOnValidCatalog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalog-info.yaml")
	if err := os.WriteFile(path, []byte(validCatalog), 0644); err != nil {
		t.Fatal(err)
	}

	err := runStatus(dir)
	if err != nil {
		t.Errorf("Status deveria passar para catalog válido: %v", err)
	}
}

func TestStatusFailsOnMissingCatalog(t *testing.T) {
	dir := t.TempDir()

	err := runStatus(dir)
	if err == nil {
		t.Error("Status deveria falhar sem catalog-info.yaml")
	}
}
