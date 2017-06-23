package main

import (
	"runtime"
	"sync/atomic"
	"time"
)

type ContentType string

const (
	ArticleContentType    ContentType = "article"
	PhotoAlbumContentType ContentType = "photoalbum"
	VideoContentType      ContentType = "video"
)

type (
	StatKey struct {
		ContentId   int
		ContentType ContentType
		Time        time.Time
	}
	StatRaw struct {
		StatKey
		StatCounters
	}
	StatCounters struct {
		ViewCount int64
	}
)

const timeFormat = "2006-01-02 15:04"

func NewStatRaw(id int, t ContentType) StatRaw {
	return StatRaw{
		StatKey{
			ContentId:   id,
			ContentType: t,
			Time:        time.Now(),
		},
		StatCounters{},
	}
}

func (s *StatRaw) Inc(what StatRaw) {
	atomic.AddInt64(&s.ViewCount, what.ViewCount)
	runtime.Gosched()
}

func (s *StatRaw) TimeFormat() string {
	return s.Time.Format(timeFormat)
}

func (s *StatRaw) TimeUnixMs() int64 {
	return s.Time.Unix() * 1e3
}
