package common

import (
	"testing"
)

func Test_componentNames(t *testing.T) {
	components := GetComponents()
	componentNames := ComponentNames()
	for _, cmp := range components {
		if !ArrContains(componentNames, cmp.Name) {
			t.Errorf("Component %s had not been found in list of component names", cmp.Name)
		}
	}
}

func Test_componentMap(t *testing.T) {
	components := GetComponents()
	componentMap := ComponentMap()
	for _, cmp := range components {
		if _, ok := componentMap[cmp.Name]; !ok {
			t.Errorf("Component %s had not been found in the map", cmp.Name)
		}
	}

}

func Test_getComponents(t *testing.T) {
	components := GetComponents()
	if len(components) == 0 {
		t.Errorf("Expected to get list of components, get empty array")
	}

	for _, cmp := range components {
		// Test mandatory fields
		if cmp.Name == "" || cmp.DockerId == "" || cmp.Image == "" {
			t.Errorf("Component Name, DockerId or Image is empty for (%s, %s, %s)", cmp.Name, cmp.DockerId, cmp.Image)
		}
	}
}
