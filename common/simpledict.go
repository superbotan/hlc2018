package common

import (
	"strings"
	"sync"
)

type SimpleDict16 struct {
	mx     sync.Mutex
	nextID uint16
	byID   map[uint16]string
	byName map[string]uint16
}

func (sd *SimpleDict16) Append(n string) (id uint16) {

	if n == "" {
		return 0
	}

	v, ok := sd.byName[n]
	if ok {
		return v
	}

	sd.mx.Lock()
	defer sd.mx.Unlock()
	v, ok = sd.byName[n]
	if ok {
		return v
	}

	id = sd.nextID
	sd.nextID = id + 1

	sd.byID[id] = n
	sd.byName[n] = id

	return id
}

func (sd *SimpleDict16) GetByName(n string) (id uint16) {
	v, ok := sd.byName[n]
	if ok {
		return v
	}

	return 0
}

func (sd *SimpleDict16) GetByID(id uint16) (n string) {
	v, ok := sd.byID[id]
	if ok {
		return v
	}

	return ""
}

func (sd *SimpleDict16) GetAllKeys() []uint16 {
	res := make([]uint16, 0, len(sd.byID))

	for i := range sd.byID {
		res = append(res, i)
	}

	return res
}

func SimpleDict16Create() (sd *SimpleDict16) {
	sd = &SimpleDict16{
		nextID: 1,
		byID:   make(map[uint16]string),
		byName: make(map[string]uint16),
	}

	return sd
}

func (sd *SimpleDict16) GetLikeValue(n string) map[uint16]bool {
	r := make(map[uint16]bool)

	for k, v := range sd.byID {
		if strings.Index(v, n) == 0 {
			r[k] = true
		}
	}

	return r
}

func (sd *SimpleDict16) GetList(ns string) map[uint16]bool {
	r := make(map[uint16]bool)

	for _, v := range strings.Split(ns, ",") {
		k, ok := sd.byName[v]
		if ok {
			r[k] = true
		}
	}

	return r
}

type SimpleDict8 struct {
	mx     sync.Mutex
	nextID uint8
	byID   map[uint8]string
	byName map[string]uint8
}

func (sd *SimpleDict8) Append(n string) (id uint8) {

	if n == "" {
		return 0
	}

	v, ok := sd.byName[n]
	if ok {
		return v
	}

	sd.mx.Lock()
	defer sd.mx.Unlock()
	v, ok = sd.byName[n]
	if ok {
		return v
	}

	id = sd.nextID
	sd.nextID = id + 1

	sd.byID[id] = n
	sd.byName[n] = id

	return id
}

func (sd *SimpleDict8) GetByName(n string) (id uint8) {
	v, ok := sd.byName[n]
	if ok {
		return v
	}

	return 0
}

func (sd *SimpleDict8) GetByID(id uint8) (n string) {
	v, ok := sd.byID[id]
	if ok {
		return v
	}

	return ""
}

func (sd *SimpleDict8) GetAllKeys() []uint8 {
	res := make([]uint8, 0, len(sd.byID))

	for i := range sd.byID {
		res = append(res, i)
	}

	return res
}

func SimpleDict8Create() (sd *SimpleDict8) {
	sd = &SimpleDict8{
		nextID: 1,
		byID:   make(map[uint8]string),
		byName: make(map[string]uint8),
	}

	return sd
}

func (sd *SimpleDict8) GetLikeValue(n string) map[uint8]bool {
	r := make(map[uint8]bool)

	for k, v := range sd.byID {
		if strings.Index(v, n) == 0 {
			r[k] = true
		}
	}

	return r
}

func (sd *SimpleDict8) GetList(ns string) map[uint8]bool {
	r := make(map[uint8]bool)

	for _, v := range strings.Split(ns, ",") {
		k, ok := sd.byName[v]
		if ok {
			r[k] = true
		}
	}

	return r
}
func (sd *SimpleDict8) GetListArray(ns string) []uint8 {
	r := make([]uint8, 0)

	for _, v := range strings.Split(ns, ",") {
		k, ok := sd.byName[v]
		if ok {
			r = append(r, k)
		}
	}

	return r
}
