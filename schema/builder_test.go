/*
   Copyright 2025 The  Authors

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

package schema

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandBytes(t *testing.T) {
	t.Run("RandBytes", func(t *testing.T) {
		out := RandBytes(2)
		assert.Equal(t, len(out), 2)
		out1 := RandBytes(2)
		assert.NotEqual(t, out, out1)
		t.Log("out:", string(out))

		out2 := RandBytes(10)
		assert.NotEqual(t, out, out2)
		t.Log("out:", string(out2))
	})
}

func TestDefinitionsWithProcess(t *testing.T) {
	db := NewDefinitionsBuilder()
	pb := NewProcessBuilder()
	task := Task{}
	pb.AddActivity(&task)
	script := ScriptTask{}
	pb.AddActivity(&script)

	db.AddProcess(*pb.Out())

	def := db.Out()

	assert.Equal(t, len(def.ProcessField), 1)

	data, _ := xml.MarshalIndent(def, "", " ")
	t.Log(string(data))
}
