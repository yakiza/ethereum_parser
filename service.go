package ethereum_parser

import (
	"context"
	"log"
)

// Ensuring that we are implementing the service interface
var _ Service = &service{}

type service struct {
	repo         Repository
	ethClient    EthereumClient
	newSubSignal chan bool
}

func (s *service) GetCurrentBlock(ctx context.Context) (int64, error) {
	currentBlock, err := s.repo.GetCurrentBlock(ctx)
	if err != nil {
		log.Println("There was an error trying to retrieve the current block from local repo")
		return 0, err
	}

	if currentBlock == 0 {
		block, err := s.ethClient.GetCurrentBlock(ctx, ID)
		if err != nil {
			return 0, err
		}

		err = s.repo.SetCurrentBlock(ctx, block)
		if err != nil {
			return 0, err
		}
		currentBlock = block
	}

	log.Printf("Retrieved current block %d", currentBlock)

	return currentBlock, err
}

func (s *service) Subscribe(ctx context.Context, address string) (bool, error) {
	if err := s.repo.Subscribe(ctx, address); err != nil {
		log.Printf("There was an issue trying to subscribe for %v", address)
		return false, err
	}

	s.newSubSignal <- true

	log.Printf("Subscribed for %v", address)

	return true, nil
}

func (s *service) GetTransactions(ctx context.Context, address string) ([]Transaction, error) {
	log.Printf("Retrieving transactions for %v", address)

	transactions, err := s.repo.GetTransactions(ctx, address)
	if err != nil {
		log.Printf("There was an issue trying to retrieve the transactions for %v", address)
		return nil, err
	}

	return transactions, err
}

func NewService(repo Repository, client EthereumClient, newSub chan bool) service {
	return service{
		repo:         repo,
		ethClient:    client,
		newSubSignal: newSub,
	}
}

type Service interface {
	// GetCurrentBlock last parsed block
	GetCurrentBlock(ctx context.Context) (int64, error)

	// Subscribe add address to observer
	Subscribe(ctx context.Context, address string) (bool, error)

	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) ([]Transaction, error)
}
