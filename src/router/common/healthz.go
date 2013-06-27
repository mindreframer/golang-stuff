package common

type Lockable interface {
	Lock()
	Unlock()
}

type Healthz struct {
	LockableObject Lockable
}

func (v *Healthz) Value() string {
	return "ok"
}
