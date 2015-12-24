package board

import (
	"io"
	"strconv"
	"strings"
	"time"
	"log"
)

type HackerNews struct {
	Site
}

func newHackerNewsSite() *Site {
	fc := &HackerNews{Site: *newDefaultSite("4chan")}
	fc.Browser = fc
	return &fc.Site
}

func (s *HackerNews) BoardDirectory() []string {
	return []string{"main"}
}

func (s *HackerNews) OpenBoardReader(b *Board) (io.Reader, error) {
	return httpReader("http://news.ycombinator.com/")// + b.name)
}

func (s *HackerNews) OpenThreadReader(t *Thread) (io.Reader, error) {
	return httpReader(
		"http://news.ycombinator.com/" + t.site_key)
}

func (s *HackerNews) ReadThread(t *Thread, rc io.Reader) <-chan *Comment {
	out := make(chan *Comment)
	go func() {
		for rec := range recsFromClasses(rc, [][2]string{
			{"ind", ""},
			{"comhead", ""}, // user, time
			{"c00", ""}, // comment
			{"reply", "attrs"}, // site key
		}) {
			out <- s.rec2comm(&t.Comment, rec)
		}
		close(out)
	}()
	return out
}

func (s *HackerNews) ReadBoard(b *Board, rc io.Reader) <-chan *Thread {
	out := make(chan *Thread)
	go func() {
		rcs, err := dupReader(rc)
		if err != nil {
			log.Println(err)
			return
		}

		rec1c := recsFromClasses(rcs[1], [][2]string{
			{"rank", ""},
			{"subtext", "attrs"},})
		for rec0 := range recsFromClasses(rcs[0], [][2]string{
			{"rank", ""}, // synchronize the stream on rank
			{"title", ""},
			{"subtext", ""},
		}) {
			if rec1, ok := <-rec1c; ok {
				out <- s.rec2thr(b, append(rec0, rec1...))
			}
		}
		close(out)
	}()
	return out
}

func (s *HackerNews) rec2thr(b *Board, rec []string) *Thread {
	subtext := strings.Split(rec[2]," ")
	// XXX this is kind of crummy in ls -l. Maybe I should randomize minutes
	num_hours, _ := strconv.ParseInt(strings.Trim(subtext[12], "\n"), 10, 64)
	hour_dur := time.Duration(time.Duration(num_hours) * time.Hour)
	when := time.Now().Add(-hour_dur)
	// item?id=123345356
	site_key := strings.Split(strings.Split(rec[4], " ")[5], "\"")[1]
	comm := &Comment{
		author:   subtext[11],
		title:    strings.Split(rec[1], "\n")[0],
		when:     when,
		body:     rec[2],
		parent:   nil,
		site_key: site_key,}
	num_replies, _ := strconv.ParseInt(
		strings.Trim(subtext[17], "\n"), 10, 64)
	return newThread(b, comm, make([]*Comment, num_replies))
}

func (s* HackerNews) rec2comm(p *Comment, rec []string) *Comment {
	rec1 := strings.Split(rec[1], " ")
	num_hours, _ := strconv.ParseInt(strings.Trim(rec1[11], "\n"), 10, 64)
	hour_dur := time.Duration(time.Duration(num_hours) * time.Hour)
	when := time.Now().Add(-hour_dur)

	body := strings.Trim(rec[2], "\n ")
	body = strings.Replace(body, "   \nreply", "", -1)
	body = strings.Trim(body, "\n ")

	return &Comment{
		author: strings.Trim(rec1[10], "\n"),
		title: "",
		when:     when,
		body:     body,
		parent:   p,
		site_key: "???",}
}