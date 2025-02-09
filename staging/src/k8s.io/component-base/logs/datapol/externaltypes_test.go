/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package datapol

import (
	"crypto/x509"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	testcases := []struct {
		value  any
		expect []string
	}{{
		value:  http.Header{},
		expect: []string{"password", "token"},
	}, {
		value:  http.Cookie{},
		expect: []string{"token"},
	}, {
		value:  x509.Certificate{},
		expect: []string{"security-key"},
	}}
	for _, tc := range testcases {
		types := GlobalDatapolicyMapping(tc.value)
		if !assert.ElementsMatch(t, tc.expect, types) {
			t.Errorf("Wrong set of datatypes detected for %T, want: %v, got %v", tc.value, tc.expect, types)
		}
	}
}
