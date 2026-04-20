// Comando `foral status` — exibe o status de compliance do repo atual.
// Output tabular semelhante ao `kubectl get pods`.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibe o status de compliance do repositório atual",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := os.Getwd()
		return runStatus(dir)
	},
}

// runStatus executa a lógica do status. Exportada para testes.
func runStatus(dir string) error {
	file := filepath.Join(dir, "catalog-info.yaml")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("❌ catalog-info.yaml não encontrado.")
		fmt.Println()
		fmt.Println("Execute 'foral init .' para criar um manifesto.")
		return fmt.Errorf("catalog-info.yaml não encontrado")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("erro ao ler %s: %w", file, err)
	}

	var catalog CatalogInfo
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return fmt.Errorf("erro ao parsear YAML: %w", err)
	}

	// Cabeçalho
	fmt.Println("┌─────────────────────────────────────────────────────┐")
	fmt.Println("│              Foral Compliance Status                │")
	fmt.Println("├─────────────────────────────────────────────────────┤")

	// Dados do manifesto
	printRow("Nome", catalog.Metadata.Name)
	printRow("Kind", catalog.Kind)
	printRow("Archetype", catalog.Metadata.Annotations["foral.dev/archetype"])
	printRow("Owner", catalog.Spec.Owner)
	printRow("Lifecycle", catalog.Spec.Lifecycle)
	printRow("Type", catalog.Spec.Type)
	printRow("API Version", catalog.APIVersion)

	fmt.Println("├─────────────────────────────────────────────────────┤")
	fmt.Println("│              Validation Checks                     │")
	fmt.Println("├─────────────────────────────────────────────────────┤")

	checks := []struct {
		name   string
		status bool
	}{
		{"@context presente", catalog.Context != ""},
		{"apiVersion válido", catalog.APIVersion == "backstage.io/v1alpha1" || catalog.APIVersion == "backstage.io/v1beta1"},
		{"kind válido", isValidKind(catalog.Kind)},
		{"metadata.name presente", catalog.Metadata.Name != ""},
		{"metadata.name RFC 1123", dnsLabelRegex.MatchString(catalog.Metadata.Name)},
		{"spec.type presente", catalog.Spec.Type != ""},
		{"spec.lifecycle válido", catalog.Spec.Lifecycle == "experimental" || catalog.Spec.Lifecycle == "production" || catalog.Spec.Lifecycle == "deprecated"},
		{"spec.owner presente", catalog.Spec.Owner != ""},
		{"spec.owner RFC 1123", dnsLabelRegex.MatchString(catalog.Spec.Owner)},
	}

	passed := 0
	for _, c := range checks {
		icon := "✅"
		if !c.status {
			icon = "❌"
		} else {
			passed++
		}
		fmt.Printf("│  %s  %-45s │\n", icon, c.name)
	}

	fmt.Println("├─────────────────────────────────────────────────────┤")
	score := float64(passed) / float64(len(checks)) * 100
	bar := renderBar(score, 30)
	fmt.Printf("│  Score: %s %3.0f%%  (%d/%d)        │\n", bar, score, passed, len(checks))
	fmt.Println("└─────────────────────────────────────────────────────┘")

	return nil
}

func printRow(label, value string) {
	if value == "" {
		value = "(não definido)"
	}
	fmt.Printf("│  %-15s %-35s │\n", label+":", value)
}

func isValidKind(kind string) bool {
	validKinds := map[string]bool{
		"Component": true, "API": true, "Resource": true,
		"System": true, "Domain": true, "Group": true,
		"User": true, "Template": true, "Location": true,
	}
	return validKinds[kind]
}

func renderBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
