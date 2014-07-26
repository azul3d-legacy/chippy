// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
)

var defaultIcon image.Image

func init() {
	var err error

	buf := bytes.NewBuffer(defaultIconBytes)
	defaultIcon, _, err = image.Decode(buf)
	if err != nil {
		panic(fmt.Sprintf("Unable to decode default icon", err))
	}
}
