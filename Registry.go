package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Registry struct {
	id int32
    globalListeners []chan RegistryGlobal
	globalRemoveListeners []chan RegistryGlobalRemove
	events []func(r *Registry, msg message) error
    name string
}

func NewRegistry() (r *Registry) {
	r = new(Registry)
    r.name = "Registry"
    r.globalListeners = make([]chan RegistryGlobal, 0)
	r.globalRemoveListeners = make([]chan RegistryGlobalRemove, 0)

    r.events = append(r.events, registryGlobal)
	r.events = append(r.events, registryGlobalRemove)
	return
}

func (r *Registry) HandleEvent(msg message) {
	if r.events[msg.opcode] != nil {
		r.events[msg.opcode](r, msg)
	}
}

func (r *Registry) SetId(id int32) {
	r.id = id
}

func (r *Registry) Id() int32 {
	return r.id
}

func (r *Registry) Name() string {
    return r.name
}

////
//// REQUESTS
////

func (r *Registry) Bind(name uint32, iface string, version uint32, id Object) {
    sendrequest(r, "wl_registry_bind", name, iface, version, id)
}

////
//// EVENTS
////

type RegistryGlobal struct {
    Name uint32
	Iface string
	Version uint32
}

func (r *Registry) AddGlobalListener(channel chan RegistryGlobal) {
    r.globalListeners = append(r.globalListeners, channel)
}

func registryGlobal(r *Registry, msg message) (err error) {
    var data RegistryGlobal

    // Read name
    name,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Name = name

    // Read iface
    iface,err := readString(msg.buf)
    if err != nil {
        return
    }
    data.Iface = iface

    // Read version
    version,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Version = version

    // Dispatch events
    for _,channel := range r.globalListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type RegistryGlobalRemove struct {
    Name uint32
}

func (r *Registry) AddGlobalRemoveListener(channel chan RegistryGlobalRemove) {
    r.globalRemoveListeners = append(r.globalRemoveListeners, channel)
}

func registryGlobalRemove(r *Registry, msg message) (err error) {
    var data RegistryGlobalRemove

    // Read name
    name,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Name = name

    // Dispatch events
    for _,channel := range r.globalRemoveListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
