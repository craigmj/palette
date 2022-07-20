package palette

import (
	`os`
	`flag`
	`log`
	// `fmt`
	`strings`
	`html/template`

	colorful `github.com/lucasb-eyer/go-colorful`
	`gopkg.in/yaml.v3`
)

type Aspect struct {
	Name string
	Values []string
}

type Config struct {
	Prefix string  // String prefix for all our palette variables
	// Name of the generated CSS file
	Css string
	// 2N+1 colors in the palette
	N int
	// Colours to use to generate the 'big' palette
	Colors []string

	Aspects []*Aspect
	Areas map[string]map[string]string

	Variables map[string]string
}

type Var struct {
	Var string
	Value string
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

func Main() {
	configFile := flag.String(`config`,`palette.yml`,`palette config file`)
	flag.Parse()
	log.Printf(`Reading config from %s`, *configFile)
	cin, err := os.Open(*configFile)
	if nil!=err {
		log.Fatal(err)
	}
	defer cin.Close()

	var config Config
	if err := yaml.NewDecoder(cin).Decode(&config); nil!=err {
		log.Fatal(err)
	}
	config.Print()

	if 0==len(config.Colors) {
		log.Fatalf(`No colors read from config %s`, *configFile)
	}

	// GENERATE OUR PALETTE, WHICH ARE THE BASIC CSS VARIABLES WE MIGHT USE FOR COLOURING
	N := config.N
	palette := make([][]colorful.Color, 0)
	for line, text := range config.Colors {
		text := strings.TrimSpace(text)
		if ``==text {
			continue
		}
		color, err := colorful.Hex(text)
		if nil!=err {
			log.Fatalf(`Failed to parse color '%s' on line %d: %s`, text, line, err.Error())
		}
		// l,a,b := color.Lab()
		a,b,l := color.Hsl()

		colorSet := make([]colorful.Color, 0)
		// I want to scale l s.t. 1 => close to black and N+1=>given color and 2N+1=>close to white
		// I do this by doing darks first, and setting my 

		ldelta := l / float64(N+1)
		for i := 0; i<N; i++ {
			ls := ldelta*float64(i+1)
			log.Printf(`Color (%f,%f,%f), step %d = (%f,%f,%f)`, l,a,b,i,ls,a,b)
			colorSet = append(colorSet, colorful.Hsl(a,b,ls))
		}
		colorSet = append(colorSet,color)

		ldelta = (1-l)/float64(N+1)
		for i:=0; i<N; i++ {
			colorSet = append(colorSet, colorful.Hsl(a,b,l + ldelta*float64(i+1)))
		}
		palette = append(palette, colorSet)
	}

	// GENEATE OUR TREE OF ASPECTS
	cssTree := config.AspectsTree()
	// cssTree.Print()
	cssVars := cssTree.Flatten()

	htmlout, err := os.Create(`palette.html`)
	if nil!=err {
		log.Fatal(err)
	}
	defer htmlout.Close()
	if err := _t.Execute(htmlout, map[string]interface{}{
		`M`: 2*N+1,
		`Prefix`: config.Prefix,		
		`Palette`: palette,
		`Areas`: config.Areas,
	}); nil!=err {
		log.Fatal(err)
	}

	log.Printf(`Writing css to %s`, config.Css)
	out, err := os.Create(config.Css)
	if nil!=err {
		log.Fatalf(`Creating output file '%s': %s`, config.Css, err.Error())
	}
	defer out.Close()
	if err := _css.Execute(out, map[string]interface{}{
		`M`: 2*N+1,
		`Palette`: palette,
		`Prefix`: config.Prefix,
		`CssVars`: cssVars,
		`Variables`: config.Variables,
		`Areas`: config.Areas,
	}); nil!=err {
		log.Fatal(err)
	}
}

var _css = template.Must(template.New(``).Parse(`
:root,.palette1 {
{{ range $i,$p := .Palette}}{{range $j, $c := .}}
--{{$.Prefix}}-c{{$i}}h{{$j}}: {{$c.Hex}};
{{- end}}{{end}}
{{range $v := .CssVars}}{{range $var, $val := $.Variables}}
{{$v.Name $var}}: {{$v.Value $var $val}};
{{- end}}{{end}}
}
{{ range $area, $vars := .Areas }}
.{{$area}} {
{{ range $prop, $val := $vars}}  --{{$.Prefix}}-{{$prop}}: {{$val}};
{{end}}
}
{{end}}
`))

var _t = template.Must(template.New(``).Parse(`<!doctype html>
<html>
<head><title>Generated Template</title>
<style>
:root {
	--box-size: 10em;	
}
body {
	width: 100%;
	height: 100vh;
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
}
#palette {
	display: grid;
	
	grid-template-rows: auto;
	grid-template-columns: repeat({{.M}}, var(--box-size));
	color: white;
	font-family: sans-serif;
	font-size: 8pt;
}
#palette>div {
	height: var(--box-size);
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: center;
	text-align: center;
}
</style>
<body>
	<div id="palette">
{{ range $i,$p := .Palette}}{{range $j, $c := .}}
		<div style="background-color: {{$c.Hex}};">--{{$.Prefix}}-c{{$i}}h{{$j}}<br />{{$c.Hex}}</div>
{{- end}}{{end}}
	</div>
</body>
</html>`))