package core

import (
	"github.com/Trapesys/go-ibft/messages/proto"
)

// messageConstructor defines a message constructor interface
type messageConstructor interface {
	// BuildPrePrepareMessage builds a PREPREPARE message based on the passed in proposal
	BuildPrePrepareMessage(proposal []byte, view *proto.View) *proto.Message

	// BuildPrepareMessage builds a PREPARE message based on the passed in proposal
	BuildPrepareMessage(proposal []byte, view *proto.View) *proto.Message

	// BuildCommitMessage builds a COMMIT message based on the passed in proposal
	BuildCommitMessage(proposal []byte, view *proto.View) *proto.Message

	// BuildRoundChangeMessage builds a ROUND_CHANGE message based on the passed in proposal
	BuildRoundChangeMessage(height, round uint64) *proto.Message
}

// Backend defines an interface all backend implementations
// need to implement
type Backend interface {
	messageConstructor

	// IsValidBlock checks if the proposed block is child of parent
	IsValidBlock(block []byte) bool

	// IsValidSender checks if signature is from sender
	IsValidSender(msg *proto.Message) bool

	// IsProposer checks if the passed in ID is the Proposer for current view (sequence, round)
	IsProposer(id []byte, height, round uint64) bool

	// BuildProposal builds a new block proposal
	BuildProposal(blockNumber uint64) ([]byte, error)

	// VerifyProposalHash checks if the hash matches the proposal
	VerifyProposalHash(proposal, hash []byte) error

	// IsValidCommittedSeal checks if the seal for the proposal is valid
	IsValidCommittedSeal(proposal, seal []byte) bool

	InsertBlock(proposal []byte, committedSeals [][]byte) error

	// ID returns the validator's ID
	ID() []byte

	//	ValidatorCount returns the number of validators for the given block
	ValidatorCount(blockNumber uint64) uint64

	// AllowedFaulty returns the maximum number of faulty nodes based
	// on the validator set
	AllowedFaulty() uint64
}