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

package mock

import (
	_ "github.com/lib/pq"

	"context"
	"database/sql"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/registry/pkg/storage/storage"
	"github.com/lastbackend/registry/pkg/storage/store"
	"github.com/lib/pq"
	"time"
)

const (
	logLevel  = 5
	logPrefix = "storage:mock"
)

var client store.IDB

type Storage struct {
	conn string

	context.Context
	context.CancelFunc

	*BuildStorage
	*BuilderStorage
	*ImageStorage
	*SystemStorage
}

const sql_connection_name = "sqlc"

func (s *Storage) Begin(ctx context.Context) (context.Context, error) {
	s.Context, s.CancelFunc = context.WithCancel(ctx)
	connect, err := client.Begin()
	if err != nil {
		return s.Context, err
	}
	return context.WithValue(s.Context, sql_connection_name, connect), nil
}

func (s *Storage) Rollback(ctx context.Context) (context.Context, error) {
	if ctx.Value(sql_connection_name) == nil {
		return context.WithValue(s.Context, sql_connection_name, nil), nil
	}
	c := ctx.Value(sql_connection_name).(*sql.Tx)
	if err := c.Rollback(); err != nil {
		return context.WithValue(s.Context, sql_connection_name, nil), err
	}
	s.CancelFunc()
	return context.WithValue(s.Context, sql_connection_name, nil), nil
}

func (s *Storage) Commit(ctx context.Context) (context.Context, error) {
	if ctx.Value(sql_connection_name) == nil {
		return context.WithValue(s.Context, sql_connection_name, nil), nil
	}
	c := ctx.Value(sql_connection_name).(*sql.Tx)
	if err := c.Commit(); err != nil {
		return context.WithValue(s.Context, sql_connection_name, nil), err
	}
	s.Context.Done()
	return context.WithValue(s.Context, sql_connection_name, nil), nil
}

func (s *Storage) Listen(ctx context.Context, key string, listener chan string) error {

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err)
		}
	}

	l := pq.NewListener(s.conn, 10*time.Second, time.Minute, reportProblem)
	err := l.Listen(key)
	if err != nil {
		panic(err)
	}

	log.Debugf("Start monitoring PostgreSQL on key: %s", key)

	for {
		select {
		case n := <-l.Notify:
			log.Debug("Received data from channel [", n.Channel, "] :", n.Extra)
			listener <- n.Extra
			log.Debug("end send info to channel")
		case <-time.After(60 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				if err := l.Ping(); err != nil {
					fmt.Println(err)
				}
			}()
		}
	}

}

func (s *Storage) Build() storage.Build {
	if s == nil {
		return nil
	}
	return s.BuildStorage
}

func (s *Storage) Builder() storage.Builder {
	if s == nil {
		return nil
	}
	return s.BuilderStorage
}

func (s *Storage) Image() storage.Image {
	if s == nil {
		return nil
	}
	return s.ImageStorage
}

func (s *Storage) System() storage.System {
	if s == nil {
		return nil
	}
	return s.SystemStorage
}

func New() (*Storage, error) {
	s := new(Storage)
	s.conn = ""

	s.BuildStorage = newBuildStorage()
	s.BuilderStorage = newBuilderStorage()
	s.ImageStorage = newImageStorage()
	s.SystemStorage = newSystemStorage()

	return s, nil
}
