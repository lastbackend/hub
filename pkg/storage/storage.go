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

package storage

import (
	"github.com/lastbackend/registry/pkg/storage/mock"
	"github.com/lastbackend/registry/pkg/storage/pgsql"
	"github.com/lastbackend/registry/pkg/storage/types/filter"
)

func Get(c string) (IStorage, error) {
	var driver = ""
	switch driver {
	case "mock":
		return mock.New()
	default:
		return pgsql.New(c)
	}
}

func Filter() IFilter {
	return filter.NewFilter()
}
