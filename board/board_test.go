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
