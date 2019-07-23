/*
 * Copyright 2019 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ristretto

import (
	"fmt"
	"testing"
)

type PolicyCreator func(uint64, bool) Policy

func GeneratePolicyTest(create PolicyCreator) func(*testing.T) {
	iterations := uint64(1024)
	return func(t *testing.T) {
		t.Run("uniform-push", func(t *testing.T) {
			policy := create(iterations, false)
			values := make([]Element, iterations)
			for i := range values {
				values[i] = Element(fmt.Sprintf("%d", i))
			}
			policy.Add("0", 1)
			policy.Push(values)
			if !policy.Has("0") || policy.Has("*") {
				t.Fatal("add/push error")
			}
		})
		t.Run("uniform-add", func(t *testing.T) {
			policy := create(iterations, false)
			for i := uint64(0); i < iterations; i++ {
				policy.Add(fmt.Sprintf("%d", i), 1)
			}
			if victims, added := policy.Add("*", 1); victims == nil || !added {
				t.Fatal("add/eviction error")
			}
		})
		t.Run("variable-push", func(t *testing.T) {
			policy := create(iterations, true)
			values := make([]Element, iterations)
			for i := range values {
				values[i] = Element(fmt.Sprintf("%d", i))
			}
			policy.Add("0", 1)
			policy.Push(values)
			if !policy.Has("0") || policy.Has("*") {
				t.Fatal("add/push error")
			}
		})
		t.Run("variable-add", func(t *testing.T) {
			policy := create(iterations, true)
			for i := uint64(0); i < iterations; i++ {
				policy.Add(fmt.Sprintf("%d", i), i)
			}
			if victims, added := policy.Add("*", 1); victims == nil || !added {
				t.Fatal("add/eviction error")
			}
		})
	}
}

func TestPolicy(t *testing.T) {
	policies := []PolicyCreator{NewPolicy}
	for _, policy := range policies {
		GeneratePolicyTest(policy)(t)
	}
}
