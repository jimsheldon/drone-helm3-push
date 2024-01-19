// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testRegistryNamespace = "example"
	testRegistryPassword  = "correct-horse-battery-staple"
	testRegistryUsername  = "kevinbacon"
)

func TestVerifyArgs(t *testing.T) {
	var err error

	// RegistryUsername unset should return error
	err = verifyArgs(&Args{
		RegistryPassword:  testRegistryPassword,
		RegistryNamespace: testRegistryNamespace,
	})
	if err == nil {
		t.Error(err)
	}

	// RegistryNamespace unset should return error
	err = verifyArgs(&Args{
		RegistryPassword: testRegistryPassword,
		RegistryUsername: testRegistryUsername,
	})
	if err == nil {
		t.Error(err)
	}

	// RegistryPassword unset should return error
	err = verifyArgs(&Args{
		RegistryNamespace: testRegistryNamespace,
		RegistryUsername:  testRegistryUsername,
	})
	if err == nil {
		t.Error(err)
	}

	// RegistryNamespace, RegistryPassword and RegistryUsername set should not return error
	err = verifyArgs(&Args{
		RegistryNamespace: testRegistryNamespace,
		RegistryPassword:  testRegistryPassword,
		RegistryUsername:  testRegistryUsername,
	})
	if err != nil {
		t.Error(err)
	}
}

// these tests are based on https://github.com/helm/helm/blob/v3.13.3/cmd/helm/package_test.go
func TestPackackageChart(t *testing.T) {
	tests := []struct {
		name      string
		chartPath string
		chartDest string
		err       bool
	}{
		{
			name:      "package testdata/testcharts/alpine",
			chartPath: "testdata/testcharts/alpine",
			chartDest: "alpine-0.1.0.tgz",
		},
		{
			name:      "package testdata/testcharts/issue1979",
			chartPath: "testdata/testcharts/issue1979",
			chartDest: "alpine-0.1.0.tgz",
		},
		{
			name:      "chart with missing repo dependencies",
			chartPath: "testdata/testcharts/chart-missing-deps",
			err:       true,
		},
		{
			name:      "chart with bad type",
			chartPath: "testdata/testcharts/chart-bad-type",
			err:       true,
		},
	}

	for _, tt := range tests {
		tempDir := t.TempDir()

		args := &Args{
			ChartPath:        tt.chartPath,
			ChartDestination: tempDir,
		}

		got, err := packageChart(args)
		if err != nil {
			// return if an error was expected
			if tt.err {
				return
			}
			t.Fatal(err)
		}

		want := filepath.Join(tempDir, tt.chartDest)
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	}
}
