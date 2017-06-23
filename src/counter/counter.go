package main

import (
	"encoding/gob"
	"os"
	"sync"
	"time"
)

type (
	Counter struct {
		stats    *StatMap
		rawQueue chan StatRaw
		wg       sync.WaitGroup
	}
)

const rawQueueSize int = 1000
const backupFilePath = "./counters.db"

// Counter

func NewCounter() *Counter {
	c := &Counter{
		stats:    NewStatMap(),
		rawQueue: make(chan StatRaw, rawQueueSize),
	}
	go c.count()
	return c
}

func (s *Counter) Pull() StatSlice {
	return s.stats.Slice()
}

func (s *Counter) Push(sr StatRaw) {
	sr.Time = time.Now()
	s.wg.Add(1)
	s.rawQueue <- sr
}

func (s *Counter) count() {
	for raw := range s.rawQueue {
		raw.Time = raw.Time.Local().Truncate(24 * time.Hour)
		s.stats.AddRaw(raw)
		s.wg.Done()
	}
}

func (s *Counter) Backup() {
	stats := s.stats.Slice()
	f, err := os.OpenFile(backupFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm|os.ModeExclusive)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	enc.Encode(stats)
}

func (s *Counter) Restore() {
	f, err := os.OpenFile(backupFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	var stats StatSlice
	if err = decoder.Decode(&stats); err != nil {
		log.Println(err)
		return
	}
	s.stats.Lock()
	defer s.stats.Unlock()

	for _, raw := range stats {
		r := &StatRaw{}
		*r = raw
		s.stats.data[raw.StatKey] = r
	}
}
