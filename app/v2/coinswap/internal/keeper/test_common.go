package keeper

import (
	"github.com/irisnet/irishub/app/protocol"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
)

// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput(t *testing.T, amt sdk.Int, nAccs int64) (sdk.Context, Keeper, []auth.Account) {
	keyAcc := protocol.KeyAccount
	keyParams := protocol.KeyParams
	tkeyParams := protocol.TkeyParams
	keyCoinswap := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyCoinswap, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	cdc := makeTestCodec()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "coinswap-chain"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(cdc, keyAcc, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)

	initialCoins := sdk.Coins{
		sdk.NewCoin(sdk.IrisAtto, amt),
		sdk.NewCoin("btc-min", amt),
	}
	initialCoins = initialCoins.Sort()
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)

	keeper := NewKeeper(cdc, keyCoinswap, bk, ak, pk.Subspace(types.DefaultParamSpace))
	keeper.SetParams(ctx, types.DefaultParams())

	return ctx, keeper, accs
}

func createTestAccs(ctx sdk.Context, numAccs int, initialCoins sdk.Coins, ak *auth.AccountKeeper) (accs []auth.Account) {
	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())
		acc := auth.NewBaseAccountWithAddress(addr)
		acc.Coins = initialCoins
		acc.PubKey = pubKey
		acc.AccountNumber = uint64(i)
		ak.SetAccount(ctx, &acc)
		accs = append(accs, &acc)
	}
	return
}
