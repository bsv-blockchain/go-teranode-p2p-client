package p2p

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	msgbus "github.com/bsv-blockchain/go-p2p-message-bus"
	teranode "github.com/bsv-blockchain/teranode/services/p2p"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMsgBusClient implements a mock message bus client for testing
type mockMsgBusClient struct {
	id            string
	peers         []msgbus.PeerInfo
	subscriptions map[string]chan msgbus.Message
	mu            sync.Mutex
	closed        bool
}

func newMockMsgBusClient(id string) *mockMsgBusClient {
	return &mockMsgBusClient{
		id:            id,
		peers:         []msgbus.PeerInfo{},
		subscriptions: make(map[string]chan msgbus.Message),
	}
}

func (m *mockMsgBusClient) GetID() string {
	return m.id
}

func (m *mockMsgBusClient) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true

	for _, ch := range m.subscriptions {
		close(ch)
	}

	return nil
}

func (m *mockMsgBusClient) GetPeers() []msgbus.PeerInfo {
	return m.peers
}

func (m *mockMsgBusClient) Subscribe(topic string) <-chan msgbus.Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch := make(chan msgbus.Message, 100)
	m.subscriptions[topic] = ch

	return ch
}

func (m *mockMsgBusClient) Publish(_ context.Context, _ string, _ []byte) error {
	return nil
}

// sendMessage sends a message to all subscribers of the given topic
func (m *mockMsgBusClient) sendMessage(topic string, data []byte) {
	m.mu.Lock()
	ch, ok := m.subscriptions[topic]
	m.mu.Unlock()

	if ok {
		ch <- msgbus.Message{Data: data}
	}
}

// createTestClient creates a Client with a mock message bus for testing
func createTestClient(mockClient *mockMsgBusClient, network string) *Client {
	return &Client{
		msgbus:  mockClient,
		logger:  slog.Default(),
		network: network,
	}
}

func TestClient_GetID(t *testing.T) {
	mock := newMockMsgBusClient("12D3KooWTestPeerID")
	client := createTestClient(mock, "main")

	assert.Equal(t, "12D3KooWTestPeerID", client.GetID())
}

func TestClient_GetNetwork(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "testnet")

	assert.Equal(t, "testnet", client.GetNetwork())
}

func TestClient_GetNetwork_AllNetworks(t *testing.T) {
	networks := []string{NetworkMainnet, NetworkTestnet, NetworkSTN, NetworkTeratestnet}

	for _, network := range networks {
		t.Run(network, func(t *testing.T) {
			mock := newMockMsgBusClient("test-id")
			client := createTestClient(mock, network)
			assert.Equal(t, network, client.GetNetwork())
		})
	}
}

func TestClient_Close(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "main")

	ctx := context.Background()

	// Subscribe to create channels
	blocks := client.SubscribeBlocks(ctx)
	subtrees := client.SubscribeSubtrees(ctx)
	rejected := client.SubscribeRejectedTxs(ctx)
	status := client.SubscribeNodeStatus(ctx)

	// Close the client
	err := client.Close()
	require.NoError(t, err)

	// Verify all channels are closed by trying to receive (should return zero value + closed)
	_, open := <-blocks
	assert.False(t, open, "blocks channel should be closed")

	_, open = <-subtrees
	assert.False(t, open, "subtrees channel should be closed")

	_, open = <-rejected
	assert.False(t, open, "rejected channel should be closed")

	_, open = <-status
	assert.False(t, open, "status channel should be closed")
}

func TestSubscribeBlocks(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "main")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	blocks := client.SubscribeBlocks(ctx)
	require.NotNil(t, blocks)

	// Send a test message
	blockMsg := teranode.BlockMessage{
		Height: 12345,
		Hash:   "0000000000000000000123456789abcdef",
	}

	data, err := json.Marshal(blockMsg)
	require.NoError(t, err)

	topic := TopicName("main", TopicBlock)
	mock.sendMessage(topic, data)

	// Receive the message
	select {
	case received := <-blocks:
		assert.Equal(t, uint32(12345), received.Height)
		assert.Equal(t, "0000000000000000000123456789abcdef", received.Hash)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for block message")
	}

	cancel()

	_ = client.Close()
}

func TestSubscribeBlocks_ContextCancel(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "main")

	defer func() { _ = client.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	blocks := client.SubscribeBlocks(ctx)

	// Cancel the context
	cancel()

	// The channel should eventually be closed
	select {
	case _, open := <-blocks:
		if open {
			t.Log("received message before channel closed")
		}
		// Channel was closed, which is expected
	case <-time.After(time.Second):
		// Give more time for the goroutine to process
		time.Sleep(100 * time.Millisecond)
	}

	// Verify subscriber was removed
	client.mu.RLock()
	count := len(client.blockSubs)
	client.mu.RUnlock()

	assert.Equal(t, 0, count, "subscriber should be removed after context cancel")
}

func TestSubscribeBlocks_FanOut(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "main")

	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())

	defer cancel1()
	defer cancel2()

	// Create two subscribers
	blocks1 := client.SubscribeBlocks(ctx1)
	blocks2 := client.SubscribeBlocks(ctx2)

	// Send a test message
	blockMsg := teranode.BlockMessage{
		Height: 99999,
		Hash:   "test-hash",
	}

	data, err := json.Marshal(blockMsg)
	require.NoError(t, err)

	topic := TopicName("main", TopicBlock)
	mock.sendMessage(topic, data)

	// Both subscribers should receive the message
	var received1, received2 teranode.BlockMessage

	select {
	case received1 = <-blocks1:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for block1")
	}

	select {
	case received2 = <-blocks2:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for block2")
	}

	assert.Equal(t, uint32(99999), received1.Height)
	assert.Equal(t, uint32(99999), received2.Height)
	assert.Equal(t, "test-hash", received1.Hash)
	assert.Equal(t, "test-hash", received2.Hash)

	_ = client.Close()
}

func TestSubscribeBlocks_SlowSubscriber(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "main")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Don't read from this channel - simulates slow subscriber
	_ = client.SubscribeBlocks(ctx)

	// Create a second subscriber that we will read from
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	fastSub := client.SubscribeBlocks(ctx2)
	topic := TopicName("main", TopicBlock)

	// Send many messages - slow subscriber should be skipped, fast one should still work
	for i := range 150 {
		blockMsg := teranode.BlockMessage{Height: uint32(i)} //nolint:gosec // test code, i is small

		data, err := json.Marshal(blockMsg)
		require.NoError(t, err)

		mock.sendMessage(topic, data)
	}

	// Fast subscriber should receive messages (some may be dropped due to channel buffer)
	receivedCount := 0
	timeout := time.After(2 * time.Second)

readLoop:
	for {
		select {
		case <-fastSub:
			receivedCount++
		case <-timeout:
			break readLoop
		default:
			// No more messages immediately available
			time.Sleep(10 * time.Millisecond)

			select {
			case <-fastSub:
				receivedCount++
			case <-time.After(100 * time.Millisecond):
				break readLoop
			}
		}
	}

	// Should have received at least some messages
	assert.Positive(t, receivedCount, "fast subscriber should receive messages even with slow subscriber")

	_ = client.Close()
}

func TestSubscribeSubtrees(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subtrees := client.SubscribeSubtrees(ctx)
	require.NotNil(t, subtrees)

	// Send a test message
	msg := teranode.SubtreeMessage{
		Hash: "subtree-hash",
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	topic := TopicName("test", TopicSubtree)
	mock.sendMessage(topic, data)

	select {
	case received := <-subtrees:
		assert.Equal(t, "subtree-hash", received.Hash)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for subtree message")
	}

	_ = client.Close()
}

func TestSubscribeRejectedTxs(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "stn")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rejected := client.SubscribeRejectedTxs(ctx)
	require.NotNil(t, rejected)

	// Send a test message
	msg := teranode.RejectedTxMessage{
		TxID:   "rejected-tx-id",
		Reason: "double-spend",
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	topic := TopicName("stn", TopicRejectedTx)
	mock.sendMessage(topic, data)

	select {
	case received := <-rejected:
		assert.Equal(t, "rejected-tx-id", received.TxID)
		assert.Equal(t, "double-spend", received.Reason)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for rejected tx message")
	}

	_ = client.Close()
}

func TestSubscribeNodeStatus(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, NetworkTeratestnet)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	status := client.SubscribeNodeStatus(ctx)
	require.NotNil(t, status)

	// Send a test message
	msg := teranode.NodeStatusMessage{
		PeerID: "node-peer-id",
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	topic := TopicName(NetworkTeratestnet, TopicNodeStatus)
	mock.sendMessage(topic, data)

	select {
	case received := <-status:
		assert.Equal(t, "node-peer-id", received.PeerID)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for node status message")
	}

	_ = client.Close()
}

func TestSubscribe_OnlyStartsOnce(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "main")

	ctx := context.Background()

	// Subscribe multiple times
	_ = client.SubscribeBlocks(ctx)
	_ = client.SubscribeBlocks(ctx)
	_ = client.SubscribeBlocks(ctx)

	client.mu.RLock()
	started := client.blockStarted
	subCount := len(client.blockSubs)
	client.mu.RUnlock()

	assert.True(t, started, "block topic should be started")
	assert.Equal(t, 3, subCount, "should have 3 subscribers")

	// Verify only one subscription was made to the underlying msgbus
	mock.mu.Lock()
	subscriptionCount := len(mock.subscriptions)
	mock.mu.Unlock()

	assert.Equal(t, 1, subscriptionCount, "should only have one underlying subscription")

	_ = client.Close()
}

func TestGetPeers(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	mock.peers = []msgbus.PeerInfo{
		{ID: "peer1"},
		{ID: "peer2"},
	}
	client := createTestClient(mock, "main")

	peers := client.GetPeers()
	require.Len(t, peers, 2)
	assert.Equal(t, "peer1", peers[0].ID)
	assert.Equal(t, "peer2", peers[1].ID)
}

func TestFanOut_InvalidJSON(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "main")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	blocks := client.SubscribeBlocks(ctx)

	// Send invalid JSON
	topic := TopicName("main", TopicBlock)
	mock.sendMessage(topic, []byte("not-valid-json"))

	// Send a valid message after
	validMsg := teranode.BlockMessage{Height: 42}

	data, err := json.Marshal(validMsg)
	require.NoError(t, err)

	mock.sendMessage(topic, data)

	// Should receive the valid message (invalid one is skipped with error log)
	select {
	case received := <-blocks:
		assert.Equal(t, uint32(42), received.Height)
	case <-time.After(time.Second):
		t.Fatal("timeout - valid message should still be received after invalid one")
	}

	_ = client.Close()
}

func TestConcurrentSubscriptions(t *testing.T) {
	mock := newMockMsgBusClient("test-id")
	client := createTestClient(mock, "main")

	var wg sync.WaitGroup

	ctx := context.Background()
	subscriberCount := 10

	var received atomic.Int32

	// Create many subscribers concurrently
	for range subscriberCount {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ch := client.SubscribeBlocks(ctx)
			// Read one message
			select {
			case <-ch:
				received.Add(1)
			case <-time.After(2 * time.Second):
			}
		}()
	}

	// Wait for all subscribers to be registered
	time.Sleep(100 * time.Millisecond)

	// Send a message
	blockMsg := teranode.BlockMessage{Height: 1}

	data, err := json.Marshal(blockMsg)
	require.NoError(t, err)

	topic := TopicName("main", TopicBlock)
	mock.sendMessage(topic, data)

	// Wait for all goroutines
	wg.Wait()

	// All subscribers should have received the message
	assert.Equal(t, int32(subscriberCount), received.Load())

	_ = client.Close()
}
