package main
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"./boardfs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <mountpoint>\n", os.Args[0])
		os.Exit(1)
	}

	mountOptions := &fuse.MountOptions{
		AllowOther: true,
		Name:       "boardfs",
		Options:    []string{"default_permissions"},
	}
	mountpoint := os.Args[1]

	root := boardfs.NewRootNode()
	conn := nodefs.NewFileSystemConnector(root, &nodefs.Options{})

	server, err := fuse.NewServer(conn.RawFS(), mountpoint, mountOptions)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	// shutdown fuseserver on SIGINT
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill)
	go func() {
		sig := <-sigchan
		fmt.Print("\nExiting on ", sig, "\n")
		server.Unmount()
	}()
	server.Serve()
	signal.Stop(sigchan)
}