// Testes dos types compartilhados.
package cmd

import "testing"

func TestDNSLabelRegex(t *testing.T) {
	valid := []string{
		"my-project",
		"a",
		"test-123",
		"hello-world-app",
		"a1b2c3",
	}
	for _, v := range valid {
		if !dnsLabelRegex.MatchString(v) {
			t.Errorf("'%s' deveria ser DNS label válido", v)
		}
	}

	invalid := []string{
		"MyProject",
		"my_project",
		"-starts-with-dash",
		"ends-with-dash-",
		"has spaces",
		"has.dots",
		"UPPER",
		"",
		"a-very-long-name-that-exceeds-sixty-three-characters-in-total-and-should-fail",
	}
	for _, v := range invalid {
		if dnsLabelRegex.MatchString(v) {
			t.Errorf("'%s' NÃO deveria ser DNS label válido", v)
		}
	}
}

func TestIsValidKind(t *testing.T) {
	validKinds := []string{
		"Component", "API", "Resource", "System",
		"Domain", "Group", "User", "Template", "Location",
	}
	for _, k := range validKinds {
		if !isValidKind(k) {
			t.Errorf("'%s' deveria ser kind válido", k)
		}
	}

	invalidKinds := []string{"", "Invalid", "component", "SERVICE"}
	for _, k := range invalidKinds {
		if isValidKind(k) {
			t.Errorf("'%s' NÃO deveria ser kind válido", k)
		}
	}
}
