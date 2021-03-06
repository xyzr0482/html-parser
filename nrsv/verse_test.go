package nrsv

import (
	"io"
	"os"
	"testing"
)

const (
	chapFilename = "genesis_1.html"
)

var (
	chapters = []Chapter{Chapter{
		"Genesis",
		1,
		"https://www.biblegateway.com/passage/?search=Genesis+1&version=NRSV"},
		Chapter{
			"1 Corinthians",
			12,
			"https://www.biblegateway.com/passage/?search=1+Corinthians+12&version=NRSV"},
		Chapter{"Deuteronomy",
			2,
			"https://www.biblegateway.com/passage/?search=Deuteronomy+2&version=NRSV"}}
)

func TestGetChapterText(t *testing.T) {
	records, done := make(chan VerseRecord), make(chan bool)
	nDone := 0
	for _, chap := range chapters {
		go GetVerseRecordsFromWeb(chap, records, done)
	}
	for nDone < len(chapters) {
		select {
		case vr := <-records:
			t.Log(vr)
		case <-done:
			nDone++
		}
	}
}

func TestGetTextNode(t *testing.T) {
	f := getChapterFile(t)
	defer f.Close()
	node, err := getTextNode(f)
	checkError(t, err)
	t.Logf("%#v\n", node)
}

func getChapterFile(t *testing.T) io.ReadCloser {
	f, err := os.Open(chapFilename)
	if err != nil {
		t.Error(err)
	}
	return f
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func getChapterVerses(t *testing.T, chap Chapter) {
	textNode, err := getRawVerseTextNodeFromWeb(chap)
	checkError(t, err)
	verses := getVersesFromPassageTextNode(textNode)
	for _, v := range verses {
		verseRecord, err := v.getRecord(chap)
		checkError(t, err)
		t.Log(verseRecord)
	}
	// t.Log(verses)
}
