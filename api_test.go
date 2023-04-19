package ethereum_parser_test

//
//type APITestSuite struct {
//	suite.Suite
//	handler http.Handler
//	service ethereum_parser.Service
//}
//
//func (suite *APITestSuite) SetupSuite() {
//	suite.service = ServiceTestDouble{}
//}
//
//var _ ethereum_parser.Service = ServiceTestDouble{}
//
//type ServiceTestDouble struct {
//	GetCurrentBlockTD func(ctx context.Context) (int64, error)
//
//	// Subscribe add address to observer
//	SubscribeTD func(ctx context.Context, address string) (bool, error)
//
//	// GetTransactions list of inbound or outbound transactions for an address
//	GetTransactionsTD func(ctx context.Context, address string) ([]Transaction, error)
//}
//
//func (s ServiceTestDouble) GetCurrentBlock(ctx context.Context) (int64, error) {
//	return s.GetCurrentBlockTD(ctx)
//}
//
//func (s ServiceTestDouble) Subscribe(ctx context.Context, address string) (bool, error) {
//	return s.SubscribeTD(ctx, address)
//}
//
//func (s ServiceTestDouble) GetTransactions(ctx context.Context, address string) ([]ethereum_parser.Transaction, error) {
//	return s.GetTransactionsTD(ctx, address)
//}
