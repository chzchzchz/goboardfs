package boardfs

import (
	"syscall"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	"github.com/chzchzchz/goboardfs/board"
)

type ThreadNode struct {
	nodefs.Node
	thread *board.Thread
}

// XXX be smarter about dispatch later
func (tn *ThreadNode) OpenDir(
	context *fuse.Context) (ret []fuse.DirEntry, st fuse.Status) {
	ret = []fuse.DirEntry{
		{Name: "comments", Mode: fuse.S_IFREG},
	}
	return ret, fuse.OK
}

func (tn *ThreadNode) Lookup(
	out *fuse.Attr,
	name string,
	context *fuse.Context) (*nodefs.Inode, fuse.Status) {
	if name != "comments" {
		return nil, fuse.ENOENT
	}

	// XXX should look at comments for this info
	//t := tn.thread.Time()
	// out.SetTimes(&t, &t, &t)
	strn := newStrNode(tn.thread.Print())
	return tn.Inode().NewChild(name, false, strn), fuse.OK
}

func (tn *ThreadNode) GetAttr(
	out *fuse.Attr,
	file nodefs.File,
	context *fuse.Context) (code fuse.Status) {
	out.Mode = 0755 | syscall.S_IFDIR
	out.Nlink = 1
	owner := getDefaultOwner()
	out.Uid = owner.Uid
	out.Gid = owner.Gid
	out.Size = uint64(tn.thread.Size())
	t := tn.thread.Time()
	out.SetTimes(&t, &t, &t)
	return fuse.OK
}
