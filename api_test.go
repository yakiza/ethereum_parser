package ethereum_parser_test

import (
	"context"
	"encoding/json"
	"ethereum_parser"
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type APITestSuite struct {
	suite.Suite
	handler http.Handler
	service ServiceTestDouble
}

func (suite *APITestSuite) SetupSuite() {
	suite.service = ServiceTestDouble{}
	suite.handler = ethereum_parser.CreateAPIMux(ethereum_parser.NewHTTPHandlers(&suite.service))
}

func (suite *APITestSuite) TestSubscriberSuccess() {
	suite.service.SubscribeTD = func(ctx context.Context, receivedAddress string) (bool, error) {
		suite.Equal(address, receivedAddress)
		return true, nil
	}

	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/subscribe?address=%v", address), nil)
	suite.Require().NoError(err)

	w := httptest.NewRecorder()
	suite.handler.ServeHTTP(w, r)

	var actual string
	err = json.Unmarshal(w.Body.Bytes(), &actual)
	suite.Require().NoError(err)

	suite.Equal(subscribedTrue, actual)
	suite.Require().Equal(http.StatusOK, w.Code)
}

func (suite *APITestSuite) TestSubscriberUnSuccessful() {
	suite.service.SubscribeTD = func(ctx context.Context, receivedAddress string) (bool, error) {
		suite.Equal(address, receivedAddress)
		return false, fmt.Errorf("there was an error trying to subscribe")
	}

	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/subscribe?address=%v", address), nil)
	suite.Require().NoError(err)

	w := httptest.NewRecorder()
	suite.handler.ServeHTTP(w, r)

	var actual string
	err = json.Unmarshal(w.Body.Bytes(), &actual)
	suite.Require().Error(err)
	suite.Require().Equal(http.StatusInternalServerError, w.Code)
}

func TestAPI(t *testing.T) {
	suite.Run(t, &APITestSuite{})
}

// Implement test double, only the subscribe has been tested above the rest are left for demonstration purposes
var _ ethereum_parser.Service = ServiceTestDouble{}

type ServiceTestDouble struct {
	GetCurrentBlockTD func(ctx context.Context) (int64, error)

	// Subscribe add address to observer
	SubscribeTD func(ctx context.Context, address string) (bool, error)

	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactionsTD func(ctx context.Context, address string) ([]ethereum_parser.Transaction, error)
}

func (s ServiceTestDouble) GetCurrentBlock(ctx context.Context) (int64, error) {
	return s.GetCurrentBlockTD(ctx)
}

func (s ServiceTestDouble) Subscribe(ctx context.Context, address string) (bool, error) {
	return s.SubscribeTD(ctx, address)
}

func (s ServiceTestDouble) GetTransactions(ctx context.Context, address string) ([]ethereum_parser.Transaction, error) {
	return s.GetTransactionsTD(ctx, address)
}

const address = "0xae2fc483527b8ef99eb5d9b44875f005ba1fae13"
const subscribedTrue = "subscribed true"
