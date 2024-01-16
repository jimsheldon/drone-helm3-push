// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import "testing"

func TestPlugin(t *testing.T) {
	t.Skip()
}

func TestVerifyArgs(t *testing.T) {
	var err error

	// RegistryUsername unset should return error
	err = verifyArgs(&Args{
		RegistryPassword:  "correct-horse-battery-staple",
		RegistryNamespace: "example",
	})
	if err == nil {
		t.Error(err)
		return
	}

	// RegistryNamespace unset should return error
	err = verifyArgs(&Args{
		RegistryPassword: "correct-horse-battery-staple",
		RegistryUsername: "kevinbacon",
	})
	if err == nil {
		t.Error(err)
		return
	}

	// RegistryPassword unset should return error
	err = verifyArgs(&Args{
		RegistryNamespace: "example",
		RegistryUsername:  "kevinbacon",
	})
	if err == nil {
		t.Error(err)
		return
	}
}
