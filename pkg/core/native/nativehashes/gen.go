//go:build ignore

// This program generates hashes.go. It can be invoked by running
// go generate.
package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/nspcc-dev/neo-go/pkg/core/native/nativenames"
	"github.com/nspcc-dev/neo-go/pkg/core/state"
)

// srcTmpl is a nativehashes package template.
const srcTmpl = `// Code generated by "go generate go run gen.go"; DO NOT EDIT.

//go:generate go run gen.go

// package nativehashes contains hashes of all native contracts in their LE and Uint160 representation.
package nativehashes

import "github.com/nspcc-dev/neo-go/pkg/util"

// Hashes of all native contracts.
var (
{{- range .Natives }}
	// {{ .Name }} is a hash of native {{ .Name }} contract.
	{{ .Name }} = {{ .Hash }}
{{- end }}
)
`

type (
	// Config contains parameters for the nativehashes package generation.
	Config struct {
		Natives []NativeInfo
	}

	// NativeInfo contains information about native contract needed for
	// nativehashes package generation.
	NativeInfo struct {
		Name string
		Hash string
	}
)

// srcTemplate is a parsed nativehashes package template.
var srcTemplate = template.Must(template.New("nativehashes").Parse(srcTmpl))

func main() {
	f, err := os.Create("hashes.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg = Config{Natives: make([]NativeInfo, len(nativenames.All))}
	for i, name := range nativenames.All {
		var hash = state.CreateNativeContractHash(name)
		cfg.Natives[i] = NativeInfo{
			Name: name,
			Hash: fmt.Sprintf("%#v", hash),
		}
	}

	srcTemplate.Execute(f, cfg)
}
