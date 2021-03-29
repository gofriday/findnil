package findnil

import (
	"fmt"
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = "findnil is ..."

var errtyp = types.Universe.Lookup("error").Type()

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "findnil",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	analysisutil.InspectFuncs(s.SrcFuncs, func(i int, instr ssa.Instruction) bool {
		sig := instr.Parent().Signature
		res := sig.Results()
		if res.Len() != 1 {
			return false
		}

		v := res.At(0)
		if !types.Identical(v.Type(), errtyp) {
			return false
		}

		_, ok := instr.(*ssa.Return)
		if !ok {
			return true
		}

		fmt.Printf("%T\n", instr)

		return true
	})

	return nil, nil
}
