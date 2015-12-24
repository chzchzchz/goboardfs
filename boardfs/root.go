package boardfs

import (
	"errors"

	"github.com/hanwen/go-fuse/fuse/nodefs"

	"github.com/chzchzchz/goboardfs/board"
)

type RootNode struct {
}

func NewRootNode() nodefs.Node {
	rn := &RootNode{}
	return newDirNode(nil, nil, rn)
}

func (rn *RootNode) NodeFromName(name string) (nodefs.Node, error) {
	sn := &SiteNode{Site: board.NewSite(name)}
	if sn.Site != nil {
		return newDirNode(nil, nil, sn), nil
	}
	return nil, errors.New("no site " + name)
}

func (rn *RootNode) List() []string {
	return []string{"reddit", "4chan", "slashdot"}
}
