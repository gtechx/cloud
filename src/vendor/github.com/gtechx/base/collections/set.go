package collections

import (
	"sync"
)

type Set struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func NewSet() *Set {
	return &Set{
		m: map[interface{}]interface{}{},
	}
}

func (s *Set) Add(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = nil
}

func (s *Set) Remove(item interface{}) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) Len() interface{} {
	return len(s.m)
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]interface{}{}
}

func (s *Set) IsEmpty() bool {
	s.RLock()
	defer s.RUnlock()
	if len(s.m) == 0 {
		return true
	}
	return false
}

// func (s *Set) List() []int {
// 	s.RLock()
// 	defer s.RUnlock()
// 	list := []int{}
// 	for item := range s.m {
// 		list = append(list, item)
// 	}
// 	return list
// }
