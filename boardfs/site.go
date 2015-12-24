package boardfs

import (
	"errors"

	"github.com/hanwen/go-fuse/fuse/nodefs"

	"github.com/chzchzchz/goboardfs/board"
)

type SiteNode struct {
	*board.Site
}

func (sn *SiteNode) NodeFromName(name string) (nodefs.Node, error) {
	bn := &BoardNode{Board: sn.Open(name)}
	if bn.Board == nil {
		return nil, errors.New("bad board " + name)
	}
	return newDirNode(nil, nil, bn), nil
}
