// Copyright 2019 Intel Corporation. All Rights Reserved.
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

package memtier

import (
	"fmt"
	"time"

	config "github.com/intel/cri-resource-manager/pkg/config"
	"github.com/intel/cri-resource-manager/pkg/topology"
)

// Duration is an alias for time.Duration.
type Duration time.Duration

// Options captures our configurable policy parameters.
type options struct {
	// PinCPU controls CPU pinning in the memtier policy.
	PinCPU bool
	// PinMemory controls memory pinning in the memtier policy.
	PinMemory bool
	// PreferIsolated controls whether isolated CPUs are preferred for isolated allocations.
	PreferIsolated bool `json:"PreferIsolatedCPUs"`
	// PreferShared controls whether shared CPU allocation is always preferred by default.
	PreferShared bool `json:"PreferSharedCPUs"`
	// FakeHints are the set of fake TopologyHints to use for testing purposes.
	FakeHints fakehints `json:",omitempty"`

	DirtyBitScanPeriod Duration `json:"DirtyBitScanPeriod"`
	PageMovePeriod     Duration `json:"PageMovePeriod"`
	PageMoveCount      uint     `json:"PageMoveCount"`
}

// MarshalJSON converts Duration to JSON string.
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Duration(d).String() + "\""), nil
}

// UnmarshalJSON converts JSON string to Duration.
func (d *Duration) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid Duration data")
	}
	parsed, err := time.ParseDuration(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*d = Duration(parsed)
	return nil
}

// String returns the value of Duration as a string.
func (d *Duration) String() string {
	return time.Duration(*d).String()
}

// Our runtime configuration.
var opt = defaultOptions().(*options)

// fakeHints is our flag.Value for per-pod or per-container faked topology.Hints.
type fakehints map[string]topology.Hints

// newFakeHints creates a new set of fake hints.
func newFakeHints() fakehints {
	return make(fakehints)
}

// merge merges the given hints to the existing set.
func (fh *fakehints) merge(hints fakehints) {
	if fh == nil {
		*fh = newFakeHints()
	}
}

// defaultOptions returns a new options instance, all initialized to defaults.
func defaultOptions() interface{} {
	return &options{
		PinCPU:             true,
		PinMemory:          true,
		PreferIsolated:     true,
		PreferShared:       false,
		FakeHints:          make(fakehints),
		DirtyBitScanPeriod: 0,
		PageMovePeriod:     0,
		PageMoveCount:      0,
	}
}

// Register us for configuration handling.
func init() {
	config.Register(PolicyPath, PolicyDescription, opt, defaultOptions)
}