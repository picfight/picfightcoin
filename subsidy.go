package picfightcoin

import (
	"github.com/jfixby/bignum"
	"github.com/jfixby/coin"
)

type SubsidyCalculator interface {
	ExpectedTotalNetworkSubsidy() coin.Amount
	NumberOfGeneratingBlocks() int64
	PreminedCoins() coin.Amount
	FirstGeneratingBlockIndex() int64
	CalcBlockWorkSubsidy(height int64, voters uint16) int64
	CalcStakeVoteSubsidy(height int64) int64
	CalcBlockTaxSubsidy(height int64, voters uint16) int64
	CalcBlockSubsidy(height int64) int64
	TicketsPerBlock() uint16
	SetEngine(engine bignum.BigNumEngine)
}

//type PicFightCoinSubsidy struct {
//	// TargetTotalSubsidy is the the expected total subsidy (in coins)
//	// produced by the network.
//	TargetTotalSubsidy coin.Amount
//
//	// SubsidyProductionPeriod is the the estimated time-period during which
//	// all the subsidy should be produced.
//	SubsidyProductionPeriod time.Duration
//
//	TargetTimePerBlock time.Duration
//
//	Premine map[string]coin.Amount
//
//}

// CalcBlockSubsidy returns the subsidy amount a block at the provided height
// should have. This is mainly used for determining how much the coinbase for
// newly generated blocks awards as well as validating the coinbase for blocks
// has the expected value.
//func CalcBlockSubsidy(height int32) int64 {
//	if bignumEngine == nil {
//		// use the default float64
//		bignumEngine = bignum.Float64Engine{}
//	}
//	satoshiBigNum := CalcBlockSubsidy(bignumEngine, height)
//	return int64(satoshiBigNum.ToFloat64())
//}
//
//func (s *PicFightCoinSubsidy) CalcBlockSubsidy(engine bignum.BigNumEngine, height int32, subsidySettings *PicFightCoinSubsidy) bignum.BigNum {
//	period := subsidySettings.SubsidyProductionPeriod
//	blockTime := subsidySettings.TargetTimePerBlock
//	totalSubsidy := subsidySettings.TargetTotalSubsidy
//	premine := chaincfg.Sum(subsidySettings.Premine)
//	subsidyBlocksNumber := int64(period / blockTime)
//	subsidyCoins := lineardown.LinearDownGenerate(engine, subsidyBlocksNumber, int64(height), totalSubsidy.AtomsValue-premine)
//	satoshi := engine.NewBigNum(chaincfg.SatoshiPerPicfightcoin)
//	satoshi = satoshi.Mul(satoshi, subsidyCoins)
//	return satoshi
//}
