package main

import (
	"os"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/raintank/snap-plugin-collector-procnum/procnum"
)

func main() {

	plugin.Start(
		procnum.Meta(),
		procnum.New(),
		os.Args[1],
	)
}
