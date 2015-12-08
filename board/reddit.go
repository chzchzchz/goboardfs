package board

import (
	"strings"
	"time"
	"strconv"
	"io"
)

type Reddit struct {
	Site
}

func newRedditSite() *Site {
	rs := &Reddit{ Site : *newDefaultSite("reddit" ) }
	rs.Browser = rs
	return &rs.Site
}

func (s *Reddit) BoardDirectory() []string {
	return []string{"news", "funny"}
}

func (s *Reddit) OpenBoardReader(b *Board) (io.Reader, error) {
	return httpReader("http://reddit.com/r/" + b.name)
}

func (s *Reddit) OpenThreadReader(t *Thread) (io.Reader, error) {
	return httpReader(t.site_key)
}

func (s *Reddit) ReadBoard(b *Board, rc io.Reader) (<-chan *Thread) {
	out := make(chan *Thread)
	go func() {
		defer close(out)
		rcs, err := dupReader(rc)
		if err != nil {
			return
		}

		recc := [2]<-chan []string{
			recsFromClasses(rcs[0], [][2]string {
				{"title may-blank", ""},
				{"live-timestamp", "title"},
				{"author may-blank", ""},
				{"comments may-blank", ""},}),
			// to get the site key ugh
			recsFromClasses(rcs[1], [][2]string {
				{"comments may-blank", "href"},}),
		}
		for {
			rec1, r0_ok := <-recc[0]
			rec2, r1_ok := <-recc[1]
			if r0_ok && r1_ok {
				out <- s.rec2thr(b, append(rec1, rec2...))
			} else if !r0_ok && !r1_ok {
				return
			}
		}
	}()
	return out
}

func (s* Reddit) ReadThread(t *Thread, rc io.Reader) (<-chan *Comment) {
	out := make(chan *Comment)
	go func() {
		for rec := range recsFromClasses(rc, [][2]string{
			{"author may-blank", ""},
			{"live-timestamp", "title"},
			{"usertext-body may-blank", ""},
		}) {
			out <- s.rec2comm(&t.Comment, rec)
		}
		close(out)
	}()
	return out
}

func reddit2time(s string) time.Time {
	const timeTitleFmt = "Mon Jan 2 15:04:05 2006 MST"
	when, _ := time.Parse(timeTitleFmt, s)
	return when
}

func (s *Reddit) rec2comm(parent *Comment, rec []string) (*Comment) {
	return &Comment{
		parent : parent,
		author : rec[0],
		title : "",
		when : reddit2time(rec[1]),
		body : rec[2],
		}
}

func (s *Reddit) rec2thr(b *Board, rec []string) (*Thread) {
	num_comments := 0
	if !strings.Contains(rec[3], "empty") {
		v := strings.Split(rec[3], " ")[0]
		n , _ := strconv.ParseInt(v, 10, 32)
		num_comments = int(n)
	}
	return newThread(
		b,
		&Comment{
			author : rec[2],
			title : rec[0],
			when : reddit2time(rec[1]),
			site_key : rec[4],
			},
		make([]*Comment, num_comments))
}
