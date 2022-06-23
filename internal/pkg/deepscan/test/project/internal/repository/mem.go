package repository

type Memory struct {
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m Memory) Fetch() {
	panic("implement me")
}
