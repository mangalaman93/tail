package tail

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Tail struct {
	Filename   string
	BufferSize int64
	Lines      chan string
	cmd        *exec.Cmd
	wait       chan bool
}

func (t *Tail) String() string {
	return fmt.Sprintf("&Tail{Filename:%s, BufferSize:%d}", t.Filename, t.BufferSize)
}

func TailFile(filepath string, buffersize int64) (*Tail, error) {
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
		Lines:      make(chan string, buffersize),
		cmd:        cmd,
		wait:       make(chan bool, 1),
	}

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			t.Lines <- scanner.Text()
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
