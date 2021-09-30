package hosts

import (
	"myhosts/app"
	"myhosts/model"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type HostsFileManager struct {
	current string
	watcher *fsnotify.Watcher
	txt     chan string
	pwd     string
	err     chan string
	// db data
	dbHosts []*model.Host
}

func NewHostFileManager() *HostsFileManager {
	h := &HostsFileManager{
		pwd: "",
		txt: make(chan string, 1),
		err: make(chan string, 1),
	}
	h.ReLoad()
	var err error
	h.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	if err := h.watcher.Add(app.HostsFile()); err != nil {
		panic(err)
	}

	go h.process()
	return h
}

func (h *HostsFileManager) SetPwd(pwd string) {
	h.pwd = pwd
}

func (h *HostsFileManager) GetPwd() string {
	return h.pwd
}
func (h *HostsFileManager) CurrentHostFile() string {
	return h.current
}

func (h *HostsFileManager) ReadHostFile() {
	b, err := os.ReadFile(app.HostsFile())
	if err != nil {
		h.err <- err.Error()
	}
	h.current = string(b)
}
func (h *HostsFileManager) WriteHostFile(txt string) {
	h.txt <- txt
}

func (h *HostsFileManager) Err() chan string {
	return h.err
}

func (h *HostsFileManager) All() []*model.Host {
	return h.dbHosts
}

func (h *HostsFileManager) ReLoad() error {
	h.ReadHostFile()
	list, err := model.HostListAll()
	if err != nil {
		return err
	}
	h.dbHosts = list
	return nil
}

func (m *HostsFileManager) Create(h *model.Host) error {
	return model.HostCreate(h)
}

func (m *HostsFileManager) Update(h *model.Host) error {
	return model.HostUpdate(h)

}

func (m *HostsFileManager) Delete(h *model.Host) error {
	return model.HostDelete(h)

}

func (m *HostsFileManager) UpdateHostFile() {

	m.ReLoad()
	var content string
	var enable []string
	for _, v := range m.All() {
		if v.Enable && v.Content != "" {
			enable = append(enable, v.Content)
		}
	}
	content = strings.Join(enable, "\n")
	if content != m.current {
		m.WriteHostFile(content)
	}
}

func (h *HostsFileManager) process() {

	for {
		select {
		case content := <-h.txt:
			if err := app.SaveHosts(content, h.pwd); err != nil {
				h.err <- err.Error()
			}
		case ev := <-h.watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				h.ReadHostFile()
			}
		}
	}

}
