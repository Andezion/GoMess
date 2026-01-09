package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ChatClient struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewChatClient(address string) (*ChatClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Не удалось подключиться к серверу: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ChatClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (c *ChatClient) Start() error {
	defer c.conn.Close()

	namePrompt, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Ошибка чтения от сервера: %w", err)
	}
	fmt.Print(namePrompt)

	nameScanner := bufio.NewScanner(os.Stdin)
	if !nameScanner.Scan() {
		return fmt.Errorf("Ошибка чтения имени")
	}
	name := nameScanner.Text()

	_, err = c.writer.WriteString(name + "\n")
	if err != nil {
		return fmt.Errorf("Ошибка отправки имени: %w", err)
	}
	c.writer.Flush()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	c.wg.Add(1)
	go c.readMessages()

	c.wg.Add(1)
	go c.writeMessages()

	select {
	case <-c.ctx.Done():
	case <-sigChan:
		fmt.Println("\nОтключение!")
		c.cancel()
	}

	c.wg.Wait()
	return nil
}

func (c *ChatClient) readMessages() {
	defer c.wg.Done()
	defer c.cancel()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))

			message, err := c.reader.ReadString('\n')
			if err != nil {
				if c.ctx.Err() == nil {
					fmt.Printf("\nСоединение с сервером потеряно: %v\n", err)
				}
				return
			}

			message = strings.TrimSpace(message)
			if message != "" {
				fmt.Printf("\r%s\n", message)
			}
		}
	}
}

func (c *ChatClient) writeMessages() {
	defer c.wg.Done()
	defer c.cancel()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		select {
		case <-c.ctx.Done():
			return
		default:
			input := scanner.Text()
			input = strings.TrimSpace(input)

			if input == "" {
				continue
			}

			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			_, err := c.writer.WriteString(input + "\n")
			if err != nil {
				fmt.Printf("\nОшибка отправки: %v\n", err)
				return
			}

			err = c.writer.Flush()
			if err != nil {
				fmt.Printf("\nОшибка отправки: %v\n", err)
				return
			}

			if input == "/quit" {
				time.Sleep(100 * time.Millisecond)
				return
			}
		}
	}
}

func main() {
	address := "localhost:8080"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	fmt.Printf("Подключение к серверу %s...\n", address)

	client, err := NewChatClient(address)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	if err := client.Start(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Отключено от сервера")
}
