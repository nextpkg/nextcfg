package configmap

import (
	"errors"
	"time"

	"github.com/nextpkg/nextcfg/source"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type watcher struct {
	group     string
	name      string
	namespace string
	format    string
	client    *kubernetes.Clientset
	st        cache.Store
	ct        cache.Controller
	ch        chan *source.ChangeSet

	exit chan bool
	stop chan struct{}
}

func newWatcher(group, name, ns, format string, c *kubernetes.Clientset) (source.Watcher, error) {
	w := &watcher{
		group:     group,
		name:      name,
		namespace: ns,
		format:    format,
		client:    c,
		ch:        make(chan *source.ChangeSet),
		exit:      make(chan bool),
		stop:      make(chan struct{}),
	}

	lw := cache.NewListWatchFromClient(w.client.CoreV1().RESTClient(), "configmaps", w.namespace,
		fields.OneTermEqualSelector("metadata.name", w.group))
	st, ct := cache.NewInformer(
		lw,
		&v12.ConfigMap{},
		time.Second*30,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: w.handle,
		},
	)

	go ct.Run(w.stop)

	w.ct = ct
	w.st = st

	return w, nil
}

func (w *watcher) handle(_ interface{}, newCmp interface{}) {
	if newCmp == nil {
		return
	}

	cm := newCmp.(*v12.ConfigMap)
	cs := &source.ChangeSet{
		Format:    w.format,
		Source:    w.name,
		Data:      []byte(cm.Data[w.name]),
		Timestamp: cm.CreationTimestamp.Time,
	}
	cs.Checksum = cs.Sum()

	w.ch <- cs
}

// Next 处理新配置
func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-w.ch:
		return cs, nil
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

// Stop 关闭监听器
func (w *watcher) Stop() error {
	select {
	case <-w.exit:
	case <-w.stop:
	default:
		close(w.exit)
	}

	return nil
}
