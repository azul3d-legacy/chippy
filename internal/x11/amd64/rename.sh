# Copyright 2014 The Azul3D Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

ar x libxkbcommon-x11.a
ar x libxkbcommon.a

# Get rid of duplicate symbols
rm libxkbcommon_x11_la-atom.o
rm libxkbcommon_x11_la-context-priv.o
rm libxkbcommon_x11_la-keymap-priv.o

ld -r *.o -o ../xkbcommon_amd64.syso
rm -rf *.o
