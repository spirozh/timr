package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/spirozh/timr"
)

type timerMessage struct {
	Id    *int             `json:"id"`
	Name  *string          `json:"name"`
	State *timr.TimerState `json:"state"`
}

type TimerHandler struct {
	path string
	timr.TimerService
	idMethods   map[string]http.HandlerFunc
	noIdMethods map[string]http.HandlerFunc
}

func newTimerHandler(path string, ts timr.TimerService) *TimerHandler {
	th := &TimerHandler{
		path:         path,
		TimerService: ts,
	}

	th.noIdMethods = map[string]http.HandlerFunc{
		http.MethodPost: th.createTimer,
		http.MethodGet:  th.listTimers,
	}

	th.idMethods = map[string]http.HandlerFunc{
		http.MethodGet:    th.getTimer,
		http.MethodHead:   th.getTimer,
		http.MethodPatch:  th.updateTimer,
		http.MethodDelete: th.deleteTimer,
	}
	return th
}

func (th *TimerHandler) createTimer(w http.ResponseWriter, r *http.Request) {
	// createTimer POST /api/timer/
	//
	// receives {"name":"zzz","state":{"duration":1,"remaining":1,"running":true}}
	// emits {"id":1,"name":"zzz","state":{"duration":1,"remaining":1,"running":true}}
	createReq, err := timerBody(w, r)
	if err != nil {
		return
	}
	if createReq.Id != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if createReq.Name == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, _ := th.Create(*createReq.Name, *createReq.State)
	_, timer, _ := th.Get(id)
	s := timer.State()

	res, err := json.Marshal(timerMessage{&id, createReq.Name, &s})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (th *TimerHandler) listTimers(w http.ResponseWriter, r *http.Request) {
	// list timers
}

func (th *TimerHandler) getTimer(w http.ResponseWriter, r *http.Request) {
	// get the context values
	id, name, timer := getIdNameTimer(r.Context())
	state := timer.State()
	b, err := json.Marshal(timerMessage{&id, &name, &state})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func getIdNameTimer(ctx context.Context) (id int, name string, timer timr.Timer) {
	var ok bool

	if id, ok = ctx.Value(keyId).(int); !ok {
		panic("id not in context")
	}
	if name, ok = ctx.Value(keyName).(string); !ok {
		panic("name not in context")
	}
	if timer, ok = ctx.Value(keyTimer).(timr.Timer); !ok {
		panic("name not in context")
	}

	return
}

func (th *TimerHandler) updateTimer(w http.ResponseWriter, r *http.Request) {
	// updateTimer PATCH /api/timer/:id
	//
	// receives {"name":"mmm","state":{"remaining":5}}
	// emits {"id":1,"name":"mmm","state":{"duration":1,"remaining":5,"running":false}}
	id, name, timer := getIdNameTimer(r.Context())

	updateReq, err := timerBody(w, r)
	if err != nil || updateReq.Id == nil || *updateReq.Id != id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if updateReq.Name == nil {
		updateReq.Name = &name
	}

	if updateReq.State == nil {
		state := timer.State()
		updateReq.State = &state
	}

	// update the timer...
	th.TimerService.Update(id, *updateReq.Name, *updateReq.State)
	s := timer.State()

	res, err := json.Marshal(timerMessage{&id, updateReq.Name, &s})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func (th *TimerHandler) deleteTimer(w http.ResponseWriter, r *http.Request) {
	id, _, _ := getIdNameTimer(r.Context())
	err := th.TimerService.Remove(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(204)
}

type timerContextKey struct{}

var (
	keyId    = timerContextKey{}
	keyName  = timerContextKey{}
	keyTimer = timerContextKey{}
)

func handle405(methodMap map[string]http.HandlerFunc, w http.ResponseWriter) {
	var methodsAllowed []string

	for k := range methodMap {
		methodsAllowed = append(methodsAllowed, k)
	}

	w.Header().Add("Allow", strings.Join(methodsAllowed, ", "))
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (th TimerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// get id
	idStr, prefixFound := strings.CutPrefix(th.path, r.URL.Path)
	if !prefixFound {
		panic(fmt.Sprintf("bad configuration, routed '%s' to handler of '%s'", r.URL.Path, th.path))
	}

	if len(idStr) > 0 {
		id, err := strconv.Atoi(idStr)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		name, t, err := th.TimerService.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if handler, exists := th.idMethods[r.Method]; exists {
			ctx := r.Context()
			ctx = context.WithValue(ctx, keyId, id)
			ctx = context.WithValue(ctx, keyName, name)
			ctx = context.WithValue(ctx, keyTimer, t)

			handler(w, r.WithContext(ctx))
			return
		}

		handle405(th.idMethods, w)
		return

	}

	if handler, exists := th.noIdMethods[r.Method]; exists {
		handler(w, r)
		return
	}

	handle405(th.noIdMethods, w)
}

func timerBody(w http.ResponseWriter, r *http.Request) (timerMessage, error) {
	var (
		msg timerMessage
		b   bytes.Buffer
	)
	io.Copy(&b, r.Body)
	err := json.Unmarshal(b.Bytes(), &msg)
	return msg, err
}

// SSE handles /api/sse/
func SSE(ts timr.TimerService) http.HandlerFunc {
	// TODO: handle 405s

	count := 0

	return func(w http.ResponseWriter, r *http.Request) {
		count++
		defer func() { count-- }()
		timr.INFO("starting SSE handler (", count, ")")

		flusher, ok := sseSetup(w, r)
		if !ok {
			return
		}

		// first connection: send json for all timers
		//m := map[string]timr.TimerState{}
		ts.ForAll(func(id int, name string, state timr.TimerState) {
			outputTimer(w, ts, id)
		})
		flusher.Flush()

		// after each state change, send json for just one timer
		//
		// either:
		//  {'name': {'running': bool, 'remaining': milliseconds}}
		// or (for delete)
		//  {'name': null}
		mu := sync.Mutex{}
		tsEventHandler := func(eventType timr.TimrEventType, id int, name string, timer timr.Timer) {
			mu.Lock()
			defer mu.Unlock()

			outputTimer(w, ts, id)

			flusher.Flush()
		}
		sub := ts.Subscribe(tsEventHandler)
		defer ts.Unsubscribe(sub)

		<-r.Context().Done()

		timr.INFO("ending SSE handler")
	}
}

func outputTimer(w io.Writer, ts timr.TimerService, id int) {
	name, t, _ := ts.Get(id)
	state := t.State()
	msg := timr.TimerMessage{Name: name, State: state}
	j, err := json.Marshal(msg)
	if err == nil {
		output(w, string(j))
	}
}

func output(w io.Writer, s string) {
	io.WriteString(w, "data: ")
	io.WriteString(w, s)
	io.WriteString(w, "\n\n")
}

func sseSetup(w http.ResponseWriter, r *http.Request) (flusher http.Flusher, ok bool) {
	// Make sure that the writer supports flushing.
	flusher, ok = w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher.Flush()
	return
}
