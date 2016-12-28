package outguess

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"testing"
)

// create testfiles with original outguess (for now):
// outguess -d msg.txt hi.jpg hi_default_key.jpg
// outguess -k test -d msg.txt hi.jpg hi_test_key.jpg

func TestRetrieve(t *testing.T) {
	file, err := os.Open("testfiles/hi_default_key.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := Get(file, w, nil); err != nil {
		t.Error(err)
	}
	w.Flush()
	if b.String() != "hi" {
		t.Fail()
	}
}

func TestRetrieveCorrectKey(t *testing.T) {
	file, err := os.Open("testfiles/hi_test_key.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := Get(file, w, []byte("test")); err != nil {
		t.Error(err)
	}
	w.Flush()
	if b.String() != "hi" {
		t.Fail()
	}
}

func TestRetrieveWrongKey(t *testing.T) {
	file, err := os.Open("testfiles/hi_test_key.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := Get(file, w, []byte("wrong")); err == nil {
		t.Fail()
	}
}
