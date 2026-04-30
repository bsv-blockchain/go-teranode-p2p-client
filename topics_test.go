package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTopicName_Mainnet(t *testing.T) {
	topic := TopicName("main", TopicBlock)
	assert.Equal(t, "teranode/bitcoin/1.0.0/mainnet-block", topic)
}

func TestTopicName_AllNetworks(t *testing.T) {
	tests := []struct {
		name     string
		network  string
		topic    string
		expected string
	}{
		{
			name:     "mainnet block using short name",
			network:  "main",
			topic:    TopicBlock,
			expected: "teranode/bitcoin/1.0.0/mainnet-block",
		},
		{
			name:     "mainnet block using full name",
			network:  NetworkMainnet,
			topic:    TopicBlock,
			expected: "teranode/bitcoin/1.0.0/mainnet-block",
		},
		{
			name:     "testnet subtree using short name",
			network:  "test",
			topic:    TopicSubtree,
			expected: "teranode/bitcoin/1.0.0/testnet-subtree",
		},
		{
			name:     "testnet subtree using full name",
			network:  NetworkTestnet,
			topic:    TopicSubtree,
			expected: "teranode/bitcoin/1.0.0/testnet-subtree",
		},
		{
			name:     "stn rejected-tx",
			network:  NetworkSTN,
			topic:    TopicRejectedTx,
			expected: "teranode/bitcoin/1.0.0/stn-rejected-tx",
		},
		{
			name:     "teratestnet node_status using short name",
			network:  "teratest",
			topic:    TopicNodeStatus,
			expected: "teranode/bitcoin/1.0.0/teratestnet-node_status",
		},
		{
			name:     "teratestnet node_status using full name",
			network:  NetworkTeratestnet,
			topic:    TopicNodeStatus,
			expected: "teranode/bitcoin/1.0.0/teratestnet-node_status",
		},
		{
			name:     "regtest block",
			network:  NetworkRegtest,
			topic:    TopicBlock,
			expected: "teranode/bitcoin/1.0.0/regtest-block",
		},
		{
			name:     "unknown network passes through unchanged",
			network:  "custom",
			topic:    TopicBlock,
			expected: "teranode/bitcoin/1.0.0/custom-block",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TopicName(tt.network, tt.topic)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAllTopics(t *testing.T) {
	topics := AllTopics("main")

	require.Len(t, topics, 4)
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/mainnet-block")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/mainnet-subtree")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/mainnet-rejected-tx")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/mainnet-node_status")
}

func TestAllTopics_Testnet(t *testing.T) {
	topics := AllTopics(NetworkTestnet)

	require.Len(t, topics, 4)
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/testnet-block")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/testnet-subtree")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/testnet-rejected-tx")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/testnet-node_status")
}

func TestAllTopics_Regtest(t *testing.T) {
	topics := AllTopics(NetworkRegtest)

	require.Len(t, topics, 4)
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/regtest-block")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/regtest-subtree")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/regtest-rejected-tx")
	assert.Contains(t, topics, "teranode/bitcoin/1.0.0/regtest-node_status")
}

func TestNetworkMapping(t *testing.T) {
	mapping := getNetworkToTopic()

	// Test that short names map correctly
	assert.Equal(t, NetworkMainnet, mapping["main"])
	assert.Equal(t, NetworkTestnet, mapping["test"])
	assert.Equal(t, NetworkTeratestnet, mapping["teratest"])

	// Test that full names map to themselves
	assert.Equal(t, NetworkMainnet, mapping[NetworkMainnet])
	assert.Equal(t, NetworkTestnet, mapping[NetworkTestnet])
	assert.Equal(t, NetworkSTN, mapping[NetworkSTN])
	assert.Equal(t, NetworkTeratestnet, mapping[NetworkTeratestnet])
	assert.Equal(t, NetworkRegtest, mapping[NetworkRegtest])
}

func TestTopicConstants(t *testing.T) {
	// Verify topic constants have expected values
	assert.Equal(t, "block", TopicBlock)
	assert.Equal(t, "subtree", TopicSubtree)
	assert.Equal(t, "rejected-tx", TopicRejectedTx)
	assert.Equal(t, "node_status", TopicNodeStatus)
}

func TestNetworkConstants(t *testing.T) {
	// Verify network constants have expected values
	assert.Equal(t, "mainnet", NetworkMainnet)
	assert.Equal(t, "testnet", NetworkTestnet)
	assert.Equal(t, "stn", NetworkSTN)
	assert.Equal(t, "teratestnet", NetworkTeratestnet)
	assert.Equal(t, "regtest", NetworkRegtest)
}
