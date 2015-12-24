package board

import (
	"fmt"
	"log"
	"time"
)

// root comment
type Thread struct {
	board *Board
	Comment
	comments     []*Comment
	last_refresh time.Time
}

func newThread(b *Board, first *Comment, comms []*Comment) *Thread {
	return &Thread{
		board:        b,
		Comment:      *first,
		comments:     comms,
		last_refresh: time.Unix(0, 0),
	}
}

func (t *Thread) Read() []*Comment {
	since_refresh := time.Since(t.last_refresh).Seconds()
	if since_refresh < REFRESH_TIMEOUT {
		log.Println("Using thread cache since", since_refresh)
		return t.comments
	}

	rc, err := t.board.site.OpenThreadReader(t)
	if err != nil {
		log.Println("Error reading ", err)
		return t.comments
	}

	t.comments = make([]*Comment, 0)
	for comm := range t.board.site.ReadThread(t, rc) {
		t.comments = append(t.comments, comm)
	}
	t.last_refresh = time.Now()
	return t.comments
}

func (t *Thread) Size() uint {
	return uint(len(t.comments))
}

func (t *Thread) Print() (ret string) {
	ret = t.String()
	t.Read()
	for i := range t.comments {
		ret = ret + fmt.Sprintf("\n%s", t.comments[i])
	}
	return ret
}
