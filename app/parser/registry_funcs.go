package main

import (
	"sort"
	"sync"
)

type dslFnRegistry struct {
	mu    *sync.Mutex
	data  map[string]*dslFnType
	state *dslRegistryState
}

func (r *dslFnRegistry) storeState() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.update()
}

func (r *dslFnRegistry) restoreState() {
	r.mu.Lock()
	defer r.mu.Unlock()
	toRemove := r.state.get()
	for _, name := range toRemove {
		delete(r.data, name)
	}
	r.state.reset()
}

func (r *dslFnRegistry) register(name, description string, parameters []dslParamMeta, returns []dslParamMeta, function func(...any) (any, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
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
	r.mu.Lock()
	defer r.mu.Unlock()
	fn, ok := r.data[name]

	if !ok {
		return nil
	}
	return fn
}

func (r *dslFnRegistry) names() []string {
	r.mu.Lock()
	names := make([]string, 0, len(r.data))
	for name := range r.data {
		names = append(names, name)
	}
	r.mu.Unlock()
	sort.Strings(names)
	return names
}
