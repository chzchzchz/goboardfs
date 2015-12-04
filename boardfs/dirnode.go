package boardfs

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse"
	"log"
	"os/user"
	"strings"
	"strconv"
	"syscall"
	"time"
)

type NodeMaker interface {
	NodeFromName(name string) (nodefs.Node, error)
	List() []string
}

type DirNode struct {
	nodefs.Node
	NodeMaker
	time	time.Time
	owner	fuse.Owner
}

var (
	defaultOwner *fuse.Owner
)

func isBadName(name string)  bool {
	// XXX hack to avoid apache freakouts since lookups may be too permissive
	return name == ".htaccess" ||
		strings.Contains(name, ".htm") ||
		strings.Contains(name, ".phtm") ||
		strings.Contains(name, ".php") ||
		strings.Contains(name, ".var")
}

func newDirNode(o *fuse.Owner, t *time.Time, nm NodeMaker) (dn *DirNode) {
	if nm == nil {
		return nil
	}
	if o == nil {
		o = getDefaultOwner()
	}
	if t == nil {
		now := time.Now()
		t = &now
	}
	return &DirNode{
		Node : nodefs.NewDefaultNode(),
		NodeMaker : nm,
		time : *t,
		owner : *o}
}

func getDefaultOwner() (*fuse.Owner) {
	if defaultOwner != nil {
		return defaultOwner
	}
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Couldn't get current user %v", err)
	}
	uid, err := strconv.ParseInt(u.Uid, 10, 32)
	if err != nil {
		log.Fatalf("Couldn't convert uid %v", err)
	}
	gid, err := strconv.ParseInt(u.Gid, 10, 32)
	if err != nil {
		log.Fatalf("Couldn't convert gid %v", err)
	}
	defaultOwner = &fuse.Owner{Uid : uint32(uid), Gid : uint32(gid)}
	return defaultOwner
}

func (dn *DirNode) GetAttr(
	out *fuse.Attr,
	file nodefs.File,
	context *fuse.Context) (code fuse.Status) {
	out.Mode = 0755 | syscall.S_IFDIR
	out.Nlink = 2
	out.Size = 1
	out.Uid = dn.owner.Uid
	out.Gid = dn.owner.Gid
	out.SetTimes(&dn.time, &dn.time, &dn.time)
	return fuse.OK
}

func (dn *DirNode) Mkdir(
	name string,
	mode uint32,
	context *fuse.Context) (*nodefs.Inode, fuse.Status) {
	if isBadName(name) {
		return nil, fuse.ENOENT
	}
	n, err:= dn.NodeFromName(name)
	if err != nil {
		return nil, fuse.ENOENT
	}
	return dn.Inode().NewChild(name, true, n), fuse.OK
}

func (dn *DirNode) Lookup(
	out *fuse.Attr,
	name string,
	context *fuse.Context) (*nodefs.Inode, fuse.Status) {
	if (isBadName(name)) {
		return nil, fuse.ENOENT
	}

	n, err := dn.NodeFromName(name)
	if err != nil {
		return nil, fuse.ENOENT
	}

	// assumes n overrides getattr from defaultnode
	errs := n.GetAttr(out, nil, context)
	if errs != fuse.OK {
		return nil, errs
	}

	return dn.Inode().NewChild(name, true, n), fuse.OK
}

func (dn *DirNode) OpenDir(
	context *fuse.Context) (ret []fuse.DirEntry, code fuse.Status) {
	l := dn.List()
	for i := range l {
		ret = append(ret, fuse.DirEntry{Name: l[i], Mode: fuse.S_IFDIR})
	}
	return ret, fuse.OK
}
