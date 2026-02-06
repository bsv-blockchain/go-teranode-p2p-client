// Package p2p provides a high-level client for connecting to the Teranode P2P network.
package p2p

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	msgbus "github.com/bsv-blockchain/go-p2p-message-bus"
	teranode "github.com/bsv-blockchain/teranode/services/p2p"
)

// Client provides a high-level interface for subscribing to Teranode P2P messages.
// It supports multiple subscribers per topic via internal fan-out.
type Client struct {
	msgbus  msgbus.Client
	logger  *slog.Logger
	network string // P2P network name (e.g., "mainnet", "testnet")

	// Fan-out support: multiple subscribers per topic
	mu           sync.RWMutex
	blockSubs    []chan teranode.BlockMessage
	subtreeSubs  []chan teranode.SubtreeMessage
	rejectedSubs []chan teranode.RejectedTxMessage
	statusSubs   []chan teranode.NodeStatusMessage

	// Track if we've started listening to each topic
	blockStarted    bool
	subtreeStarted  bool
	rejectedStarted bool
	statusStarted   bool

	done     chan struct{}  // signals shutdown to fan-out goroutines
	fanOutWg sync.WaitGroup // tracks running fan-out goroutines
}

// NewClient creates a new Teranode P2P client.
// Prefer using Config.Initialize() which applies defaults before creating the client.
func NewClient(cfg Config) (*Client, error) {
	p2pClient, err := msgbus.NewClient(cfg.MsgBus)
	if err != nil {
		return nil, err
	}

	return &Client{
		msgbus:  p2pClient,
		logger:  slog.Default(),
		network: cfg.Network,
		done:    make(chan struct{}),
	}, nil
}

// GetID returns this client's peer ID.
func (c *Client) GetID() string {
	return c.msgbus.GetID()
}

// GetNetwork returns the network this client is connected to.
func (c *Client) GetNetwork() string {
	return c.network
}

// Close shuts down the P2P client and closes all subscriber channels.
func (c *Client) Close() error {
	c.mu.Lock()

	// Signal shutdown first (before closing channels)
	select {
	case <-c.done:
		// Already closed
	default:
		close(c.done)
	}

	c.mu.Unlock()

	// Close msgbus to stop raw channels, causing fan-out goroutines to exit
	err := c.msgbus.Close()

	// Wait for all fan-out goroutines to finish — no goroutine is sending
	// to subscriber channels after this returns
	c.fanOutWg.Wait()

	// Now safe to close subscriber channels
	c.mu.Lock()

	for _, ch := range c.blockSubs {
		close(ch)
	}

	c.blockSubs = nil

	for _, ch := range c.subtreeSubs {
		close(ch)
	}

	c.subtreeSubs = nil

	for _, ch := range c.rejectedSubs {
		close(ch)
	}

	c.rejectedSubs = nil

	for _, ch := range c.statusSubs {
		close(ch)
	}

	c.statusSubs = nil
	c.mu.Unlock()

	return err
}

// GetPeers returns information about all known peers.
func (c *Client) GetPeers() []msgbus.PeerInfo {
	return c.msgbus.GetPeers()
}

// SubscribeBlocks subscribes to block announcements.
// Multiple callers can subscribe; each receives all messages (fan-out).
// The returned channel is closed when the client is closed or context is canceled.
func (c *Client) SubscribeBlocks(ctx context.Context) <-chan teranode.BlockMessage {
	out := make(chan teranode.BlockMessage, 100)

	c.mu.Lock()
	c.blockSubs = append(c.blockSubs, out)

	// Start the topic listener only once
	if !c.blockStarted {
		c.blockStarted = true
		topic := TopicName(c.network, TopicBlock)

		rawChan := c.msgbus.Subscribe(topic)

		c.fanOutWg.Add(1)

		go func() {
			defer c.fanOutWg.Done()

			c.fanOutBlocks(rawChan, topic)
		}()
	}

	c.mu.Unlock()

	// Handle context cancellation
	go func() {
		<-ctx.Done()

		c.mu.Lock()

		for i, ch := range c.blockSubs {
			if ch == out {
				c.blockSubs = append(c.blockSubs[:i], c.blockSubs[i+1:]...)

				close(out)

				break
			}
		}

		c.mu.Unlock()
	}()

	return out
}

// SubscribeSubtrees subscribes to subtree (transaction batch) announcements.
// Multiple callers can subscribe; each receives all messages (fan-out).
func (c *Client) SubscribeSubtrees(ctx context.Context) <-chan teranode.SubtreeMessage {
	out := make(chan teranode.SubtreeMessage, 100)

	c.mu.Lock()
	c.subtreeSubs = append(c.subtreeSubs, out)

	if !c.subtreeStarted {
		c.subtreeStarted = true
		topic := TopicName(c.network, TopicSubtree)

		rawChan := c.msgbus.Subscribe(topic)

		c.fanOutWg.Add(1)

		go func() {
			defer c.fanOutWg.Done()

			c.fanOutSubtrees(rawChan, topic)
		}()
	}

	c.mu.Unlock()

	go func() {
		<-ctx.Done()

		c.mu.Lock()

		for i, ch := range c.subtreeSubs {
			if ch == out {
				c.subtreeSubs = append(c.subtreeSubs[:i], c.subtreeSubs[i+1:]...)

				close(out)

				break
			}
		}

		c.mu.Unlock()
	}()

	return out
}

// SubscribeRejectedTxs subscribes to rejected transaction notifications.
// Multiple callers can subscribe; each receives all messages (fan-out).
func (c *Client) SubscribeRejectedTxs(ctx context.Context) <-chan teranode.RejectedTxMessage {
	out := make(chan teranode.RejectedTxMessage, 100)

	c.mu.Lock()
	c.rejectedSubs = append(c.rejectedSubs, out)

	if !c.rejectedStarted {
		c.rejectedStarted = true
		topic := TopicName(c.network, TopicRejectedTx)

		rawChan := c.msgbus.Subscribe(topic)

		c.fanOutWg.Add(1)

		go func() {
			defer c.fanOutWg.Done()

			c.fanOutRejectedTxs(rawChan, topic)
		}()
	}

	c.mu.Unlock()

	go func() {
		<-ctx.Done()

		c.mu.Lock()

		for i, ch := range c.rejectedSubs {
			if ch == out {
				c.rejectedSubs = append(c.rejectedSubs[:i], c.rejectedSubs[i+1:]...)

				close(out)

				break
			}
		}

		c.mu.Unlock()
	}()

	return out
}

// SubscribeNodeStatus subscribes to node status updates.
// Multiple callers can subscribe; each receives all messages (fan-out).
func (c *Client) SubscribeNodeStatus(ctx context.Context) <-chan teranode.NodeStatusMessage {
	out := make(chan teranode.NodeStatusMessage, 100)

	c.mu.Lock()
	c.statusSubs = append(c.statusSubs, out)

	if !c.statusStarted {
		c.statusStarted = true
		topic := TopicName(c.network, TopicNodeStatus)

		rawChan := c.msgbus.Subscribe(topic)

		c.fanOutWg.Add(1)

		go func() {
			defer c.fanOutWg.Done()

			c.fanOutNodeStatus(rawChan, topic)
		}()
	}

	c.mu.Unlock()

	go func() {
		<-ctx.Done()

		c.mu.Lock()

		for i, ch := range c.statusSubs {
			if ch == out {
				c.statusSubs = append(c.statusSubs[:i], c.statusSubs[i+1:]...)

				close(out)

				break
			}
		}

		c.mu.Unlock()
	}()

	return out
}

// Fan-out goroutines - one per topic type

// isShuttingDown checks if the client is shutting down without blocking.
func (c *Client) isShuttingDown() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

// getBlockSubs returns a snapshot of current block subscribers, or nil if shutting down.
func (c *Client) getBlockSubs() []chan teranode.BlockMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isShuttingDown() {
		return nil
	}

	subs := make([]chan teranode.BlockMessage, len(c.blockSubs))
	copy(subs, c.blockSubs)

	return subs
}

// getSubtreeSubs returns a snapshot of current subtree subscribers, or nil if shutting down.
func (c *Client) getSubtreeSubs() []chan teranode.SubtreeMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isShuttingDown() {
		return nil
	}

	subs := make([]chan teranode.SubtreeMessage, len(c.subtreeSubs))
	copy(subs, c.subtreeSubs)

	return subs
}

// getRejectedSubs returns a snapshot of current rejected tx subscribers, or nil if shutting down.
func (c *Client) getRejectedSubs() []chan teranode.RejectedTxMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isShuttingDown() {
		return nil
	}

	subs := make([]chan teranode.RejectedTxMessage, len(c.rejectedSubs))
	copy(subs, c.rejectedSubs)

	return subs
}

// getStatusSubs returns a snapshot of current status subscribers, or nil if shutting down.
func (c *Client) getStatusSubs() []chan teranode.NodeStatusMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isShuttingDown() {
		return nil
	}

	subs := make([]chan teranode.NodeStatusMessage, len(c.statusSubs))
	copy(subs, c.statusSubs)

	return subs
}

func (c *Client) fanOutBlocks(rawChan <-chan msgbus.Message, topic string) {
	for msg := range rawChan {
		if c.isShuttingDown() {
			return
		}

		var typed teranode.BlockMessage

		if err := json.Unmarshal(msg.Data, &typed); err != nil {
			c.logger.Error("failed to unmarshal block message",
				slog.String("topic", topic),
				slog.String("error", err.Error()))

			continue
		}

		subs := c.getBlockSubs()
		if subs == nil {
			return
		}

		c.sendToBlockSubs(subs, typed)
	}
}

func (c *Client) sendToBlockSubs(subs []chan teranode.BlockMessage, msg teranode.BlockMessage) {
	for _, ch := range subs {
		select {
		case <-c.done:
			return
		case ch <- msg:
		default:
			// Subscriber is slow, skip to avoid blocking
		}
	}
}

func (c *Client) fanOutSubtrees(rawChan <-chan msgbus.Message, topic string) {
	for msg := range rawChan {
		if c.isShuttingDown() {
			return
		}

		var typed teranode.SubtreeMessage

		if err := json.Unmarshal(msg.Data, &typed); err != nil {
			c.logger.Error("failed to unmarshal subtree message",
				slog.String("topic", topic),
				slog.String("error", err.Error()))

			continue
		}

		subs := c.getSubtreeSubs()
		if subs == nil {
			return
		}

		c.sendToSubtreeSubs(subs, typed)
	}
}

func (c *Client) sendToSubtreeSubs(subs []chan teranode.SubtreeMessage, msg teranode.SubtreeMessage) {
	for _, ch := range subs {
		select {
		case <-c.done:
			return
		case ch <- msg:
		default:
			// Subscriber is slow, skip to avoid blocking
		}
	}
}

func (c *Client) fanOutRejectedTxs(rawChan <-chan msgbus.Message, topic string) {
	for msg := range rawChan {
		if c.isShuttingDown() {
			return
		}

		var typed teranode.RejectedTxMessage

		if err := json.Unmarshal(msg.Data, &typed); err != nil {
			c.logger.Error("failed to unmarshal rejected-tx message",
				slog.String("topic", topic),
				slog.String("error", err.Error()))

			continue
		}

		subs := c.getRejectedSubs()
		if subs == nil {
			return
		}

		c.sendToRejectedSubs(subs, typed)
	}
}

func (c *Client) sendToRejectedSubs(subs []chan teranode.RejectedTxMessage, msg teranode.RejectedTxMessage) {
	for _, ch := range subs {
		select {
		case <-c.done:
			return
		case ch <- msg:
		default:
			// Subscriber is slow, skip to avoid blocking
		}
	}
}

func (c *Client) fanOutNodeStatus(rawChan <-chan msgbus.Message, topic string) {
	for msg := range rawChan {
		if c.isShuttingDown() {
			return
		}

		var typed teranode.NodeStatusMessage

		if err := json.Unmarshal(msg.Data, &typed); err != nil {
			c.logger.Error("failed to unmarshal node_status message",
				slog.String("topic", topic),
				slog.String("error", err.Error()))

			continue
		}

		subs := c.getStatusSubs()
		if subs == nil {
			return
		}

		c.sendToStatusSubs(subs, typed)
	}
}

func (c *Client) sendToStatusSubs(subs []chan teranode.NodeStatusMessage, msg teranode.NodeStatusMessage) {
	for _, ch := range subs {
		select {
		case <-c.done:
			return
		case ch <- msg:
		default:
			// Subscriber is slow, skip to avoid blocking
		}
	}
}
