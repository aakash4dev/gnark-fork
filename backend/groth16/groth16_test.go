package groth16_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/aakash4dev/gnark2/backend"
	"github.com/aakash4dev/gnark2/backend/groth16"
	"github.com/aakash4dev/gnark2/constraint"
	"github.com/aakash4dev/gnark2/frontend"
	"github.com/aakash4dev/gnark2/frontend/cs/r1cs"
	"github.com/aakash4dev/gnark2/test"
	"github.com/consensys/gnark"
	"github.com/consensys/gnark-crypto/ecc"
)

func TestCustomHashToField(t *testing.T) {
	assert := test.NewAssert(t)
	assignment := &commitmentCircuit{X: 1}
	for _, curve := range getCurves() {
		assert.Run(func(assert *test.Assert) {
			ccs, err := frontend.Compile(curve.ScalarField(), r1cs.NewBuilder, &commitmentCircuit{})
			assert.NoError(err)
			pk, vk, err := groth16.Setup(ccs)
			assert.NoError(err)
			witness, err := frontend.NewWitness(assignment, curve.ScalarField())
			assert.NoError(err)
			assert.Run(func(assert *test.Assert) {
				proof, err := groth16.Prove(ccs, pk, witness, backend.WithProverHashToFieldFunction(constantHash{}))
				assert.NoError(err)
				pubWitness, err := witness.Public()
				assert.NoError(err)
				err = groth16.Verify(proof, vk, pubWitness, backend.WithVerifierHashToFieldFunction(constantHash{}))
				assert.NoError(err)
			}, "custom success")
			assert.Run(func(assert *test.Assert) {
				proof, err := groth16.Prove(ccs, pk, witness, backend.WithProverHashToFieldFunction(constantHash{}))
				assert.NoError(err)
				pubWitness, err := witness.Public()
				assert.NoError(err)
				err = groth16.Verify(proof, vk, pubWitness)
				assert.Error(err)
			}, "prover_only")
			assert.Run(func(assert *test.Assert) {
				proof, err := groth16.Prove(ccs, pk, witness)
				assert.Error(err)
				_ = proof
			}, "verifier_only")
		}, curve.String())
	}
}

//--------------------//
//     benches		  //
//--------------------//

func BenchmarkSetup(b *testing.B) {
	for _, curve := range getCurves() {
		b.Run(curve.String(), func(b *testing.B) {
			r1cs, _ := referenceCircuit(curve)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = groth16.Setup(r1cs)
			}
		})
	}
}

func BenchmarkProver(b *testing.B) {
	for _, curve := range getCurves() {
		b.Run(curve.String(), func(b *testing.B) {
			r1cs, _solution := referenceCircuit(curve)
			fullWitness, err := frontend.NewWitness(_solution, curve.ScalarField())
			if err != nil {
				b.Fatal(err)
			}
			pk, err := groth16.DummySetup(r1cs)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = groth16.Prove(r1cs, pk, fullWitness)
			}
		})
	}
}

func BenchmarkVerifier(b *testing.B) {
	for _, curve := range getCurves() {
		b.Run(curve.String(), func(b *testing.B) {
			r1cs, _solution := referenceCircuit(curve)
			fullWitness, err := frontend.NewWitness(_solution, curve.ScalarField())
			if err != nil {
				b.Fatal(err)
			}
			publicWitness, err := fullWitness.Public()
			if err != nil {
				b.Fatal(err)
			}

			pk, vk, err := groth16.Setup(r1cs)
			if err != nil {
				b.Fatal(err)
			}
			proof, err := groth16.Prove(r1cs, pk, fullWitness)
			if err != nil {
				panic(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = groth16.Verify(proof, vk, publicWitness)
			}
		})
	}
}

type refCircuit struct {
	nbConstraints int
	X             frontend.Variable
	Y             frontend.Variable `gnark:",public"`
}

func (circuit *refCircuit) Define(api frontend.API) error {
	for i := 0; i < circuit.nbConstraints; i++ {
		circuit.X = api.Mul(circuit.X, circuit.X)
	}
	api.AssertIsEqual(circuit.X, circuit.Y)
	return nil
}

func referenceCircuit(curve ecc.ID) (constraint.ConstraintSystem, frontend.Circuit) {
	const nbConstraints = 40000
	circuit := refCircuit{
		nbConstraints: nbConstraints,
	}
	r1cs, err := frontend.Compile(curve.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	var good refCircuit
	good.X = 2

	// compute expected Y
	expectedY := new(big.Int).SetUint64(2)
	exp := big.NewInt(1)
	exp.Lsh(exp, nbConstraints)
	expectedY.Exp(expectedY, exp, curve.ScalarField())

	good.Y = expectedY

	return r1cs, &good
}

type commitmentCircuit struct {
	X frontend.Variable
}

func (c *commitmentCircuit) Define(api frontend.API) error {
	cmt, err := api.(frontend.Committer).Commit(c.X)
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	api.AssertIsEqual(cmt, "0xaabbcc")
	return nil
}

type constantHash struct{}

func (h constantHash) Write(p []byte) (n int, err error) { return len(p), nil }
func (h constantHash) Sum(b []byte) []byte               { return []byte{0xaa, 0xbb, 0xcc} }
func (h constantHash) Reset()                            {}
func (h constantHash) Size() int                         { return 3 }
func (h constantHash) BlockSize() int                    { return 32 }

func getCurves() []ecc.ID {
	if testing.Short() {
		return []ecc.ID{ecc.BN254}
	}
	return gnark.Curves()
}
