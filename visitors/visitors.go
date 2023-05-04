package visitors

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"zartekAssignment/variables"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

type visitors struct {
	IpChan chan string

	VisitorChan chan *visitor

	Wg sync.WaitGroup
}

func NewVisitors() *visitors {
	v := &visitors{
		IpChan:      make(chan string),
		VisitorChan: make(chan *visitor),
	}
	v.Wg.Add(2)
	go v.ProcessRequests()
	go v.UpdateVisitors()
	return v
}

func (v *visitors) ProcessRequests() {
	counts := make(map[string]int)

	t := time.NewTicker(variables.Duration)
	defer t.Stop()

	for {
		select {
		case ip := <-v.IpChan:
			counts[ip]++

			if counts[ip] > variables.MaxRequests {
				v.VisitorChan <- nil
			} else {
				v.VisitorChan <- &visitor{lastSeen: time.Now(), count: counts[ip]}
			}
		case <-t.C:
			counts = make(map[string]int)
		}
	}
}

func (v *visitors) UpdateVisitors() {
	visitors := make(map[string]*visitor)

	for {
		select {
		case ip := <-v.IpChan:
			if visitor, exists := visitors[ip]; exists {
				v.VisitorChan <- visitor
				continue
			}

			visitor := &visitor{lastSeen: time.Now(), count: 1}
			visitors[ip] = visitor
			v.VisitorChan <- visitor
		}
	}
}

func (v *visitors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.IpChan <- r.RemoteAddr

	visitor := <-v.VisitorChan

	if visitor == nil {
		reponse := map[string]string{
			"error": "Too Many Requests",
		}
		jsonData, err := json.Marshal(reponse)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write(jsonData)
		return
	}

	visitor.lastSeen = time.Now()

	response := map[string]interface{}{
		"message": "Hello, Gopher!",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
}
