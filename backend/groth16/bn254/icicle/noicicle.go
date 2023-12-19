//go:build !icicle

package icicle_bn254

import (
	"fmt"

	"github.com/aakash4dev/gnark2/backend"
	groth16_bn254 "github.com/aakash4dev/gnark2/backend/groth16/bn254"
	"github.com/aakash4dev/gnark2/backend/witness"
	cs "github.com/aakash4dev/gnark2/constraint/bn254"
)

const HasIcicle = false

func Prove(r1cs *cs.R1CS, pk *ProvingKey, fullWitness witness.Witness, opts ...backend.ProverOption) (*groth16_bn254.Proof, error) {
	return nil, fmt.Errorf("icicle backend requested but program compiled without 'icicle' build tag")
}
