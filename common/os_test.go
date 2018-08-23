// Copyright(c) 2018 DSiSc Group All Rights Reserved.
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

package common

import (
	"os"
	"testing"
)

func TestGoPath(t *testing.T) {
	// restore original gopath upon exit
	path := os.Getenv("GOPATH")
	defer func() {
		_ = os.Setenv("GOPATH", path)
	}()

	err := os.Setenv("GOPATH", "~/testgopath")
	if err != nil {
		t.Fatal(err)
	}
	path = GoPath()
	if path != "~/testgopath" {
		t.Fatalf("should get GOPATH env var value, got %v", path)
	}
	os.Unsetenv("GOPATH")

	path = GoPath()
	if path != "~/testgopath" {
		t.Fatalf("subsequent calls should return the same value, got %v", path)
	}
}

func TestGoPathWithoutEnvVar(t *testing.T) {
	// restore original gopath upon exit
	path := os.Getenv("GOPATH")
	defer func() {
		_ = os.Setenv("GOPATH", path)
	}()

	os.Unsetenv("GOPATH")
	// reset cache
	gopath = ""

	path = GoPath()
	if path == "" || path == "~/testgopath" {
		t.Fatalf("should get nonempty result of calling go env GOPATH, got %v", path)
	}
}
