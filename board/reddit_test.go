package board

import (
	"testing"
)

func TestRefreshRedditBoard(t *testing.T) {
	doTestRefreshBoard(t, newRedditSite(), "reddit", 25)
}

func TestRefreshRedditThread(t *testing.T) {
	doTestRefreshThread(t, newRedditSite(), "reddit")
}
