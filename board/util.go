package board

import (
	"log"
	"strings"
	"io"
	"io/ioutil"
	"net/http"
	"golang.org/x/net/html"
)

var (
	closedChan chan []string
)

func getClosedChan() <-chan []string{
	if closedChan == nil{
		closedChan = make(chan []string)
		close(closedChan)
	}
	return closedChan
}

func httpReader(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("Couldn't get URL ", url, err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return strings.NewReader(string(body)), err
}

func recsFromClasses(rc io.Reader, matchClasses [][2]string)  <-chan []string {
	out, err := newRecChan(rc, matchClasses)
	if err != nil {
		log.Println("err: ", err)
		return getClosedChan()
	}
	return out
}

func newParseChan(rc io.Reader) (<-chan *html.Node, error) {
	doc, err := html.Parse(rc)
	if err != nil {
		return nil, err
	}
	return newNodeChan(doc), nil
}

func newNodeChan(doc *html.Node) (<-chan *html.Node) {
	out := make(chan *html.Node)
	var f func(*html.Node)
	f = func(n *html.Node) {
		out <- n
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
		if n == doc {
			close(out)
		}
	}
	go f(doc)
	return out
}

func newRecChan(rc io.Reader, classes[][2]string) (<-chan []string, error) {
	nc, err := newParseChan(rc)
	if err != nil {
		return nil, err
	}
	mnc := make(chan []string)
	go func() {
		for {
			if rec := matchRec(nc, classes); rec != nil {
				mnc <- rec
			} else {
				break
			}
		}
		close(mnc)
	}()
	return mnc, nil
}

func isClassMatch(n *html.Node, cls string) bool {
	for i := range n.Attr {
		if n.Attr[i].Key == "class" {
			return strings.Contains(n.Attr[i].Val, cls)
		}
	}
	return false
}

func findAttr(n *html.Node, attr string) (ret string) {
	if attr == "" {
		for nn := range newNodeChan(n) {
			if nn.Type == html.TextNode {
				ret = ret + nn.Data + "\n"
			}
		}
		return strings.Trim(ret, "\n")
	}
	for i := range n.Attr {
		if n.Attr[i].Key == attr {
			return n.Attr[i].Val
		}
	}
	return ret
}

func matchRec(nodec <-chan *html.Node, classes [][2]string) (rec []string) {
	for i := range classes {
		for n := range nodec {
			if isClassMatch(n, classes[i][0]) {
				attr := findAttr(n, classes[i][1])
				rec = append(rec, attr)
				break;
			}
		}
		if len(rec)-1 != i {
			return nil // no more nodes to read
		}
	}
	return rec
}