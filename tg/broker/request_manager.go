package broker

import (
	"errors"
	"sync"
	"time"

	"github.com/MrReality255/turbo-go/tg/utils"
)

var (
	ErrHandleConflict = errors.New("request with same ID is still in processing")
	ErrRequestTimeout = errors.New("request timeout")
)

type requestWrapper[Command ICommand] struct {
	refID         Handle
	timeout       time.Time
	handlerLocked func(responseCmd Command, err error) bool

	mx sync.Mutex
}

type requestManager[Command ICommand] struct {
	activeRequests map[Handle]*requestWrapper[Command]
	mx             sync.Mutex
	descriptor     CommandDescriptor[Command]
	senderFct      func(cmd Command) error
	receiverFct    func(cmd Command)
	timeout        time.Duration
}

func newRequestManager[Command ICommand](
	descriptor CommandDescriptor[Command],
	senderFct func(cmd Command) error,
	receiverFct func(Cmd Command),
	timeout time.Duration,
) *requestManager[Command] {
	return &requestManager[Command]{
		activeRequests: make(map[Handle]*requestWrapper[Command]),
		descriptor:     descriptor,
		senderFct:      senderFct,
		receiverFct:    receiverFct,
		timeout:        timeout,
	}
}

func (m *requestManager[Command]) Accept(msg Command) {
	refID := m.descriptor.GetRef(msg)

	m.mx.Lock()
	defer m.mx.Unlock()
	rec := m.activeRequests[refID]
	if rec == nil {
		go m.receiverFct(msg)
		return
	}

	go m.handleResponse(msg, refID, rec)
}

func (m *requestManager[Command]) Request(req Command) (Command, error) {
	ch := make(chan utils.ItemWithErr[Command], 1)
	go m.RequestMultiple(req, func(responseCmd Command, err error) bool {
		response := utils.ItemWithErr[Command]{Data: responseCmd, Err: err}
		ch <- response
		close(ch)
		return true
	})
	r := <-ch
	return r.Data, r.Err
}

func (m *requestManager[Command]) RequestMultiple(
	req Command,
	handler func(responseCmd Command, err error) bool,
) {
	var (
		dummy  Command
		isDone bool

		handle = m.descriptor.GetID(req)
		newRec = &requestWrapper[Command]{
			refID:         handle,
			handlerLocked: handler,
			timeout:       time.Now().Add(m.timeout),
		}
	)

	utils.ExecLocked(&m.mx, func() {
		if m.activeRequests[handle] != nil {
			_ = handler(dummy, ErrHandleConflict)
			isDone = true
			return
		}
		m.activeRequests[handle] = newRec
	})
	if isDone {
		return
	}
	go m.checkTimeout(handle, m.timeout+time.Millisecond, newRec)
}

func (m *requestManager[Command]) checkTimeout(
	handle Handle, sleepTime time.Duration, ref *requestWrapper[Command],
) {
	time.Sleep(sleepTime)
	m.mx.Lock()
	defer m.mx.Unlock()
	rec, ok := m.activeRequests[handle]
	if !ok || ref != rec {
		return
	}

	if n := time.Now(); n.Before(rec.timeout) {
		go m.checkTimeout(handle, rec.timeout.Sub(n)+time.Millisecond, ref)
		return
	}

	m.removeActiveRequest(rec.refID, rec, true)
	go func() {
		var dummy Command
		rec.mx.Lock()
		defer rec.mx.Unlock()
		rec.handlerLocked(dummy, ErrRequestTimeout)
	}()
}

func (m *requestManager[Command]) handleResponse(
	msg Command, id Handle, rec *requestWrapper[Command],
) {
	rec.mx.Lock()
	defer rec.mx.Unlock()
	isDone := rec.handlerLocked(msg, nil)
	if isDone {
		go m.removeActiveRequest(id, rec, false)
		return
	}

	// not done yet? fix the timeout
	rec.timeout = time.Now().Add(m.timeout)
}

func (m *requestManager[Command]) removeActiveRequest(
	id Handle, rec *requestWrapper[Command], isLocked bool,
) {
	if !isLocked {
		m.mx.Lock()
		defer m.mx.Unlock()
	}
	if m.activeRequests[id] == rec {
		delete(m.activeRequests, id)
	}
}
