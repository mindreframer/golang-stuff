package spec

type Component interface {
	Start() error
	Stop()
	Running() bool
}
