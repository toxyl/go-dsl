package main

import (
	"fmt"
	"sort"
	"sync"
)

type dslVarRegistry struct {
	mu    *sync.Mutex
	data  map[string]*dslMetaVarType
	state *dslRegistryState
}

func (r *dslVarRegistry) storeState() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.update()
}

func (r *dslVarRegistry) restoreState() {
	r.mu.Lock()
	defer r.mu.Unlock()
	toRemove := r.state.get()
	for _, name := range toRemove {
		if varRef, exists := r.data[name]; exists {
			// Explicitly clear the data field to release references
			varRef.data = nil
			// Clear closures to break circular references
			varRef.get = nil
			varRef.set = nil
		}
		delete(r.data, name)
	}
	r.state.reset()
}

func (r *dslVarRegistry) register(name, typ, unit, description string, min, max, def any, fnGet func() any, fnSet func(any)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.add(name, def)
	// Create a local copy of the variable for the closure
	varRef := &dslMetaVarType{
		meta: dslMetaVar{
			name: name,
			desc: description,
			typ:  typ,
			min:  min,
			max:  max,
			def:  def,
			unit: unit,
		},
		get: fnGet,
	}

	// Set up Set method without nested locking
	varRef.set = func(a any) error {
		// Validate without acquiring the lock
		if err := varRef.validate(a); err != nil {
			return err
		}
		fnSet(a)
		return nil
	}

	r.data[name] = varRef
}

func (r *dslVarRegistry) has(name string) bool {
	r.mu.Lock()
	_, exists := r.data[name]
	r.mu.Unlock()
	return exists
}

func (r *dslVarRegistry) get(name string) *dslMetaVarType {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.data[name]
}

func (r *dslVarRegistry) set(name string, value any) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if variable exists and create if it doesn't
	if _, exists := r.data[name]; !exists {
		newVar := &dslMetaVarType{
			meta: dslMetaVar{
				name: name,
				typ:  fmt.Sprintf("%T", value),
			},
			data: value,
		}
		newVar.get = func() any { return newVar.data }
		newVar.set = func(v any) error { newVar.data = v; return nil }
		r.data[name] = newVar
		r.state.add(name, value)
		return nil
	}

	// Update existing variable
	return r.data[name].set(value)
}

func (r *dslVarRegistry) names() []string {
	r.mu.Lock()

	names := make([]string, 0, len(r.data))
	for name := range r.data {
		names = append(names, name)
	}
	r.mu.Unlock()
	sort.Strings(names)
	return names
}
