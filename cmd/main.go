package main

import (
	"fmt"
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
	makefile, err := makego.ParseMakefile("testdata/makefile")
	if err != nil {
		panic(err)
	}
	PrintRules(makefile.Stages)
}
