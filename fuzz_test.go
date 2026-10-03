//go:build go1.18
// +build go1.18

// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testify_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzAssertEqual(f *testing.F) {
	f.Add("hello", "hello")
	f.Add("", "")
	f.Add("hello", "world")
	f.Fuzz(func(t *testing.T, expected, actual string) {
		if len(expected) > 10000 || len(actual) > 10000 {
			return
		}
		func() {
			defer func() { recover() }()
			mockT := new(testing.T)
			assert.Equal(mockT, expected, actual)
			assert.NotEqual(mockT, expected, "different-"+expected)
			assert.Contains(mockT, expected+actual, expected)
		}()
	})
}

func FuzzAssertJSON(f *testing.F) {
	f.Add(`{"a":1}`, `{"a":1}`)
	f.Add(`{"a":1}`, `{"a":2}`)
	f.Add(`[1,2,3]`, `[1,2,3]`)
	f.Add(`null`, `null`)
	f.Add(`"hello"`, `"hello"`)
	f.Fuzz(func(t *testing.T, expectedJSON, actualJSON string) {
		if len(expectedJSON) > 1<<16 || len(actualJSON) > 1<<16 {
			return
		}
		func() {
			defer func() { recover() }()
			mockT := new(testing.T)
			assert.JSONEq(mockT, expectedJSON, actualJSON)
		}()
	})
}

func FuzzAssertYAML(f *testing.F) {
	f.Add("key: value", "key: value")
	f.Add("list:\n  - a\n  - b", "list:\n  - a\n  - b")
	f.Add("", "")
	f.Fuzz(func(t *testing.T, expectedYAML, actualYAML string) {
		if len(expectedYAML) > 1<<16 || len(actualYAML) > 1<<16 {
			return
		}
		func() {
			defer func() { recover() }()
			mockT := new(testing.T)
			assert.YAMLEq(mockT, expectedYAML, actualYAML)
		}()
	})
}

func FuzzRequireInt(f *testing.F) {
	f.Add(42, 42)
	f.Add(-1, 0)
	f.Fuzz(func(t *testing.T, val, compare int) {
		func() {
			defer func() { recover() }()
			mockT := new(testing.T)
			require.NotNil(mockT, &val)
			require.GreaterOrEqual(mockT, val, compare-1)
			require.LessOrEqual(mockT, val, compare+100)
		}()
	})
}

func FuzzElementsMatch(f *testing.F) {
	f.Add("a", "b", "c")
	f.Add("", "", "")
	f.Fuzz(func(t *testing.T, a, b, c string) {
		if len(a) > 1000 || len(b) > 1000 || len(c) > 1000 {
			return
		}
		func() {
			defer func() { recover() }()
			mockT := new(testing.T)
			assert.ElementsMatch(mockT, []string{a, b, c}, []string{c, b, a})
		}()
	})
}
