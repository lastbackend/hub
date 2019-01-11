//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package stream

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/registry/pkg/util/stream/backend"
)

type IStream interface {
	Write(p []byte) (n int, err error)
	Pipe()
	Flush()
	Close()
	Done()
	AddSocketBackend(endpoint string) *Stream
	AddHttpBackend(w io.Writer) *Stream
	AddStdoutBackend(w io.Writer) *Stream
	AddFileBackend(dir, file string) *Stream
}

type Stream struct {
	io.Writer

	end  chan bool
	done chan bool

	buffer *bytes.Buffer
	mutex  sync.Mutex

	parts   int
	written int
	limit   int

	close bool

	timer   *time.Time
	timeout time.Duration

	streams map[backend.IStreamBackend]bool
}

type part struct {
	Chunk int    `json:"chunk"`
	Data  string `json:"data"`
}

var ErrWrotePastMaxLogLength = errors.New("wrote past max length")

func (s *Stream) Pipe() {

	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()

		for {
			select {
			case <-s.end:
				s.done <- true
				return

			case <-tick.C:
				s.Flush()
			}
		}
	}()
}

func (s *Stream) Write(p []byte) (n int, err error) {

	if s.close {
		return 0, fmt.Errorf("attempted write to closed log")
	}

	s.written += len(p)
	if s.written > s.limit {
		s.mutex.Lock()
		_, _ = s.buffer.Write([]byte(
			fmt.Sprintf("\n\nThe log length has exceeded the limit of %d MB (this usually means "+
				"hat the test suite is raising the same exception over and over).\n\nThe job has been terminated\n",
				s.limit/1000/1000)))
		s.mutex.Unlock()
		return 0, ErrWrotePastMaxLogLength
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.Write(p)
}

func (s *Stream) Flush() {

	if s.buffer.Len() <= 0 {
		return
	}

	buf := make([]byte, 1024*10) // Try to find better chunk size. Start from 1024

	for s.buffer.Len() > 0 {
		s.mutex.Lock()
		c, err := s.buffer.Read(buf)
		s.mutex.Unlock()

		if err != nil {
			log.V(6).Debug("Empty buffer returns! Panic panic!!!")
			return
		}

		p := part{
			Data:  string(buf[0:c]),
			Chunk: s.parts,
		}

		s.parts++

		body, err := json.Marshal(p)
		if err != nil {
			log.V(6).Debugf("Builder: Log marshal error: %s", err)
		}

		chunk := body

		for stream := range s.streams {
			stream.Write(chunk)
		}

	}

}

func (s *Stream) Close() {
	log.Debug("close stream connection")
	if !s.close {
		log.Debug("connection needs to be closed")
		s.Flush()

		s.mutex.Lock()
		for stream := range s.streams {
			stream.Disconnect()
			delete(s.streams, stream)
		}
		s.mutex.Unlock()
	}
	s.close = true
}

func (s *Stream) Done() {
	<-s.done
}

func (s *Stream) AddSocketBackend(endpoint string) *Stream {

	stream := backend.NewSocketBackend(endpoint)

	go func() {
		stream.End()
		s.close = true
		log.Debug("stream closed")
		s.done <- true
	}()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.streams[stream] = true
	return s
}

func (s *Stream) AddHttpBackend(w io.Writer) *Stream {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.streams[backend.NewHttpBackend(w)] = true
	return s
}

func NewStream() *Stream {
	var s = new(Stream)

	s.end = make(chan bool)
	s.done = make(chan bool)

	s.timeout = time.Second
	s.buffer = new(bytes.Buffer)
	s.limit = 1024 * 1000

	s.streams = make(map[backend.IStreamBackend]bool, 0)

	return s
}
