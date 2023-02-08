package palette

import (
	`strings`
	`os`

	`gopkg.in/yaml.v3`
)

// AspectProperty is a single propety that will appear in an Aspect, structured into a tree
// so that we're able to find a particular aspect, and trace through it's parent properties to find
// the eventual value of it.
type AspectProperty struct {
	Prefix string
	Aspect string
	Parent *AspectProperty `yaml:"-"`
	Children []*AspectProperty
}

// Flatten converts the AspectProperty tree into a flat Slice, but the parent links remain valid.
func  (ap *AspectProperty) Flatten() []*AspectProperty {
	f := []*AspectProperty{ap}
	for _, c := range ap.Children{
		f = append(f, c.Flatten()...)
	}
	return f
}

func (ap *AspectProperty) IsRoot() bool {
	return nil==ap || nil==ap.Parent
}

func (ap *AspectProperty) ParentIsRoot() bool {
	return ap.Parent.IsRoot()
}

// Path returns the path from this aspect to the parent
func (ap *AspectProperty) Path() string {
	if ap.ParentIsRoot() {
		return ap.Aspect
	}
	return ap.Parent.Path() + `-` + ap.Aspect
}

func (ap *AspectProperty) Name(variable string) string {
	parts := []string{}
	if ``!=ap.Prefix {
		parts = append(parts, ap.Prefix)
	}
	path := ap.Path()
	if ``!=path {
		parts = append(parts,path)
	}
	parts = append(parts, variable)
	return `--`+strings.Join(parts, `-`)
}

func (ap *AspectProperty) Value(variable, value string) string {
	// the value of a property is either a set value: which the user does
	// in their over-ride CSS, or a sequence of previous values
	parentNames := ap.ParentNames(variable)[1:]	// The first parent name is us	
	if 0<len(parentNames) {
		return `var(` + strings.Join(parentNames, `, var(`) + `, ` + value + strings.Repeat(`)`,len(parentNames))
	}
	return value
}

func (ap *AspectProperty) ParentNames(variable string) []string {
	if ap.IsRoot() {
		return []string{ ap.Name(variable) }
	}
	return append([]string{ap.Name(variable)}, ap.Parent.ParentNames(variable)...)
}

func (ap *AspectProperty) Print() {
	yaml.NewEncoder(os.Stdout).Encode(ap)
}