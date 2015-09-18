package tail

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

const (
	TESTDIR = ".test"
	FILE    = "tailed.txt"
)

func check(message string, err error, t *testing.T) {
	if err != nil {
		t.Log(message)
		t.Fatal(err)
	}
}

func setup(t *testing.T, numlines int64, line string) {
	_, err := os.Stat(TESTDIR)
	if err == nil {
		t.Fatal("test dir already exists!")
	}

	err = os.Mkdir(TESTDIR, 0777)
	check("unable to create test dir!", err, t)
	file, err := os.Create(path.Join(TESTDIR, FILE))
	check("unable to create test file!", err, t)

	var i int64 = 0
	for i = 0; i < numlines; i++ {
		_, err = file.WriteString(line + "\n")
		check("unable to write to test file!", err, t)
	}

	fmt.Println("setup done!")
}

func tear(t *testing.T) {
	err := os.RemoveAll(TESTDIR)
	check("unable to delete test dir!", err, t)
	fmt.Println("teared down!\n")
}

func TestOne(t *testing.T) {
	myline := "This is a simple line"
	var mycount int64 = 10000000
	setup(t, mycount, myline)
	defer tear(t)

	tail, err := TailFile(path.Join(TESTDIR, FILE), 100)
	check("unable to tail file!", err, t)
	fmt.Println(tail)

	var count int64 = 0
	for line := range tail.Lines {
		count++
		if line != myline {
			tail.Stop()
			t.Fatal("line does not match!")
		}

		if count == mycount {
			tail.Stop()
		}
	}
}

func TestTwo(t *testing.T) {
	myline := "This is a simple line"
	var mycount int64 = 10000000
	setup(t, mycount, myline)
	defer tear(t)

	tail, err := TailFile(path.Join(TESTDIR, FILE), 100)
	check("unable to tail file!", err, t)
	fmt.Println(tail)
	fmt.Println("waiting for 10 seconds...")
	time.Sleep(10 * time.Second)
	fmt.Println("done waiting!")

	var count int64 = 0
	for line := range tail.Lines {
		count++
		if line != myline {
			tail.Stop()
			t.Fatal("line does not match!")
		}

		if count == mycount {
			tail.Stop()
		}
	}
}
