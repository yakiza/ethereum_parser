package ethereum_parser

import (
	"context"
	"log"
	"sync"
)

type InMemStorage struct {
	mux sync.Mutex

	transactions      map[string][]Transaction
	TransactionByHash map[string]Transaction
	subscribers       map[string]bool

	currentBlock int64
}

func (s *InMemStorage) GetCurrentBlock(_ context.Context) (int64, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.currentBlock, nil
}

func (s *InMemStorage) SetCurrentBlock(_ context.Context, currentBlock int64) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.currentBlock = currentBlock

	return nil
}

func (s *InMemStorage) GetTransactions(_ context.Context, address string) ([]Transaction, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.transactions[address], nil
}

func (s *InMemStorage) GetTransactionByHash(_ context.Context, hash string) (Transaction, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.TransactionByHash[hash]; ok {
		return s.TransactionByHash[hash], nil
	}

	return Transaction{}, nil
}

func (s *InMemStorage) AddTransaction(_ context.Context, transaction Transaction) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.TransactionByHash[transaction.Hash] = transaction
	s.transactions[transaction.From] = append(s.transactions[transaction.From], transaction)
	s.transactions[transaction.To] = append(s.transactions[transaction.To], transaction)

	return nil
}

func (s *InMemStorage) Subscribe(_ context.Context, address string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.subscribers[address]; !ok {
		s.subscribers[address] = true
	} else {
		// Not really an error perhaps could log a message, or we might not even care and just say nothing
		log.Printf("%v is already subscribed", address)
	}

	return nil
}

func (s *InMemStorage) GetSubscribers(ctx context.Context) ([]string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var subscribers []string
	for sub, _ := range s.subscribers {
		subscribers = append(subscribers, sub)
	}

	return subscribers, nil
}

func NewMemStorage() InMemStorage {
	return InMemStorage{
		transactions:      make(map[string][]Transaction),
		TransactionByHash: make(map[string]Transaction),
		subscribers:       make(map[string]bool),
	}
}

type Repository interface {
	// GetCurrentBlock retrieving current block from locally, if it's missing it goes and fetches it from ethereum ethClient
	GetCurrentBlock(ctx context.Context) (int64, error)

	// SetCurrentBlock responsible for setting the current block locally
	SetCurrentBlock(ctx context.Context, currentBlock int64) error

	// Subscribe subscribes an address
	Subscribe(ctx context.Context, address string) error

	// GetSubscribers retrieves all subscribed addressed
	GetSubscribers(ctx context.Context) ([]string, error)

	// GetTransactions retrieves all parsed transactions from repo
	GetTransactions(ctx context.Context, address string) ([]Transaction, error)

	// AddTransaction responsible for inserting a single transaction in repo
	AddTransaction(ctx context.Context, transaction Transaction) error

	// GetTransactionByHash retrieves transaction data for given hash
	GetTransactionByHash(_ context.Context, hash string) (Transaction, error)
}
