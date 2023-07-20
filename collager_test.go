package collager

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func getFileBytes(file string) (r []byte, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return
	}

	r = make([]byte, stat.Size())
	_, err = bufio.NewReader(f).Read(r)
	if err != nil && err != io.EOF {
		return
	}

	return
}

func TestBase(t *testing.T) {
	file1, err := getFileBytes("path/to/file_1.jpg")
	if err != nil {
		t.Error(err)
	}
	file2, err := getFileBytes("path/to/file_2.jpg")
	if err != nil {
		t.Error(err)
	}
	c := NewCollager()
	c.FromJPG(file1)
	c.FromJPG(file2)
	c.Collage(2, 1, &SaveTo{Name: "output.jpg", Type: JPG})
}
