package process

// ============================================================================
// PROCESS STATE HELPERS
// ============================================================================

// State provides a simple key-value store for process-local state.
//
// In a distributed algorithm, each process maintains its own local state.
// This state is:
//   - Private: no other process can read it
//   - Volatile: lost on crash (unless using persistent storage)
//   - Updated atomically during each Handle() step

type State struct {
	store map[string]any
}

// NewState creates an empty state store.
func NewState() *State {
	return &State{store: make(map[string]any)}
}

// Set stores a value.
func (s *State) Set(key string, value any) {
	s.store[key] = value
}

// Get retrieves a value (returns nil if not found).
func (s *State) Get(key string) any {
	return s.store[key]
}

// GetInt retrieves an integer value (returns 0 if not found or wrong type).
func (s *State) GetInt(key string) int {
	v, ok := s.store[key].(int)
	if !ok {
		return 0
	}
	return v
}

// GetString retrieves a string value.
func (s *State) GetString(key string) string {
	v, ok := s.store[key].(string)
	if !ok {
		return ""
	}
	return v
}

// GetBool retrieves a boolean value.
func (s *State) GetBool(key string) bool {
	v, ok := s.store[key].(bool)
	if !ok {
		return false
	}
	return v
}

// Has checks if a key exists.
func (s *State) Has(key string) bool {
	_, ok := s.store[key]
	return ok
}

// Delete removes a key.
func (s *State) Delete(key string) {
	delete(s.store, key)
}

// Clear removes all state (simulates volatile state loss on crash).
func (s *State) Clear() {
	s.store = make(map[string]any)
}

// Clone returns a copy of the state.
func (s *State) Clone() *State {
	ns := NewState()
	for k, v := range s.store {
		ns.store[k] = v
	}
	return ns
}

// Keys returns all keys in the state.
func (s *State) Keys() []string {
	keys := make([]string, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys
}
