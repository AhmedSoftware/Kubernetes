/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// pointuh is a trivial gengo/v2 program which consider its inputs, and emits
// to new packages the same types, except for structs, where all fields are
// pointers.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"k8s.io/gengo/v2/args"
	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/types"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	stdArgs, myArgs := getArgs()

	// Collect and parse flags.
	stdArgs.AddFlags(pflag.CommandLine)
	myArgs.AddFlags(pflag.CommandLine)
	flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := validateArgs(stdArgs); err != nil {
		klog.ErrorS(err, "fatal error")
		os.Exit(1)
	}

	// Run the tool.
	if err := stdArgs.Execute(
		getNameSystems(),
		getDefaultNameSystem(),
		getPackages,
	); err != nil {
		klog.ErrorS(err, "fatal error")
		os.Exit(1)
	}
	klog.V(2).InfoS("Completed successfully.")
}

// toolArgs is used by the gengo framework to pass args specific to this generator.
type toolArgs struct {
}

// getArgs returns default arguments for the generator.
func getArgs() (*args.GeneratorArgs, *toolArgs) {
	stdArgs := args.Default().WithoutDefaultFlagParsing()
	stdArgs.OutputFileBaseName = "pointuh_generated"
	toolArgs := &toolArgs{}
	stdArgs.CustomArgs = toolArgs
	return stdArgs, toolArgs
}

// AddFlags add the generator flags to the flag set.
func (ta *toolArgs) AddFlags(fs *pflag.FlagSet) {
}

// validateArgs checks the given arguments.
func validateArgs(stdArgs *args.GeneratorArgs) error {
	if len(stdArgs.OutputPackagePath) == 0 {
		return fmt.Errorf("output package must be specified")
	}
	if len(stdArgs.OutputFileBaseName) == 0 {
		return fmt.Errorf("output file base name must be specified")
	}

	_ = stdArgs.CustomArgs.(*toolArgs)

	return nil
}

// getNameSystems returns the name system used by the generators in this package.
func getNameSystems() namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer("", nil),
	}
}

// getDefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func getDefaultNameSystem() string {
	return "public"
}

// getPackages is called after the inputs have been loaded.  It is expected to
// examine the provided context and return a list of Packages which will be
// executed further.
func getPackages(c *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	header := []byte(fmt.Sprintf("//go:build !%s\n// +build !%s\n\n", arguments.GeneratedBuildTag, arguments.GeneratedBuildTag))

	var pkgs generator.Packages
	for _, input := range c.Inputs {
		klog.V(2).InfoS("processing", "pkg", input)

		pkg := c.Universe[input]
		if pkg == nil { // e.g. the input had no Go files
			continue
		}

		pkgs = append(pkgs, &generator.DefaultPackage{
			PackageName: pkg.Name,
			PackagePath: filepath.Join(arguments.OutputPackagePath, pkg.Name),
			HeaderText:  header,
			// FilterFunc returns true if this Package cares about this type.
			// Each Generator has its own Filter method which will be checked
			// subsequently.  This will be called for every type in every
			// loaded package, not just things in our inputs.
			FilterFunc: func(c *generator.Context, t *types.Type) bool {
				// Only consider types in our inputs
				return t.Name.Package == pkg.Path
			},
			// GeneratorFunc returns a list of Generators, each of which is
			// responsible for a single output file (though multiple generators
			// may write to the same one).
			GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
				return []generator.Generator{
					newPointuhGenerator(arguments.OutputFileBaseName, pkg),
				}
			},
		})
	}

	return pkgs
}

// pointuhGenerator produces a file with autogenerated functions.
type pointuhGenerator struct {
	generator.DefaultGen
	myPackage *types.Package
}

func newPointuhGenerator(sanitizedName string, pkg *types.Package) generator.Generator {
	return &pointuhGenerator{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		myPackage: pkg,
	}
}

// Namers returns a set of NameSystems which will be merged with the namers
// provided when executing this package. In case of a name collision, the
// values produced here will win.
func (g *pointuhGenerator) Namers(*generator.Context) namer.NameSystems {
	return namer.NameSystems{
		// This elides package names when the name is in "this" package.
		"localraw": namer.NewRawNamer(g.myPackage.Path, nil),
	}
}

// GenerateType should emit code for the specified type.  This will be called
// for every type which made it through this Generator's Filter method.
func (g *pointuhGenerator) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	if namer.IsPrivateGoName(t.Name.Name) {
		return nil
	}

	klog.V(2).InfoS("generating pointerful type", "type", t.String())

	sw := generator.NewSnippetWriter(w, c, "$", "$")
	// Only modify structs.
	if t.Kind == types.Struct {
		emitModifiedStruct(t, sw)
	} else {
		emitUnmodifiedType(t, sw)
	}
	return sw.Error()
}

func emitUnmodifiedType(t *types.Type, sw *generator.SnippetWriter) {
	if t.Kind == types.DeclarationOf || t.Kind == types.Interface {
		return
	}

	args := argsFromType(t)
	sw.Do("// $.type|localraw$ is an autogenerated clone of $.type|raw$\n", args)
	sw.Do("type $.type|localraw$ ", args)
	for {
		if t.Kind != types.Pointer {
			break
		}
		sw.Do("*", nil)
		t = t.Elem
	}
	switch t.Kind {
	case types.Builtin:
		sw.Do("$.type.Name.Name$\n", args)
	case types.Map:
		sw.Do("map[$.type.Key$]$.type.Elem$\n", args)
	case types.Slice:
		sw.Do("[]$.type.Elem$\n", args)
	case types.Array:
		sw.Do("[$.type.Len$]$.type.Elem$\n", args)
	case types.Alias:
		sw.Do("$.type.Underlying|localraw$\n", args)
	case types.Struct:
		// must be non-exported
		sw.Do("struct {\n", args)
		sw.Do("}\n", nil)
	default:
		sw.Do("ERROR_Unhandled_input_type // $.type|raw$ ($.type.Kind$)\n", args)
	}
	sw.Do("\n", nil)
}

func emitModifiedStruct(t *types.Type, sw *generator.SnippetWriter) {
	args := argsFromType(t)
	sw.Do("// $.type|localraw$ is an autogenerated type.\n", args)
	sw.Do("type $.type|localraw$ struct {\n", args)
	for _, field := range t.Members {
		args := argsFromType(field.Type)
		if field.Type.Kind == types.Pointer {
			sw.Do(fmt.Sprintf("%s $.type|raw$\n", field.Name), args)
		} else {
			sw.Do(fmt.Sprintf("%s *$.type|raw$\n", field.Name), args)
		}
	}
	sw.Do("}\n", args)
}

func argsFromType(ts ...*types.Type) generator.Args {
	a := generator.Args{
		"type": ts[0],
	}
	for i, t := range ts {
		a[fmt.Sprintf("type%d", i+1)] = t
	}
	return a
}
