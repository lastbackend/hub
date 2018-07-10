//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package v1

import (
	"github.com/lastbackend/registry/pkg/api/client/http/request"
	"github.com/lastbackend/registry/pkg/api/client/types"
)

type Client struct {
	client *request.RESTClient
}

func New(client *request.RESTClient) *Client {
	return &Client{client: client}
}

func (s *Client) Build() types.BuildClientV1 {
	if s == nil {
		return nil
	}
	return newBuildClient(s.client)
}

func (s *Client) Builder() types.BuilderClientV1 {
	if s == nil {
		return nil
	}
	return newBuilderClient(s.client)
}

func (s *Client) Image(owner, name string) types.ImageClientV1 {
	if s == nil {
		return nil
	}
	return newImageClient(s.client, owner, name)
}

func (s *Client) Registry() types.RegistryClientV1 {
	if s == nil {
		return nil
	}
	return newRegistryClient(s.client)
}