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

package request

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"io"
	"io/ioutil"
)

type BuilderRequest struct{}

func (BuilderRequest) BuilderUpdateManifestOptions() *BuilderUpdateManifestOptions {
	return new(BuilderUpdateManifestOptions)
}

func (b *BuilderUpdateManifestOptions) Validate() *errors.Err {
	return nil
}

func (b *BuilderUpdateManifestOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("builder").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("builder").Unknown(err)
	}

	err = json.Unmarshal(body, b)
	if err != nil {
		return errors.New("builder").IncorrectJSON(err)
	}

	if err := b.Validate(); err != nil {
		return err
	}

	return nil
}

func (b BuilderUpdateManifestOptions) ToJson() ([]byte, error) {
	return json.Marshal(b)
}
