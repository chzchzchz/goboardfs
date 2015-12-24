package board

import (
	"io"
	"strconv"
	"strings"
	"time"
)

type FourChan struct {
	Site
}

func newFourChanSite() *Site {
	fc := &FourChan{Site: *newDefaultSite("4chan")}
	fc.Browser = fc
	return &fc.Site
}

func (s *FourChan) BoardDirectory() []string {
	return []string{"g", "ck"}
}

func (s *FourChan) OpenBoardReader(b *Board) (io.Reader, error) {
	return httpReader("http://boards.4chan.org/" + b.name)
}

func (s *FourChan) OpenThreadReader(t *Thread) (io.Reader, error) {
	return httpReader(
		"http://boards.4chan.org/" + t.board.name + "/" + t.site_key)
}

func (s *FourChan) ReadThread(t *Thread, rc io.Reader) <-chan *Comment {
	out := make(chan *Comment)
	go func() {
		for rec := range recsFromClasses(rc, [][2]string{
			{"postInfoM mobile", "id"},
			{"dateTime postNum", "data-utc"},
			{"postMessage", ""},
		}) {
			out <- s.rec2comm(&t.Comment, rec)
		}
		close(out)
	}()
	return out
}

func (s *FourChan) ReadBoard(b *Board, rc io.Reader) <-chan *Thread {
	out := make(chan *Thread)
	go func() {
		for rec := range recsFromClasses(rc, [][2]string{
			{"postInfoM mobile", "id"},
			{"subject", ""},
			{"dateTime postNum", "data-utc"},
			{"postMessage", ""},
			{"info", ""},
		}) {
			out <- s.rec2thr(b, rec)
		}
		close(out)
	}()
	return out
}

func (s *FourChan) rec2thr(b *Board, rec []string) *Thread {
	// no parent since root of thread
	comm := s.rec2comm(nil, []string{rec[0], rec[2], rec[3]})
	title := rec[1]
	if title == "" {
		v := 64
		if v > len(comm.body) {
			v = len(comm.body)
		}
		title = comm.body[:v] + "..."
	}
	comm.title = title
	num_replies, _ := strconv.ParseInt(strings.Split(rec[4], " ")[0], 10, 64)
	return newThread(b, comm, make([]*Comment, num_replies))
}

// {id, date, message}
func (s *FourChan) rec2comm(p *Comment, rec []string) *Comment {
	dateNum, _ := strconv.ParseInt(rec[1], 10, 64)
	when := time.Unix(dateNum, 0)
	return &Comment{
		author:   "Anonymous",
		title:    "",
		when:     when,
		body:     rec[2],
		parent:   p,
		site_key: rec[0][3:len(rec)]}
}
