package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type MessageType int

const (
	TextMessage MessageType = iota
	JoinMessage
	LeaveMessage
	UserListMessage
)

type Message struct {
	Type      MessageType
	From      string
	Content   string
	Timestamp time.Time
}

type Client struct {
	conn     net.Conn
	name     string
	outgoing chan Message
	server   *Server
	ctx      context.Context
	cancel   context.CancelFunc
}

type Server struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewServer(ctx context.Context) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		clients:    make(map[string]*Client),
		broadcast:  make(chan Message, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (s *Server) Run() {
	defer s.cancel()

	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("Сервер останавливается!")
			s.closeAllClients()
			return

		case client := <-s.register:
			s.registerClient(client)

		case client := <-s.unregister:
			s.unregisterClient(client)

		case msg := <-s.broadcast:
			s.broadcastMessage(msg)
		}
	}
}

func (s *Server) registerClient(client *Client) {
	s.mu.Lock()
	s.clients[client.name] = client
	count := len(s.clients)
	s.mu.Unlock()

	fmt.Printf("[Сервер] Клиент '%s' подключился. Всего онлайн: %d\n", client.name, count)

	s.broadcast <- Message{
		Type:      JoinMessage,
		From:      "Сервер",
		Content:   fmt.Sprintf("%s присоединился к чату", client.name),
		Timestamp: time.Now(),
	}

	go s.sendUserList(client)
}

func (s *Server) unregisterClient(client *Client) {
	s.mu.Lock()
	if _, ok := s.clients[client.name]; ok {
		delete(s.clients, client.name)
		close(client.outgoing)
	}
	count := len(s.clients)
	s.mu.Unlock()

	fmt.Printf("[Сервер] Клиент '%s' отключился. Всего онлайн: %d\n", client.name, count)

	s.broadcast <- Message{
		Type:      LeaveMessage,
		From:      "Сервер",
		Content:   fmt.Sprintf("%s покинул чат", client.name),
		Timestamp: time.Now(),
	}
}

func (s *Server) broadcastMessage(msg Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		select {
		case client.outgoing <- msg:
		case <-time.After(1 * time.Second):
			fmt.Printf("[Предупреждение] Не удалось отправить сообщение клиенту %s\n", client.name)
		}
	}
}

func (s *Server) sendUserList(client *Client) {
	s.mu.RLock()
	var users []string
	for name := range s.clients {
		users = append(users, name)
	}
	s.mu.RUnlock()

	msg := Message{
		Type:      UserListMessage,
		From:      "Сервер",
		Content:   fmt.Sprintf("Онлайн: %s", strings.Join(users, ", ")),
		Timestamp: time.Now(),
	}

	select {
	case client.outgoing <- msg:
	case <-time.After(1 * time.Second):
	}
}

func (s *Server) closeAllClients() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, client := range s.clients {
		client.cancel()
		client.conn.Close()
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	fmt.Println("[DEBUG] Новое подключение от", conn.RemoteAddr())

	writer := bufio.NewWriter(conn)
	writer.WriteString("Введите ваше имя: ")
	writer.Flush()

	fmt.Println("[DEBUG] Отправили запрос имени")

	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[DEBUG] Ошибка чтения имени:", err)
		conn.Close()
		return
	}
	name = strings.TrimSpace(name)
	fmt.Println("[DEBUG] Получили имя:", name)

	s.mu.RLock()
	_, exists := s.clients[name]
	s.mu.RUnlock()

	if exists {
		writer := bufio.NewWriter(conn)
		writer.WriteString("Это имя уже занято. Попробуйте другое.\n")
		writer.Flush()
		conn.Close()
		return
	}

	fmt.Println("[DEBUG] Регистрируем клиента:", name)

	ctx, cancel := context.WithCancel(s.ctx)
	client := &Client{
		conn:     conn,
		name:     name,
		outgoing: make(chan Message, 10),
		server:   s,
		ctx:      ctx,
		cancel:   cancel,
	}

	s.register <- client

	fmt.Println("[DEBUG] Клиент зарегистрирован, отправляем приветствие")

	writer = bufio.NewWriter(conn)
	writer.WriteString("\n=== Добро пожаловать в GoMess! ===\n")
	writer.WriteString("Команды: /users, /help, /quit\n\n")
	writer.Flush()

	fmt.Println("[DEBUG] Запускаем горутины для клиента:", name)

	go client.writeMessages()
	go client.readMessages()
}

func (c *Client) readMessages() {
	defer func() {
		c.server.unregister <- c
		c.cancel()
		c.conn.Close()
	}()

	reader := bufio.NewReader(c.conn)
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "/") {
				c.handleCommand(line)
				continue
			}

			msg := Message{
				Type:      TextMessage,
				From:      c.name,
				Content:   line,
				Timestamp: time.Now(),
			}

			c.server.broadcast <- msg
		}
	}
}

func (c *Client) writeMessages() {
	defer c.conn.Close()

	for {
		select {
		case <-c.ctx.Done():
			return
		case msg, ok := <-c.outgoing:
			if !ok {
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			formatted := c.formatMessage(msg)
			_, err := c.conn.Write([]byte(formatted + "\n"))
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) formatMessage(msg Message) string {
	timeStr := msg.Timestamp.Format("15:04:05")

	switch msg.Type {
	case TextMessage:
		return fmt.Sprintf("[%s] %s: %s", timeStr, msg.From, msg.Content)
	case JoinMessage, LeaveMessage, UserListMessage:
		return fmt.Sprintf("[%s] *** %s ***", timeStr, msg.Content)
	default:
		return msg.Content
	}
}

func (c *Client) handleCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "/users":
		c.server.sendUserList(c)
	case "/help":
		help := `Доступные команды:
/users - показать список пользователей онлайн
/help - показать эту справку
/quit - выйти из чата`
		c.outgoing <- Message{
			Type:      UserListMessage,
			From:      "Сервер",
			Content:   help,
			Timestamp: time.Now(),
		}
	case "/quit":
		c.cancel()
	default:
		c.outgoing <- Message{
			Type:      UserListMessage,
			From:      "Сервер",
			Content:   "Неизвестная команда. Используйте /help для справки",
			Timestamp: time.Now(),
		}
	}
}

func (s *Server) Listen(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("Ошибка запуска сервера: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Сервер запущен на %s\n", address)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-s.ctx.Done():
					return
				default:
					fmt.Printf("Ошибка подключения: %v\n", err)
					continue
				}
			}

			go s.HandleConnection(conn)
		}
	}()

	<-s.ctx.Done()
	return nil
}
