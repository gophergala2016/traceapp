package rta

import (
	"fmt"
	"strings"

	"github.com/extemporalgenome/slug"
	"github.com/gophergala2016/traceapp/node"
	"golang.org/x/tools/go/callgraph/rta"
	"golang.org/x/tools/go/ssa"
)

// GetNodes returns the D3-compatible nodes by rta analysis
func GetNodes(prog *ssa.Program) ([]node.Node, error) {
	main, err := mainPackage(prog, false)
	if err != nil {
		return nil, err
	}

	roots := []*ssa.Function{
		main.Func("init"),
		main.Func("main"),
	}

	added := make(map[string]struct{})
	extraImports := make(map[string][]string)
	var nodes []node.Node

	callgraph := rta.Analyze(roots, true)
	callgraph.CallGraph.DeleteSyntheticNodes()
	for _, bar := range callgraph.CallGraph.Nodes {
		_, ok := added[bar.Func.Package().String()]
		if !ok {
			for _, edge := range bar.Out {
				if extraImports[edge.Callee.Func.Package().String()] == nil {
					extraImports[edge.Callee.Func.Package().String()] = make([]string, 0)
				}
				extraImports[edge.Callee.Func.Package().String()] = append(extraImports[edge.Callee.Func.Package().String()], strings.TrimPrefix(slug.Slug(bar.Func.Package().String()), "package-"))
			}

			added[bar.Func.Package().String()] = struct{}{}
		}
	}

	added = make(map[string]struct{})
	for _, bar := range callgraph.CallGraph.Nodes {
		_, ok := added[bar.Func.Package().String()]
		if !ok {
			var imports []string
			for _, edge := range bar.In {
				imports = append(imports, strings.TrimPrefix(slug.Slug(edge.Caller.Func.Package().String()), "package-"))
			}
			imports = append(imports, extraImports[bar.Func.Package().String()]...)

			nodes = append(nodes, node.Node{Name: strings.TrimPrefix(slug.Slug(bar.Func.Package().String()), "package-"), Imports: imports, Size: len(imports) * 100})

			added[bar.Func.Package().String()] = struct{}{}
		}
	}

	return nodes, nil
}

// The resulting package has a main() function.
func mainPackage(prog *ssa.Program, tests bool) (*ssa.Package, error) {
	pkgs := prog.AllPackages()

	// TODO(adonovan): allow independent control over tests, mains and libraries.
	// TODO(adonovan): put this logic in a library; we keep reinventing it.

	if tests {
		// If -test, use all packages' tests.
		if len(pkgs) > 0 {
			if main := prog.CreateTestMainPackage(pkgs...); main != nil {
				return main, nil
			}
		}
		return nil, fmt.Errorf("no tests")
	}

	// Otherwise, use the first package named main.
	for _, pkg := range pkgs {
		if pkg.Pkg.Name() == "main" {
			if pkg.Func("main") == nil {
				return nil, fmt.Errorf("no func main() in main package")
			}
			return pkg, nil
		}
	}

	return nil, fmt.Errorf("no main package")
}
