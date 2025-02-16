// +build !relayer

package integration_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestServeStargateWithWasm(t *testing.T) {
	var (
		env     = newEnv(t)
		apath   = env.Scaffold("sgblog")
		servers = env.RandomizeServerPorts(apath, "")
	)

	env.Must(env.Exec("add Wasm module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "wasm"),
			step.Workdir(apath),
		)),
	))

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "", "", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithCustomHome(t *testing.T) {
	var (
		env     = newEnv(t)
		apath   = env.Scaffold("sgblog2")
		servers = env.RandomizeServerPorts(apath, "")
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "./home", "", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithConfigHome(t *testing.T) {
	var (
		env     = newEnv(t)
		apath   = env.Scaffold("sgblog3")
		servers = env.RandomizeServerPorts(apath, "")
	)

	// Set config homes
	env.SetRandomHomeConfig(apath, "")

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "", "", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithCustomConfigFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "starporttest")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	var (
		env   = newEnv(t)
		apath = env.Scaffold("sgblog4")
	)
	// Move config
	newConfig := "new_config.yml"
	newConfigPath := filepath.Join(tmpDir, newConfig)
	err = os.Rename(filepath.Join(apath, "config.yml"), newConfigPath)
	require.NoError(t, err)

	servers := env.RandomizeServerPorts(tmpDir, newConfig)

	// Set config homes
	env.SetRandomHomeConfig(tmpDir, newConfig)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "", newConfigPath, ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
