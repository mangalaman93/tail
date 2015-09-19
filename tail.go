package tail

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	QUEUE_SIZE = 100
)

type Tail struct {
	Filename   string
	BufferSize int
	Lines      chan string
	cmd        *exec.Cmd
	wait       chan bool
}

func (t *Tail) string() string {
	return fmt.Sprintf("&Tail{Filename:%s, BufferSize:%d}", t.Filename, t.BufferSize)
}

func TailFile(filepath string, buffersize int) (*Tail, error) {
	// check whether the file exists
	_, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("tail", "-c", "+1", "-f", filepath)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	t := &Tail{
		Filename:   filepath,
		BufferSize: buffersize,
		Lines:      make(chan string, QUEUE_SIZE),
		cmd:        cmd,
		wait:       make(chan bool, 1),
	}

	go func() {
		bigreader := bufio.NewReaderSize(reader, buffersize)
		line, isPrefix, err := bigreader.ReadLine()
		for err == nil && !isPrefix {
			t.Lines <- string(line)
			line, isPrefix, err = bigreader.ReadLine()
		}

		if isPrefix {
			log.Println("buffer size is too small!")
		}

		if err != io.EOF {
			log.Println(err)
		}

		close(t.Lines)
		t.wait <- true
	}()

	return t, nil
}

func (t *Tail) Stop() {
	t.cmd.Process.Signal(syscall.SIGINT)
	timeout := time.After(2 * time.Second)
	select {
	case <-t.wait:
	case <-timeout:
		t.cmd.Process.Kill()
		<-t.wait
	}

	close(t.wait)
}
