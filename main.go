// Copyright © 2019 Jay Liu <liujianping.itech@qq.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/liujianping/job/build"
	"github.com/liujianping/job/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

//go:generate protoc pb/job.proto --go_out=${GOPATH}/src
func main() {
	build.Info(version, commit, date)
	cmd.Execute()
}
