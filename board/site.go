package board

import (
	"io"
	"log"
)

type Browser interface {
	BoardDirectory() []string
	OpenBoardReader(*Board) (io.Reader, error)
	ReadBoard(*Board, io.Reader) <-chan *Thread
	OpenThreadReader(*Thread) (io.Reader, error)
	ReadThread(*Thread, io.Reader) <-chan *Comment
}

type Site struct {
	Browser
	name        string
	list_boards []string // cached ListBoards
	boards      map[string]*Board
}

type siteInfo struct {
	name string
	init func() (*Site)
}

var (
	sites = []siteInfo {
		siteInfo {name : "4chan", init : newRedditSite, },
		siteInfo {name : "hackernews", init : newHackerNewsSite, },
		siteInfo {name : "reddit", init : newRedditSite, },
		siteInfo {name : "slashdot", init : newSlashdotSite, },
	}
)

func siteByName(name string) (s *Site) {
	for i := range sites {
		if sites[i].name == name {
			return sites[i].init()
		}
	}
	return nil
}

func Sites() (v []string) {
	for i := range sites {
		v = append(v, sites[i].name)
	}
	return v
}

// create board by name
func NewSite(name string) (s *Site) {
	s = siteByName(name)
	if s != nil {
		s.list_boards = s.Browser.BoardDirectory()
	} else {
		log.Println("no board for " + name)
	}
	return s
}

func (s *Site) NewBoard(boardname string) *Board {
	return newDefaultBoard(s, boardname)
}

func (s *Site) Open(name string) *Board {
	b := s.boards[name]
	if b != nil {
		log.Println("Already have board ", name)
		return b
	}
	b = s.NewBoard(name)
	if b == nil {
		log.Println("Could not get board ", name)
		return nil
	}
	s.boards[name] = b
	return b
}

func (s *Site) Close(name string) {
	delete(s.boards, name)
}

func (s *Site) List() []string {
	return s.list_boards
}

func newDefaultSite(n string) *Site {
	return &Site{
		name:   n,
		boards: make(map[string]*Board)}
}
