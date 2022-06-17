package main

import (
	"reflect"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func NewFunctions(MaxParamLen int) *Functions {
	t := new(Functions)
	t.fns = make(map[string]interface{})
	t.mu = &sync.RWMutex{}
	t.mpl = MaxParamLen
	return t
}

type Functions struct {
	fns map[string]interface{}
	mu  *sync.RWMutex
	mpl int
}

func (t *Functions) Run(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	cmd := t.sw(&s)

	t.mu.RLock()
	fn, ok := t.fns[cmd]
	t.mu.RUnlock()

	if !ok {
		return
	}

	funcArgs := reflect.ValueOf(fn).Type().NumIn() - 2
	funcValue := reflect.ValueOf(fn)

	v := make([]reflect.Value, funcArgs+2)

	if funcArgs == 1 {
		v[0] = reflect.ValueOf(s)
		v[1] = reflect.ValueOf(ds)
		v[2] = reflect.ValueOf(dm)
	}
	if funcArgs > 1 {
		for i := 0; i < funcArgs-1; i++ {
			v[i] = reflect.ValueOf(t.sw(&s))
		}
		v[funcArgs-1] = reflect.ValueOf(s)
		v[funcArgs] = reflect.ValueOf(ds)
		v[funcArgs+1] = reflect.ValueOf(dm)
	}

	funcValue.Call(v)
	// tw(&s)
}

func (t *Functions) Bind(cmd string, fn interface{}) {
	t.mu.Lock()
	t.fns[cmd] = fn
	t.mu.Unlock()
}

func (t *Functions) sw(s *string) string {
	l, rt, i := t.mpl, "", 0

	for _, v := range *s {
		if v != ' ' {
			rt = rt + string(v)
		} else {
			break
		}
		i++
		if i >= l {
			break
		}
	}
	if len(*s) <= i {
		*s = ""
		return rt
	}
	*s = (*s)[i+1:]
	return rt
}
