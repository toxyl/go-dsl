package main

import (
	"maps"
)

type dslRegistryState struct {
	data      map[string]any
	new       map[string]any
	protected bool
}

func (s *dslRegistryState) add(key string, value any) {
	if s.protected {
		// it's protected, so we need to add it to the new map
		s.new[key] = value
	} else {
		// it's not protected, so we need to add it to the data map
		s.data[key] = value
	}
}

func (s *dslRegistryState) reset() {
	s.new = make(map[string]any)
	s.protected = true
}

func (s *dslRegistryState) update() {
	// update the data map with the new values
	maps.Copy(s.data, s.new)
	// clear the new map
	s.new = make(map[string]any)
	s.protected = true
}

func (s *dslRegistryState) get() []string {
	keys := make([]string, 0, len(s.new))
	for key := range s.new {
		keys = append(keys, key)
	}
	return keys
}
