package p2p

import (
	"encoding/json"
	"strings"
	"testing"

	msgbus "github.com/bsv-blockchain/go-p2p-message-bus"
	teranode "github.com/bsv-blockchain/teranode/services/p2p"
)

// FuzzBlockMessageUnmarshal tests JSON unmarshaling of BlockMessage.
// This mirrors the parsing in client.go:266 (fanOutBlocks).
func FuzzBlockMessageUnmarshal(f *testing.F) {
	// Seed corpus: valid inputs
	f.Add([]byte(`{"height":0}`))
	f.Add([]byte(`{"height":12345,"hash":"abc123"}`))
	f.Add([]byte(`{"height":9223372036854775807}`)) // max int64

	// Edge cases
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`123`))
	f.Add([]byte(`true`))
	f.Add([]byte(`false`))

	// Malformed JSON
	f.Add([]byte(`{`))
	f.Add([]byte(`{"height":`))
	f.Add([]byte(`{"height":}`))
	f.Add([]byte(`{height:1}`))
	f.Add([]byte(``))
	f.Add([]byte(`{"height": "not a number"}`))
	f.Add([]byte(`{"height": null}`))
	f.Add([]byte(`{"height": -1}`))
	f.Add([]byte(`{"height": 1.5}`))
	f.Add([]byte(`{"height": 1e100}`))

	// Unicode and special characters
	f.Add([]byte(`{"hash": "\u0000"}`))
	f.Add([]byte(`{"hash": "\n\r\t"}`))
	f.Add([]byte("{\"hash\": \"\x00\x01\x02\"}"))

	f.Fuzz(func(_ *testing.T, data []byte) {
		var msg teranode.BlockMessage

		_ = json.Unmarshal(data, &msg) // Only panics are failures
	})
}

// FuzzSubtreeMessageUnmarshal tests JSON unmarshaling of SubtreeMessage.
// This mirrors the parsing in client.go:293 (fanOutSubtrees).
func FuzzSubtreeMessageUnmarshal(f *testing.F) {
	// Seed corpus: valid inputs
	f.Add([]byte(`{"txid":"abc123"}`))
	f.Add([]byte(`{"txid":"","count":0}`))
	f.Add([]byte(`{"count":1000000}`))

	// Edge cases
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`123`))

	// Malformed JSON
	f.Add([]byte(`{`))
	f.Add([]byte(`{"txid":`))
	f.Add([]byte(``))
	f.Add([]byte(`{"txid": 123}`))
	f.Add([]byte(`{"txid": null}`))
	f.Add([]byte(`{"txid": []}`))

	// Large values
	f.Add([]byte(`{"count": 9223372036854775807}`))
	f.Add([]byte(`{"count": -9223372036854775808}`))

	// Unicode
	f.Add([]byte(`{"txid": "\u0000\u0001"}`))

	f.Fuzz(func(_ *testing.T, data []byte) {
		var msg teranode.SubtreeMessage

		_ = json.Unmarshal(data, &msg) // Only panics are failures
	})
}

// FuzzRejectedTxMessageUnmarshal tests JSON unmarshaling of RejectedTxMessage.
// This mirrors the parsing in client.go:319 (fanOutRejectedTxs).
func FuzzRejectedTxMessageUnmarshal(f *testing.F) {
	// Seed corpus: valid inputs
	f.Add([]byte(`{"txid":"abc123","reason":"invalid"}`))
	f.Add([]byte(`{"txid":"","reason":""}`))
	f.Add([]byte(`{"reason":"double-spend"}`))

	// Edge cases
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`123`))

	// Malformed JSON
	f.Add([]byte(`{`))
	f.Add([]byte(`{"txid":`))
	f.Add([]byte(``))
	f.Add([]byte(`{"txid": 123, "reason": 456}`))
	f.Add([]byte(`{"txid": null, "reason": null}`))

	// Long strings
	longStr := strings.Repeat("a", 10000)
	f.Add([]byte(`{"reason": "` + longStr + `"}`))

	// Unicode and special characters
	f.Add([]byte(`{"reason": "\u0000\n\r\t"}`))
	f.Add([]byte(`{"reason": "emoji: 🚀"}`))

	f.Fuzz(func(_ *testing.T, data []byte) {
		var msg teranode.RejectedTxMessage

		_ = json.Unmarshal(data, &msg) // Only panics are failures
	})
}

// FuzzNodeStatusMessageUnmarshal tests JSON unmarshaling of NodeStatusMessage.
// This mirrors the parsing in client.go:345 (fanOutNodeStatus).
func FuzzNodeStatusMessageUnmarshal(f *testing.F) {
	// Seed corpus: valid inputs
	f.Add([]byte(`{"peer_id":"12D3KooW..."}`))
	f.Add([]byte(`{"peer_id":"","connected":true}`))
	f.Add([]byte(`{"connected":false}`))

	// Edge cases
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`123`))
	f.Add([]byte(`true`))

	// Malformed JSON
	f.Add([]byte(`{`))
	f.Add([]byte(`{"peer_id":`))
	f.Add([]byte(``))
	f.Add([]byte(`{"connected": "yes"}`))
	f.Add([]byte(`{"connected": 1}`))
	f.Add([]byte(`{"connected": null}`))

	// Nested objects
	f.Add([]byte(`{"status": {"key": "value"}}`))
	f.Add([]byte(`{"metrics": [1, 2, 3]}`))

	// Unicode
	f.Add([]byte(`{"peer_id": "\u0000"}`))

	f.Fuzz(func(_ *testing.T, data []byte) {
		var msg teranode.NodeStatusMessage

		_ = json.Unmarshal(data, &msg) // Only panics are failures
	})
}

// FuzzPrivateKeyFromHex tests msgbus.PrivateKeyFromHex with arbitrary input.
// This mirrors the parsing in config.go:40 (LoadOrGeneratePrivateKey).
func FuzzPrivateKeyFromHex(f *testing.F) {
	// Seed corpus: valid hex (wrong length but valid hex chars)
	f.Add("0123456789abcdef")
	f.Add("ABCDEF0123456789")

	// Empty and short inputs
	f.Add("")
	f.Add("0")
	f.Add("00")
	f.Add("ff")

	// Invalid hex characters
	f.Add("ghijklmn")
	f.Add("0x1234")
	f.Add("zzzz")

	// Whitespace
	f.Add(" ")
	f.Add("\t\n\r")
	f.Add("  1234  ")

	// Special characters
	f.Add("!@#$%^&*()")
	f.Add("-")
	f.Add("+")

	// Unicode
	f.Add("\u0000")
	f.Add("🚀")

	// Long inputs
	f.Add(strings.Repeat("a", 1000))
	f.Add(strings.Repeat("0", 64))  // 32 bytes in hex
	f.Add(strings.Repeat("f", 128)) // 64 bytes in hex

	// Valid-looking but wrong length
	f.Add(strings.Repeat("0", 63))
	f.Add(strings.Repeat("0", 65))

	f.Fuzz(func(_ *testing.T, data string) {
		_, _ = msgbus.PrivateKeyFromHex(data) // Only panics are failures
	})
}

// FuzzTopicName tests the TopicName function with arbitrary input.
// This tests topics.go:39 (TopicName).
func FuzzTopicName(f *testing.F) {
	// Seed corpus: valid networks and topics
	f.Add("main", "block")
	f.Add("test", "subtree")
	f.Add("stn", "rejected-tx")
	f.Add("mainnet", "node_status")
	f.Add("testnet", "block")
	f.Add("teratestnet", "subtree")
	f.Add("teratest", "block")

	// Empty inputs
	f.Add("", "")
	f.Add("main", "")
	f.Add("", "block")

	// Unknown networks/topics
	f.Add("unknown", "unknown")
	f.Add("production", "status")
	f.Add("dev", "test")

	// Special characters
	f.Add("/", "/")
	f.Add("../", "../")
	f.Add("main\n", "block\n")
	f.Add("main\x00", "block\x00")

	// Unicode (using escape sequences to avoid gosmopolitan lint)
	f.Add("\u4e3b\u7f51", "\u533a\u5757") // Chinese characters
	f.Add("\U0001F680", "\U0001F525")     // Emoji: rocket, fire

	// Long inputs
	f.Add(strings.Repeat("a", 1000), strings.Repeat("b", 1000))

	// Whitespace
	f.Add(" main ", " block ")
	f.Add("\t", "\n")

	f.Fuzz(func(t *testing.T, network, topic string) {
		result := TopicName(network, topic)

		// Invariant: result must always start with protocol prefix
		if !strings.HasPrefix(result, protocolPrefix) {
			t.Errorf("TopicName(%q, %q) = %q, does not have prefix %q",
				network, topic, result, protocolPrefix)
		}

		// Invariant: result format is always "prefix/network-topic"
		if !strings.Contains(result, "-") {
			t.Errorf("TopicName(%q, %q) = %q, missing hyphen separator",
				network, topic, result)
		}
	})
}

// FuzzResolveStoragePath tests the resolveStoragePath function with arbitrary input.
// This tests config.go:70 (resolveStoragePath).
func FuzzResolveStoragePath(f *testing.F) {
	// Seed corpus: valid paths
	f.Add("")
	f.Add(".")
	f.Add("..")
	f.Add("./data")
	f.Add("../data")
	f.Add("/tmp/data")
	f.Add("/absolute/path")

	// Home directory expansion
	f.Add("~/")
	f.Add("~/.teranode")
	f.Add("~/subdir/data")

	// Edge cases with ~
	f.Add("~")     // Just tilde, not ~/
	f.Add("~user") // Not supported, should pass through
	f.Add("~~")    // Double tilde
	f.Add("~/../..")

	// Special characters
	f.Add("/path with spaces/")
	f.Add("/path\twith\ttabs/")
	f.Add("/path\nwith\nnewlines/")
	f.Add("/path\x00with\x00nulls/")

	// Unicode paths (using escape sequences to avoid gosmopolitan lint)
	f.Add("/\u8def\u5f84/\u6570\u636e") // Chinese path
	f.Add("/\u043f\u0443\u0442\u044c")  // Cyrillic path
	f.Add("/\U0001F680/\U0001F4C1")     // Emoji path

	// Long paths
	f.Add("/" + strings.Repeat("a", 1000))
	f.Add("~/" + strings.Repeat("b", 1000))

	// Relative paths
	f.Add("relative/path")
	f.Add("./relative/./path")
	f.Add("../../../escape")

	// Path traversal attempts
	f.Add("~/../../etc/passwd")
	f.Add("/tmp/../etc/passwd")

	f.Fuzz(func(t *testing.T, path string) {
		result, err := resolveStoragePath(path)
		if err != nil {
			return // Errors are acceptable, only panics are failures
		}

		// Validate invariants based on input type
		validateResolveStoragePathInvariants(t, path, result)
	})
}

// validateResolveStoragePathInvariants checks the invariants for resolveStoragePath results.
func validateResolveStoragePathInvariants(t *testing.T, path, result string) {
	t.Helper()

	startsWithTilde := len(path) >= 2 && path[:2] == "~/"

	switch {
	case path == "":
		// Empty input should use default path
		if !strings.HasSuffix(result, ".teranode-p2p") {
			t.Errorf("resolveStoragePath(%q) = %q, expected suffix .teranode-p2p", path, result)
		}
	case startsWithTilde:
		// Tilde paths should be expanded (result should not start with ~/)
		if len(result) >= 2 && result[:2] == "~/" {
			t.Errorf("resolveStoragePath(%q) = %q, should have expanded ~/", path, result)
		}
	default:
		// Other paths should pass through unchanged
		if result != path {
			t.Errorf("resolveStoragePath(%q) = %q, expected %q", path, result, path)
		}
	}
}
