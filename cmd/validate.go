// Comando `foral validate` — valida o repo atual contra o Foral Protocol.
// Executa: JSON Schema, OPA policies, naming conventions.
package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	schemaOnly bool
	policyOnly bool
	namingOnly bool
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Valida catalog-info.yaml contra o Foral Protocol",
	Long: `Executa validação offline contra o Foral Protocol:
  1. Campos obrigatórios (schema)
  2. Valores válidos (enums)
  3. Naming conventions (RFC 1123)
  4. Tags (kebab-case)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file := "catalog-info.yaml"
		if len(args) > 0 {
			file = args[0]
		}
		return runValidate(file, schemaOnly, policyOnly, namingOnly)
	},
}

// runValidate executa a lógica de validação. Exportada para testes.
func runValidate(file string, schema, policy, naming bool) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("❌ Não foi possível ler %s: %w", file, err)
	}

	var catalog CatalogInfo
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return fmt.Errorf("❌ Erro ao parsear YAML: %w", err)
	}

	fmt.Printf("Validando %s...\n\n", file)

	errors := 0
	passed := 0

	// Determinar quais checks rodar (nenhum flag = todos)
	runAll := !schema && !policy && !naming

	// --- Schema (campos obrigatórios) ---
	if runAll || schema {
		checks := []struct {
			name string
			ok   bool
			msg  string
		}{
			{"@context", catalog.Context != "", "campo '@context' é obrigatório"},
			{"apiVersion", catalog.APIVersion != "", "campo 'apiVersion' é obrigatório"},
			{"kind", catalog.Kind != "", "campo 'kind' é obrigatório"},
			{"metadata.name", catalog.Metadata.Name != "", "campo 'metadata.name' é obrigatório"},
			{"spec.type", catalog.Spec.Type != "", "campo 'spec.type' é obrigatório"},
			{"spec.lifecycle", catalog.Spec.Lifecycle != "", "campo 'spec.lifecycle' é obrigatório"},
			{"spec.owner", catalog.Spec.Owner != "", "campo 'spec.owner' é obrigatório"},
		}

		for _, c := range checks {
			if c.ok {
				fmt.Printf("  ✅ %-20s presente\n", c.name)
				passed++
			} else {
				fmt.Printf("  ❌ %-20s %s\n", c.name, c.msg)
				errors++
			}
		}
	}

	// --- Policy (valores válidos) ---
	if runAll || policy {
		fmt.Println()

		validAPIVersions := map[string]bool{
			"backstage.io/v1alpha1": true,
			"backstage.io/v1beta1":  true,
		}
		if validAPIVersions[catalog.APIVersion] {
			fmt.Printf("  ✅ %-20s '%s' válido\n", "apiVersion", catalog.APIVersion)
			passed++
		} else if catalog.APIVersion != "" {
			fmt.Printf("  ❌ %-20s '%s' inválido\n", "apiVersion", catalog.APIVersion)
			errors++
		}

		validKinds := map[string]bool{
			"Component": true, "API": true, "Resource": true,
			"System": true, "Domain": true, "Group": true,
			"User": true, "Template": true, "Location": true,
		}
		if validKinds[catalog.Kind] {
			fmt.Printf("  ✅ %-20s '%s' válido\n", "kind", catalog.Kind)
			passed++
		} else if catalog.Kind != "" {
			fmt.Printf("  ❌ %-20s '%s' inválido\n", "kind", catalog.Kind)
			errors++
		}

		validLifecycles := map[string]bool{
			"experimental": true, "production": true, "deprecated": true,
		}
		if validLifecycles[catalog.Spec.Lifecycle] {
			fmt.Printf("  ✅ %-20s '%s' válido\n", "spec.lifecycle", catalog.Spec.Lifecycle)
			passed++
		} else if catalog.Spec.Lifecycle != "" {
			fmt.Printf("  ❌ %-20s '%s' inválido\n", "spec.lifecycle", catalog.Spec.Lifecycle)
			errors++
		}
	}

	// --- Naming (RFC 1123) ---
	if runAll || naming {
		fmt.Println()
		dns := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)
		kebab := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

		if dns.MatchString(catalog.Metadata.Name) {
			fmt.Printf("  ✅ %-20s '%s' é DNS label válido (RFC 1123)\n", "metadata.name", catalog.Metadata.Name)
			passed++
		} else if catalog.Metadata.Name != "" {
			fmt.Printf("  ❌ %-20s '%s' não é DNS label válido\n", "metadata.name", catalog.Metadata.Name)
			errors++
		}

		if dns.MatchString(catalog.Spec.Owner) {
			fmt.Printf("  ✅ %-20s '%s' é DNS label válido\n", "spec.owner", catalog.Spec.Owner)
			passed++
		} else if catalog.Spec.Owner != "" {
			fmt.Printf("  ❌ %-20s '%s' não é DNS label válido\n", "spec.owner", catalog.Spec.Owner)
			errors++
		}

		for _, tag := range catalog.Metadata.Tags {
			if kebab.MatchString(tag) {
				fmt.Printf("  ✅ %-20s '%s' é kebab-case\n", "tag", tag)
				passed++
			} else {
				fmt.Printf("  ❌ %-20s '%s' não é kebab-case\n", "tag", tag)
				errors++
			}
		}
	}

	fmt.Println()
	fmt.Printf("Resultado: %d passed, %d errors\n", passed, errors)

	if errors > 0 {
		return fmt.Errorf("validação falhou com %d erros", errors)
	}

	fmt.Println("✅ Todas as validações passaram.")
	return nil
}

func init() {
	validateCmd.Flags().BoolVar(&schemaOnly, "schema", false, "Apenas validação de schema")
	validateCmd.Flags().BoolVar(&policyOnly, "policy", false, "Apenas validação de policies")
	validateCmd.Flags().BoolVar(&namingOnly, "naming", false, "Apenas validação de naming")
	rootCmd.AddCommand(validateCmd)
}
