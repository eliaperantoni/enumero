# Enumero
A little tool to generate namespaced, type-checked and polyglot enums in Go.

## Installation
```shell
$ go get github.com/eliaperantoni/enumero
```

## Usage
```shell
$ enumero --help
Usage of enumero:
  -logrus
        use logrus.Fatal instead of panic
  -marshalJSON
        generate implementation for json.Marshaler
  -marshalText
        generate implementation for encoding.TextMarshaler
  -name string
        name of the enum (default "Enum")
  -output string
        output file path (default "./%s_enumero.go")
  -package string
        name of the generated package (default "enums")
  -source string
        source file to use for values and variants
  -unmarshalJSON
        generate implementation for json.Unmarshaler
  -unmarshalText
        generate implementation for encoding.TextUnmarshaler
  -values string
        values comma separated
  -variant value
        VariantName:Value1,Value2,...,VariantN
```

```shell
$ enumero -package=main -name Color -values=Red,Green,Blue -variant=German:Rot,Grün,Blau
```

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	var color ColorT
	color = Color.Blue

	fmt.Println(color.String())
	// Ordinal number assigned to enum value, starting from 1
	fmt.Println(color.Ordinal())
	fmt.Println(color.German())

	// Iterate enum values
	for i, color := range Color.Values {
		fmt.Printf("%d: [%d] %s %s\n", i, color.Ordinal(), color.String(), color.German())
	}

	// Returns non-nil error if string doesn't belong to any value
	if color, err := Color.FromString("Green"); err == nil {
		fmt.Println(color.String())
	} else {
		log.Fatal("invalid color")
	}

	// Returns no error, panics on failure
	color = Color.MustFromOrdinal(1)
	fmt.Println(color.German())

	color = Color.MustFromGerman("Blau")
	fmt.Println(color.String())
}

```

```shell
$ go run color_enumero.go test.go 
Blue
3
Blau
0: [1] Red Rot
1: [2] Green Grün
2: [3] Blue Blau
Green
Rot
Blue
```

### Source
For very big enums, a source file can be provided instead of listing all values and variants.

It is expected to look like this:

```
,German
Red,Rot
Green,Grün
Blue,Blau
```

That is: the first line is the header and contains an empty string (marking the column of primary values) followed by
the names of the variants.

Each following line contains a single value and its corresponding variants.
