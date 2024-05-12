package server

type Router interface {
	RegisterHandler(*Server)
}
