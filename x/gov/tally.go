package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Address         sdk.AccAddress // sdk.AccAddress of the validator owner
	Power           sdk.Rat        // Power of a Validator
	DelegatorShares sdk.Rat        // Total outstanding delegator shares
	Minus           sdk.Rat        // Minus of validator, used to compute validator's voting power
	Vote            VoteOption     // Vote of the validator
}

func tally(ctx sdk.Context, keeper Keeper, proposal Proposal) (passes bool, nonVoting []sdk.AccAddress) {
	results := make(map[VoteOption]sdk.Rat)
	results[OptionYes] = sdk.ZeroRat()
	results[OptionAbstain] = sdk.ZeroRat()
	results[OptionNo] = sdk.ZeroRat()
	results[OptionNoWithVeto] = sdk.ZeroRat()

	totalVotingPower := sdk.ZeroRat()
	currValidators := make(map[string]validatorGovInfo)

	for _, valAddr := range keeper.ds.GetBondedValidatorOwnerAddresses(ctx) {
		power, _ := keeper.ds.GetValidatorPower(ctx, valAddr)
		delegatorShares, _ := keeper.ds.GetValidatorTotalDelegationShares(ctx, valAddr)
		currValidators[valAddr.String()] = validatorGovInfo{
			Address:         valAddr,
			Power:           power,
			DelegatorShares: delegatorShares,
			Minus:           sdk.ZeroRat(),
			Vote:            OptionEmpty,
		}
	}

	// iterate over all the votes
	votesIterator := keeper.GetVotes(ctx, proposal.GetProposalID())
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := &Vote{}
		keeper.cdc.MustUnmarshalBinary(votesIterator.Value(), vote)

		// if validator, just record it in the map
		// if delegator tally voting power
		if val, ok := currValidators[vote.Voter.String()]; ok {
			val.Vote = vote.Option
			currValidators[vote.Voter.String()] = val
		} else {

			for _, valAddr := range keeper.ds.GetDelegatorDelegations(ctx, vote.Voter) {
				if val, ok := currValidators[valAddr.String()]; ok {
					bondShares, _ := keeper.ds.GetDelegatorDelegationShares(ctx, vote.Voter, val.Address)
					val.Minus = val.Minus.Add(bondShares)
					currValidators[valAddr.String()] = val

					votingPower := bondShares.Quo(val.DelegatorShares).Mul(val.Power)
					results[vote.Option] = results[vote.Option].Add(votingPower)
					totalVotingPower = totalVotingPower.Add(votingPower)

				}

			}
		}

		keeper.deleteVote(ctx, vote.ProposalID, vote.Voter)
	}
	votesIterator.Close()

	// Iterate over the validators again to tally their voting power and see who didn't vote
	nonVoting = []sdk.AccAddress{}
	for _, val := range currValidators {
		if val.Vote == OptionEmpty {
			nonVoting = append(nonVoting, val.Address)
			continue
		}
		sharesAfterMinus := val.DelegatorShares.Sub(val.Minus)
		percentAfterMinus := sharesAfterMinus.Quo(val.DelegatorShares)
		votingPower := val.Power.Mul(percentAfterMinus)

		results[val.Vote] = results[val.Vote].Add(votingPower)
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyingProcedure := keeper.GetTallyingProcedure()

	// If no one votes, proposal fails
	if totalVotingPower.Sub(results[OptionAbstain]).Equal(sdk.ZeroRat()) {
		return false, nonVoting
	}
	// If more than 1/3 of voters veto, proposal fails
	if results[OptionNoWithVeto].Quo(totalVotingPower).GT(tallyingProcedure.Veto) {
		return false, nonVoting
	}
	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[OptionYes].Quo(totalVotingPower.Sub(results[OptionAbstain])).GT(tallyingProcedure.Threshold) {
		return true, nonVoting
	}
	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, nonVoting
}
