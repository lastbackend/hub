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

package request

type BuildListOptions struct {
	Active *bool  `json:"active"`
	Page   *int64 `json:"page"`
	Limit  *int64 `json:"limit"`
}

type BuildCreateOptions struct {
	Tag        string            `json:"tag"`
	Auth       string            `json:"auth"`
	DockerFile string            `json:"dockerfile"`
	Context    string            `json:"context"`
	EnvVars    []string          `json:"environments"`
	Command    string            `json:"command"`
	Workdir    string            `json:"workdir"`
	Source     ImageSource       `json:"source"`
	Labels     map[string]string `json:"labels"`
}

type BuildUpdateStatusOptions struct {
	Step     string `json:"step"`
	Message  string `json:"message"`
	Error    bool   `json:"error"`
	Canceled bool   `json:"canceled"`
}

type BuildSetImageInfoOptions struct {
	Hash        string `json:"id"`
	Size        int64  `json:"size"`
	VirtualSize int64  `json:"virtual_size"`
}

type BuildLogsOptions struct {
	Follow bool `json:"follow"`
}
