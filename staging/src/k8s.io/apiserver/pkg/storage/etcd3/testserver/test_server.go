/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testserver

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
)

// getAvailablePort returns a TCP port that is available for binding.
func getAvailablePorts(count int) ([]int, error) {
	ports := []int{}
	for i := 0; i < count; i++ {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return nil, fmt.Errorf("could not bind to a port: %v", err)
		}
		// It is possible but unlikely that someone else will bind this port before we get a chance to use it.
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}

// NewTestConfig returns a configuration for an embedded etcd server.
// The configuration is based on embed.NewConfig(), with the following adjustments:
//   - sets UnsafeNoFsync = true to improve test performance (only reasonable in a test-only
//     single-member server we never intend to restart or keep data from)
//   - uses free ports for client and peer listeners
//   - cleans up the data directory on test termination
//   - silences server logs other than errors
func NewTestConfig(t *testing.T) *embed.Config {
	cfg := embed.NewConfig()

	cfg.UnsafeNoFsync = true

	ports, err := getAvailablePorts(2)
	if err != nil {
		t.Fatal(err)
	}
	clientURL := url.URL{Scheme: "http", Host: net.JoinHostPort("localhost", strconv.Itoa(ports[0]))}
	peerURL := url.URL{Scheme: "http", Host: net.JoinHostPort("localhost", strconv.Itoa(ports[1]))}

	cfg.LPUrls = []url.URL{peerURL}
	cfg.APUrls = []url.URL{peerURL}
	cfg.LCUrls = []url.URL{clientURL}
	cfg.ACUrls = []url.URL{clientURL}
	cfg.InitialCluster = cfg.InitialClusterFromName(cfg.Name)

	cfg.ZapLoggerBuilder = embed.NewZapLoggerBuilder(zaptest.NewLogger(t, zaptest.Level(zapcore.ErrorLevel)).Named("etcd-server"))
	cfg.Dir = t.TempDir()
	os.Chmod(cfg.Dir, 0700)
	return cfg
}

func startEtcd(t *testing.T, cfg *embed.Config) (e *embed.Etcd, cfgR *embed.Config, err error) {
	cfgR = NewTestConfig(t)
	if cfg != nil {
		cfgR.ExperimentalWatchProgressNotifyInterval = cfg.ExperimentalWatchProgressNotifyInterval
		cfgR.ClientTLSInfo = cfg.ClientTLSInfo
		if len(cfg.ClientTLSInfo.CertFile) > 0 && len(cfg.ClientTLSInfo.KeyFile) > 0 {
			for i := range cfgR.LCUrls {
				cfgR.LCUrls[i].Scheme = "https"
			}
			for i := range cfgR.ACUrls {
				cfgR.ACUrls[i].Scheme = "https"
			}
		}
	}

	e, err = embed.StartEtcd(cfgR)
	return e, cfgR, err
}

// RunEtcd starts an embedded etcd server with the provided config
// (or NewTestConfig(t) if nil), and returns a client connected to the server.
// The server is terminated when the test ends.
func RunEtcd(t *testing.T, cfg *embed.Config) *clientv3.Client {
	t.Helper()
	var e *embed.Etcd
	var err error
	step := 0
	for {
		e, cfg, err = startEtcd(t, cfg)
		if err != nil {
			if strings.Contains(err.Error(), "bind: address already in use") {
				time.Sleep(100 * time.Millisecond)
				if step >= 5 {
					break
				}
				step = step + 1
				continue
			}
			t.Fatal(err)
		}
		break
	}
	t.Cleanup(e.Close)

	select {
	case <-e.Server.ReadyNotify():
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		t.Fatal("server took too long to start")
	}
	go func() {
		err := <-e.Err()
		if err != nil {
			t.Error(err)
		}
	}()

	tlsConfig, err := cfg.ClientTLSInfo.ClientConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := clientv3.New(clientv3.Config{
		TLS:         tlsConfig,
		Endpoints:   e.Server.Cluster().ClientURLs(),
		DialTimeout: 10 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		Logger:      zaptest.NewLogger(t, zaptest.Level(zapcore.ErrorLevel)).Named("etcd-client"),
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}
