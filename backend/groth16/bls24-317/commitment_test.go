// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gnark DO NOT EDIT

package groth16_test

import (
	"fmt"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/aakash4dev/gnark2/backend/groth16"
	"github.com/aakash4dev/gnark2/backend/witness"
	"github.com/aakash4dev/gnark2/constraint"
	"github.com/aakash4dev/gnark2/frontend"
	"github.com/aakash4dev/gnark2/frontend/cs/r1cs"
	"github.com/stretchr/testify/assert"
)

type singleSecretCommittedCircuit struct {
	One frontend.Variable
}

func (c *singleSecretCommittedCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(c.One, 1)
	commitCompiler, ok := api.Compiler().(frontend.Committer)
	if !ok {
		return fmt.Errorf("compiler does not commit")
	}
	commit, err := commitCompiler.Commit(c.One)
	if err != nil {
		return err
	}
	api.AssertIsDifferent(commit, 0)
	return nil
}

func setup(t *testing.T, circuit frontend.Circuit) (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey) {
	_r1cs, err := frontend.Compile(ecc.BLS24_317.ScalarField(), r1cs.NewBuilder, circuit)
	assert.NoError(t, err)

	pk, vk, err := groth16.Setup(_r1cs)
	assert.NoError(t, err)

	return _r1cs, pk, vk
}

func prove(t *testing.T, assignment frontend.Circuit, cs constraint.ConstraintSystem, pk groth16.ProvingKey) (witness.Witness, groth16.Proof) {
	_witness, err := frontend.NewWitness(assignment, ecc.BLS24_317.ScalarField())
	assert.NoError(t, err)

	proof, err := groth16.Prove(cs, pk, _witness)
	assert.NoError(t, err)

	public, err := _witness.Public()
	assert.NoError(t, err)
	return public, proof
}

func test(t *testing.T, circuit frontend.Circuit, assignment frontend.Circuit) {

	_r1cs, pk, vk := setup(t, circuit)

	public, proof := prove(t, assignment, _r1cs, pk)

	assert.NoError(t, groth16.Verify(proof, vk, public))
}

func TestSingleSecretCommitted(t *testing.T) {
	circuit := singleSecretCommittedCircuit{}
	assignment := singleSecretCommittedCircuit{One: 1}
	test(t, &circuit, &assignment)
}

type noCommitmentCircuit struct { // to see if unadulterated groth16 is still correct
	One frontend.Variable
}

func (c *noCommitmentCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(c.One, 1)
	return nil
}

func TestNoCommitmentCircuit(t *testing.T) {
	circuit := noCommitmentCircuit{}
	assignment := noCommitmentCircuit{One: 1}

	test(t, &circuit, &assignment)
}

// Just to see if the A,B,C values are computed correctly
type singleSecretFauxCommitmentCircuit struct {
	One        frontend.Variable `gnark:",public"`
	Commitment frontend.Variable `gnark:",public"`
}

func (c *singleSecretFauxCommitmentCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(c.One, 1)
	api.AssertIsDifferent(c.Commitment, 0)
	return nil
}

func TestSingleSecretFauxCommitmentCircuit(t *testing.T) {
	test(t, &singleSecretFauxCommitmentCircuit{}, &singleSecretFauxCommitmentCircuit{
		One:        1,
		Commitment: 2,
	})
}

type oneSecretOnePublicCommittedCircuit struct {
	One frontend.Variable
	Two frontend.Variable `gnark:",public"`
}

func (c *oneSecretOnePublicCommittedCircuit) Define(api frontend.API) error {
	commitCompiler, ok := api.Compiler().(frontend.Committer)
	if !ok {
		return fmt.Errorf("compiler does not commit")
	}
	commit, err := commitCompiler.Commit(c.One, c.Two)
	if err != nil {
		return err
	}

	// constrain vars
	api.AssertIsDifferent(commit, 0)
	api.AssertIsEqual(c.One, 1)
	api.AssertIsEqual(c.Two, 2)

	return nil
}

func TestOneSecretOnePublicCommitted(t *testing.T) {
	test(t, &oneSecretOnePublicCommittedCircuit{}, &oneSecretOnePublicCommittedCircuit{
		One: 1,
		Two: 2,
	})
}
