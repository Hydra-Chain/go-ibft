package core

import (
	"bytes"
	"errors"
	"math/big"
	"sync"

	"github.com/0xPolygon/go-ibft/messages/proto"
)

const (
	// rcMinQuorumThreshold is the threshold above which the rcMinQuorum is applied
	rcMinQuorumThreshold uint64 = 5
)

var (
	errVotingPowerNotCorrect = errors.New("total voting power is zero or less")
)

// ValidatorBackend defines interface that has GetVotingPower
type ValidatorBackend interface {
	// GetVotingPowers returns map of validators addresses and their voting powers for the specified height.
	GetVotingPowers(height uint64) (map[string]*big.Int, error)
}

// ValidatorManager keeps voting power and other information about validators
type ValidatorManager struct {
	vpLock *sync.RWMutex

	// quorumSize represents quorum for the height specified in the current View
	quorumSize *big.Int

	// rcMinQuorum represents the voting power needed to create a round change certificate
	// when round is above rcMinQuorumThreshold for the height specified in the current View
	rcMinQuorum *big.Int

	// validatorsVotingPower is a map of the validator addresses on their voting power for
	// the height specified in the current View
	validatorsVotingPower map[string]*big.Int

	backend ValidatorBackend

	log Logger
}

// NewValidatorManager creates new ValidatorManager
func NewValidatorManager(backend ValidatorBackend, log Logger) *ValidatorManager {
	return &ValidatorManager{
		quorumSize:            big.NewInt(0),
		rcMinQuorum:           big.NewInt(0),
		backend:               backend,
		validatorsVotingPower: nil,
		log:                   log,
		vpLock:                &sync.RWMutex{},
	}
}

// Init sets voting power and quorum size
func (vm *ValidatorManager) Init(height uint64) error {
	validatorsVotingPower, err := vm.backend.GetVotingPowers(height)
	if err != nil {
		return err
	}

	return vm.setCurrentVotingPower(validatorsVotingPower)
}

// setCurrentVotingPower sets the current total voting power and quorum size
// based on current validators voting power
func (vm *ValidatorManager) setCurrentVotingPower(validatorsVotingPower map[string]*big.Int) error {
	vm.vpLock.Lock()
	defer vm.vpLock.Unlock()

	totalVotingPower := calculateTotalVotingPower(validatorsVotingPower)
	if totalVotingPower.Cmp(big.NewInt(0)) <= 0 {
		return errVotingPowerNotCorrect
	}

	vm.validatorsVotingPower = validatorsVotingPower
	vm.quorumSize = calculateQuorum(totalVotingPower)
	vm.rcMinQuorum = calculateRCMinQuorum(totalVotingPower)

	return nil
}

// HasQuorum provides information on whether messages have reached the quorum
func (vm *ValidatorManager) HasQuorum(sendersAddrs map[string]struct{}) bool {
	vm.vpLock.RLock()
	defer vm.vpLock.RUnlock()

	// if not initialized correctly return false
	if vm.validatorsVotingPower == nil {
		return false
	}

	messageVotePower := big.NewInt(0)

	for from := range sendersAddrs {
		if vote, ok := vm.validatorsVotingPower[from]; ok {
			messageVotePower.Add(messageVotePower, vote)
		}
	}

	// aggVotingPower >= (2 * totalVotingPower / 3) + 1
	return messageVotePower.Cmp(vm.quorumSize) >= 0
}

// HasPrepareQuorum provides information on whether prepared messages have reached the quorum
func (vm *ValidatorManager) HasPrepareQuorum(stateName stateType, proposalMessage *proto.Message,
	msgs []*proto.Message) bool {
	if proposalMessage == nil {
		// If the state is in prepare phase, the proposal must be set. Otherwise, just return false since
		// this is a valid scenario e.g. proposal msg is received before prepare msg for the same view
		if stateName == prepare {
			vm.log.Error("HasPrepareQuorum - proposalMessage is not set")
		}

		return false
	}

	proposerAddress := proposalMessage.From
	sendersAddressesMap := map[string]struct{}{
		string(proposerAddress): {},
	}

	for _, message := range msgs {
		if bytes.Equal(message.From, proposerAddress) {
			vm.log.Error("HasPrepareQuorum - proposer is among signers but it is not expected to be")

			return false
		}

		sendersAddressesMap[string(message.From)] = struct{}{}
	}

	return vm.HasQuorum(sendersAddressesMap)
}

// HasRoundChangeQuorum provides information on whether round-change messages have reached the needed quorum.
// When round is above rcMinQuorumThreshold we allow for easier RC quorum
// to achieve faster restore in case of network stall. Otherwise we use the default quorum.
func (vm *ValidatorManager) HasRoundChangeQuorum(currentRound uint64, sendersAddrs map[string]struct{}) bool {
	if currentRound <= rcMinQuorumThreshold {
		return vm.HasQuorum(sendersAddrs)
	}

	vm.vpLock.RLock()
	defer vm.vpLock.RUnlock()

	if vm.validatorsVotingPower == nil {
		return false
	}

	messageVotePower := big.NewInt(0)

	for from := range sendersAddrs {
		if vote, ok := vm.validatorsVotingPower[from]; ok {
			messageVotePower.Add(messageVotePower, vote)
		}
	}

	return messageVotePower.Cmp(vm.rcMinQuorum) >= 0
}

// calculateQuorum calculates quorum size which is FLOOR(2 * totalVotingPower / 3) + 1
func calculateQuorum(totalVotingPower *big.Int) *big.Int {
	quorum := new(big.Int).Mul(totalVotingPower, big.NewInt(2))

	// this will floor the (2 * totalVotingPower/3) and add 1
	return quorum.Div(quorum, big.NewInt(3)).Add(quorum, big.NewInt(1))
}

// calculateRCMinQuorum calculates a special RC quorum size
// Quorum is FLOOR(30%)
func calculateRCMinQuorum(totalVotingPower *big.Int) *big.Int {
	quorum := new(big.Int).Mul(totalVotingPower, big.NewInt(30))

	return quorum.Div(quorum, big.NewInt(100))
}

func calculateTotalVotingPower(validatorsVotingPower map[string]*big.Int) *big.Int {
	totalVotingPower := big.NewInt(0)
	for _, validatorVotingPower := range validatorsVotingPower {
		totalVotingPower.Add(totalVotingPower, validatorVotingPower)
	}

	return totalVotingPower
}

// convertMessageToAddressSet converts messages slice to addresses map
func convertMessageToAddressSet(messages []*proto.Message) map[string]struct{} {
	result := make(map[string]struct{}, len(messages))

	for _, x := range messages {
		result[string(x.From)] = struct{}{}
	}

	return result
}
