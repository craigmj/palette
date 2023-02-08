package palette

import (
	`os`

	`gopkg.in/yaml.v3`
)

type Config struct {
	Prefix string  // String prefix for all our palette variables
	// Name of the generated CSS file
	Css string
	// 2N+1 colors in the palette
	N int
	// Colours to use to generate the 'big' palette
	Colors []string

	Aspects []*Aspect
	Areas map[string]interface{}

	Variables map[string]string
}


// NewAspectProperty creates a new Property as a descendant of the parent property
func (c *Config) NewAspectProperty(parent *AspectProperty, aspectValue string, nextAspectIndex int) *AspectProperty {
	ap := &AspectProperty {
		Prefix: c.Prefix,
		Aspect: aspectValue,
		Parent: parent,
		Children: []*AspectProperty{},
	}
	if nextAspectIndex == len(c.Aspects) {
		return ap
	}
	// Step through each of the values of the next child aspect, building the tree of each aspect level, with
	// every sub-branch a sub-branch
	for _, child := range c.Aspects[nextAspectIndex].Values {
		ap.Children = append(ap.Children, c.NewAspectProperty(ap, child, nextAspectIndex+1))
	}
	return ap
}

// func (c *Config) ScanForVars() ([]string, error) {
// 	m := map[string]bool{}
// 	var pattern string
// 	if ``!=c.Prefix {
// 		pattern = fmt.Sprintf(`var\(--(%s-[a-zA-Z_0-9]+)`)
// 	} else {
// 		pattern = `var\(--([a-zA-Z_0-9]+)`
// 	}
// 	reg, err := regexp.Compile(pattern)
// 	if nil!=err {
// 		return nil, fmt.Errorf(`Failed to compile search pattern '%s': %w`, pattern, err)
// 	}
	
// }

func (c *Config) AspectsTree() *AspectProperty {
	return c.NewAspectProperty(nil, ``, 0)
}

func (c *Config) Print() {
	yaml.NewEncoder(os.Stderr).Encode(c)
}

