package gpadstats

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
)

const (
	qCount    = "SELECT COUNT(DISTINCT DB_Object_ID) FROM gpad"
	qEcoCount = `SELECT COUNT(DISTINCT DB_Object_ID) FROM gpad WHERE Evidence_Code IN (
				'ECO:0000269', 'ECO:0000314', 'ECO:0000353', 
				'ECO:0000315', 'ECO:0000316',
				'ECO:0000270', 'ECO:0006056', 
				'ECO:0007005', 'ECO:0007001', 
				'ECO:0007003','ECO:0007007', 'ECO:0005581'
			)`
)

func queryCount(config StatsLoaderConfig) IOE.IOEither[error, int] {
	return IOE.TryCatch(func() (int, error) {
		var count int
		err := config.DB.QueryRow(
			qCount,
		).Scan(&count)
		return count, err
	}, func(err error) error {
		return fmt.Errorf("error in running gpad query %w", err)
	})
}

func queryEcoCount(config StatsLoaderConfig) IOE.IOEither[error, int] {
	return IOE.TryCatch(func() (int, error) {
		var count int
		err := config.DB.QueryRow(
			qEcoCount,
		).Scan(&count)
		return count, err
	}, func(err error) error {
		return fmt.Errorf("error in running gpad eco query %w", err)
	})
}

func queryCounts(config StatsLoaderConfig) IOE.IOEither[error, GeneCountStats] {
	return F.Pipe1(
		IOE.SequenceArraySeq([]IOE.IOEither[error, int]{
			queryCount(config),
			queryEcoCount(config),
		}),
		IOE.Map[error](func(results []int) GeneCountStats {
			return GeneCountStats{
				Count:    results[0],
				EcoCount: results[1],
			}
		}),
	)
}
