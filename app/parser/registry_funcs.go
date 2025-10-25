package main

import (
	"sort"
	"sync"
)

type dslFnRegistry struct {
	lock  *sync.Mutex
	data  map[string]*dslFnType
	state *dslRegistryState
}

func (r *dslFnRegistry) storeState() {
	r.state.update()
}

func (r *dslFnRegistry) restoreState() {
	toRemove := r.state.get()
	for _, name := range toRemove {
		delete(r.data, name)
	}
	r.state.reset()
}

func (r *dslFnRegistry) register(name, description string, parameters []dslParamMeta, returns []dslParamMeta, function func(...any) (any, error)) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.state.add(name, nil) // for functions, we don't need to store the function itself, just the name

	// Create a local copy of the function for the registry
	r.data[name] = &dslFnType{
		meta: dslFnMeta{
			name:    name,
			desc:    description,
			params:  parameters,
			returns: returns,
		},
		data: function,
	}
}

func (r *dslFnRegistry) get(name string) *dslFnType {
	r.lock.Lock()
	fn, ok := r.data[name]
	r.lock.Unlock()

	if !ok {
		return nil
	}
	return fn
}

func (r *dslFnRegistry) names() []string {
	r.lock.Lock()
	names := make([]string, 0, len(r.data))
	for name := range r.data {
		names = append(names, name)
	}
	r.lock.Unlock()
	sort.Strings(names)
	return names
}
