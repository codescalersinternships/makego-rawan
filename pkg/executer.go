package makego

import "fmt"

func ExecuteMakefile(path string) error {
	makefile, err := parseMakefile(path)
	if err != nil {
		return err
	}

	ordered, err := dependencyResolver(makefile.Stages)
	if err != nil {
		return err
	}
	fmt.Println("Execution order:", ordered)

	return nil

}
