package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AnatoliyBr/go-clean-project-example/internal/entity"
	"github.com/gorilla/mux"
)

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
	cmds   chan entity.Command
}

func NewHTTPServer(config *Config, uc CounterUseCase) *HTTPServer {
	return &HTTPServer{
		config: config,
		router: mux.NewRouter(),
		uc:     uc,
		cmds:   make(chan entity.Command),
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
		case entity.SetCommand:
			val := s.uc.Set(cmd.Name, cmd.Val)
			cmd.ReplyChan <- val
		case entity.GetCommand:
			val := s.uc.Get(cmd.Name)
			cmd.ReplyChan <- val
		case entity.IncCommand:
			val := s.uc.Inc(cmd.Name)
			cmd.ReplyChan <- val
		case entity.DecCommand:
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
		s.cmds <- entity.Command{
			Ty:        entity.SetCommand,
			Name:      name,
			Val:       intval,
			ReplyChan: replyChan,
		}
		reply := <-replyChan
		fmt.Fprintf(w, "Counter '%s' with val '%d' was set\n", name, reply)
	}
}

func (s *HTTPServer) get(w http.ResponseWriter, r *http.Request) {
	log.Printf("get %v\n", r)
	name := r.URL.Query().Get("name")
	replyChan := make(chan int)
	s.cmds <- entity.Command{
		Ty:        entity.GetCommand,
		Name:      name,
		ReplyChan: replyChan,
	}

	reply := <-replyChan
	if reply >= 0 {
		fmt.Fprintf(w, "%s: %d\n", name, reply)
	} else {
		fmt.Fprintf(w, "%s not found\n", name)
	}
}

func (s *HTTPServer) inc(w http.ResponseWriter, r *http.Request) {
	log.Printf("inc %v\n", r)
	name := r.URL.Query().Get("name")
	replyChan := make(chan int)
	s.cmds <- entity.Command{
		Ty:        entity.IncCommand,
		Name:      name,
		ReplyChan: replyChan,
	}

	reply := <-replyChan
	if reply >= 0 {
		fmt.Fprintf(w, "ok, %s: %d\n", name, reply)
	} else {
		fmt.Fprintf(w, "%s not found\n", name)
	}
}

func (s *HTTPServer) dec(w http.ResponseWriter, r *http.Request) {
	log.Printf("inc %v\n", r)
	name := r.URL.Query().Get("name")
	replyChan := make(chan int)
	s.cmds <- entity.Command{
		Ty:        entity.DecCommand,
		Name:      name,
		ReplyChan: replyChan,
	}

	reply := <-replyChan
	if reply >= 0 {
		fmt.Fprintf(w, "ok, %s: %d\n", name, reply)
	} else {
		fmt.Fprintf(w, "%s not found\n", name)
	}
}
