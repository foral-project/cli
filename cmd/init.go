// Comando `foral init` — scaffold de um novo projeto federado.
// Gera catalog-info.yaml e CI workflow.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	archetype  string
	owner      string
	lifecycle  string
	ciPlatform string
)


var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Scaffold de um novo projeto federado",
	Long: `Cria a estrutura mínima para um projeto governado pelo Foral Protocol:
  - catalog-info.yaml (manifesto Backstage)
  - .github/workflows/foral.yml (CI validation)
  - .gitignore`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit(args[0], archetype, owner, lifecycle, ciPlatform)
	},
}

// runInit executa a lógica do scaffold. Exportada para testes.
func runInit(projectPath, arch, own, life, ci string) error {
	projectName := filepath.Base(projectPath)

	// Validação RFC 1123
	if !dnsLabelRegex.MatchString(projectName) {
		return fmt.Errorf(
			"nome '%s' não é um DNS label válido (RFC 1123).\n"+
				"Use: lowercase, hyphens, max 63 caracteres. Ex: my-project",
			projectName,
		)
	}

	// Validar archetype
	validArchetypes := map[string]bool{
		"application": true, "infrastructure": true, "bot": true,
		"library": true, "service": true,
	}
	if !validArchetypes[arch] {
		return fmt.Errorf("archetype '%s' inválido. Válidos: application, infrastructure, bot, library, service", arch)
	}

	// Validar lifecycle
	validLifecycles := map[string]bool{
		"experimental": true, "production": true, "deprecated": true,
	}
	if !validLifecycles[life] {
		return fmt.Errorf("lifecycle '%s' inválido. Válidos: experimental, production, deprecated", life)
	}

	// Validar owner
	if !dnsLabelRegex.MatchString(own) {
		return fmt.Errorf("owner '%s' não é um DNS label válido (RFC 1123)", own)
	}

	// Criar diretórios
	dirs := []string{projectPath}
	if ci == "github" {
		dirs = append(dirs, filepath.Join(projectPath, ".github", "workflows"))
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
		}
	}

	// Dados do template
	data := map[string]string{
		"Name":      projectName,
		"Archetype": arch,
		"Owner":     own,
		"Lifecycle": life,
		"Type":      arch,
	}
	if arch == "infrastructure" {
		data["Type"] = "resource"
	}

	// Gerar catalog-info.yaml
	if err := writeTemplate(
		filepath.Join(projectPath, "catalog-info.yaml"),
		catalogInfoTpl, data,
	); err != nil {
		return err
	}
	fmt.Printf("  ✅ %s/catalog-info.yaml\n", projectPath)

	// Gerar CI workflow
	if ci == "github" {
		if err := writeTemplate(
			filepath.Join(projectPath, ".github", "workflows", "foral.yml"),
			githubWorkflowTpl, data,
		); err != nil {
			return err
		}
		fmt.Printf("  ✅ %s/.github/workflows/foral.yml\n", projectPath)
	}

	// Gerar .gitignore
	if err := writeTemplate(
		filepath.Join(projectPath, ".gitignore"),
		gitignoreTpl, data,
	); err != nil {
		return err
	}
	fmt.Printf("  ✅ %s/.gitignore\n", projectPath)

	fmt.Println()
	fmt.Printf("🏛️  Projeto '%s' criado com governança Foral.\n", projectName)
	fmt.Println()
	fmt.Println("Próximos passos:")
	fmt.Printf("  cd %s\n", projectPath)
	fmt.Println("  git init && git add -A && git commit -m \"feat: initial commit\"")
	fmt.Println("  foral validate")

	return nil
}

func writeTemplate(path, tplContent string, data map[string]string) error {
	tpl, err := template.New(filepath.Base(path)).Parse(tplContent)
	if err != nil {
		return fmt.Errorf("erro ao parsear template: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro ao criar %s: %w", path, err)
	}
	defer f.Close()

	return tpl.Execute(f, data)
}

func init() {
	initCmd.Flags().StringVarP(&archetype, "archetype", "a", "application",
		"Archetype do projeto (application, infrastructure, bot, library, service)")
	initCmd.Flags().StringVarP(&owner, "owner", "o", "foral-project",
		"Owner do projeto (RFC 1123 DNS label)")
	initCmd.Flags().StringVarP(&lifecycle, "lifecycle", "l", "experimental",
		"Lifecycle do projeto (experimental, production, deprecated)")
	initCmd.Flags().StringVar(&ciPlatform, "ci", "github",
		"Plataforma de CI (github, gitlab, none)")
	rootCmd.AddCommand(initCmd)
}

// Templates embarcados — sempre atualizados com o build.

var catalogInfoTpl = strings.TrimSpace(`
"@context": "https://foral-project.github.io/protocol/context/v1/catalog.jsonld"
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: {{.Name}}
  description: "TODO: Adicione uma descrição."
  annotations:
    foral.dev/archetype: {{.Archetype}}
  tags: []
spec:
  type: {{.Type}}
  lifecycle: {{.Lifecycle}}
  owner: {{.Owner}}
`) + "\n"

var githubWorkflowTpl = strings.TrimSpace(`
# Foral Validation — gerado por 'foral init'
# Consome reusable workflows do foral-project/governance.
name: Foral Validation

on: [push, pull_request]

jobs:
  catalog:
    uses: foral-project/governance/.github/workflows/validate-catalog.yml@main

  naming:
    uses: foral-project/governance/.github/workflows/validate-naming.yml@main

  commits:
    uses: foral-project/governance/.github/workflows/validate-conventional.yml@main
`) + "\n"

var gitignoreTpl = strings.TrimSpace(`
# OS
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/

# AI
.agents/
.gemini/
.claude/

# Env
.env
`) + "\n"
