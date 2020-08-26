package util

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"os"
	"time"
)

var (
	pid = os.Getpid()
)

// NewTraceID 创建追踪ID
func NewTraceID() string {
	return fmt.Sprintf("trace-id-%d-%s",
		pid,
		time.Now().Format("2006.01.02.15.04.05.999999"))
}

func RandomInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}

func GetSelect(doc *goquery.Document, sel, attr string) (string, bool) {
	return doc.Find(sel).Attr(attr)
}
