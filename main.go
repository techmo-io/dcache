package main

import (
	"time"
)

type clock interface {
	Now() time.Time
}

type Store struct {
	d     map[string][]byte
	exp   map[string]int64
	clock clock
}

func (s *Store) set(key string, value []byte, TTLMilliseconds ...int64) error {
	var ttl int64 = 0
	if len(TTLMilliseconds) > 0 {
		ttl = TTLMilliseconds[0]
	}

	s.d[key] = value
	if ttl > 0 {
		s.exp[key] = s.clock.Now().UnixMilli() + ttl
	} else {
		delete(s.exp, key)
	}

	return nil
}

func (s *Store) get(key string) []byte {
	v, e := s.d[key]
	if !e {
		return nil
	}

	if exp := s.exp[key]; exp != 0 && s.clock.Now().UnixMilli() >= exp {
		// consider key for removal
		return nil
	}

	return v
}

func (s *Store) delete(key string) {
	delete(s.d, key)
	delete(s.exp, key)
}

func (s *Store) has(key string) bool {

	_, e := s.d[key]
	if !e {
		return false
	}

	if exp, e := s.exp[key]; e && exp <= s.clock.Now().UnixMilli() {
		return false
	}

	return true
}

func (s *Store) clear() {
	s.d = make(map[string][]byte)
	s.exp = make(map[string]int64)
}

func NewStore(clock clock) *Store {
	s := Store{
		clock: clock,
	}
	s.clear()
	return &s
}

func main() {

}
