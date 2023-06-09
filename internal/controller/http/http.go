package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CommandType int

const (
	SetCommand = iota
	GetCommand
	IncCommand
	DecCommand
)

type Command struct {
	Ty        CommandType
	Name      string
	Val       int
	ReplyChan chan int
}

type CounterUseCase interface {
	Set(string, int) int
	Get(string) int
	Inc(string) int
	Dec(string) int
}

type HTTPServer struct {
	config *Config
	router *mux.Router
	uc     CounterUseCase
	cmds   chan Command
}

func NewHTTPServer(config *Config, uc CounterUseCase) *HTTPServer {
	return &HTTPServer{
		config: config,
		router: mux.NewRouter(),
		uc:     uc,
		cmds:   make(chan Command),
	}
}

func (s *HTTPServer) StartHTTPServer() error {
	log.Println("start server...")
	s.configureRouter()
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *HTTPServer) StartCounterManager() {
	for cmd := range s.cmds {
		switch cmd.Ty {
		case SetCommand:
			val := s.uc.Set(cmd.Name, cmd.Val)
			cmd.ReplyChan <- val
		case GetCommand:
			val := s.uc.Get(cmd.Name)
			cmd.ReplyChan <- val
		case IncCommand:
			val := s.uc.Inc(cmd.Name)
			cmd.ReplyChan <- val
		case DecCommand:
			val := s.uc.Dec(cmd.Name)
			cmd.ReplyChan <- val
		default:
			log.Fatal("unknown command type", cmd.Ty)
		}
	}
}

func (s *HTTPServer) configureRouter() {
	s.router.HandleFunc("/set", s.set)
	s.router.HandleFunc("/get", s.get)
	s.router.HandleFunc("/inc", s.inc)
	s.router.HandleFunc("/dec", s.dec)
}

func (s *HTTPServer) set(w http.ResponseWriter, r *http.Request) {
	log.Printf("set %v\n", r)
	name := r.URL.Query().Get("name")
	val := r.URL.Query().Get("val")
	intval, err := strconv.Atoi(val)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	} else {
		replyChan := make(chan int)
		s.cmds <- Command{
			Ty:        SetCommand,
			Name:      name,
			Val:       intval,
			ReplyChan: replyChan,
		}
		reply := <-replyChan
		if reply == -1 {
			fmt.Fprintf(w, "counter '%s' can't be negative\n", name)
		} else {
			fmt.Fprintf(w, "counter '%s' with val '%d' was set\n", name, reply)
		}
	}
}

func (s *HTTPServer) get(w http.ResponseWriter, r *http.Request) {
	log.Printf("get %v\n", r)
	name := r.URL.Query().Get("name")
	replyChan := make(chan int)
	s.cmds <- Command{
		Ty:        GetCommand,
		Name:      name,
		ReplyChan: replyChan,
	}

	reply := <-replyChan
	if reply >= 0 {
		fmt.Fprintf(w, "'%s': %d\n", name, reply)
	} else {
		fmt.Fprintf(w, "'%s' not found\n", name)
	}
}

func (s *HTTPServer) inc(w http.ResponseWriter, r *http.Request) {
	log.Printf("inc %v\n", r)
	name := r.URL.Query().Get("name")
	replyChan := make(chan int)
	s.cmds <- Command{
		Ty:        IncCommand,
		Name:      name,
		ReplyChan: replyChan,
	}

	reply := <-replyChan
	if reply >= 0 {
		fmt.Fprintf(w, "ok, '%s': %d\n", name, reply)
	} else if reply == -1 {
		fmt.Fprintf(w, "'%s' not found\n", name)
	} else {
		fmt.Fprintf(w, "'%s' has reached its maximum\n", name)
	}
}

func (s *HTTPServer) dec(w http.ResponseWriter, r *http.Request) {
	log.Printf("inc %v\n", r)
	name := r.URL.Query().Get("name")
	replyChan := make(chan int)
	s.cmds <- Command{
		Ty:        DecCommand,
		Name:      name,
		ReplyChan: replyChan,
	}

	reply := <-replyChan
	if reply >= 0 {
		fmt.Fprintf(w, "ok, '%s': %d\n", name, reply)
	} else if reply == -1 {
		fmt.Fprintf(w, "'%s' not found\n", name)
	} else {
		fmt.Fprintf(w, "'%s' can't be negative\n", name)
	}
}
