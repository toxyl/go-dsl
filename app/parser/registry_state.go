package main

import (
	"maps"
	"sync"
)

type dslRegistryState struct {
	data      map[string]any
	new       map[string]any
	protected bool
	mu        *sync.Mutex
}

func (s *dslRegistryState) add(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.protected {
		s.new[key] = value // it's protected, so we need to add it to the new map
	} else {
		s.data[key] = value // it's not protected, so we need to add it to the data map
	}
}

func (s *dslRegistryState) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.new = make(map[string]any)
	s.protected = true
}

func (s *dslRegistryState) update() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// update the data map with the new values
	maps.Copy(s.data, s.new)
	// clear the new map
	s.new = make(map[string]any)
	s.protected = true
}

func (s *dslRegistryState) get() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.new))
	for key := range s.new {
		keys = append(keys, key)
	}
	return keys
}
