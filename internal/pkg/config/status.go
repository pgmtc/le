package config

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

type statusAction struct{}

func (statusAction) Run(log common.Logger, args ...string) error {
	color.HiWhite("Current profile: %s", common.CONFIG.Profile)
	color.HiWhite("Available profiles: %s", common.GetAvailableProfiles())
	if len(args) > 0 && args[0] == "-v" {
		// Verbose output
		s, _ := json.MarshalIndent(common.GetComponents(), "", "  ")
		color.White("Components: \n%s\n", s)
	} else {
		color.White("Components: (for more verbose output, add '-v' parameter)")
		for i, cmp := range common.GetComponents() {
			color.White("   %02d | Name: %s, DockerId: %s, Image: %s", i, cmp.Name, cmp.DockerId, cmp.Image)
		}
	}
	return nil
}
