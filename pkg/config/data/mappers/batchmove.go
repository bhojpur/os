package mappers

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"path"

	types "github.com/bhojpur/os/pkg/config/data"
)

type BatchMove struct {
	From              []string
	To                string
	DestDefined       bool
	NoDeleteFromField bool
	moves             []Move
}

func (b *BatchMove) FromInternal(data map[string]interface{}) {
	for _, m := range b.moves {
		m.FromInternal(data)
	}
}

func (b *BatchMove) ToInternal(data map[string]interface{}) error {
	errors := types.Errors{}
	for i := len(b.moves) - 1; i >= 0; i-- {
		errors = append(errors, b.moves[i].ToInternal(data))
	}
	return errors.Err()
}

func (b *BatchMove) ModifySchema(s *types.Schema, schemas *types.Schemas) error {
	for _, from := range b.From {
		b.moves = append(b.moves, Move{
			From:              from,
			To:                path.Join(b.To, from),
			DestDefined:       b.DestDefined,
			NoDeleteFromField: b.NoDeleteFromField,
		})
	}

	for _, m := range b.moves {
		if err := m.ModifySchema(s, schemas); err != nil {
			return err
		}
	}

	return nil
}
