// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var typesTmpl = `
{{define "GoTypePrefix" -}}
	{{if eq .MaxOccurs "unbounded"}}[]{{end}}{{if or (eq .MinOccurs "0") (.Nillable)}}*{{end}}
{{- end}}

{{define "SimpleType"}}
	{{$type := replaceReservedWords .Name | makePublic}}
	{{if .Doc}} {{.Doc | comment}} {{end}}
	{{if ne .List.ItemType ""}}
		type {{$type}} []{{toGoType .List.ItemType }}
	{{else if ne .Union.MemberTypes ""}}
		type {{$type}} string
	{{else if .Union.SimpleType}}
		type {{$type}} string
	{{else if .Restriction.Base}}
		type {{$type}} {{toGoType .Restriction.Base}}
    {{else}}
		type {{$type}} interface{}
	{{end}}

	{{if .Restriction.Enumeration}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = "{{goString .Value}}" {{end}}
		{{end}}
	)
	{{end}}
{{end}}

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base}}
	{{ if $baseType }}
		{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.Sequence}}
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "Attributes"}}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{ if ne .Type "" }}
			{{ normalize .Name | makeFieldPublic}} {{toGoType .Type}} ` + "`" + `xml:"{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ else }}
			{{ normalize .Name | makeFieldPublic}} string ` + "`" + `xml:"{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ end }}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base}} ` + "`xml:\",chardata\" json:\"-,\"`" + `
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "ComplexTypeInline"}}
	{{replaceReservedWords .Name | makePublic}} {{template "GoTypePrefix" .}}struct {
	{{with .ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
		{{end}}
	{{end}}
	} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{template "GoTypePrefix" .}}{{.Ref | toGoType | removePointerFromType}} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty" json:"{{.Ref | removeNS}},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ normalize .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{else}}
					{{ normalize .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{end}}
			{{else}}
				{{template "ComplexTypeInline" .}}
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceAttrReservedWords .Name | makeFieldPublic}} {{template "GoTypePrefix" .}}{{.Type | toGoType | removePointerFromType}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}

{{define "Any"}}
	{{range .}}
		Items     []string ` + "`" + `xml:",any" json:"items,omitempty"` + "`" + `
	{{end}}
{{end}}

{{range .Schemas}}
	{{ $targetNamespace := .TargetNamespace }}

	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}

	{{range .Elements}}
		{{$name := .Name}}
		{{if not .Type}}
			{{if .ComplexType}}
				{{/* ComplexTypeLocal */}}
				{{with .ComplexType}}
					type {{$name | replaceReservedWords | makePublic}} struct {
						XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + `
						{{if ne .ComplexContent.Extension.Base ""}}
							{{template "ComplexContent" .ComplexContent}}
						{{else if ne .SimpleContent.Extension.Base ""}}
							{{template "SimpleContent" .SimpleContent}}
						{{else}}
							{{template "Elements" .Sequence}}
							{{template "Any" .Any}}
							{{template "Elements" .Choice}}
							{{template "Elements" .SequenceChoice}}
							{{template "Elements" .All}}
							{{template "Attributes" .Attributes}}
						{{end}}
					}
				{{end}}
			{{else if .SimpleType}}
				{{/* SimpleTypeLocal */}}
				{{with .SimpleType}}
					{{$type := replaceReservedWords $name | makePublic}}
					{{if .Doc}} {{.Doc | comment}} {{end}}
					{{if ne .List.ItemType ""}}
						type {{$type}} []{{toGoType .List.ItemType }}
					{{else if ne .Union.MemberTypes ""}}
						type {{$type}} string
					{{else if .Union.SimpleType}}
						type {{$type}} string
					{{else if .Restriction.Base}}
						type {{$type}} {{toGoType .Restriction.Base}}
				    {{else}}
						type {{$type}} interface{}
					{{end}}

					{{if .Restriction.Enumeration}}
					const (
						{{with .Restriction}}
							{{range .Enumeration}}
								{{if .Doc}} {{.Doc | comment}} {{end}}
								{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = "{{goString .Value}}" {{end}}
						{{end}}
					)
					{{end}}
				{{end}}
			{{end}}
		{{else}}
			{{if ne ($name | replaceReservedWords | makePublic) (toGoType .Type | removePointerFromType)}}
				type {{$name | replaceReservedWords | makePublic}} {{toGoType .Type | removePointerFromType}}
			{{end}}
		{{end}}
	{{end}}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$name := replaceReservedWords .Name | makePublic}}
		{{if eq (toGoType .SimpleContent.Extension.Base) "string"}}
			type {{$name}} string
		{{else}}
			type {{$name}} struct {
				{{$typ := findNameByType .Name}}
				{{if ne $name $typ}}
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$typ}}\"`" + `
				{{end}}

				{{if ne .ComplexContent.Extension.Base ""}}
					{{template "ComplexContent" .ComplexContent}}
				{{else if ne .SimpleContent.Extension.Base ""}}
					{{template "SimpleContent" .SimpleContent}}
				{{else}}
					{{template "Elements" .Sequence}}
					{{template "Any" .Any}}
					{{template "Elements" .Choice}}
					{{template "Elements" .SequenceChoice}}
					{{template "Elements" .All}}
					{{template "Attributes" .Attributes}}
				{{end}}
			}
		{{end}}
	{{end}}
{{end}}
`
