package system

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

import "path/filepath"

const (
	// DefaultRootDir represents where persistent installations are located
	DefaultRootDir = "/bhojpur/system"
	// DefaultDataDir represents where persistent state is located
	DefaultDataDir = "/bhojpur/data"
	// DefaultLocalDir represents where local, persistent configuration is located
	DefaultLocalDir = "/var/lib/bhojpur/os"
	// DefaultStateDir represents where ephemeral state is located
	DefaultStateDir = "/run/bos"
)

var (
	rootDirectory  = DefaultRootDir
	dataDirectory  = DefaultDataDir
	localDirectory = DefaultLocalDir
	stateDirectory = DefaultStateDir
)

// RootPath joins any number of elements into a single path underneath the persistent installation root, by default `DefaultRootDir`
func RootPath(elem ...string) string {
	return filepath.Join(rootDirectory, filepath.Join(elem...))
}

// DataPath joins any number of elements into a single path underneath the persistent state root, by default `DefaultDataDir`
func DataPath(elem ...string) string {
	return filepath.Join(dataDirectory, filepath.Join(elem...))
}

// LocalPath joins any number of elements into a single path underneath the persistent configuration root, by default `DefaultLocalDir`
func LocalPath(elem ...string) string {
	return filepath.Join(localDirectory, filepath.Join(elem...))
}

// StatePath joins any number of elements into a single path underneath the ephemeral state root, by default `DefaultStateDir`
func StatePath(elem ...string) string {
	return filepath.Join(stateDirectory, filepath.Join(elem...))
}
