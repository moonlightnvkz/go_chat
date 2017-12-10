package common

import (
	"bufio"
	"io"
	"fmt"
)

type Terminal struct {
	In     chan string
	Out    chan string
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func NewTerminal(rd io.Reader, w io.Writer) *Terminal {
	return &Terminal{
		make(chan string),
		make(chan string),
		bufio.NewReader(rd),
		bufio.NewWriter(w),
	}
}

func (term *Terminal) Read() {
	for {
		line, err := term.Reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				//fmt.Println("EOF")
				break
			}
			fmt.Println(err)
			continue
		}
		//fmt.Println("Read: " + line)
		term.In <- line
	}
}

func (term *Terminal) Write() {
	for data := range term.Out {
		//fmt.Println("Write: " + data)
		term.Writer.WriteString(data)
		term.Writer.Flush()
	}
}

func (term *Terminal) Listen() {
	go term.Read()
	go term.Write()
}
