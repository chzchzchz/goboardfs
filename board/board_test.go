package board

import (
	"os"
	"testing"
	"time"
)

func doTestRefreshBoard(t *testing.T, s *Site, fname string, expected_c uint) {
	rc, err := os.Open("testfiles/" + fname + ".board.html")
	if err != nil {
		t.Error("Couldn't open board html")
		return
	}

	thrs := s.ReadBoard(s.Open("abc"), rc)
	t_c := uint(0)
	for thr := range thrs {
		if thr.when.Year() < 2015 {
			t.Error("Bad year on ", *thr)
		}
		t_c++
	}

	if t_c != expected_c {
		t.Error("Bad thread count. Expected ", expected_c, " got ", t_c)
		return
	}
}

func doTestRefreshThread(t *testing.T, s *Site, fname string) {
	path := "testfiles/" + fname + ".thread.html"
	rc, err := os.Open(path)
	if err != nil {
		t.Error("Couldn't open thread html " + path)
		return
	}

	thr := newThread(nil, &Comment{when: time.Now()}, nil)
	c_c := 0
	for c := range s.Browser.ReadThread(thr, rc) {
		if c.when.Year() < 2015 {
			t.Error("Bad year on ", *c)
		}
		c_c++
	}

	if c_c == 0 {
		t.Error("Expected some comments, got none")
		return
	}
}


func TestRefreshBoard(t *testing.T) {
	tests := []struct {
		name string
		boardTotal uint
	}{
		{"4chan", 15},
		{"hackernews", 30},
		{"reddit", 25},
		{"slashdot", 15},
	}
	for i := range tests  {
		test := &tests[i]
		t.Log("Testing...", test.name)
		doTestRefreshBoard(
			t, NewSite(test.name), test.name, test.boardTotal)
	}
}

func TestRefreshThread(t *testing.T) {
	sites := Sites()
	for i := range sites {
		t.Log("Testing...", sites[i])
		doTestRefreshThread(t, NewSite(sites[i]), sites[i])
	}
}
