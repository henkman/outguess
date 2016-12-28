package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/henkman/outguess"
)

var (
	_help bool
	_in   string
	_out  string
	_msg  string
	_key  string
)

func init() {
	flag.BoolVar(&_help, "h", false, "help")
	flag.StringVar(&_in, "i", "", "the input jpeg")
	flag.StringVar(&_out, "o", "", "output jpeg or retrieved message")
	flag.StringVar(&_msg, "m", "", "file containing the message to embed")
	flag.StringVar(&_key, "k", "", "key")
	flag.Parse()
}

func main() {
	if _help {
		flag.Usage()
		fmt.Println("to embed message: outguess -m msg.txt -i in.jpg -o steg.jpg")
		fmt.Println("to retrieve message: outguess -i steg.jpg -o msg.txt")
		return
	}
	var in io.Reader
	if _in != "" {
		temp, err := os.Open(_in)
		if err != nil {
			log.Fatal(err)
		}
		defer temp.Close()
		in = temp
	} else {
		in = os.Stdin
	}
	var out io.Writer
	if _out != "" {
		temp, err := os.OpenFile(_out,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0750)
		if err != nil {
			log.Fatal(err)
		}
		defer temp.Close()
		out = temp
	} else {
		out = os.Stdout
	}
	var key []byte = nil
	if _key != "" {
		key = []byte(_key)
	}
	if _msg == "" {
		if err := outguess.Get(in, out, key); err != nil {
			log.Fatal(err)
		}
	} else {
		msg, err := os.Open(_msg)
		if err != nil {
			log.Fatal(err)
		}
		defer msg.Close()
		if err := outguess.Put(in, msg, key, out); err != nil {
			log.Fatal(err)
		}
	}
}
