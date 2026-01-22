package p2p

import "fmt"

// Protocol prefix for Teranode P2P topics
const protocolPrefix = "teranode/bitcoin/1.0.0"

// Topic type constants
const (
	TopicBlock      = "block"
	TopicSubtree    = "subtree"
	TopicRejectedTx = "rejected-tx"
	TopicNodeStatus = "node_status"
)

// Network constants for use with TopicName
const (
	NetworkMainnet     = "mainnet"
	NetworkTestnet     = "testnet"
	NetworkSTN         = "stn"
	NetworkTeratestnet = "teratestnet"
)

// getNetworkToTopic returns a map from config network names to topic network names.
func getNetworkToTopic() map[string]string {
	return map[string]string{
		"main":             NetworkMainnet,
		NetworkMainnet:     NetworkMainnet,
		"test":             NetworkTestnet,
		NetworkTestnet:     NetworkTestnet,
		NetworkSTN:         NetworkSTN,
		"teratest":         NetworkTeratestnet,
		NetworkTeratestnet: NetworkTeratestnet,
	}
}

// TopicName constructs a full topic name for subscribing to Teranode P2P messages.
// Example: TopicName("main", TopicBlock) returns "teranode/bitcoin/1.0.0/mainnet-block"
func TopicName(network, topic string) string {
	if mapped, ok := getNetworkToTopic()[network]; ok {
		network = mapped
	}

	return fmt.Sprintf("%s/%s-%s", protocolPrefix, network, topic)
}

// AllTopics returns all topic names for a given network.
func AllTopics(network string) []string {
	return []string{
		TopicName(network, TopicBlock),
		TopicName(network, TopicSubtree),
		TopicName(network, TopicRejectedTx),
		TopicName(network, TopicNodeStatus),
	}
}
