package mux

func (m *Mux) Group() *Mux {
	copy := *m
	return &copy
}
