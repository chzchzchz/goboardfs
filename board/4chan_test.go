package board

import (
	"testing"
)

func TestRefreshFourChanBoard(t *testing.T) {
	doTestRefreshBoard(t, newFourChanSite(), "4chan", 15)
}

func TestRefreshFourChanThread(t *testing.T) {
	doTestRefreshThread(t, newFourChanSite(), "4chan")
}
