package boardfs

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"syscall"
)

type StrNode struct {
	nodefs.Node
	s string
}

func newStrNode(in_s string) (*StrNode) {
	return &StrNode{
		Node : nodefs.NewDefaultNode(),
		s : in_s}
}

func (sn *StrNode) Open(flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	return nodefs.NewDataFile([]byte(sn.s)), fuse.OK
}

func (sn* StrNode) GetAttr(
	out *fuse.Attr,
	file nodefs.File,
	context *fuse.Context) (code fuse.Status) {
	out.Mode = 0444 | syscall.S_IFREG
	owner := getDefaultOwner()
	out.Uid = owner.Uid
	out.Gid = owner.Gid
	out.Nlink = 1
	out.Size = uint64(len(sn.s))
	return fuse.OK
}
