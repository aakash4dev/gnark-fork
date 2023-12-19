package groth16

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"testing"

	"github.com/aakash4dev/gnark2/backend/witness"
	"github.com/consensys/gnark-crypto/ecc"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tests adapted from https://github.com/esuwu/groth16-verifier-bls12381
func TestVerifyBellmanProof(t *testing.T) {
	for _, test := range []struct {
		vk     string
		proof  string
		inputs string
		ok     bool
	}{
		{
			"hwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bhwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bo5ViaDBdO7ZBxAhLSe5k/5TFQyF5Lv7KN2tLKnwgoWMqB16OL8WdbePIwTCuPtJNAFKoTZylLDbSf02kckMcZQDPF9iGh+JC99Pio74vDpwTEjUx5tQ99gNQwxULtztsqDRsPnEvKvLmsxHt8LQVBkEBm2PBJFY+OXf1MNW021viDBpR10mX4WQ6zrsGL5L0GY4cwf4tlbh+Obit+LnN/SQTnREf8fPpdKZ1sa/ui3pGi8lMT6io4D7Ujlwx2RdChwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bkBF+isfMf77HCEGsZANw0hSrO2FGg14Sl26xLAIohdaW8O7gEaag8JdVAZ3OVLd5Df1NkZBEr753Xb8WwaXsJjE7qxwINL1KdqA4+EiYW4edb7+a9bbBeOPtb67ZxmFqAAAAAoMkzUv+KG8WoXszZI5NNMrbMLBDYP/xHunVgSWcix/kBrGlNozv1uFr0cmYZiij3YqToYs+EZa3dl2ILHx7H1n+b+Bjky/td2QduHVtf5t/Z9sKCfr+vOn12zVvOVz/6wAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			"lvQLU/KqgFhsLkt/5C/scqs7nWR+eYtyPdWiLVBux9GblT4AhHYMdCgwQfSJcudvsgV6fXoK+DUSRgJ++Nqt+Wvb7GlYlHpxCysQhz26TTu8Nyo7zpmVPH92+UYmbvbQCSvX2BhWtvkfHmqDVjmSIQ4RUMfeveA1KZbSf999NE4qKK8Do+8oXcmTM4LZVmh1rlyqznIdFXPN7x3pD4E0gb6/y69xtWMChv9654FMg05bAdueKt9uA4BEcAbpkdHF", "LcMT3OOlkHLzJBKCKjjzzVMg+r+FVgd52LlhZPB4RFg=",
			true,
		},
		{
			"hwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bhwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bo5ViaDBdO7ZBxAhLSe5k/5TFQyF5Lv7KN2tLKnwgoWMqB16OL8WdbePIwTCuPtJNAFKoTZylLDbSf02kckMcZQDPF9iGh+JC99Pio74vDpwTEjUx5tQ99gNQwxULtztsqDRsPnEvKvLmsxHt8LQVBkEBm2PBJFY+OXf1MNW021viDBpR10mX4WQ6zrsGL5L0GY4cwf4tlbh+Obit+LnN/SQTnREf8fPpdKZ1sa/ui3pGi8lMT6io4D7Ujlwx2RdChwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bkBF+isfMf77HCEGsZANw0hSrO2FGg14Sl26xLAIohdaW8O7gEaag8JdVAZ3OVLd5Df1NkZBEr753Xb8WwaXsJjE7qxwINL1KdqA4+EiYW4edb7+a9bbBeOPtb67ZxmFqAAAAAoMkzUv+KG8WoXszZI5NNMrbMLBDYP/xHunVgSWcix/kBrGlNozv1uFr0cmYZiij3YqToYs+EZa3dl2ILHx7H1n+b+Bjky/td2QduHVtf5t/Z9sKCfr+vOn12zVvOVz/6wAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			"lvQLU/KqgFhsLkt/5C/scqs7nWR+eYtyPdWiLVBux9GblT4AhHYMdCgwQfSJcudvsgV6fXoK+DUSRgJ++Nqt+Wvb7GlYlHpxCysQhz26TTu8Nyo7zpmVPH92+UYmbvbQCSvX2BhWtvkfHmqDVjmSIQ4RUMfeveA1KZbSf999NE4qKK8Do+8oXcmTM4LZVmh1rlyqznIdFXPN7x3pD4E0gb6/y69xtWMChv9654FMg05bAdueKt9uA4BEcAbpkdHF", "cmzVCcRVnckw3QUPhmG4Bkppeg4K50oDQwQ9EH+Fq1s=",
			false,
		},
		{
			"hwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bhwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bo5ViaDBdO7ZBxAhLSe5k/5TFQyF5Lv7KN2tLKnwgoWMqB16OL8WdbePIwTCuPtJNAFKoTZylLDbSf02kckMcZQDPF9iGh+JC99Pio74vDpwTEjUx5tQ99gNQwxULtztsqDRsPnEvKvLmsxHt8LQVBkEBm2PBJFY+OXf1MNW021viDBpR10mX4WQ6zrsGL5L0GY4cwf4tlbh+Obit+LnN/SQTnREf8fPpdKZ1sa/ui3pGi8lMT6io4D7Ujlwx2RdChwk883gUlTKCyXYA6XWZa8H9/xKIYZaJ0xEs0M5hQOMxiGpxocuX/8maSDmeCk3bkBF+isfMf77HCEGsZANw0hSrO2FGg14Sl26xLAIohdaW8O7gEaag8JdVAZ3OVLd5Df1NkZBEr753Xb8WwaXsJjE7qxwINL1KdqA4+EiYW4edb7+a9bbBeOPtb67ZxmFqAAAAAoMkzUv+KG8WoXszZI5NNMrbMLBDYP/xHunVgSWcix/kBrGlNozv1uFr0cmYZiij3YqToYs+EZa3dl2ILHx7H1n+b+Bjky/td2QduHVtf5t/Z9sKCfr+vOn12zVvOVz/6wAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			"lvQLU/KqgFhsLkt/5C/scqs7nWR+eYtyPdWiLVBux9GblT4AhHYMdCgwQfSJcudvsgV6fXoK+DUSRgJ++Nqt+Wvb7GlYlHpxCysQhz26TTu8Nyo7zpmVPH92+UYmbvbQCSvX2BhWtvkfHmqDVjmSIQ4RUMfeveA1KZbSf999NE4qKK8Do+8oXcmTM4LZVmh1rlyqznIdFXPN7x3pD4E0gb6/y69xtWMChv9654FMg05bAdueKt9uA4BEcAbpkdHF", "cmzVCcRVnckw3QUPhmG4Bkppeg4K50oDQwQ9EH+Fq1s=",
			false,
		},
		{
			"kYYCAS8vM2T99GeCr4toQ+iQzvl5fI89mPrncYqx3C1d75BQbFk8LMtcnLWwntd6kYYCAS8vM2T99GeCr4toQ+iQzvl5fI89mPrncYqx3C1d75BQbFk8LMtcnLWwntd6knkzSwcsialcheg69eZYPK8EzKRVI5FrRHKi8rgB+R5jyPV70ejmYEx1neTmfYKODRmARr/ld6pZTzBWYDfrCkiS1QB+3q3M08OQgYcLzs/vjW4epetDCmk0K1CEGcWdh7yLzdqr7HHQNOpZI8mdj/7lR0IBqB9zvRfyTr+guUG22kZo4y2KINDp272xGglKEeTglTxyDUriZJNF/+T6F8w70MR/rV+flvuo6EJ0+HA+A2ZnBbTjOIl9wjisBV+0kYYCAS8vM2T99GeCr4toQ+iQzvl5fI89mPrncYqx3C1d75BQbFk8LMtcnLWwntd6jgld4oAppAOzvQ7eoIx2tbuuKVSdbJm65KDxl/T+boaYnjRm3omdETYnYRk3HAhrAeWpefX+dM/k7PrcheInnxHUyjzSzqlN03xYjg28kdda9FZJaVsQKqdEJ/St9ivXAAAAAZae/nTwyDn5u+4WkhZ76991cGB/ymyGpXziT0bwS86pRw/AcbpzXmzK+hq+kvrvpwAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			"sStVLdyxqInmv76iaNnRFB464lGq48iVeqYWSi2linE9DST0fTNhxSnvSXAoPpt8tFsanj5vPafC+ij/Fh98dOUlMbO42bf280pOZ4lm+zr63AWUpOOIugST+S6pq9zeB0OHp2NY8XFmriOEKhxeabhuV89ljqCDjlhXBeNZwM5zti4zg89Hd8TbKcw46jAsjIJe2Siw3Th7ELQQKR5ucX50f0GISmnOSceePPdvjbGJ8fSFOnSmSp8dK7uyehrU", "",
			true,
		},
		{
			"mY//hEITCBCZUJUN/wsOlw1iUSSOESL6PFSbN1abGK80t5jPNICNlPuSorio4mmWmY//hEITCBCZUJUN/wsOlw1iUSSOESL6PFSbN1abGK80t5jPNICNlPuSorio4mmWpf+4uOyv3gPZe54SYGM4pfhteqJpwFQxdlpwXWyYxMTNaSLDj8VtSn/EJaSu+P6nFmWsda3mTYUPYMZzWE4hMqpDgFPcJhw3prArMThDPbR3Hx7E6NRAAR0LqcrdtsbDqu2T0tto1rpnFILdvHL4PqEUfTmF2mkM+DKj7lKwvvZUbukqBwLrnnbdfyqZJryzGAMIa2JvMEMYszGsYyiPXZvYx6Luk54oWOlOrwEKrCY4NMPwch6DbFq6KpnNSQwOmY//hEITCBCZUJUN/wsOlw1iUSSOESL6PFSbN1abGK80t5jPNICNlPuSorio4mmWpgRYCz7wpjk57X+NGJmo85tYKc+TNa1rT4/DxG9v6SHkpXmmPeHhzIIW8MOdkFjxB5o6Qn8Fa0c6Tt6br2gzkrGr1eK5/+RiIgEzVhcRrqdY/p7PLmKXqawrEvIv9QZ3AAAAAoo8rTzcIp5QvF3USzv2Lz99z43CPVkjHB1ejzj/SjzKNa54GiDzHoCoAL0xKLjRSqeL98AF0V1+cRI8FwJjOcMgf0gDmjzwiv3ppbPZKqJR7Go+57k02670lfG6s1MM0AAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			"g53N8ecorvG2sDgNv8D7quVhKMIIpdP9Bqk/8gmV5cJ5Rhk9gKvb4F0ll8J/ZZJVqa27OyciJwx6lym6QpVK9q1ASrqio7rD5POMDGm64Iay/ixXXn+//F+uKgDXADj9AySri2J1j3qEkqqe3kxKthw94DzAfUBPncHfTPazVtE48AfzB1KWZA7Vf/x/3phYs4ckcP7ZrdVViJVLbUgFy543dpKfEH2MD30ZLLYRhw8SatRCyIJuTZcMlluEKG+d", "aZ8tqrOeEJKt4AMqiRF/WJhIKTDC0HeDTgiJVLZ8OEs=",
			true,
		},
		{
			"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",
			true,
		},
		{
			"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",
			true,
		},
		{
			"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",
			true,
		},
		{
			"lp7+dPDIOfm77haSFnvr33VwYH/KbIalfOJPRvBLzqlHD8BxunNebMr6Gr6S+u+nlp7+dPDIOfm77haSFnvr33VwYH/KbIalfOJPRvBLzqlHD8BxunNebMr6Gr6S+u+nh7yLzdqr7HHQNOpZI8mdj/7lR0IBqB9zvRfyTr+guUG22kZo4y2KINDp272xGglKEeTglTxyDUriZJNF/+T6F8w70MR/rV+flvuo6EJ0+HA+A2ZnBbTjOIl9wjisBV+0iISo2JdNY1vPXlpwhlL2fVpW/WlREkF0bKlBadDIbNJBgM4niJGuEZDru3wqrGueETKHPv7hQ8em+p6vQolp7c0iknjXrGnvlpf4QtUtpg3z/D+snWjRPbVqRgKXWtihlp7+dPDIOfm77haSFnvr33VwYH/KbIalfOJPRvBLzqlHD8BxunNebMr6Gr6S+u+nuIvPFaM6dt7HZEbkeMnXWwSINeYC/j3lqYnce8Jq+XkuF42stVNiooI+TuXECnFdFi9Ib25b9wtyz3H/oKg48He1ftntj5uIRCOBvzkFHGUF6Ty214v3JYvXJjdS4uS2AAAAAY3pKZWWKFdGlxJ2BDiL3xe3eMWst6pMdjbKaKDOB3maYh5JyFpFFYRlSQcwQy4ywQAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			"jiGBK+TGHfH8Oadexhdet7ExyIWibSmamWQvffZkyl3WnMoVbTQ3lOks4Mca3sU5qgcaLyQQ1FjFW4g6vtoMapZ43hTGKaWO7bQHsOCvdwHCdwJDulVH16cMTyS9F0BfBJxa88F+JKZc4qMTJjQhspmq755SrKhN9Jf+7uPUhgB4hJTSrmlOkTatgW+/HAf5kZKhv2oRK5p5kS4sU48oqlG1azhMtcHEXDQdcwf9ANel4Z9cb+MQyp2RzI/3hlIx", "",
			false,
		},
		{
			"lp7+dPDIOfm77haSFnvr33VwYH/KbIalfOJPRvBLzqlHD8BxunNebMr6Gr6S+u+nlp7+dPDIOfm77haSFnvr33VwYH/KbIalfOJPRvBLzqlHD8BxunNebMr6Gr6S+u+nh7yLzdqr7HHQNOpZI8mdj/7lR0IBqB9zvRfyTr+guUG22kZo4y2KINDp272xGglKEeTglTxyDUriZJNF/+T6F8w70MR/rV+flvuo6EJ0+HA+A2ZnBbTjOIl9wjisBV+0iISo2JdNY1vPXlpwhlL2fVpW/WlREkF0bKlBadDIbNJBgM4niJGuEZDru3wqrGueETKHPv7hQ8em+p6vQolp7c0iknjXrGnvlpf4QtUtpg3z/D+snWjRPbVqRgKXWtihlp7+dPDIOfm77haSFnvr33VwYH/KbIalfOJPRvBLzqlHD8BxunNebMr6Gr6S+u+nuIvPFaM6dt7HZEbkeMnXWwSINeYC/j3lqYnce8Jq+XkuF42stVNiooI+TuXECnFdFi9Ib25b9wtyz3H/oKg48He1ftntj5uIRCOBvzkFHGUF6Ty214v3JYvXJjdS4uS2AAAAAo3pKZWWKFdGlxJ2BDiL3xe3eMWst6pMdjbKaKDOB3maYh5JyFpFFYRlSQcwQy4ywY4hgSvkxh3x/DmnXsYXXrexMciFom0pmplkL332ZMpd1pzKFW00N5TpLODHGt7FOQAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			"hp1iMepdu0rKoBh0NXcw9F9hkiggDIkRNINq2rlvUypPiSmp8U8tDSMeG0YVSovFteecr3THhBJj0qNeEe9jA2Ci64fKG9WT1heMYzEAQKebOErYXYCm9d72n97mYn1XBq+g1Y730XEDv4BIDI1hBDntJcgcj/cSvcILB1+60axJvtyMyuizxUr1JUBUq9njtmJ9m8zK6QZLNqMiKh0f2jokQb5mVhu6v5guW3KIjwQc/oFK/l5ehKAOPKUUggNh", "c9BSUPtO0xjPxWVNkEMfXe7O4UZKpaH/nLIyQJj7iA4=",
			false,
		},
		{
			"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",
			false,
		},
		{
			"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",
			false,
		},
	} {
		// decode verifying key
		vk := NewVerifyingKey(ecc.BLS12_381)

		vkBytes, err := base64.StdEncoding.DecodeString(test.vk)
		require.NoError(t, err)

		_, err = vk.ReadFrom(bytes.NewReader(vkBytes))
		require.NoError(t, err)

		// decode proof
		proofBytes, err := base64.StdEncoding.DecodeString(test.proof)
		require.NoError(t, err)

		// pad with 0 bytes to account for commitment stuff
		proofBytes = append(proofBytes, make([]byte, bls12381.SizeOfG1AffineUncompressed+4)...)

		proof := NewProof(ecc.BLS12_381)
		_, err = proof.ReadFrom(bytes.NewReader(proofBytes))
		require.NoError(t, err)

		// decode inputs
		inputsBytes, err := base64.StdEncoding.DecodeString(test.inputs)
		require.NoError(t, err)

		// verify groth16 proof
		// we need to prepend the number of elements in the witness.
		// witness package expects [nbPublic nbSecret] followed by [n | elements];
		// note that n is redundant with nbPublic + nbSecret
		var buf bytes.Buffer
		_ = binary.Write(&buf, binary.BigEndian, uint32(len(inputsBytes)/(fr.Limbs*8)))
		_ = binary.Write(&buf, binary.BigEndian, uint32(0))
		_ = binary.Write(&buf, binary.BigEndian, uint32(len(inputsBytes)/(fr.Limbs*8)))
		buf.Write(inputsBytes)

		witness, err := witness.New(ecc.BLS12_381.ScalarField())
		require.NoError(t, err)

		err = witness.UnmarshalBinary(buf.Bytes())
		require.NoError(t, err)

		err = Verify(proof, vk, witness)
		if test.ok {
			assert.NoError(t, err)
		}
	}
}
