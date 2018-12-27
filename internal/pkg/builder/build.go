package builder

import (
	"fmt"
)

func build(args []string) error {
	fmt.Printf("Build an image, %s\n", args)
	//if len(args) == 0 {
	//	//return errors.New(fmt.Sprintf("Missing component name. Available components = %s", componentNames()))
	//}
	//componentId := args[0]
	return nil
}
