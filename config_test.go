package p2p

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadOrGeneratePrivateKey_Generate(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Generate a new key
	privKey, err := LoadOrGeneratePrivateKey(tempDir)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	// Verify key file was created
	keyPath := filepath.Join(tempDir, "p2p_key.hex")
	_, err = os.Stat(keyPath)
	require.NoError(t, err, "key file should exist")

	// Verify file permissions (should be 0600)
	info, err := os.Stat(keyPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestLoadOrGeneratePrivateKey_Load(t *testing.T) {
	tempDir := t.TempDir()

	// Generate a key first
	originalKey, err := LoadOrGeneratePrivateKey(tempDir)
	require.NoError(t, err)

	// Load the same key
	loadedKey, err := LoadOrGeneratePrivateKey(tempDir)
	require.NoError(t, err)

	// The keys should be the same (same raw bytes)
	originalRaw, err := originalKey.Raw()
	require.NoError(t, err)
	loadedRaw, err := loadedKey.Raw()
	require.NoError(t, err)

	assert.Equal(t, originalRaw, loadedRaw, "loaded key should match original")
}

func TestLoadOrGeneratePrivateKey_InvalidKey(t *testing.T) {
	tempDir := t.TempDir()

	// Write an invalid key file
	keyPath := filepath.Join(tempDir, "p2p_key.hex")
	err := os.WriteFile(keyPath, []byte("invalid-hex-data!!!"), 0o600)
	require.NoError(t, err)

	// Try to load the invalid key
	_, err = LoadOrGeneratePrivateKey(tempDir)
	require.Error(t, err, "should fail on invalid key")
	assert.Contains(t, err.Error(), "failed to parse private key")
}

func TestResolveStoragePath_Default(t *testing.T) {
	resolved, err := resolveStoragePath("")
	require.NoError(t, err)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If home dir fails, should fallback to ./.teranode-p2p
		assert.Equal(t, "./.teranode-p2p", resolved)
	} else {
		expected := filepath.Join(homeDir, ".teranode-p2p")
		assert.Equal(t, expected, resolved)
	}
}

func TestResolveStoragePath_Tilde(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	resolved, err := resolveStoragePath("~/custom-p2p-dir")
	require.NoError(t, err)

	expected := filepath.Join(homeDir, "custom-p2p-dir")
	assert.Equal(t, expected, resolved)
}

func TestResolveStoragePath_Absolute(t *testing.T) {
	absPath := "/tmp/test-p2p-storage"
	resolved, err := resolveStoragePath(absPath)
	require.NoError(t, err)
	assert.Equal(t, absPath, resolved)
}

func TestResolveStoragePath_Relative(t *testing.T) {
	relPath := "./data/p2p"
	resolved, err := resolveStoragePath(relPath)
	require.NoError(t, err)
	assert.Equal(t, relPath, resolved)
}

func TestConfig_SetDefaults(t *testing.T) {
	v := viper.New()
	cfg := &Config{}

	cfg.SetDefaults(v, "")

	assert.Equal(t, "main", v.GetString("network"))
	assert.Equal(t, "off", v.GetString("msgbus.dht_mode"))
	assert.Equal(t, 35, v.GetInt("msgbus.max_connections"))
	assert.Equal(t, 25, v.GetInt("msgbus.min_connections"))
	assert.Equal(t, 20*time.Second, v.GetDuration("msgbus.connection_grace_period"))
	assert.Equal(t, 24*time.Hour, v.GetDuration("msgbus.peer_cache_ttl"))
}

func TestConfig_SetDefaults_WithPrefix(t *testing.T) {
	v := viper.New()
	cfg := &Config{}

	cfg.SetDefaults(v, "p2p")

	assert.Equal(t, "main", v.GetString("p2p.network"))
	assert.Equal(t, "off", v.GetString("p2p.msgbus.dht_mode"))
	assert.Equal(t, 35, v.GetInt("p2p.msgbus.max_connections"))
}

func TestGetDefaultBootstrapPeers_Main(t *testing.T) {
	peers := getDefaultBootstrapPeers()

	mainPeers := peers["main"]
	require.Len(t, mainPeers, 1)
	assert.Contains(t, mainPeers[0], "mainnet")
}

func TestGetDefaultBootstrapPeers_Test(t *testing.T) {
	peers := getDefaultBootstrapPeers()

	testPeers := peers["test"]
	require.Len(t, testPeers, 1)
	assert.Contains(t, testPeers[0], "testnet")
}

func TestGetDefaultBootstrapPeers_STN(t *testing.T) {
	peers := getDefaultBootstrapPeers()

	stnPeers := peers["stn"]
	require.Len(t, stnPeers, 1)
	assert.Contains(t, stnPeers[0], "teratestnet")
}

func TestGetDefaultBootstrapPeers_Teratestnet_NotIncluded(t *testing.T) {
	peers := getDefaultBootstrapPeers()

	// Teratestnet should not have default bootstrap peers
	teratestPeers := peers["teratestnet"]
	assert.Empty(t, teratestPeers, "teratestnet should not have default bootstrap peers")

	// Also check with short name
	teratestPeers2 := peers["teratest"]
	assert.Empty(t, teratestPeers2, "teratest should not have default bootstrap peers")
}

func TestGetDefaultBootstrapPeers_Regtest_NotIncluded(t *testing.T) {
	peers := getDefaultBootstrapPeers()

	regtestPeers := peers["regtest"]
	assert.Empty(t, regtestPeers, "regtest should not have default bootstrap peers")
}

func TestConfig_Initialize_SetsStoragePath(t *testing.T) {
	// This test only verifies the storage path resolution,
	// not the full initialization which requires network connectivity.
	tempDir := t.TempDir()

	cfg := &Config{
		Network:     "main",
		StoragePath: tempDir,
	}

	// We can't fully test Initialize without network, but we can test the
	// storage path resolution by checking it directly
	resolved, err := resolveStoragePath(cfg.StoragePath)
	require.NoError(t, err)
	assert.Equal(t, tempDir, resolved)
}

func TestBootstrapPeerAddresses(t *testing.T) {
	// Verify bootstrap peers are valid dnsaddr multiaddr format
	peers := getDefaultBootstrapPeers()

	for network, addrs := range peers {
		for i, addr := range addrs {
			assert.Contains(t, addr, "/dnsaddr/", "peer %d for %s should have dnsaddr", i, network)
			assert.Contains(t, addr, "bootstrap.teranode.bsvb.tech", "peer %d for %s should have bootstrap domain", i, network)
		}
	}
}
