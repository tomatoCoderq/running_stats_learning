package repositories

import (
	"context"
	"database/sql"
)

func BeginTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error{
	ctx := context.Background()
	transaction, err := resultsRepository.dbHandler.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	runnersRepository.transaction = transaction
	resultsRepository.transaction = transaction
	return nil
}

func CommitTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	transaction := runnersRepository.transaction
	runnersRepository.transaction = nil
	resultsRepository.transaction = nil
	return transaction.Commit()
   }

func RollBackTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	transaction := runnersRepository.transaction
	runnersRepository.transaction = nil
	resultsRepository.transaction = nil

	return transaction.Rollback()

}