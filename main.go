// Ponto de entrada do Foral CLI.
// Delega execução para o comando root (Cobra).
package main

import "github.com/foral-project/cli/cmd"

func main() {
	cmd.Execute()
}
