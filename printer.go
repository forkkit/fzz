package main

import (
	"fmt"
	"io"
	"sync"
)

type PrintResetter interface {
	Print(string) (int, error)
	Reset()
}

func NewPrinter(target io.Writer, maxCol, maxRow int) *Printer {
	return &Printer{
		target:  target,
		maxCol:  maxCol,
		maxRow:  maxRow,
		printed: 0,
		mutex:   &sync.Mutex{},
	}
}

type Printer struct {
	target  io.Writer
	maxCol  int
	maxRow  int
	printed int
	mutex   *sync.Mutex
}

func (p *Printer) Print(line string) (n int, err error) {
	p.mutex.Lock()
	if p.printed == p.maxRow {
		p.mutex.Unlock()
		return 0, nil
	}

	if p.printed == 0 {
		fmt.Fprintf(p.target, "\n")
	}

	if len(line) > p.maxCol {
		n, err = fmt.Fprintf(p.target, "%s\n", line[:p.maxCol])
	} else {
		n, err = fmt.Fprint(p.target, line)
	}

	if err == nil {
		p.printed++
	}

	p.mutex.Unlock()

	return
}

func (p *Printer) Reset() {
	p.mutex.Lock()
	p.printed = 0
	p.mutex.Unlock()
}
