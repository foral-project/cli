// Testes do comando `foral validate`.
package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// Fixture: catalog-info.yaml válido.
const validCatalog = `"@context": "https://foral-project.github.io/protocol/context/v1/catalog.jsonld"
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: test-project
  description: "Test project."
  annotations:
    foral.dev/archetype: application
  tags:
    - test
    - governance
spec:
  type: service
  lifecycle: experimental
  owner: test-org
`

// Fixture: catalog-info.yaml sem @context.
const missingContextCatalog = `apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: test-project
spec:
  type: service
  lifecycle: experimental
  owner: test-org
`

// Fixture: catalog-info.yaml com naming inválido.
const invalidNamingCatalog = `"@context": "https://foral-project.github.io/protocol/context/v1/catalog.jsonld"
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: MyProject_Bad
  tags:
    - Invalid_Tag
spec:
  type: service
  lifecycle: experimental
  owner: Bad_Owner
`

func TestValidatePassesOnValidCatalog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalog-info.yaml")
	if err := os.WriteFile(path, []byte(validCatalog), 0644); err != nil {
		t.Fatal(err)
	}

	err := runValidate(path, false, false, false)
	if err != nil {
		t.Errorf("Validação deveria passar para catalog válido: %v", err)
	}
}

func TestValidateFailsOnMissingContext(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalog-info.yaml")
	if err := os.WriteFile(path, []byte(missingContextCatalog), 0644); err != nil {
		t.Fatal(err)
	}

	err := runValidate(path, false, false, false)
	if err == nil {
		t.Error("Validação deveria falhar para catalog sem @context")
	}
}

func TestValidateFailsOnInvalidNaming(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalog-info.yaml")
	if err := os.WriteFile(path, []byte(invalidNamingCatalog), 0644); err != nil {
		t.Fatal(err)
	}

	err := runValidate(path, false, false, false)
	if err == nil {
		t.Error("Validação deveria falhar para naming inválido")
	}
}

func TestValidateFailsOnMissingFile(t *testing.T) {
	err := runValidate("/nonexistent/catalog-info.yaml", false, false, false)
	if err == nil {
		t.Error("Validação deveria falhar para arquivo inexistente")
	}
}
