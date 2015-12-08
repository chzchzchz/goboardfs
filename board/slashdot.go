package board

import (
	"time"
	"log"
	"strings"
	"strconv"
	"io"
)

type Slashdot struct {
	Site
}

func newSlashdotSite() *Site {
	rs := &Slashdot{ Site : *newDefaultSite("slashdot" ) }
	rs.Browser = rs
	return &rs.Site
}

func (s *Slashdot) BoardDirectory() []string {
	return []string{"m", "science", "tech", "devices", "yro", "developers"}
}

func (s *Slashdot) OpenBoardReader(b *Board) (io.Reader, error) {
	return httpReader("http://" + b.name + ".slashdot.org/")
}

func (s *Slashdot) OpenThreadReader(t *Thread) (io.Reader, error) {
	return httpReader("http:" + t.site_key)
}

func (s *Slashdot) ReadBoard(b *Board, rc io.Reader) (<-chan *Thread) {
	out := make(chan *Thread)
	go func() {
		defer close(out)
		rcs, err := dupReader(rc)
		if err != nil {
			log.Println(err)
			return
		}

		recc := [2]<-chan []string{
			recsFromClasses(rcs[0], [][2]string{
				{"story-title", ""},
				{"comment-bubble", ""},
				{"story-byline", ""},
				{"body", ""},}),
			recsFromClasses(rcs[1], [][2]string{
				{"story-title", "attrs"},}),
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

func (s* Slashdot) ReadThread(t *Thread, rc io.Reader) (<-chan *Comment) {
	out := make(chan *Comment)
	go func() {
		for rec := range recsFromClasses(rc, [][2]string{
			{"title", ""},
			{"by", ""},
			{"otherdetails", ""},
			{"commentBody", ""},
		}) {
			out <- s.rec2comm(&t.Comment, rec)
		}
		close(out)
	}()
	return out
}

func slashdot2time(s string) (time.Time, error) {
	const timeTitleFmt = "Monday January 2, 2006 @3:04PM"
	return time.Parse(timeTitleFmt, s)
}

func (s *Slashdot) rec2thr(b *Board, rec []string) (*Thread) {
	title := strings.Join(strings.Fields(rec[0]), " ")
	num_comments, _ := strconv.ParseInt(rec[1], 10, 64)
	author := strings.Fields(rec[2])[2]
	when, _ := slashdot2time(strings.Join(strings.Fields(rec[2])[4:9], " "))
	site_key := strings.Split(strings.Split(rec[4], "href=")[1], "\"")[1]
	return newThread(
		b,
		&Comment{
			author : author,
			title : title,
			when : when,
			site_key : site_key,
			},
		make([]*Comment, num_comments))
}

func (s *Slashdot) rec2comm(parent *Comment, rec []string) (*Comment) {
	on_split := strings.Split(rec[2], "on ")
	when := parent.when
	if len(on_split) > 1 {
		past_on := on_split[1]
		timestr := strings.Join(strings.Fields(past_on)[0:5], " ")
		t, err := slashdot2time(timestr)
		if err != nil {
			log.Println(err)
		} else {
			when = t
		}
	}
	return &Comment{
		parent : parent,
		author : strings.Fields(rec[1])[1],
		title : strings.Split(rec[0], "\n")[0],
		when : when,
		body : rec[3],
		}
}
