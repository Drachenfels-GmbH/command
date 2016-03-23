// Copyright 2016 Drachenfels GmbH. All Rights Reserved.
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

package command

import (
	"testing"
)

func TestMissingCommand(t *testing.T) {
	t.Skip("Not implemented.")
}

func TestMissingParams(t *testing.T) {
	t.Skip("Not implemented.")
}

func TestInvalidParams(t *testing.T) {
	t.Skip("Not implemented.")
}

func TestCommandFunc(t *testing.T) {
	var val string
	Add(
		"hello",
		"simple testing function",
		CmdFunc(func(args []string) error {
			val = args[0]
			return nil

		}),
	)

	_, err := Run("hello", "world")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if val != "world" {
		t.Fatalf("Command should set val to %q but was %q.", "world", val)
	}
}
