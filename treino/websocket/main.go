/*
    Abordagem para gerenciar WebSockets (será a melhor do mundo?):

	A cada intervalo de tempo atualizamos a "cópia" do mapa de conexões que o loop de broadcast
	(que no mundo real estaria obtendo as mensagens de um broker).

	Evitamos assim interagir com um mapa que está sendo constantemente atualizado, mas também
	evitamos o bloqueio do mapa a cada vez que ele é alterado. A cópia do mapa é feita em um
	atomic.Value, mais rápido que o Mutex.

	Novas conexões demoram um intervalo de 5 segundos para serem consideradas; conexões removidas
	demoram esse mesmo período para serem desconsideradas (tentativas de envio de mensagens sujeitas
	a tratamento de erro).
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	mainMap      map[*websocket.Conn]struct{}
	mutex        sync.Mutex
	broadcastRef atomic.Value
}

func NewServer() *Server {
	s := &Server{
		mainMap: make(map[*websocket.Conn]struct{}),
	}
	s.broadcastRef.Store(make(map[*websocket.Conn]struct{})) // Cópia atômica do mapa de conexões
	return s
}

func (s *Server) AddConnection(conn *websocket.Conn) {
	s.mainMap[conn] = struct{}{}
}

func (s *Server) RemoveConnection(conn *websocket.Conn) {
	delete(s.mainMap, conn)
}

// Cria uma nova cópia para o broadcast
func (s *Server) UpdateBroadcastMap() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newCopy := make(map[*websocket.Conn]struct{}, len(s.mainMap))
	for conn := range s.mainMap {
		newCopy[conn] = struct{}{}
	}
	s.broadcastRef.Store(newCopy)
}

// Envio das mensagens (simula leitura em broker)
func (s *Server) BroadcastLoop() {
	count := 1

	for {
		currentMap := s.broadcastRef.Load().(map[*websocket.Conn]struct{})
		msg := fmt.Sprintf("Mensagem %d do broker", count)
		count++

		for conn := range currentMap {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Println("Erro ao enviar mensagem:", err)
				s.RemoveConnection(conn)
			}
		}
		time.Sleep(1 * time.Second) // Intervalo de envio
	}
}

func main() {
	server := NewServer()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Erro ao criar conexão WebSocket:", err)
			return
		}
		server.AddConnection(conn)
	})

	go server.BroadcastLoop()
	go server.PeriodicUpdate()

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

