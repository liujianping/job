package routine

import (
	"io"
)

type ReadCounter struct {
	cnt int
	rd  io.Reader
}

type WriteCounter struct {
	cnt int
	wr  io.Writer
}

func NewReadCounter(rd io.Reader) io.Reader {
	return &ReadCounter{rd: rd}
}

func (rc *ReadCounter) Read(p []byte) (int, error) {
	cnt, err := rc.rd.Read(p)
	rc.cnt += cnt
	return cnt, err
}

func (rc *ReadCounter) Count() int {
	return rc.cnt
}

func NewWriteCounter(wr io.Writer) io.Writer {
	return &WriteCounter{wr: wr}
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	cnt, err := wc.wr.Write(p)
	wc.cnt += cnt
	return cnt, err
}

func (wc *WriteCounter) Count() int {
	return wc.cnt
}
