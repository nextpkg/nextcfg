package configmap

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeMap(kv map[string]string) map[string]interface{} {
	data := make(map[string]interface{})

	for k, v := range kv {
		data[k] = make(map[string]interface{})

		mp := make(map[string]interface{})

		vals := strings.Split(v, "\n")
		for _, h := range vals {
			m, n := split(h, "=")
			mp[m] = n
		}

		data[k] = mp
	}

	return data
}

var localCfg = os.Getenv("HOME") + "/.kube/config"

func TestGetClient(t *testing.T) {
	if tr := os.Getenv("TRAVIS"); len(tr) > 0 {
		return
	}

	tt := []struct {
		name    string
		cfgPath string
		error   string
		assert  string
	}{
		{name: "fail loading incluster kubeconfig, from external",
			error: "unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST " +
				"and KUBERNETES_SERVICE_PORT must be defined", cfgPath: localCfg},
		{name: "fail loading external kubeconfig", cfgPath: "/invalid/path",
			error: "stat /invalid/path: no such file or directory"},
		{name: "loading an incluster kubeconfig", cfgPath: localCfg, error: "",
			assert: "open /var/run/secrets/kubernetes.io/serviceaccount/token: no such file or directory"},
		{name: "loading kubeconfig from external", cfgPath: localCfg},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getClient(tc.cfgPath)
			if err != nil {
				if err.Error() == tc.error {
					return
				}
				if err.Error() == tc.assert {
					return
				}

				t.Errorf("found an unhandled error: %v", err)
			}
		})
	}
}

func TestMakeMap(t *testing.T) {
	at := require.New(t)

	if tr := os.Getenv("TRAVIS"); len(tr) > 0 {
		return
	}

	tt := []struct {
		name  string
		din   map[string]string
		dout  map[string]interface{}
		jdout []byte
	}{
		{
			name: "simple valid data",
			din: map[string]string{
				"foo": "bar=bazz",
			},
			dout: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "bazz",
				},
			},
			jdout: []byte(`{"foo":{"bar":"bazz"}}`),
		},
		{
			name: "complex valid data",
			din: map[string]string{
				"mongodb": "host=127.0.0.1\nport=27017\nuser=user\npassword=password",
				"config":  "host=0.0.0.0\nport=1337",
				"redis":   "url=redis://127.0.0.1:6379/db01",
				"sql":     "username=user\npassword=password=1234",
			},
			dout: map[string]interface{}{
				"mongodb": map[string]interface{}{
					"host":     "127.0.0.1",
					"port":     "27017",
					"user":     "user",
					"password": "password",
				},
				"config": map[string]interface{}{
					"host": "0.0.0.0",
					"port": "1337",
				},
				"redis": map[string]interface{}{
					"url": "redis://127.0.0.1:6379/db01",
				},
				"sql": map[string]interface{}{
					"username": "user",
					"password": "password=1234",
				},
			},
			jdout: []byte(`{"config":{"host":"0.0.0.0","port":"1337"},"mongodb":{"host":"127.0.0.1","password":"password","port":"27017","user":"user"},"redis":{"url":"redis://127.0.0.1:6379/db01"},"sql":{"password":"password=1234","username":"user"}}`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dout := makeMap(tc.din)
			jdout, _ := json.Marshal(dout)

			at.Equal(tc.dout, dout)
			at.Equal(string(tc.jdout), string(jdout))
		})
	}
}

// TestConfigmap_Read 需在容器内执行
func TestConfigmap_Read(t *testing.T) {
	at := assert.New(t)

	if tr := os.Getenv("TRAVIS"); len(tr) > 0 {
		return
	}

	data := []byte(`host=127.0.0.1
port=27017
user=user
password=password`)
	tt := []struct {
		name      string
		sname     string
		namespace string
	}{
		{name: "read data with source default values", sname: DefaultGroup, namespace: DefaultNamespace},
		{name: "read data with source with custom configmap name", sname: "config", namespace: DefaultNamespace},
		{name: "read data with source with custom namespace", sname: DefaultGroup, namespace: "kube-public"},
		{name: "read data with source with custom configmap name and namespace",
			sname: "config", namespace: "kube-public"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			source := NewSource(
				WithGroup(tc.sname),
				WithName("mongodb"),
				WithNamespace(tc.namespace),
				WithConfigPath(localCfg),
			)

			r, err := source.Read()
			at.Nil(err)

			at.Equal(data, r.Data)
		})
	}
}

func TestConfigmap_String(t *testing.T) {
	if tr := os.Getenv("TRAVIS"); len(tr) > 0 {
		return
	}

	source := NewSource()

	if source.String() != "configmap" {
		t.Errorf("expecting to get %v and instead got %v", "configmap", source)
	}
}
