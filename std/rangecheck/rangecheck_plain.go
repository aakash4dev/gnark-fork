package rangecheck

import (
	"github.com/aakash4dev/gnark2/frontend"
	"github.com/aakash4dev/gnark2/std/math/bits"
)

type plainChecker struct {
	api frontend.API
}

func (pl plainChecker) Check(v frontend.Variable, nbBits int) {
	bits.ToBinary(pl.api, v, bits.WithNbDigits(nbBits))
}
