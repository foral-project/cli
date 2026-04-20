// Tipos e constantes compartilhados entre os commands do Foral CLI.
// Centraliza structs de domínio e regex reutilizáveis.
package cmd

import "regexp"

// Regex RFC 1123 DNS label — compartilhado por init, validate e status.
var dnsLabelRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

// CatalogInfo representa a estrutura mínima do catalog-info.yaml (Backstage).
type CatalogInfo struct {
	Context    string   `yaml:"@context"`
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Metadata do catalog-info.yaml.
type Metadata struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Annotations map[string]string `yaml:"annotations"`
	Tags        []string          `yaml:"tags"`
}

// Spec do catalog-info.yaml.
type Spec struct {
	Type      string `yaml:"type"`
	Lifecycle string `yaml:"lifecycle"`
	Owner     string `yaml:"owner"`
}
