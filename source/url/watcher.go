package url

import (
	"errors"
	"time"

	"github.com/nextpkg/nextcfg/source"
	"log"
)

type urlWatcher struct {
	u    *urlSource
	ch   chan *source.ChangeSet
	exit chan bool
}

func newWatcher(u *urlSource) (*urlWatcher, error) {
	w := &urlWatcher{
		u:    u,
		ch:   make(chan *source.ChangeSet),
		exit: make(chan bool),
	}
	go w.run()

	return w, nil
}

func (u *urlWatcher) run() {
	for {
		select {
		case <-u.exit:
		case <-time.After(30 * time.Second):
			cs, err := u.u.Read()
			if err != nil {
				log.Println("watch failed: ", err)
				continue
			}

			u.ch <- cs
		}
	}
}

// Next 处理新配置
func (u *urlWatcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-u.ch:
		return cs, nil
	case <-u.exit:
		return nil, errors.New("url watcher stopped")
	}
}

// Stop 关闭监听器
func (u *urlWatcher) Stop() error {
	select {
	case <-u.exit:
	default:
		close(u.exit)
	}

	return nil
}
