# Kubernetes ConfigMap Source (configmap)

The configmap source reads config from a kubernetes configmap key/values

## Kubernetes ConfigMap Format

The configmap source expects keys under a namespace default to `default` and a confimap default to `config`

```shell
// we recommend to setup your variables from multiples files example:
kubectl create configmap config --namespace default --from-file=./testdata
```

## Kubernetes wrights

Since Kubernetes 1.9 the app must have wrights to be able to access configmaps. You must provide Role and RoleBinding so
that your app can access configmaps.

```bash
cat << EOF > cm-role.yaml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cm-role
  labels:
    app: tools-rbac
rules:
- apiGroups: [""]
  resources: ["configmaps"]
  verbs: ["get", "update", "list", "watch"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: global-rolebinding
  labels:
    app: tools-rbac
subjects:
- kind: Group
  name: system:serviceaccounts
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: cm-role
  apiGroup: ""
EOF
```

To configure your Kubernetes cluster just apply the file:

```bash
kubectl apply -n default -f cm-role.yaml
```

查看已生效的权限:

```bash
 kubectl get role cm-role -oyaml
```

## New Source

Specify source with data

```go
configmapSource := configmap.NewSource(
  // optionally specify name for ConfigMap; defaults config
  configmap.WithGroup("config").
	// optionally specify a namespace; default to default
	configmap.WithNamespace("kube-public"),
	// 分组内的Key
	configmap.WithName("example.yaml"),
  // optionally strip the provided path to a kube config file mostly used outside of a cluster, defaults to "" for in cluster support.
  configmap.WithConfigPath($HOME/.kube/config),
)
```

## Load Source

Load the source into config

```go
// Create new config
conf := nextcfg.NewConfig()

// Load file source
conf.Load(configmapSource)
```

## Running Go Tests

### Requirements

Have a kubernetes cluster running (external or minikube) have a valid `kubeconfig` file.

```shell
// Setup testing configmaps feel free to remove them after testing.
cd source/configmap

# 给default和kube-public授权（测试需要）
kubectl apply -n default -f role.yaml
kubectl apply -n kube-public -f role.yaml

# 给default和kube-public命名空间创建配置项
kubectl create configmap config --from-file=./testdata
kubectl create configmap config --from-file=./testdata --namespace kube-public

go test -v -run  -cover
```

测试后清理环境

```shell
# To clean up the testing configmaps
kubectl delete configmap config --all-namespaces

# 删除测试用的权限组
kubectl delete roles cm-role --namespace default
kubectl delete roles cm-role --namespace kube-public
```