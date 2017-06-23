package main

import (
	"sync"
	"time"
)

type (
	StatMap struct {
		sync.RWMutex
		data statsMap
	}

	statsMap map[StatKey]*StatRaw
)

func NewStatMap() *StatMap {
	sm := &StatMap{
		data: make(statsMap),
	}
	go sm.putToDb()
	return sm
}

func (s *StatMap) Slice() StatSlice {
	s.RLock()

	stats := s.data
	list := make(StatSlice, len(stats))
	i := 0

	for _, st := range stats {
		list[i] = *st
		i++
	}

	s.RUnlock()
	list.Sort()

	return list
}

func (s *StatMap) AddRaw(sr StatRaw) {
	key := sr.StatKey
	raw := s.Get(key)

	if raw != nil {
		raw.Inc(sr)
		return
	}

	s.Lock()
	s.data[key] = &sr
	s.Unlock()
}

func (s *StatMap) Get(key StatKey) *StatRaw {
	s.RLock()
	sr, _ := s.data[key]
	s.RUnlock()
	return sr
}

func (s *StatMap) Len() int {
	s.RLock()
	l := len(s.data)
	s.RUnlock()
	return l
}

func (s *StatMap) putToDb() {
	for {
		time.Sleep(1 * time.Minute)
		s.Lock()
		old := s.data
		s.data = make(statsMap)
		s.Unlock()

		s.inc(old)
	}
}

func (s *StatMap) inc(old statsMap) {
	fail := make(statsMap)

	for k, v := range old {
		switch v.ContentType {
		case VideoContentType:
			if err := s.incVideo(k.ContentId, v.ViewCount); err != nil {
				log.Println(err)
				fail[k] = v
			}
		case ArticleContentType:
			if err := s.incArticle(k.ContentId, v.ViewCount); err != nil {
				log.Println(err)
				fail[k] = v
			}
		case PhotoAlbumContentType:
			if err := s.incPhotoAlbum(k.ContentId, v.ViewCount); err != nil {
				log.Println(err)
				fail[k] = v
			}
		}
	}

	if len(fail) > 0 {
		s.inc(fail)
	}
}

func (s *StatMap) printVideo(id int) {
	const q = `SELECT ViewCount, Rank FROM tbh_videos WHERE VideoId = ?`
	row := Context.DB.QueryRow(q, id)
	var (
		view int
		rank int64
	)
	if err := row.Scan(&view, &rank); err != nil {
		log.Println("Scan error", err)
	}
	//log.Println("We went data of video", id, view, rank)
}

func (s *StatMap) incVideo(id int, count int64) error {
	if err := Context.DB.Ping(); err != nil {
		return err
	}
	s.printVideo(id)
	const q = `UPDATE tbh_videos
	SET ViewCount = ViewCount + ?,
		Rank = (FavoriteCount / (ViewCount + ?)) * 1000000000,
		ViewDate = ?
	WHERE VideoId = ?`
	_, err := Context.DB.Exec(q, count, count, time.Now(), id)
	if err != nil {
		return err
	}

	log.Println("video", id, count)
	return nil
}

func (s *StatMap) incPhotoAlbum(id int, count int64) error {
	if err := Context.DB.Ping(); err != nil {
		return err
	}
	const q = `UPDATE tbh_photo_albums
	SET ViewCount = ViewCount + ?
	WHERE PhotoAlbumId = ?`
	_, err := Context.DB.Exec(q, count, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *StatMap) incArticle(id int, count int64) error {
	if err := Context.DB.Ping(); err != nil {
		return err
	}
	const q = `UPDATE tbh_articles
	SET ViewCount = ViewCount + ?
	WHERE ArticleId = ?`
	_, err := Context.DB.Exec(q, count, id)
	if err != nil {
		return err
	}
	return nil
}
