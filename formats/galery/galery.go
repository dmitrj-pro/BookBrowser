package galery

import (
	"crypto/sha1"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
//	"sort"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/geek1011/BookBrowser/booklist"
	"github.com/geek1011/BookBrowser/formats"
	"github.com/geek1011/BookBrowser/util"
	"github.com/pkg/errors"
)

type galery struct {
	book *booklist.Book
	folder string
	coverpath string
}

func (e *galery) Book() *booklist.Book {
	return e.book
}

func (e *galery) HasCover() bool {
	return true
}

func (e *galery) GetCover() (i image.Image, err error) {
	f, err:= os.Open(e.coverpath)
	defer f.Close()
	i, _, _ = image.Decode(f)
	return i, nil
}

func load(filename string) (bi formats.BookInfo, ferr error) {
	defer func() {
		if r := recover(); r != nil {
			bi = nil
			ferr = fmt.Errorf("unknown error: %s", r)
		}
	}()

	p := &galery{book: &booklist.Book{}}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, errors.Wrapf(err, "could not stat book")
	}
	p.book.FilePath = filename
	p.book.FileSize = fi.Size()
	p.book.ModTime = fi.ModTime()

	s := sha1.New()
	i, err := io.Copy(s, f)
	if err == nil && i != fi.Size() {
		err = errors.New("could not read whole file")
	}
	if err != nil {
		f.Close()
		return nil, errors.Wrap(err, "could not hash book")
	}
	p.book.Hash = fmt.Sprintf("%x", s.Sum(nil))

	f.Close()

	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	str := string(c)
	c = []byte{}

	folder := ""//util.StringBetween(str, "<folder>", "</folder>")
	name := util.StringBetween(str, "<name>", "</name>")
	cover := util.StringBetween(str, "<cover>", "</cover>")
	author := util.StringBetween(str, "<author>", "</author>")


	p.book.Title = name
	p.book.Author = author
	if (len(cover) > 0) {
		p.coverpath = cover
	}
	p.book.HasCover = true
	p.folder = folder
	if (len(p.folder) == 0) {
		p.folder = strings.ReplaceAll(filename, ".galery", "")
	} else {
		p.folder = filepath.Dir(filename) + "/" + p.folder
	}
	if stat_,err_ := os.Stat(p.folder); err_==nil && stat_.IsDir() {

	} else {
		return nil, errors.Wrapf(err, "could not find galery folder")
	}

	existsCover := false
	if (len(cover) > 0) {
		p.coverpath = p.folder + "/" + cover
		_, err:= os.Open(p.coverpath)
		if (err == nil) {
			existsCover = true
		}
	}	
	if !existsCover {
		lst, err := ioutil.ReadDir(p.folder)
		if (err != nil) {
			return nil, errors.New("folder is not found")
		}
		//sort.Strings(lst)

		_, err = os.Open(p.folder + "/" + lst[0].Name())
		if (err != nil) {
			return nil, errors.New("galery has empty folder")
		} else {
			p.coverpath = p.folder + "/" + lst[0].Name()
		}
	}

	debug.FreeOSMemory()

	return p, nil
}

func init() {
	formats.Register("galery", load)
}
