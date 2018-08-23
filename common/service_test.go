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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testService struct {
	BaseService
}

func (testService) OnReset() error {
	return nil
}

func TestBaseServiceWait(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)
	ts.Start()

	waitFinished := make(chan struct{})
	go func() {
		ts.Wait()
		waitFinished <- struct{}{}
	}()

	go ts.Stop()

	select {
	case <-waitFinished:
		// all good
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected Wait() to finish within 100 ms.")
	}
}

func TestBaseServiceReset(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *NewBaseService(nil, "TestService", ts)
	ts.Start()

	err := ts.Reset()
	require.Error(t, err, "expected cant reset service error")

	ts.Stop()

	err = ts.Reset()
	require.NoError(t, err)

	err = ts.Start()
	require.NoError(t, err)
}
