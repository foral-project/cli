// Comando `foral version` — exibe versão do CLI e do protocol.
package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Variáveis injetadas via ldflags no build.
// Ex: go build -ldflags "-X github.com/foral-project/cli/cmd.Version=0.1.0"
var (
	Version     = "dev"
	GitCommit   = "unknown"
	BuildDate   = "unknown"
	ProtocolURL = "https://foral-project.github.io/protocol"
)

type versionInfo struct {
	Version     string `json:"version"`
	GitCommit   string `json:"gitCommit"`
	BuildDate   string `json:"buildDate"`
	GoVersion   string `json:"goVersion"`
	Platform    string `json:"platform"`
	ProtocolURL string `json:"protocolUrl"`
}

var jsonOutput bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe versão do CLI, Go, e Protocol URL",
	Run: func(cmd *cobra.Command, args []string) {
		info := versionInfo{
			Version:     Version,
			GitCommit:   GitCommit,
			BuildDate:   BuildDate,
			GoVersion:   runtime.Version(),
			Platform:    runtime.GOOS + "/" + runtime.GOARCH,
			ProtocolURL: ProtocolURL,
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("foral %s\n", info.Version)
		fmt.Printf("  commit:   %s\n", info.GitCommit)
		fmt.Printf("  built:    %s\n", info.BuildDate)
		fmt.Printf("  go:       %s\n", info.GoVersion)
		fmt.Printf("  platform: %s\n", info.Platform)
		fmt.Printf("  protocol: %s\n", info.ProtocolURL)
	},
}

func init() {
	versionCmd.Flags().BoolVar(&jsonOutput, "json", false, "Saída em JSON")
	rootCmd.AddCommand(versionCmd)
}
