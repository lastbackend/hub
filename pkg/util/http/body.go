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

package http

import (
	"bytes"
	"encoding/json"
	"io"
)

func encodeBodyJSON(bodyJSON interface{}) (io.Reader, error) {

	var buf = new(bytes.Buffer)
	if bodyJSON != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(bodyJSON)
		if err != nil {
			return nil, err
		}
		//fmt.Fprintf(os.Stdout, "JSON %s", buf.String())
	}
	return buf, nil

}
