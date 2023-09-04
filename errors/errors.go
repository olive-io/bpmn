// Copyright 2023 Lack (xingyys@gmail.com).
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

package errors

import "fmt"

type InvalidArgumentError struct {
	Expected interface{}
	Actual   interface{}
}

func (e InvalidArgumentError) Error() string {
	return fmt.Sprintf("Invalid argument: expected %v, got %v", e.Expected, e.Actual)
}

type InvalidStateError struct {
	Expected interface{}
	Actual   interface{}
}

func (e InvalidStateError) Error() string {
	return fmt.Sprintf("Invalid state: expected %v, got %v", e.Expected, e.Actual)
}

type NotFoundError struct {
	Expected interface{}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%v not found", e.Expected)
}

type RequirementExpectationError struct {
	Expected interface{}
	Actual   interface{}
}

func (e RequirementExpectationError) Error() string {
	return fmt.Sprintf("Requirement expectation failed: expected %v, got %v", e.Expected, e.Actual)
}

type NotSupportedError struct {
	What   string
	Reason string
}

func (e NotSupportedError) Error() string {
	return fmt.Sprintf("%s is not supported because %s", e.What, e.Reason)
}

type TaskCallError struct {
	Type   string
	Id     string
	Reason string
}

func (e TaskCallError) Error() string {
	return fmt.Sprintf("call task %s-%s failed: %v", e.Type, e.Id, e.Reason)
}
