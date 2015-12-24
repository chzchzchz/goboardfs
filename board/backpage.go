package board

//http://sanjose.backpage.com/FemaleEscorts/airport-upscale-incall-15-mins-80-roses-408-585-9002/28397009
// board "sanjose.backpage.com/Events"


import (
	"io"
	"fmt"
	"log"
	"errors"
	"strings"
	"time"
)

type Backpage struct {
	Site
}

func newBackpageSite() *Site {
	fc := &Backpage{Site: *newDefaultSite("")}
	fc.Browser = fc
	return &fc.Site
}

func (s *Backpage) BoardDirectory() []string {
	return []string{"sanjose-events"}
}

func board2pieces(boardName string) (city string, category string, err error) {
	// sanjose-events
	if args := strings.Split(boardName, "-"); len(args) == 2 {
		city = args[0]
		category = args[1]
	} else {
		err = errors.New("bad backpage board name " + boardName)
	}
	return city, category, err
}

func (s *Backpage) OpenBoardReader(b *Board) (io.Reader, error) {
	city, category, err := board2pieces(b.name)
	if err != nil {
		return nil, err
	}
	return httpReader("http://" + city + ".backpage.com/" + category)
}

func (s *Backpage) OpenThreadReader(t *Thread) (io.Reader, error) {
	return httpReader(t.site_key)
}

func (s *Backpage) ReadThread(t *Thread, rc io.Reader) <-chan *Comment {
	out := make(chan *Comment)
	go func() {
		for rec := range recsFromClasses(rc, [][2]string{
			{"h1link", ""},
			// Thursday, December 24, 2015 1:15 AM
			{"adInfo", ""},
			// class="replyDisplay" ahref="..."
			{"replyDisplay", "attrs"},
			{"postingBody", ""},
		}) {
			site_key := strings.Split(rec[2], "\"")[3]
			timestr := strings.Join(strings.Fields(rec[1])[1:], " ")
			when, _ := time.Parse("Monday, January 2, 2006 3:15 PM", timestr)
			out <- &Comment{
				author:  site_key,
				title:	rec[0],
				when:     when,
				body:    rec[3],
				parent:   &t.Comment,
				site_key: site_key}
		}
		close(out)
	}()
	return out
}

func (s *Backpage) ReadBoard(b *Board, rc io.Reader) <-chan *Thread {
	out := make(chan *Thread)
	go func() {
		rcs, err := dupReaderN(rc, 3)
		if err != nil {
			log.Println(err)
			return
		}

		datec := recsFromClasses(
			rcs[0],
			[][2]string{{"date",""}, {"cat","attrs"},})
		titlec := recsFromClasses(rcs[1], [][2]string{{"cat", ""}})
		urlc := recsFromClasses(rcs[2],[][2]string{{"cat", "attrs"}})


		curDate := <-datec
		nextDate := curDate
		for title := range titlec {
			url := (<-urlc)[0]
			if url == nextDate[1] {
				curDate = nextDate
				nextDate = <-datec
			}
			out <- s.rec2thr(b, title[0], url, curDate[0])
		}

		close(out)
	}()
	return out
}

func (s *Backpage) rec2thr(b *Board, title string, url string, date string) *Thread {
	site_key := strings.Trim(strings.Split(url, "\"")[3], " ")
	timestr := fmt.Sprintf("%s 15:04:05 %d",
			strings.TrimSpace(date),
			time.Now().Year())
	when, _ := time.Parse("Mon. Jan. 2 15:05:05 2006", timestr)
	comm := &Comment{
		author:   b.name,
		title:	  strings.Join(strings.Fields(title), " "),
		when:     when,
		body:     "",
		parent:   nil,
		site_key: site_key}
	return newThread(b, comm, make([]*Comment, 1))
}
