package unsigned

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-go/data/feeTx"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"math/big"
	"sync"
)

// TODO: Set MinGasPrice and MinTxFee to some positive value (TBD)
// MinGasPrice is the minimal gas price to be paid for any transaction
var MinGasPrice = uint64(0)

// MinTxFee is the minimal fee to be paid for any transaction
var MinTxFee = uint64(0)

const communityPercentage = 0.1 // 1 = 100%, 0 = 0%
const leaderPercentage = 0.4    // 1 = 100%, 0 = 0%
const burnPercentage = 0.5      // 1 = 100%, 0 = 0%

type feeTxHandler struct {
	address     process.SpecialAddressHandler
	hasher      hashing.Hasher
	marshalizer marshal.Marshalizer
	mutTxs      sync.Mutex
	feeTxs      []*feeTx.FeeTx

	feeTxsFromBlock map[string]*feeTx.FeeTx
}

// NewFeeTxHandler constructor for the fx tee handler
func NewFeeTxHandler(
	address process.SpecialAddressHandler,
	hasher hashing.Hasher,
	marshalizer marshal.Marshalizer,
) (*feeTxHandler, error) {
	if address == nil {
		return nil, process.ErrNilSpecialAddressHandler
	}
	if hasher == nil {
		return nil, process.ErrNilHasher
	}
	if marshalizer == nil {
		return nil, process.ErrNilMarshalizer
	}

	ftxh := &feeTxHandler{
		address:     address,
		hasher:      hasher,
		marshalizer: marshalizer,
	}
	ftxh.feeTxs = make([]*feeTx.FeeTx, 0)
	ftxh.feeTxsFromBlock = make(map[string]*feeTx.FeeTx)

	return ftxh, nil
}

// SaveCurrentIntermediateTxToStorage saves current cached data into storage - already saaved for txs
func (ftxh *feeTxHandler) SaveCurrentIntermediateTxToStorage() error {
	//TODO implement me - save only created feeTxs
	return nil
}

// AddIntermediateTransactions adds intermediate transactions to local cache
func (ftxh *feeTxHandler) AddIntermediateTransactions(txs []data.TransactionHandler) error {
	return nil
}

// CreateAllInterMiniBlocks creates miniblocks from process transactions
func (ftxh *feeTxHandler) CreateAllInterMiniBlocks() map[uint32]*block.MiniBlock {
	calculatedFeeTxs := ftxh.CreateAllUTxs()

	miniBlocks := make(map[uint32]*block.MiniBlock)
	for _, value := range calculatedFeeTxs {
		dstShId, err := ftxh.address.ShardIdForAddress(value.GetRecvAddress())
		if err != nil {
			log.Debug(err.Error())
			continue
		}

		txHash, err := core.CalculateHash(ftxh.marshalizer, ftxh.hasher, value)
		if err != nil {
			log.Debug(err.Error())
			continue
		}

		var ok bool
		var mb *block.MiniBlock
		if mb, ok = miniBlocks[dstShId]; !ok {
			mb = &block.MiniBlock{
				ReceiverShardID: dstShId,
			}
		}

		mb.TxHashes = append(mb.TxHashes, txHash)
		miniBlocks[dstShId] = mb
	}

	return miniBlocks
}

// VerifyInterMiniBlocks verifies if transaction fees were correctly handled for the block
func (ftxh *feeTxHandler) VerifyInterMiniBlocks(body block.Body) error {
	err := ftxh.VerifyCreatedUTxs()
	ftxh.CleanProcessedUTxs()

	return err
}

func (ftxh *feeTxHandler) CreateBlockStarted() {
	ftxh.CleanProcessedUTxs()
}

// CleanProcessedUTxs deletes the cached data
func (ftxh *feeTxHandler) CleanProcessedUTxs() {
	ftxh.mutTxs.Lock()
	ftxh.feeTxs = make([]*feeTx.FeeTx, 0)
	ftxh.feeTxsFromBlock = make(map[string]*feeTx.FeeTx)
	ftxh.mutTxs.Unlock()
}

// AddTxFeeFromBlock adds an existing txfee from block into local cache
func (ftxh *feeTxHandler) AddTxFeeFromBlock(tx data.TransactionHandler) {
	currFeeTx, ok := tx.(*feeTx.FeeTx)
	if !ok {
		log.Error(process.ErrWrongTypeAssertion.Error())
		return
	}

	ftxh.mutTxs.Lock()
	ftxh.feeTxsFromBlock[string(tx.GetRecvAddress())] = currFeeTx
	ftxh.mutTxs.Unlock()
}

// AddProcessedUTx adds a new feeTx to the cache
func (ftxh *feeTxHandler) AddProcessedUTx(tx data.TransactionHandler) {
	currFeeTx, ok := tx.(*feeTx.FeeTx)
	if !ok {
		log.Debug(process.ErrWrongTypeAssertion.Error())
		return
	}

	ftxh.mutTxs.Lock()
	ftxh.feeTxs = append(ftxh.feeTxs, currFeeTx)
	ftxh.mutTxs.Unlock()
}

func getPercentageOfValue(value *big.Int, percentage float64) *big.Int {
	x := new(big.Float).SetInt(value)
	y := big.NewFloat(percentage)

	z := new(big.Float).Mul(x, y)

	op := big.NewInt(0)
	result, _ := z.Int(op)

	return result
}

func (ftxh *feeTxHandler) createLeaderTx(totalGathered *big.Int) *feeTx.FeeTx {
	currTx := &feeTx.FeeTx{}

	currTx.Value = getPercentageOfValue(totalGathered, leaderPercentage)
	currTx.RcvAddr = ftxh.address.LeaderAddress()

	return currTx
}

func (ftxh *feeTxHandler) createBurnTx(totalGathered *big.Int) *feeTx.FeeTx {
	currTx := &feeTx.FeeTx{}

	currTx.Value = getPercentageOfValue(totalGathered, burnPercentage)
	currTx.RcvAddr = ftxh.address.BurnAddress()

	return currTx
}

func (ftxh *feeTxHandler) createCommunityTx(totalGathered *big.Int) *feeTx.FeeTx {
	currTx := &feeTx.FeeTx{}

	currTx.Value = getPercentageOfValue(totalGathered, communityPercentage)
	currTx.RcvAddr = ftxh.address.ElrondCommunityAddress()

	return currTx
}

// CreateAllUTxs creates all the needed fee transactions
// According to economic paper 50% burn, 40% to the leader, 10% to Elrond community fund
func (ftxh *feeTxHandler) CreateAllUTxs() []data.TransactionHandler {
	ftxh.mutTxs.Lock()
	defer ftxh.mutTxs.Unlock()

	totalFee := big.NewInt(0)
	for _, val := range ftxh.feeTxs {
		totalFee = totalFee.Add(totalFee, val.Value)
	}

	if totalFee.Cmp(big.NewInt(1)) < 0 {
		ftxh.feeTxs = make([]*feeTx.FeeTx, 0)
		return nil
	}

	leaderTx := ftxh.createLeaderTx(totalFee)
	communityTx := ftxh.createCommunityTx(totalFee)
	burnTx := ftxh.createBurnTx(totalFee)

	currFeeTxs := make([]data.TransactionHandler, 0)
	currFeeTxs = append(currFeeTxs, leaderTx, communityTx, burnTx)

	ftxh.feeTxs = make([]*feeTx.FeeTx, 0)

	return currFeeTxs
}

// VerifyCreatedUTxs creates all fee txs from added values, than verifies if in block the values are the same
func (ftxh *feeTxHandler) VerifyCreatedUTxs() error {
	calculatedFeeTxs := ftxh.CreateAllUTxs()

	ftxh.mutTxs.Lock()
	defer ftxh.mutTxs.Unlock()

	totalFeesFromBlock := big.NewInt(0)
	for _, value := range ftxh.feeTxsFromBlock {
		totalFeesFromBlock = totalFeesFromBlock.Add(totalFeesFromBlock, value.Value)
	}

	totalCalculatedFees := big.NewInt(0)
	for _, value := range calculatedFeeTxs {
		totalCalculatedFees = totalCalculatedFees.Add(totalCalculatedFees, value.GetValue())

		txFromBlock, ok := ftxh.feeTxsFromBlock[string(value.GetRecvAddress())]
		if !ok {
			return process.ErrTxsFeesNotFound
		}
		if txFromBlock.Value.Cmp(value.GetValue()) != 0 {
			return process.ErrTxsFeesDoesNotMatch
		}
	}

	if totalCalculatedFees.Cmp(totalFeesFromBlock) != 0 {
		return process.ErrTotalTxsFeesDoNotMatch
	}

	return nil
}

// CreateMarshalizedData creates the marshalized data for broadcasting purposes
func (ftxh *feeTxHandler) CreateMarshalizedData(txHashes [][]byte) ([][]byte, error) {
	// TODO: implement me

	return make([][]byte, 0), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ftxh *feeTxHandler) IsInterfaceNil() bool {
	if ftxh == nil {
		return true
	}
	return false
}