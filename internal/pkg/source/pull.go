package source

import (
	"fmt"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

type pullAction struct{}

func (pullAction) Run(log common.Logger, config common.Configuration, args ...string) error {
	fmt.Println("Pull latest code")
	return nil
}
