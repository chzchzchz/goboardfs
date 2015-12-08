package board

import (
	"testing"
)

func TestRefreshSlashdotBoard(t *testing.T) {
	doTestRefreshBoard(t, newSlashdotSite(), "slashdot", 15)
}

func TestRefreshSlashdotThread(t *testing.T) {
	doTestRefreshThread(t, newSlashdotSite(), "slashdot")
}
