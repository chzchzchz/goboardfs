package board

import (
	"testing"
)

func TestRefreshHackerNewsBoard(t *testing.T) {
	doTestRefreshBoard(t, newHackerNewsSite(), "hackernews", 30)
}

func TestRefreshHackerNewsThread(t *testing.T) {
	doTestRefreshThread(t, newHackerNewsSite(), "hackernews")
}
