package booklist

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	//b
	"os"
	//e
)

type Book struct {
	Hash     string
	FilePath string
	FileSize int64
	ModTime  time.Time
	//b
	AccessTime time.Time
	//e

	HasCover    bool
	Title       string
	Author      string
	Description string
	Series      string
	SeriesIndex float64
	Publisher   string

	ISBN		string
	PublishDate	time.Time
}

func (b *Book) ID() string {
	return b.Hash[:10]
}

//b
func (b *Book) UpdateAccessTime() time.Time {
	b.AccessTime = b.ModTime
	posFile, err := os.Open(b.FilePath + ".position")
	if err == nil {
		posFileStat, err := posFile.Stat()
		if err == nil {
			b.AccessTime = posFileStat.ModTime()
		}
	}
	return b.AccessTime
}
//e

func (b *Book) AuthorID() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(b.Author)))[:10]
}

func (b *Book) SeriesID() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(b.Series)))[:10]
}

func (b *Book) FileType() string {
	return strings.Replace(strings.ToLower(filepath.Ext(b.FilePath)), ".", "", -1)
}
