package std

import (
	"sync"

	"github.com/aakash4dev/gnark-fork/constraint/solver"
	"github.com/aakash4dev/gnark-fork/std/algebra/emulated/sw_emulated"
	"github.com/aakash4dev/gnark-fork/std/algebra/native/sw_bls12377"
	"github.com/aakash4dev/gnark-fork/std/algebra/native/sw_bls24315"
	"github.com/aakash4dev/gnark-fork/std/evmprecompiles"
	"github.com/aakash4dev/gnark-fork/std/internal/logderivarg"
	"github.com/aakash4dev/gnark-fork/std/math/bits"
	"github.com/aakash4dev/gnark-fork/std/math/bitslice"
	"github.com/aakash4dev/gnark-fork/std/math/cmp"
	"github.com/aakash4dev/gnark-fork/std/math/emulated"
	"github.com/aakash4dev/gnark-fork/std/rangecheck"
	"github.com/aakash4dev/gnark-fork/std/selector"
)

var registerOnce sync.Once

// RegisterHints register all gnark/std hints
// In the case where the Solver/Prover code is loaded alongside the circuit, this is not useful.
// However, if a Solver/Prover services consumes serialized constraint systems, it has no way to
// know which hints were registered; caller code may add them through backend.WithHints(...).
func RegisterHints() {
	registerOnce.Do(registerHints)
}

func registerHints() {
	// note that importing these packages may already trigger a call to solver.RegisterHint(...)
	solver.RegisterHint(sw_bls24315.DecomposeScalarG1)
	solver.RegisterHint(sw_bls12377.DecomposeScalarG1)
	solver.RegisterHint(sw_bls24315.DecomposeScalarG2)
	solver.RegisterHint(sw_bls12377.DecomposeScalarG2)
	solver.RegisterHint(bits.GetHints()...)
	solver.RegisterHint(cmp.GetHints()...)
	solver.RegisterHint(selector.GetHints()...)
	solver.RegisterHint(emulated.GetHints()...)
	solver.RegisterHint(rangecheck.GetHints()...)
	solver.RegisterHint(evmprecompiles.GetHints()...)
	solver.RegisterHint(logderivarg.GetHints()...)
	solver.RegisterHint(bitslice.GetHints()...)
	solver.RegisterHint(sw_emulated.GetHints()...)
}
