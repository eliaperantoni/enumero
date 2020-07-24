package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"text/template"
	"unicode"
)

const enumTemplateStr = `// Code generated by ono/scripts/enumGen. DO NOT EDIT.
package {{.Package}}

import "errors"

{{$nameUpper := capUpper .Name}}
{{$nameLower := capLower .Name}}
{{$type := concat $nameUpper "T"}}
{{$err := concat "ErrInvalid" (capUpper .Name)}}
{{$panic := $nameLower | printf "panic(\"invalid %s\")"}}

var {{$err}} = errors.New("invalid value is not a {{$nameLower}}")

type {{$type}} int

// Private constants
const (
	_	{{$type}}	=	iota
{{range .Values}}	{{valueName .}}
{{end}})

{{$stringMap := concat $nameLower "String" "Map"}}
// Forward string lookup
var {{$stringMap}} = map[{{$type}}]string{
{{range .Values}}	{{valueName .}}:	"{{capUpper .}}", 
{{end}}}

{{$reverseStringMap := concat $nameLower "Reverse" "String" "Map"}}
// Reverse string lookup
var {{$reverseStringMap}} = map[string]{{$type}} {
{{range .Values}}	"{{capUpper .}}":	{{valueName .}}, 
{{end}}}

{{$fromString := concat $nameLower "FromString"}}
func {{$fromString}}(str string) ({{$type}}, error) {
	v, ok := {{$reverseStringMap}}[str]
	if !ok {
		return 0, {{$err}}
	}
	return v, nil
}

{{$mustFromString := concat "must" $nameUpper "FromString"}}
func {{$mustFromString}}(str string) {{$type}} {
	v, err := {{$fromString}}(str)
	if err != nil {
		{{$panic}}
	}
	return v
}

func (v {{$type}}) String() string {
	return {{$stringMap}}[v]
}

{{$ordinalMap := concat $nameLower "Ordinal" "Map"}}
// Forward ordinal lookup
var {{$ordinalMap}} = map[{{$type}}]int{
{{range $index, $element := .Values}}	{{valueName $element}}:	{{inc $index}}, 
{{end}}}

{{$reverseOrdinalMap := concat $nameLower "Reverse" "Ordinal" "Map"}}
// Reverse ordinal lookup
var {{$reverseOrdinalMap}} = map[int]{{$type}}{
{{range $index, $element := .Values}}	{{inc $index}}:	{{valueName $element}}, 
{{end}}}

{{$fromOrdinal := concat $nameLower "FromOrdinal"}}
func {{$fromOrdinal}}(ord int) ({{$type}}, error) {
	v, ok := {{$reverseOrdinalMap}}[ord]
	if !ok {
		return 0, {{$err}}
	}
	return v, nil
}

{{$mustFromOrdinal := concat "must" $nameUpper "FromOrdinal"}}
func {{$mustFromOrdinal}}(ord int) {{$type}} {
	v, err := {{$fromOrdinal}}(ord)
	if err != nil {
		{{$panic}}
	}
	return v
}

func (v {{$type}}) Ordinal() int {
	return {{$ordinalMap}}[v]
}

{{$count := concat $nameLower "Count"}}
var {{$count}} = {{len .Values}}

{{$values := concat $nameLower "Values"}}
var {{$values}} = []{{$type}} {
{{range .Values}}	{{valueName .}},
{{end}}}

{{range .Variants}}

var {{variantMapName .}} = map[{{$type}}]string{
{{range $index, $variantVal := .Values}}{{$originalVal := index $.Values $index}}	{{valueName $originalVal}}:	"{{$variantVal}}", 
{{end}}}

var {{variantReverseMapName .}} = map[string]{{$type}} {
{{range $index, $variantVal := .Values}}{{$originalVal := index $.Values $index}}	"{{$variantVal}}":	{{valueName $originalVal}}, 
{{end}}}

func {{fromVariantName .}}(str string) ({{$type}}, error) {
	v, ok := {{variantReverseMapName .}}[str]
	if !ok {
		return 0, {{$err}}
	}
	return v, nil
}

func {{mustFromVariantName .}}(str string) {{$type}} {
	v, err := {{fromVariantName .}}(str)
	if err != nil {
		{{$panic}}
	}
	return v
}

func (v {{$type}}) {{capUpper .Name}}() string {
	return {{variantMapName .}}[v]
}
{{end}}

{{$struct := $nameLower}}

type {{$struct}} struct {
	Count int

	Values	[]{{$type}}

{{range .Values}}	{{capUpper .}}	{{$type}}
{{end}}

	FromString	func(string) ({{$type}}, error)
	MustFromString	func(string) {{$type}}

	FromOrdinal	func(int) ({{$type}}, error)
	MustFromOrdinal	func(int) {{$type}}

{{range .Variants}}
	From{{capUpper .Name}}	func(string) ({{$type}}, error)
	MustFrom{{capUpper .Name}}	func(string) {{$type}}
{{end}}
}

{{$structObj := $nameUpper}}

var {{$structObj}} = {{$struct}}{
	Count: {{$count}},

	Values:	{{$values}},

{{range .Values}}	{{capUpper .}}:	{{valueName .}},
{{end}}	

	FromString: {{$fromString}},
	MustFromString: {{$mustFromString}},

	FromOrdinal: {{$fromOrdinal}},
	MustFromOrdinal: {{$mustFromOrdinal}},
	
{{range .Variants}}
	From{{capUpper .Name}}: {{fromVariantName .}},
	MustFrom{{capUpper .Name}}: {{mustFromVariantName .}},
{{end}}
}
`

var (
	output      = flag.String("output", "./%s_enumero.go", "output file path")
	packageName = flag.String("package", "enums", "name of the generated package")
	name        = flag.String("name", "Enum", "name of the enum")
	values      = flag.String("values", "", "values comma separated")
	variants    = &VariantsFlag{}
)

func dieOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Var(variants, "variant", "VariantName:Value1,Value2,...,VariantN")
	flag.Parse()

	data := struct {
		Package  string
		Name     string
		Output   string
		Values   []string
		Variants []Variant
	}{
		Package:  *packageName,
		Name:     *name,
		Output:   *output,
		Values:   strings.Split(*values, ","),
		Variants: *variants,
	}

	capLower := func(s string) string {
		return string(unicode.ToLower([]rune(s)[0])) + string([]rune(s[1:]))
	}

	capUpper := func(s string) string {
		return string(unicode.ToUpper([]rune(s)[0])) + string([]rune(s[1:]))
	}

	funcs := map[string]interface{}{
		"capLower": capLower,
		"capUpper": capUpper,
		"inc": func(v int) int {
			return v + 1
		},
		"concat": func(ss ...string) string {
			out := ""
			for _, s := range ss {
				out += s
			}
			return out
		},
		"valueName": func(val string) string {
			return fmt.Sprintf("%s_%s_v", capLower(data.Name), capLower(val))
		},
		"variantMapName": func(v Variant) string {
			return fmt.Sprintf("%s%sMap", capLower(data.Name), capUpper(v.Name))
		},
		"variantReverseMapName": func(v Variant) string {
			return fmt.Sprintf("%sReverse%sMap", capLower(data.Name), capUpper(v.Name))
		},
		"fromVariantName": func(v Variant) string {
			return fmt.Sprintf("%sFrom%s", capLower(data.Name), capUpper(v.Name))
		},
		"mustFromVariantName": func(v Variant) string {
			return fmt.Sprintf("must%sFrom%s", capUpper(data.Name), capUpper(v.Name))
		},
	}
	pluginsTemplate := template.Must(template.New("Enum").Funcs(funcs).Parse(enumTemplateStr))

	fileName := fmt.Sprintf(data.Output, strings.ToLower(*name))

	outFile, err := os.Create(fileName)
	dieOnErr(err)

	defer outFile.Close()

	w := tabwriter.NewWriter(outFile, 4, 4, 4, ' ', 0)

	err = pluginsTemplate.Execute(w, data)
	dieOnErr(err)

	cmd := exec.Command("go", "fmt", fileName)

	err = cmd.Start()
	dieOnErr(err)

	err = cmd.Wait()
	dieOnErr(err)
}
