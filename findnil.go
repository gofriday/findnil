package findnil

import (
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

func isConstNil(value ssa.Value) bool {
	cnst, _ := value.(*ssa.Const)
	if cnst == nil {
		return false
	}
	return cnst.IsNil()
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

		ret, ok := instr.(*ssa.Return)
		if !ok {
			return true
		}

		if len(ret.Results) != res.Len() {
			return false
		}

		if !isConstNil(ret.Results[0]) {
			return true
		}

		ifinstr := analysisutil.IfInstr(instr.Block().Preds[0])
		if ifinstr == nil {
			return true
		}

		thenBlock := ifinstr.Block().Succs[0]
		if instr.Block() != thenBlock {
			return true
		}

		binOp, _ := ifinstr.Cond.(*ssa.BinOp)
		if binOp == nil {
			return true
		}
		switch {
		case types.Identical(binOp.X.Type(), errtyp) && isConstNil(binOp.Y):
			pass.Reportf(instr.Pos(), "NG")
		case types.Identical(binOp.Y.Type(), errtyp) && isConstNil(binOp.X):
			pass.Reportf(instr.Pos(), "NG")
		}

		return true
	})

	return nil, nil
}
