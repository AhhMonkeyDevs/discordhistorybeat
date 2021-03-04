package main

import (
	"os"

	"github.com/AhhMonkeyDevs/discordhistorybeat/cmd"

	_ "github.com/AhhMonkeyDevs/discordhistorybeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
