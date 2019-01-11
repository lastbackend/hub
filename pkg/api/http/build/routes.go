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

package build

import (
	"github.com/lastbackend/registry/pkg/util/http"
	"github.com/lastbackend/registry/pkg/util/http/middleware"
)

var Routes = []http.Route{
	// Build handlers
	{Path: "/image/{owner}/{name}/build", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Authenticate}, Handler: BuildCreateH},
	{Path: "/image/{owner}/{name}/build", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: BuildListH},
	{Path: "/image/{owner}/{name}/build/{build}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: BuildInfoH},
	{Path: "/image/{owner}/{name}/build/{build}/logs", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: BuildLogsH},
	{Path: "/image/{owner}/{name}/build/{build}/cancel", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: BuildCancelH},
	{Path: "/image/{owner}/{name}/build/{build}/status", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: BuildStatusUpdateH},
}
