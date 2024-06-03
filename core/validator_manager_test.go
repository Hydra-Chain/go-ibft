package core

import (
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CalculateQuorum(t *testing.T) {
	t.Parallel()

	vm := &ValidatorManager{
		vpLock: &sync.RWMutex{},
	}

	cases := []struct {
		validatorsVotingPower map[string]*big.Int
		signers               map[string]struct{}
		hasQuorum             bool
	}{
		{
			// case total voting power 4
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(1),
				"B": big.NewInt(1),
				"C": big.NewInt(1),
				"D": big.NewInt(1),
			},
			// all 4 signed, has quorum (quorum is 4)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power 4
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(1),
				"B": big.NewInt(1),
				"C": big.NewInt(1),
				"D": big.NewInt(1),
			},
			// only two signed (quorum is 4)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 9
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(2),
				"C": big.NewInt(2),
				"D": big.NewInt(2),
				"E": big.NewInt(1),
			},
			// 8 signed (quorum should be 9)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 10
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(2),
				"C": big.NewInt(2),
				"D": big.NewInt(2),
				"E": big.NewInt(1),
				"F": big.NewInt(1),
			},
			// 7 signed (quorum should be 7)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"E": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 10
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(2),
				"C": big.NewInt(2),
				"D": big.NewInt(2),
				"E": big.NewInt(1),
				"F": big.NewInt(1),
			},
			// 6 signed (quorum should be 7)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 60
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(17),
				"B": big.NewInt(10),
				"C": big.NewInt(10),
				"D": big.NewInt(10),
				"E": big.NewInt(10),
				"F": big.NewInt(3),
			},
			// only 37 signed (quorum should be 37)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 60
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(16),
				"B": big.NewInt(10),
				"C": big.NewInt(10),
				"D": big.NewInt(10),
				"E": big.NewInt(10),
				"F": big.NewInt(3),
			},
			// only 36 signed (quorum should be 37)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 90
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(20),
				"B": big.NewInt(20),
				"C": big.NewInt(16),
				"D": big.NewInt(34),
			},
			// 3 signed with voting power of 56 (quorum should be 56)
			signers: map[string]struct{}{
				"A": {},
				"C": {},
				"D": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 90
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(20),
				"B": big.NewInt(20),
				"C": big.NewInt(15),
				"D": big.NewInt(35),
			},
			// 3 signed with voting power of 55 (quorum should be 56)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 2100
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(290),
				"B": big.NewInt(810),
				"C": big.NewInt(500),
				"D": big.NewInt(500),
			},
			// 3 signed with voting power of 1290 (quorum should be 1290)
			signers: map[string]struct{}{
				"A": {},
				"C": {},
				"D": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 2100
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(289),
				"B": big.NewInt(811),
				"C": big.NewInt(500),
				"D": big.NewInt(500),
			},
			// 3 signed with voting power of 1289 (quorum should be 1290)
			signers: map[string]struct{}{
				"A": {},
				"C": {},
				"D": {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 19,128,543
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(3914974),
				"B": big.NewInt(3914976),
				"C": big.NewInt(3914975),
				"D": big.NewInt(7383617),
			},
			// 3 signed with voting power of 11,744,926 (quorum should be 11,744,926)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 19,128,543
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(3914974),
				"B": big.NewInt(3914975),
				"C": big.NewInt(3914975),
				"D": big.NewInt(7383618),
			},
			// 3 signed with voting power of 11,744,925 (quorum should be 11,744,926)
			signers: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
			hasQuorum: false,
		},
	}

	for _, c := range cases {
		require.NoError(t, vm.setCurrentVotingPower(c.validatorsVotingPower))
		require.Equal(t, c.hasQuorum, vm.HasQuorum(c.signers))
	}
}

func Test_calculateRCMinQuorum(t *testing.T) {
	t.Parallel()

	cases := []struct {
		totalVotingPower *big.Int
		expected         *big.Int
	}{
		{
			totalVotingPower: big.NewInt(9),
			expected:         nil,
		},
		{
			totalVotingPower: big.NewInt(10),
			expected:         big.NewInt(3),
		},
		{
			totalVotingPower: big.NewInt(115),
			expected:         big.NewInt(34),
		},
		{
			totalVotingPower: big.NewInt(1085),
			expected:         big.NewInt(325),
		},
		{
			totalVotingPower: big.NewInt(12763),
			expected:         big.NewInt(3828),
		},
		{
			totalVotingPower: big.NewInt(999999999),
			expected:         big.NewInt(299999999),
		},
	}

	for _, c := range cases {
		require.Equal(t, c.expected, calculateRCMinQuorum(c.totalVotingPower))
	}
}

func Test_HasRoundChangeQuorum(t *testing.T) {
	t.Parallel()

	vm := &ValidatorManager{
		vpLock: &sync.RWMutex{},
	}

	cases := []struct {
		round                 uint64
		validatorsVotingPower map[string]*big.Int
		hasQuorum             bool
		signers               map[string]struct{}
	}{
		{
			round: 5,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(1),
			},
			signers:   map[string]struct{}{},
			hasQuorum: false,
		},
		{
			round: 5,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(1),
			},
			signers: map[string]struct{}{
				"A": {},
			},
			hasQuorum: true,
		},
		{
			round: 5,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(4),
				"B": big.NewInt(4),
			},
			signers: map[string]struct{}{
				"A": {},
			},
			hasQuorum: false,
		},
		{
			round: 6,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(4),
				"B": big.NewInt(4),
			},
			signers: map[string]struct{}{
				"A": {},
			},
			hasQuorum: false,
		},
		{
			round: 5,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(5),
				"B": big.NewInt(5),
			},
			signers: map[string]struct{}{
				"A": {},
			},
			hasQuorum: false,
		},
		{
			round: 6,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(5),
				"B": big.NewInt(5),
			},
			signers: map[string]struct{}{
				"A": {},
			},
			hasQuorum: true,
		},
		{
			round: 6,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(2),
				"B": big.NewInt(7),
				"C": big.NewInt(6),
				"D": big.NewInt(6),
			},
			signers: map[string]struct{}{
				"C": {},
			},
			hasQuorum: true,
		},
		{
			round: 6,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(6),
				"B": big.NewInt(7),
				"C": big.NewInt(6),
				"D": big.NewInt(6),
			},
			signers: map[string]struct{}{
				"C": {},
			},
			hasQuorum: false,
		},
		{
			round: 6,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(15783),
				"B": big.NewInt(11432),
				"C": big.NewInt(13242),
				"D": big.NewInt(14324),
				"E": big.NewInt(32141),
			},
			signers: map[string]struct{}{
				"B": {},
				"C": {},
			},
			hasQuorum: false,
		},
		{
			round: 6,
			validatorsVotingPower: map[string]*big.Int{
				"A": big.NewInt(15783),
				"B": big.NewInt(11432),
				"C": big.NewInt(13242),
				"D": big.NewInt(14324),
				"E": big.NewInt(32141),
			},
			signers: map[string]struct{}{
				"B": {},
				"C": {},
				"D": {},
			},
			hasQuorum: true,
		},
	}

	for _, c := range cases {
		require.NoError(t, vm.setCurrentVotingPower(c.validatorsVotingPower))
		require.Equal(t, c.hasQuorum, vm.HasRoundChangeQuorum(c.round, c.signers))
	}
}
