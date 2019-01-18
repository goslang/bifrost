package engine

type Route struct {
	ch chan []byte
}

func (r *Route) Close() {
	close(r.ch)
}
