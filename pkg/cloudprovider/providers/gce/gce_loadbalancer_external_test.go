/*
Copyright 2017 The Kubernetes Authors.

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

package gce

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	computealpha "google.golang.org/api/compute/v0.alpha"
	compute "google.golang.org/api/compute/v1"

	"k8s.io/api/core/v1"
)

func TestEnsureStaticIP(t *testing.T) {
	fcas := NewFakeCloudAddressService()
	ipName := "some-static-ip"
	serviceName := ""
	region := "us-central1"

	// First ensure call
	ip, existed, err := ensureStaticIP(fcas, ipName, serviceName, region, "", NetworkTierDefault)
	if err != nil || existed || ip == "" {
		t.Fatalf(`ensureStaticIP(%v, %v, %v, %v, "") = %v, %v, %v; want valid ip, false, nil`, fcas, ipName, serviceName, region, ip, existed, err)
	}

	// Second ensure call
	var ipPrime string
	ipPrime, existed, err = ensureStaticIP(fcas, ipName, serviceName, region, ip, NetworkTierDefault)
	if err != nil || !existed || ip != ipPrime {
		t.Fatalf(`ensureStaticIP(%v, %v, %v, %v, %v) = %v, %v, %v; want %v, true, nil`, fcas, ipName, serviceName, region, ip, ipPrime, existed, err, ip)
	}
}

func TestEnsureStaticIPWithTier(t *testing.T) {
	s := NewFakeCloudAddressService()
	serviceName := ""
	region := "us-east1"

	for desc, tc := range map[string]struct {
		name     string
		netTier  NetworkTier
		expected string
	}{
		"Premium (default)": {
			name:     "foo-1",
			netTier:  NetworkTierPremium,
			expected: "PREMIUM",
		},
		"Standard": {
			name:     "foo-2",
			netTier:  NetworkTierStandard,
			expected: "STANDARD",
		},
	} {
		t.Run(desc, func(t *testing.T) {
			ip, existed, err := ensureStaticIP(s, tc.name, serviceName, region, "", tc.netTier)
			assert.NoError(t, err)
			assert.False(t, existed)
			assert.NotEqual(t, "", ip)
			// Get the Address from the fake address service and verify that the tier
			// is set correctly.
			alphaAddr, err := s.GetAlphaRegionAddress(tc.name, region)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, alphaAddr.NetworkTier)
		})
	}
}

func TestVerifyRequestedIP(t *testing.T) {
	region := "test-region"
	lbRef := "test-lb"
	s := NewFakeCloudAddressService()

	for desc, tc := range map[string]struct {
		requestedIP     string
		fwdRuleIP       string
		netTier         NetworkTier
		addrList        []*computealpha.Address
		expectErr       bool
		expectUserOwned bool
	}{
		"requested IP exists": {
			requestedIP:     "1.1.1.1",
			netTier:         NetworkTierPremium,
			addrList:        []*computealpha.Address{{Name: "foo", Address: "1.1.1.1", NetworkTier: "PREMIUM"}},
			expectErr:       false,
			expectUserOwned: true,
		},
		"requested IP is not static, but is in use by the fwd rule": {
			requestedIP: "1.1.1.1",
			fwdRuleIP:   "1.1.1.1",
			netTier:     NetworkTierPremium,
			expectErr:   false,
		},
		"requested IP is not static and is not used by the fwd rule": {
			requestedIP: "1.1.1.1",
			fwdRuleIP:   "2.2.2.2",
			netTier:     NetworkTierPremium,
			expectErr:   true,
		},
		"no requested IP": {
			netTier:   NetworkTierPremium,
			expectErr: false,
		},
		"requested IP exists, but network tier does not match": {
			requestedIP: "1.1.1.1",
			netTier:     NetworkTierStandard,
			addrList:    []*computealpha.Address{{Name: "foo", Address: "1.1.1.1", NetworkTier: "PREMIUM"}},
			expectErr:   true,
		},
	} {
		t.Run(desc, func(t *testing.T) {
			s.SetRegionalAddresses(region, tc.addrList)
			isUserOwnedIP, err := verifyUserRequestedIP(s, region, tc.requestedIP, tc.fwdRuleIP, lbRef, tc.netTier)
			assert.Equal(t, tc.expectErr, err != nil, fmt.Sprintf("err: %v", err))
			assert.Equal(t, tc.expectUserOwned, isUserOwnedIP)
		})
	}
}

func TestCreateForwardingRuleWithTier(t *testing.T) {
	s := NewFakeCloudForwardingRuleService()
	// Common variables among the tests.
	ports := []v1.ServicePort{{Name: "foo", Protocol: v1.ProtocolTCP, Port: int32(123)}}
	region := "test-region"
	target := "test-target-pool"
	svcName := "foo-svc"

	for desc, tc := range map[string]struct {
		netTier      NetworkTier
		expectedRule *computealpha.ForwardingRule
	}{
		"Premium tier": {
			netTier: NetworkTierPremium,
			expectedRule: &computealpha.ForwardingRule{
				Name:        "lb-1",
				Description: `{"kubernetes.io/service-name":"foo-svc"}`,
				IPAddress:   "1.1.1.1",
				IPProtocol:  "TCP",
				PortRange:   "123-123",
				Target:      target,
				NetworkTier: "PREMIUM",
			},
		},
		"Standard tier": {
			netTier: NetworkTierStandard,
			expectedRule: &computealpha.ForwardingRule{
				Name:        "lb-2",
				Description: `{"kubernetes.io/service-name":"foo-svc"}`,
				IPAddress:   "2.2.2.2",
				IPProtocol:  "TCP",
				PortRange:   "123-123",
				Target:      target,
				NetworkTier: "STANDARD",
			},
		},
	} {
		t.Run(desc, func(t *testing.T) {
			lbName := tc.expectedRule.Name
			ipAddr := tc.expectedRule.IPAddress

			err := createForwardingRule(s, lbName, svcName, region, ipAddr, target, ports, tc.netTier)
			assert.NoError(t, err)

			alphaRule, err := s.GetAlphaRegionForwardingRule(lbName, region)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRule, alphaRule)
		})
	}
}

func TestDeleteAddressWithWrongTier(t *testing.T) {
	region := "test-region"
	lbRef := "test-lb"
	s := NewFakeCloudAddressService()

	for desc, tc := range map[string]struct {
		addrName     string
		netTier      NetworkTier
		addrList     []*computealpha.Address
		expectDelete bool
	}{
		"Network tiers (premium) match; do nothing": {
			addrName: "foo1",
			netTier:  NetworkTierPremium,
			addrList: []*computealpha.Address{{Name: "foo1", Address: "1.1.1.1", NetworkTier: "PREMIUM"}},
		},
		"Network tiers (standard) match; do nothing": {
			addrName: "foo2",
			netTier:  NetworkTierStandard,
			addrList: []*computealpha.Address{{Name: "foo2", Address: "1.1.1.2", NetworkTier: "STANDARD"}},
		},
		"Wrong network tier (standard); delete address": {
			addrName:     "foo3",
			netTier:      NetworkTierPremium,
			addrList:     []*computealpha.Address{{Name: "foo3", Address: "1.1.1.3", NetworkTier: "STANDARD"}},
			expectDelete: true,
		},
		"Wrong network tier (preimium); delete address": {
			addrName:     "foo4",
			netTier:      NetworkTierStandard,
			addrList:     []*computealpha.Address{{Name: "foo4", Address: "1.1.1.4", NetworkTier: "PREMIUM"}},
			expectDelete: true,
		},
	} {
		t.Run(desc, func(t *testing.T) {
			s.SetRegionalAddresses(region, tc.addrList)
			// Sanity check to ensure we inject the right address.
			_, err := s.GetRegionAddress(tc.addrName, region)
			require.NoError(t, err)

			err = deleteAddressWithWrongTier(s, region, tc.addrName, lbRef, tc.netTier)
			assert.NoError(t, err)
			// Check whether the address still exists.
			_, err = s.GetRegionAddress(tc.addrName, region)
			if tc.expectDelete {
				assert.True(t, isNotFound(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMergeHttpHealthChecks(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		desc                   string
		checkIntervalSec       int64
		timeoutSec             int64
		healthyThreshold       int64
		unhealthyThreshold     int64
		wantCheckIntervalSec   int64
		wantTimeoutSec         int64
		wantHealthyThreshold   int64
		wantUnhealthyThreshold int64
	}{
		{"unchanged", gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold, gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold},
		{"interval - too small - should reconcile", gceHcCheckIntervalSeconds - 1, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold, gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold},
		{"timeout - too small - should reconcile", gceHcCheckIntervalSeconds, gceHcTimeoutSeconds - 1, gceHcHealthyThreshold, gceHcUnhealthyThreshold, gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold},
		{"healthy threshold - too small - should reconcile", gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold - 1, gceHcUnhealthyThreshold, gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold},
		{"unhealthy threshold - too small - should reconcile", gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold - 1, gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold},
		{"interval - user configured - should keep", gceHcCheckIntervalSeconds + 1, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold, gceHcCheckIntervalSeconds + 1, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold},
		{"timeout - user configured - should keep", gceHcCheckIntervalSeconds, gceHcTimeoutSeconds + 1, gceHcHealthyThreshold, gceHcUnhealthyThreshold, gceHcCheckIntervalSeconds, gceHcTimeoutSeconds + 1, gceHcHealthyThreshold, gceHcUnhealthyThreshold},
		{"healthy threshold - user configured - should keep", gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold + 1, gceHcUnhealthyThreshold, gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold + 1, gceHcUnhealthyThreshold},
		{"unhealthy threshold - user configured - should keep", gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold + 1, gceHcCheckIntervalSeconds, gceHcTimeoutSeconds, gceHcHealthyThreshold, gceHcUnhealthyThreshold + 1},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			wantHC := makeHttpHealthCheck("hc", "/", 12345)
			hc := &compute.HttpHealthCheck{
				CheckIntervalSec:   tc.checkIntervalSec,
				TimeoutSec:         tc.timeoutSec,
				HealthyThreshold:   tc.healthyThreshold,
				UnhealthyThreshold: tc.unhealthyThreshold,
			}
			mergeHttpHealthChecks(hc, wantHC)
			if wantHC.CheckIntervalSec != tc.wantCheckIntervalSec {
				t.Errorf("wantHC.CheckIntervalSec = %d; want %d", wantHC.CheckIntervalSec, tc.checkIntervalSec)
			}
			if wantHC.TimeoutSec != tc.wantTimeoutSec {
				t.Errorf("wantHC.TimeoutSec = %d; want %d", wantHC.TimeoutSec, tc.timeoutSec)
			}
			if wantHC.HealthyThreshold != tc.wantHealthyThreshold {
				t.Errorf("wantHC.HealthyThreshold = %d; want %d", wantHC.HealthyThreshold, tc.healthyThreshold)
			}
			if wantHC.UnhealthyThreshold != tc.wantUnhealthyThreshold {
				t.Errorf("wantHC.UnhealthyThreshold = %d; want %d", wantHC.UnhealthyThreshold, tc.unhealthyThreshold)
			}
		})
	}
}

func TestNeedToUpdateHttpHealthChecks(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		desc        string
		modifier    func(*compute.HttpHealthCheck)
		wantChanged bool
	}{
		{"unchanged", nil, false},
		{"desc does not match", func(hc *compute.HttpHealthCheck) { hc.Description = "bad-desc" }, true},
		{"port does not match", func(hc *compute.HttpHealthCheck) { hc.Port = 54321 }, true},
		{"requestPath does not match", func(hc *compute.HttpHealthCheck) { hc.RequestPath = "/anotherone" }, true},
		{"interval needs update", func(hc *compute.HttpHealthCheck) { hc.CheckIntervalSec = gceHcCheckIntervalSeconds - 1 }, true},
		{"timeout needs update", func(hc *compute.HttpHealthCheck) { hc.TimeoutSec = gceHcTimeoutSeconds - 1 }, true},
		{"healthy threshold needs update", func(hc *compute.HttpHealthCheck) { hc.HealthyThreshold = gceHcHealthyThreshold - 1 }, true},
		{"unhealthy threshold needs update", func(hc *compute.HttpHealthCheck) { hc.UnhealthyThreshold = gceHcUnhealthyThreshold - 1 }, true},
		{"interval does not need update", func(hc *compute.HttpHealthCheck) { hc.CheckIntervalSec = gceHcCheckIntervalSeconds + 1 }, false},
		{"timeout does not need update", func(hc *compute.HttpHealthCheck) { hc.TimeoutSec = gceHcTimeoutSeconds + 1 }, false},
		{"healthy threshold does not need update", func(hc *compute.HttpHealthCheck) { hc.HealthyThreshold = gceHcHealthyThreshold + 1 }, false},
		{"unhealthy threshold does not need update", func(hc *compute.HttpHealthCheck) { hc.UnhealthyThreshold = gceHcUnhealthyThreshold + 1 }, false},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			hc := makeHttpHealthCheck("hc", "/", 12345)
			wantHC := makeHttpHealthCheck("hc", "/", 12345)
			if tc.modifier != nil {
				tc.modifier(hc)
			}
			if gotChanged := needToUpdateHttpHealthChecks(hc, wantHC); gotChanged != tc.wantChanged {
				t.Errorf("needToUpdateHttpHealthChecks(%#v, %#v) = %t; want changed = %t", hc, wantHC, gotChanged, tc.wantChanged)
			}
		})
	}
}
