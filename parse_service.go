package ethereum_parser

import (
	"context"
	"log"
	"sync"
	"time"
)

var _ Parser = ParserService{}

type ParserService struct {
	storage Repository
	client  ethereumClient
}

func (p ParserService) Parse(ctx context.Context, newSub chan bool, wg *sync.WaitGroup) error {
	defer wg.Done()

	// Can move internal to environment
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-newSub:
			log.Println("New Sub, retrieving all transactions")
		case <-ticker.C:
			log.Println("Retrieving Transactions")
			transactions, err := p.UnsyncedTransactions(ctx)
			if err != nil {
				return err
			}

			subs, err := p.storage.GetSubscribers(ctx)
			if err != nil {
				return err
			}

			// Adding all unparsed transactions and firing up events
			for _, trans := range transactions {
				for _, sub := range subs {
					if sub == trans.From || sub == trans.To {
						err := p.FireUpEvent(sub, trans)
						if err != nil {
							return err
						}

						err = p.storage.AddTransaction(ctx, trans)
						if err != nil {
							return err
						}
					}
				}
			}
		case <-ctx.Done():
			// times up
			return nil
		}
	}
}

// UnsyncedTransactions responsible for checking if each new transaction from current block is already processed or not
func (p ParserService) UnsyncedTransactions(ctx context.Context) ([]Transaction, error) {
	currentRemoteBlock, err := p.client.GetCurrentBlock(ctx, 1)
	if err != nil {
		return nil, nil
	}

	blockTransactions, err := p.client.GetBlockBlockByNumber(ctx, currentRemoteBlock)
	if err != nil {
		return nil, err
	}

	// gathering transactions that have not been parsed
	var unprocessedTransactions []Transaction
	for _, trans := range blockTransactions.Transactions {
		retrievedTransaction, err := p.storage.GetTransactionByHash(ctx, trans.Hash)
		if err != nil {
			return nil, err
		}

		// Assuming if there is no transaction hash there is no transaction
		if retrievedTransaction.Hash == "" {
			unprocessedTransactions = append(unprocessedTransactions, trans)
		}

	}

	return unprocessedTransactions, nil
}

// FireUpEvent will trigger an event that will be sent to the notification service
func (p ParserService) FireUpEvent(address string, transaction Transaction) error {
	log.Printf("Event for address %v transaction with Hash: %v From: %v To: %v with Value: %v ", address, transaction.Hash, transaction.From, transaction.To, transaction.Value)
	return nil
}

func NewParserService(storage Repository, client ethereumClient) ParserService {
	return ParserService{
		storage: storage,
		client:  client,
	}
}

type Parser interface {
	// Parse a parser triggered by a ticker as well as new subscription
	Parse(ctx context.Context, newSub chan bool, wg *sync.WaitGroup) error

	// UnsyncedTransactions retrieves all transactions that have not been parsed ( perhaps better name ( sync local with blockchain was the thought )
	UnsyncedTransactions(ctx context.Context) ([]Transaction, error)

	// FireUpEvent responsible for sending an event to the notification service
	FireUpEvent(address string, transaction Transaction) error
}
