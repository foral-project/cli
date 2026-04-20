// Comando root do Foral CLI.
// Define flags globais e inicialização.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "foral",
	Short: "Foral — Federated governance from the command line",
	Long: `Foral é um framework de governança federada inspirado nos forais
medievais portugueses. Permite governar repositórios autônomos
dentro de limites constitucionais definidos pelo Foral Protocol.

Documentação: https://github.com/foral-project
Especificação: https://github.com/foral-project/protocol/blob/main/PROTOCOL.md`,
}

// Execute é chamado pelo main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
