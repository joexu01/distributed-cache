package http

import (
	"bytes"
	"net/http"
)

type rebalanceHandler struct {
	*Server
}

func (h *rebalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	go h.rebalanceHandler()
}

func (h *rebalanceHandler) rebalance() {
	s := h.NewScanner()
	defer s.Close()
	c := &http.Client{}
	for s.Scan() {
		k := s.Key()
		n, ok := h.ShouldProcess(k)
		if !ok {
			r, _ := http.NewRequest(http.MethodPut, "http://"+n+":12345/cache/"+k, bytes.NewReader(s.Value()))
			_, _ = c.Do(r)
			_ = h.Del(k)
		}
	}
}

func (s *Server) rebalanceHandler() http.Handler {
	return &rebalanceHandler{s}
}
