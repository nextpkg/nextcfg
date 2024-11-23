package env

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/nextpkg/nextcfg/source"
	"github.com/stretchr/testify/require"
)

func TestEnv_Read(t *testing.T) {

	at := require.New(t)

	expected := map[string]map[string]string{
		"database": {
			"host":       "localhost",
			"password":   "password",
			"datasource": "user:password@tcp(localhost:port)/db?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}

	at.Nil(os.Setenv("DATABASE_HOST", "localhost"))
	at.Nil(os.Setenv("DATABASE_PASSWORD", "password"))
	at.Nil(os.Setenv("DATABASE_DATASOURCE",
		"user:password@tcp(localhost:port)/db?charset=utf8mb4&parseTime=True&loc=Local"))

	ns := NewSource(WithPrefix("DATABASE"))
	c, err := ns.Read()
	at.Nil(err)

	var actual map[string]interface{}
	at.Nil(json.Unmarshal(c.Data, &actual))

	actualDB := actual["database"].(map[string]interface{})

	for k, v := range expected["database"] {
		a := actualDB[k]
		at.Equal(v, a)
	}
}

func TestEnvVar_Prefixes(t *testing.T) {
	os.Setenv("APP_DATABASE_HOST", "localhost")
	os.Setenv("APP_DATABASE_PASSWORD", "password")
	os.Setenv("VAULT_ADDR", "vault:1337")
	os.Setenv("MICRO_REGISTRY", "mdns")

	var prefixTests = []struct {
		prefixOpts   []source.Option
		expectedKeys []string
	}{
		{[]source.Option{WithPrefix("APP", "MICRO")}, []string{"app", "micro"}},
		{[]source.Option{WithPrefix("MICRO"), WithStrippedPrefix("APP")}, []string{"database", "micro"}},
		{[]source.Option{WithPrefix("MICRO"), WithStrippedPrefix("APP")}, []string{"database", "micro"}},
	}

	for _, pt := range prefixTests {
		source := NewSource(pt.prefixOpts...)

		c, err := source.Read()
		if err != nil {
			t.Error(err)
		}

		var actual map[string]interface{}
		if err := json.Unmarshal(c.Data, &actual); err != nil {
			t.Error(err)
		}

		// assert other prefixes ignored
		if l := len(actual); l != len(pt.expectedKeys) {
			t.Errorf("expected %v top keys, got %v", len(pt.expectedKeys), l)
		}

		for _, k := range pt.expectedKeys {
			if !containsKey(actual, k) {
				t.Errorf("expected key %v, not found", k)
			}
		}
	}
}

func TestEnvVar_WatchNextNoOpsUntilStop(t *testing.T) {
	src := NewSource(WithStrippedPrefix("GOMICRO_"))
	w, err := src.Watch()
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		w.Stop()
	}()

	if _, err := w.Next(); err != source.ErrWatcherStopped {
		t.Errorf("expected watcher stopped error, got %v", err)
	}
}

func containsKey(m map[string]interface{}, s string) bool {
	for k := range m {
		if k == s {
			return true
		}
	}
	return false
}
