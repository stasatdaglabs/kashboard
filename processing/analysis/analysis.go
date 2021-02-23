package analysis

import (
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kashboard/processing/database"
	"github.com/stasatdaglabs/kashboard/processing/database/model"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
	"time"
)

var log = logging.Logger()
var spawn = panics.GoroutineWrapperFunc(log)

func Start(database *database.Database, blockChan chan *model.Block) {
	spawn("analysis", func() {
		for block := range blockChan {
			err := handleBlock(database, block)
			if err != nil {
				panic(err)
			}
		}
	})
}

func handleBlock(database *database.Database, block *model.Block) error {
	const durationForAnalysis = 1 * time.Minute

	averageParentAmount, err := database.AverageParentAmount(block, durationForAnalysis)
	if err != nil {
		return err
	}

	blockCount, err := database.BlockCount(block, durationForAnalysis)
	if err != nil {
		return err
	}
	blockRate := float64(blockCount) / durationForAnalysis.Seconds()

	transactionCount, err := database.TransactionCount(block, durationForAnalysis)
	if err != nil {
		return err
	}
	transactionRate := float64(transactionCount) / durationForAnalysis.Seconds()

	analyzedBlock := &model.AnalyzedBlock{
		ID:                  block.ID,
		Timestamp:           block.Timestamp,
		AverageParentAmount: averageParentAmount,
		BlockRate:           blockRate,
		TransactionRate:     transactionRate,
	}
	return database.InsertAnalyzedBlock(analyzedBlock)
}