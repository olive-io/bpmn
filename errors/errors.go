/*
   Copyright 2023 The bpmn Authors

   This program is offered under a commercial and under the AGPL license.
   For AGPL licensing, see below.

   AGPL licensing:
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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

type TaskExecError struct {
	Id     string
	Reason string
}

func (e TaskExecError) Error() string {
	return fmt.Sprintf("call task %s failed: %v", e.Id, e.Reason)
}

type SubProcessError struct {
	Id     string
	Reason string
}

func (e SubProcessError) Error() string {
	return fmt.Sprintf("exec subprocess %s failed: %v", e.Id, e.Reason)
}
