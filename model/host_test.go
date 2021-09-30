package model

import (
	"fmt"
	"testing"
)

func TestHostCreate(t *testing.T) {
	Init()
	HostCreate(&Host{
		ID:      0,
		Title:   "fat2",
		Content: "127.0.0.2 bilibili.com",
		Enable:  false,
	})

}

func TestHostGet(t *testing.T) {
	Init()
	got, err := HostGet(3)
	fmt.Printf("%+v err:%+v", got, err)

}

func TestHostListAll(t *testing.T) {
	Init()
	got, err := HostListAll()
	fmt.Printf("%+v err:%+v", got, err)

}

func TestHostUpdate(t *testing.T) {
	Init()
	HostUpdate(&Host{ID: 3, Content: "123123sdsxxxf"})

}

func TestHostDelete(t *testing.T) {
	Init()
	HostDelete(&Host{ID: 1})

}
