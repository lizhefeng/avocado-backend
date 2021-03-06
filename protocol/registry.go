package protocol

import (
	"log"
	"sync"

	"github.com/pkg/errors"
)

// Registry is the hub of all protocols deployed on the chain
type Registry struct {
	protocols sync.Map
}

// Register registers the protocol with a unique ID
func (r *Registry) Register(id string, p Protocol) error {
	_, loaded := r.protocols.LoadOrStore(id, p)
	if loaded {
		return errors.Errorf("Protocol with ID %s is already registered", id)
	}
	return nil
}

// ForceRegister registers the protocol with a unique ID and force replacing the previous protocol if it exists
func (r *Registry) ForceRegister(id string, p Protocol) error {
	r.protocols.Store(id, p)
	return nil
}

// Find finds a protocol by ID
func (r *Registry) Find(id string) (Protocol, bool) {
	value, ok := r.protocols.Load(id)
	if !ok {
		return nil, false
	}
	p, ok := value.(Protocol)
	if !ok {
		log.Panic("Registry stores the item which is not a protocol")
	}
	return p, true
}

// All returns all protocols
func (r *Registry) All() []Protocol {
	all := make([]Protocol, 0)
	r.protocols.Range(func(_, value interface{}) bool {
		p, ok := value.(Protocol)
		if !ok {
			log.Panic("Registry stores the item which is not a protocol")
		}
		all = append(all, p)
		return true
	})
	return all
}
