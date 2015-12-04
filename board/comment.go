package board

import (
	"strings"
	"time"
	"fmt"
)

var (
	niceXlateMap = map[string]string {
		"&#039;" : "'",
		"/" : "|",
		"&quot;" : "\"",
		"&lt;" : "<",
		"&gt;" : ">",
	}
)

type Comment struct {
	parent	*Comment
	author	string
	when	time.Time
	title	string
	nice_title *string
	body	string
	site_key string
}

func (c *Comment) String() string {
	return fmt.Sprintf("%s @ %v\n%s\n", c.author, c.when, c.body)
}

func (c *Comment) Depth() int {
	cur := c.parent
	v := 0
	for cur != nil {
		cur = cur.parent
		v++
	}
	return v
}

func (c *Comment) Title() string {
	if c.nice_title == nil {
		nice := c.title
		for from, to := range niceXlateMap {
			nice = strings.Replace(nice, from, to, -1)
		}
		c.nice_title = &nice
	}
	return *c.nice_title
}

func (c *Comment) Time() time.Time {
	return c.when
}
