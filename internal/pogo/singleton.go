package pogo

import "sync"

type Singleton struct {
	one sync.Once
}

func (s *Singleton) Once(f func()) {
	s.one.Do(f)
}
