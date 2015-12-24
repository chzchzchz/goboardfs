package boardfs

import (
	"errors"
	"log"

	"github.com/hanwen/go-fuse/fuse/nodefs"

	"github.com/chzchzchz/goboardfs/board"
)

type BoardNode struct {
	*board.Board
}

func (bn *BoardNode) List() (ret []string) {
	thrs := bn.Read()
	for i := range thrs {
		ret = append(ret, thrs[i].Title())
	}
	return ret
}

func (bn *BoardNode) NodeFromName(name string) (nodefs.Node, error) {
	t := bn.LookupByTitle(name)
	if t == nil {
		log.Println("could not open thrnode", name)
		return nil, errors.New("could not open threads")
	}
	return &ThreadNode{
		Node:   nodefs.NewDefaultNode(),
		thread: t}, nil
}
