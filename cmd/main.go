package main

import (
	"fmt"
	"os"
	"strings"

	makego "github.com/codescalersinternships/makego-rawan/pkg"
)

func PrintRules(rules []makego.Stage) {
	for _, r := range rules {
		if len(r.Dependencies) > 0 {
			fmt.Printf("%s: %s\n", r.Target, strings.Join(r.Dependencies, " "))
		} else {
			fmt.Printf("%s:\n", r.Target)
		}

		for _, cmd := range r.Commands {
			fmt.Printf("\t%s\n", cmd)
		}

		fmt.Println()
	}
}
func main() {
	args := os.Args[1:]

	var targets []string
	targets = nil
	if len(args) > 0 {
		targets = args[0:]
	}

	err := makego.ExecuteMakefile("testdata/makefile", targets)
	if err != nil {
		panic(err)
	}
}
