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

package types

import (
	"time"
)

const EmptyString = ""
const EmptyStringSlice = "[]"

type Meta struct {
	ID string `json:"id,omitempty"`
	// Meta name
	Name string `json:"name,omitempty"`
	// Meta description
	Description string `json:"description,omitempty"`
	// Meta labels
	Labels map[string]string `json:"labels,omitempty"`
	// Meta created time
	Created time.Time `json:"created,omitempty"`
	// Meta updated time
	Updated time.Time `json:"updated,omitempty"`
}

func (m *Meta) SetDefault() {
	m.Labels = make(map[string]string, 0)
	m.Created = time.Now().UTC()
	m.Updated = time.Now().UTC()
}
