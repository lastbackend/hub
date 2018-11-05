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

package views

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/registry/pkg/distribution/types"
	"time"
	"unsafe"
)

type BuildView struct{}

func (bv *BuildView) New(obj *types.Build) *Build {
	if obj == nil {
		return nil
	}
	b := new(Build)
	b.Meta = bv.ToBuildMeta(&obj.Meta)
	b.Spec = bv.ToBuildSpec(&obj.Spec)
	b.Status = bv.ToBuildStatus(&obj.Status)
	return b
}

func (obj *Build) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (bv *BuildView) NewList(list *types.BuildList) *BuildList {
	if list == nil {
		return nil
	}

	bl := new(BuildList)
	bl.Items = make([]*Build, 0)

	bl.Limit = list.Limit
	bl.Page = list.Page
	bl.Total = list.Total

	for _, v := range list.Items {
		bl.Items = append(bl.Items, bv.New(v))
	}

	return bl
}

func (obj *BuildList) ToJson() ([]byte, error) {
	if unsafe.Sizeof(obj) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(obj)
}

func (bv *BuildView) ToBuildMeta(obj *types.BuildMeta) *BuildMeta {
	bm := &BuildMeta{
		ID:      obj.ID,
		Number:  obj.Number,
		Updated: obj.Updated,
		Created: obj.Created,
	}

	bm.Labels = obj.Labels
	if obj.Labels != nil {
		bm.Labels = make(map[string]string, 0)
	}

	return bm
}

func (bv *BuildView) ToBuildSpec(obj *types.BuildSpec) *BuildSpec {
	spec := &BuildSpec{
		Image: BuildImage{
			Name: fmt.Sprintf("%s/%s", obj.Image.Owner, obj.Image.Name),
			Tag:  obj.Image.Tag,
		},
		Source: BuildSource{
			Hub:    obj.Source.Hub,
			Owner:  obj.Source.Owner,
			Name:   obj.Source.Name,
			Branch: obj.Source.Branch,
		},
		Config: BuildConfig{
			Dockerfile: obj.Config.Dockerfile,
			Context:    obj.Config.Context,
			Workdir:    obj.Config.Workdir,
			EnvVars:    obj.Config.EnvVars,
			Command:    obj.Config.Command,
		},
	}

	if obj.Source.Commit != nil {
		spec.Source.Commit = new(BuildCommit)
		spec.Source.Commit.Hash = obj.Source.Commit.Hash
		spec.Source.Commit.Username = obj.Source.Commit.Username
		spec.Source.Commit.Message = obj.Source.Commit.Message
		spec.Source.Commit.Email = obj.Source.Commit.Email
		spec.Source.Commit.Date = obj.Source.Commit.Date
	}

	return spec
}

func (bv *BuildView) ToBuildStatus(obj *types.BuildStatus) *BuildStatus {
	started := &time.Time{}
	if obj.Started.IsZero() {
		started = nil
	} else {
		started = &obj.Started
	}

	finished := &time.Time{}
	if obj.Finished.IsZero() {
		finished = nil
	} else {
		finished = &obj.Finished
	}

	return &BuildStatus{
		Size:       obj.Size,
		Step:       obj.Step,
		Message:    obj.Message,
		Status:     obj.Status,
		Done:       obj.Done,
		Processing: obj.Processing,
		Canceled:   obj.Canceled,
		Error:      obj.Error,
		Finished:   finished,
		Started:    started,
	}
}
