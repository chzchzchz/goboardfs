package board

import (
	"log"
	"time"
)

const REFRESH_TIMEOUT float64 = 10 * 60.0

type Board struct {
	site         *Site
	name         string
	threads      []*Thread
	last_refresh time.Time
}

func newDefaultBoard(s *Site, n string) *Board {
	return &Board{
		site:         s,
		name:         n,
		threads:      nil,
		last_refresh: time.Unix(0, 0),
	}
}

func (b *Board) Read() []*Thread {
	since_refresh := time.Since(b.last_refresh).Seconds()
	if since_refresh < REFRESH_TIMEOUT {
		log.Println("Using cache due to short refresh", since_refresh)
		return b.threads
	}

	rc, err := b.site.OpenBoardReader(b)
	if err != nil {
		log.Printf("error reading (%v)", err)
		return b.threads
	}

	b.threads = make([]*Thread, 0)
	for thr := range b.site.ReadBoard(b, rc) {
		b.threads = append(b.threads, thr)
	}
	b.last_refresh = time.Now()

	return b.threads
}

// XXX use map?
func (b *Board) LookupByTitle(name string) *Thread {
	b.Read()
	for i := 0; i < len(b.threads); i++ {
		if name == b.threads[i].Title() {
			return b.threads[i]
		}
	}
	return nil
}
