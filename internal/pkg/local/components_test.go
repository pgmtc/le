package local

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

func Test_componentNames(t *testing.T) {
	components := getComponents()
	componentNames := componentNames()
	for _, cmp := range components {
		if !common.ArrContains(componentNames, cmp.name) {
			t.Errorf("Component %s had not been found in list of component names", cmp.name)
		}
	}
}

func Test_componentMap(t *testing.T) {
	components := getComponents()
	componentMap := componentMap()
	for _, cmp := range components {
		if _, ok := componentMap[cmp.name]; !ok {
			t.Errorf("Component %s had not been found in the map", cmp.name)
		}
	}

}

func Test_getComponents(t *testing.T) {
	components := getComponents()
	if len(components) == 0 {
		t.Errorf("Expected to get list of components, get empty array")
	}

	for _, cmp := range components {
		// Test mandatory fields
		if cmp.name == "" || cmp.dockerId == "" || cmp.image == "" {
			t.Errorf("Component name, dockerId or image is empty for (%s, %s, %s)", cmp.name, cmp.dockerId, cmp.image)
		}
	}
}
