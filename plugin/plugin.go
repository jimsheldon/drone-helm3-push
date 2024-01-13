// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/registry"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	ChartPath         string `envconfig:"PLUGIN_CHART_PATH"`
	ChartDestination  string `envconfig:"PLUGIN_CHART_DESTINATION"`
	RegistryNamespace string `envconfig:"PLUGIN_REGISTRY_NAMESPACE"`
	RegistryPassword  string `envconfig:"PLUGIN_REGISTRY_PASSWORD"`
	RegistryURL       string `envconfig:"PLUGIN_REGISTRY_URL"`
	RegistryUsername  string `envconfig:"PLUGIN_REGISTRY_USERNAME"`
}

type registryLoginOptions struct {
	username             string
	password             string
	passwordFromStdinOpt bool
	certFile             string
	keyFile              string
	caFile               string
	insecure             bool
}

var errConfiguration = errors.New("configuration error")

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	if err := verifyArgs(&args); err != nil {
		return err
	}

	// package the chart
	packageRun, err := packageChart(&args)
	if err != nil {
		return err
	}

	// push the chart
	err = pushChart(&args, packageRun)
	if err != nil {
		return err
	}

	// push the chart
	//client.Push()

	return nil
}

// verifyArgs verifies arguments
func verifyArgs(args *Args) error {
	if args.RegistryUsername == "" {
		return fmt.Errorf("No registry username provided: %w", errConfiguration)
	}

	if args.RegistryPassword == "" {
		return fmt.Errorf("No registry password provided: %w", errConfiguration)
	}

	if args.RegistryNamespace == "" {
		return fmt.Errorf("No registry namespace provided: %w", errConfiguration)
	}

	if args.ChartPath == "" {
		// default to workspace root
		args.ChartPath = "./"
	}

	if args.ChartDestination == "" {
		// default path to write packages
		args.ChartDestination = ".packaged_charts"
	}

	if args.RegistryURL == "" {
		// default to Docker Hub
		args.RegistryURL = "registry.hub.docker.com"
	}

	return nil
}

// packageChart packages a Helm chart
func packageChart(args *Args) (string, error) {
	helmClient := action.NewPackage()
	helmClient.DependencyUpdate = true
	helmClient.Destination = args.ChartDestination

	// minimal downloadManager settings which supports charts in the filesystem
	downloadManager := &downloader.Manager{
		Out:       io.Discard,
		ChartPath: args.ChartPath,
		Debug:     true,
	}
	if err := downloadManager.Build(); err != nil {
		return args.ChartPath, fmt.Errorf("Failed to retrieve chart in %s (%s)\n", args.ChartPath, err.Error())
	}

	packageRun, err := helmClient.Run(args.ChartPath, nil)
	if err != nil {
		return args.ChartPath, fmt.Errorf("Failed to package chart in %s (%s)\n", args.ChartPath, err.Error())
	}
	fmt.Printf("Successfully packaged chart in %s and saved it to %s\n", args.ChartPath, packageRun)

	return packageRun, nil
}

// pushChart pushes a Helm chart to a registry
func pushChart(args *Args, packageRun string) error {
	o := &registryLoginOptions{}

	opts := []registry.ClientOption{
		registry.ClientOptDebug(false),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(os.Stderr),
	}
	cfg := new(action.Configuration)
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return fmt.Errorf("Failed to create client (%s)\n", err)
	}
	cfg.RegistryClient = registryClient

	action.NewRegistryLogin(cfg).Run(os.Stdout, args.RegistryURL, args.RegistryUsername, args.RegistryPassword,
		action.WithCertFile(o.certFile),
		action.WithKeyFile(o.keyFile),
		action.WithCAFile(o.caFile),
		action.WithInsecure(o.insecure))

	client := action.NewPushWithOpts(action.WithPushConfig(cfg),
		action.WithTLSClientConfig(o.certFile, o.keyFile, o.caFile),
		action.WithInsecureSkipTLSVerify(false),
		action.WithPlainHTTP(false),
		action.WithPushOptWriter(nil))

	settings := new(cli.EnvSettings)
	client.Settings = settings
	ociURL := "oci://" + args.RegistryURL + "/" + args.RegistryNamespace

	// discard returned string since it appears to be empty
	_, err = client.Run(packageRun, ociURL)
	if err != nil {
		return fmt.Errorf("Failed to push chart %s to %s (%s)\n", packageRun, ociURL, err.Error())
	}
	fmt.Printf("Successfully pushed chart %s to %s\n", packageRun, ociURL)

	return nil
}
