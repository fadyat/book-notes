package repo

// BrokerStorage is the storage layer of the broker.
//
// By the PoC, it is an in-memory storage; however, it can be replaced with a persistent storage,
// such as a file system.
type BrokerStorage struct {

	// topics is the map of topics in the broker.
	topics map[string]*Topic

	// producers is the map of producers in the broker.
	producers map[int64]*Producer

	// consumers is the map of consumers in the broker.
	consumers map[int64]*Consumer
}
