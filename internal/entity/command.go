package entity

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
