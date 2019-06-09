package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	// StartingGatewayID is the initial number from which the gateway ids start
	StartingGatewayID = 2
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params  `json:"params"` // asset params
	Assets []Asset `json:"assets"` // issued assets
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	// set the initial gateway id
	if err := k.setInitialGatewayID(ctx, StartingGatewayID); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	// TODO: init assets with data.Assets
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	assets := []Asset{} // TODO: extract existing assets from app state
	return GenesisState{
		Params: k.GetParamSet(ctx),
		Assets: assets,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		Assets: []Asset{},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Params: DefaultParamsForTest(),
		Assets: []Asset{},
	}
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}
