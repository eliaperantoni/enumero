package main

import (
	"fmt"
	"strings"
)

type Variant struct {
	Name   string
	Values []string
}

func (v Variant) String() string {
	return fmt.Sprintf("{name: %s, values: [%s]}", v.Name, strings.Join(v.Values, ", "))
}

type VariantsFlag []Variant

func (variants *VariantsFlag) String() string {
	var joint string
	for i, variant := range *variants {
		joint += variant.String()
		if i < len(*variants)-1 {
			joint += ", "
		}
	}
	return "[" + joint + "]"
}

func (variants *VariantsFlag) Set(s string) error {
	var variant Variant
	nameVars := strings.Split(s, ":")
	variant.Name = nameVars[0]
	variant.Values = strings.Split(nameVars[1], ",")
	*variants = append(*variants, variant)
	return nil
}
