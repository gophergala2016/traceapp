package cha

import (
	"fmt"
	"strings"

	"github.com/extemporalgenome/slug"
	"github.com/gophergala2016/traceapp/node"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/ssa"
)

// GetNodes returns the D3-compatible nodes by cha analysis -
// notably, works without a main()
func GetNodes(prog *ssa.Program) ([]node.Node, error) {
	added := make(map[string]struct{})
	extraImports := make(map[string][]string)
	var nodes []node.Node

	callgraph := cha.CallGraph(prog)
	callgraph.DeleteSyntheticNodes()
	for f, bar := range callgraph.Nodes {
		fName := fmt.Sprintf("%+v", f)
		_, ok := added[fName]
		if !ok {
			for _, edge := range bar.Out {
				if extraImports[edge.Callee.Func.String()] == nil {
					extraImports[edge.Callee.Func.String()] = make([]string, 0)
				}
				extraImports[edge.Callee.Func.String()] = append(extraImports[edge.Callee.Func.String()], strings.TrimPrefix(slug.Slug(fName), "package-"))
			}

			added[fName] = struct{}{}
		}
	}

	added = make(map[string]struct{})
	for f, bar := range callgraph.Nodes {
		fName := fmt.Sprintf("%+v", f)
		_, ok := added[fName]
		if !ok {
			var imports []string
			for _, edge := range bar.In {
				imports = append(imports, strings.TrimPrefix(slug.Slug(edge.Caller.Func.String()), "package-"))
			}
			imports = append(imports, extraImports[fName]...)

			nodes = append(nodes, node.Node{Name: strings.TrimPrefix(slug.Slug(fName), "package-"), Imports: imports, Size: len(imports) * 100})

			added[fName] = struct{}{}
		}
	}

	return nodes, nil
}
