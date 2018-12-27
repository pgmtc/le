package builder

import (
	"fmt"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

func build(component common.Component) error {
	fmt.Printf("Build an image, %s\n", component.Name)
	return nil
}
