package internal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

// Client is type for tcp-client logic
type Client struct {
	conn      net.Conn
	wg        sync.WaitGroup
	cancelCtx context.Context
	cancelFun context.CancelFunc
}

// NewClient creates new Client instance
func NewClient() *Client {
	return &Client{}
}

// Start - connects to server
func (c *Client) Start(ctx context.Context, host string, port int64, inputSrc io.Reader) (<-chan struct{}, error) {
	var err error

	dialer := &net.Dialer{}

	c.conn, err = dialer.DialContext(ctx, "tcp", host+":"+strconv.FormatInt(port, 10))
	if err != nil {
		return nil, err
	}

	c.cancelCtx, c.cancelFun = context.WithCancel(context.Background())

	c.wg.Add(1)
	go c.reader()

	c.wg.Add(1)
	go c.writer(inputSrc)

	return c.cancelCtx.Done(), nil
}

func (c *Client) reader() {
	defer c.wg.Done()

	ch := make(chan string, 5)

	go c.scanner(ch, c.conn)

	var line string
	for {
		select {
		case <-c.cancelCtx.Done():
			fmt.Println("Reader closing by context")
			return
		case line = <-ch:
			if line == "" {
				fmt.Println("Connection closed")
				c.cancelFun()
				return
			}
			fmt.Println(line)
		}
	}
}

func (c *Client) writer(inputSrc io.Reader) {
	defer c.wg.Done()

	ch := make(chan string, 5)

	go c.scanner(ch, inputSrc)

	var line string
	for {
		select {
		case <-c.cancelCtx.Done():
			fmt.Println("Writer closing by context")
			return
		case line = <-ch:
			if line == "" {
				fmt.Println("Input closed")
				c.cancelFun()
				return
			}
			_, err := c.conn.Write([]byte(line + "\n"))
			if err != nil {
				log.Printf("Fail to write")
				c.cancelFun()
				return
			}
		}
	}
}

func (c *Client) scanner(ch chan<- string, src io.Reader) {
	defer close(ch)

	scanner := bufio.NewScanner(src)

	var line string
	for {
		if !scanner.Scan() {
			break
		}
		line = scanner.Text()
		if line != "" {
			ch <- line
		}
	}
}

// Stop - stops goroutines (and waits them) and closes connection
func (c *Client) Stop() {
	if c.conn == nil {
		return
	}

	c.cancelFun()

	c.wg.Wait()

	_ = c.conn.Close()
}
