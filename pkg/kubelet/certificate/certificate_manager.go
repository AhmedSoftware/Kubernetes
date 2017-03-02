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

package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/util/cert"
	certificates "k8s.io/kubernetes/pkg/apis/certificates/v1beta1"
	certificatesclient "k8s.io/kubernetes/pkg/client/clientset_generated/clientset/typed/certificates/v1beta1"
)

const (
	syncPeriod = 1 * time.Hour
)

// Manager maintains and updates the certificates in use by this certificate
// manager. In the background it communicates with the API server to get new
// certificates for certificates about to expire.
type Manager interface {
	// Start the API server status sync loop.
	Start()
	// Current returns the currently selected certificate from the
	// certificate manager.
	Current() *tls.Certificate
}

// Config is the set of configuration parameters available for a new Manager.
type Config struct {
	// CertificateSigningRequestClient will be used for signing new certificate
	// requests generated when a key rotation occurs.
	CertificateSigningRequestClient certificatesclient.CertificateSigningRequestInterface
	// Template is the CertificateRequest that will be used as a template for
	// generating certificate signing requests for all new keys generated as
	// part of rotation. It follows the same rules as the template parameter of
	// crypto.x509.CreateCertificateRequest in the Go standard libraries.
	Template *x509.CertificateRequest
	// Usages is the types of usages that certificates generated by the manager
	// can be used for.
	Usages []certificates.KeyUsage
	// CertificateStore is a persistent store where the current cert/key is
	// kept and future cert/key pairs will be persisted after they are
	// generated.
	CertificateStore Store
	// CertificateRotationPercent is the percentage of time remaining of the
	// current certificate's lifetime before it will be rotated. If this is set
	// to 100 the certificate will be rotated each time it is checked. If this
	// is set to 0, it will never be rotated. If this value is over 100, an
	// error will be returned.
	CertificateRotationPercent uint
	// BootstrapCertificatePEM is the certificate data that will be returned
	// from the Manager if the CertificateStore doesn't have any cert/key pairs
	// currently available. If the CertificateStore does have a cert/key pair,
	// this will be ignored. If the bootstrap cert/key pair are used, they will
	// be rotated at the first opportunity, possibly well in advance of
	// expiring. This is intended to allow the first boot of a component to be
	// initialized using a generic, multi-use cert/key pair which will be
	// quickly replaced with a unique cert/key pair.
	BootstrapCertificatePEM []byte
	// BootstrapKeyPEM is the key data that will be returned from the Manager
	// if the CertificateStore doesn't have any cert/key pairs currently
	// available. If the CertificateStore does have a cert/key pair, this will
	// be ignored. If the bootstrap cert/key pair are used, they will be
	// rotated at the first opportunity, possibly well in advance of expiring.
	// This is intended to allow the first boot of a component to be
	// initialized using a generic, multi-use cert/key pair which will be
	// quickly replaced with a unique cert/key pair.
	BootstrapKeyPEM []byte
}

// Store is responsible for getting and updating the current certificate.
// Depending on the concrete implementation, the backing store for this
// behavior may vary.
type Store interface {
	// Current returns the currently selected certificate. If the Store doesn't
	// have a cert/key pair currently, it should return a NoCertKeyError so
	// that the Manager can recover by using bootstrap certificates to request
	// a new cert/key pair.
	Current() (*tls.Certificate, error)
	// Update accepts the PEM data for the cert/key pair and makes the new
	// cert/key pair the 'current' pair, that will be returned by future calls
	// to Current().
	Update(cert, key []byte) (*tls.Certificate, error)
}

// NoCertKeyError indicates there is no cert/key currently available.
type NoCertKeyError struct {
	msg string
}

func (e *NoCertKeyError) Error() string { return e.msg }

type manager struct {
	certSigningRequestClient certificatesclient.CertificateSigningRequestInterface
	template                 *x509.CertificateRequest
	usages                   []certificates.KeyUsage
	certStore                Store
	certAccessLock           sync.RWMutex
	cert                     *tls.Certificate
	shouldRotatePercent      uint
	forceRotation            bool
}

// NewManager returns a new certificate manager. A certificate manager is
// responsible for being the authoritative source of certificates in the
// Kubelet and handling updates due to rotation.
func NewManager(config *Config) (Manager, error) {
	cert, forceRotation, err := getCurrentCertificateOrBootstrap(
		config.CertificateStore,
		config.BootstrapCertificatePEM,
		config.BootstrapKeyPEM)
	if err != nil {
		return nil, err
	}

	if config.CertificateRotationPercent > 100 {
		return nil, fmt.Errorf("certificate rotation percent was over 100%%")
	}

	m := manager{
		certSigningRequestClient: config.CertificateSigningRequestClient,
		template:                 config.Template,
		usages:                   config.Usages,
		certStore:                config.CertificateStore,
		cert:                     cert,
		shouldRotatePercent:      config.CertificateRotationPercent,
		forceRotation:            forceRotation,
	}

	return &m, nil
}

// Current returns the currently selected certificate from the
// certificate manager.
func (m *manager) Current() *tls.Certificate {
	m.certAccessLock.RLock()
	defer m.certAccessLock.RUnlock()
	return m.cert
}

// Start will start the background work of rotating the certificates.
func (m *manager) Start() {
	if m.shouldRotatePercent < 1 {
		glog.V(2).Infof("Certificate rotation is not enabled.")
		return
	}

	// Certificate rotation depends on access to the API server certificate
	// signing API, so don't start the certificate manager if we don't have a
	// client. This will happen on the master, where the kubelet is responsible
	// for bootstrapping the pods of the master components.
	if m.certSigningRequestClient == nil {
		glog.V(2).Infof("Certificate rotation is not enabled, no connection to the apiserver.")
		return
	}

	glog.V(2).Infof("Certificate rotation is enabled.")
	go wait.Forever(func() {
		for range time.Tick(syncPeriod) {
			err := m.rotateCerts()
			if err != nil {
				glog.Errorf("Could not rotate certificates: %v", err)
			}
		}
	}, 0)
}

func getCurrentCertificateOrBootstrap(
	store Store,
	bootstrapCertificatePEM []byte,
	bootstrapKeyPEM []byte) (cert *tls.Certificate, rotated bool, errResult error) {

	currentCert, err := store.Current()
	if err == nil {
		return currentCert, false, nil
	}

	if _, ok := err.(*NoCertKeyError); !ok {
		return nil, false, err
	}

	if bootstrapCertificatePEM == nil || bootstrapKeyPEM == nil {
		return nil, false, fmt.Errorf("no cert/key available and no bootstrap cert/key to fall back to")
	}

	bootstrapCert, err := tls.X509KeyPair(bootstrapCertificatePEM, bootstrapKeyPEM)
	if err != nil {
		return nil, false, err
	}
	certs, err := x509.ParseCertificates(bootstrapCert.Certificate[0])
	if err != nil {
		return nil, false, fmt.Errorf("unable to parse certificate data: %v", err)
	}
	bootstrapCert.Leaf = certs[0]
	return &bootstrapCert, true, nil
}

// shouldRotate looks at how close the current certificate is to expiring and
// decides if it is time to rotate or not.
func (m *manager) shouldRotate() bool {
	m.certAccessLock.RLock()
	defer m.certAccessLock.RUnlock()
	notAfter := m.cert.Leaf.NotAfter
	total := notAfter.Sub(m.cert.Leaf.NotBefore)
	remaining := notAfter.Sub(time.Now())

	passedThreshold := remaining < 0 || uint(remaining*100/total) < m.shouldRotatePercent
	return m.forceRotation || passedThreshold
}

func (m *manager) rotateCerts() error {
	if !m.shouldRotate() {
		return nil
	}

	csrPEM, keyPEM, err := m.generateCSR()
	if err != nil {
		return err
	}

	// Call the Certificate Signing Request API to get a certificate for the
	// new private key.
	crtPEM, err := requestCertificate(m.certSigningRequestClient, csrPEM, m.usages)
	if err != nil {
		return fmt.Errorf("unable to get a new key signed: %v", err)
	}

	cert, err := m.certStore.Update(crtPEM, keyPEM)
	if err != nil {
		return fmt.Errorf("unable to store the new cert/key pair: %v", err)
	}

	m.updateCached(cert)
	m.forceRotation = false
	return nil
}

func (m *manager) updateCached(cert *tls.Certificate) {
	m.certAccessLock.Lock()
	defer m.certAccessLock.Unlock()
	m.cert = cert
}

func (m *manager) generateCSR() (csrPEM []byte, keyPEM []byte, err error) {
	// Generate a new private key.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate a new private key: %v", err)
	}
	der, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to marshal the new key to DER: %v", err)
	}

	keyPEM = pem.EncodeToMemory(&pem.Block{Type: cert.ECPrivateKeyBlockType, Bytes: der})

	csrPEM, err = cert.MakeCSRFromTemplate(privateKey, m.template)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create a csr from the private key: %v", err)
	}
	return csrPEM, keyPEM, nil
}

// requestCertificate will create a certificate signing request using the PEM
// encoded CSR and send it to API server, then it will watch the object's
// status, once approved by API server, it will return the API server's issued
// certificate (pem-encoded). If there is any errors, or the watch timeouts, it
// will return an error.
//
// NOTE This is a copy of a function with the same name in
// k8s.io/kubernetes/pkg/kubelet/util/csr/csr.go, changing only the package that
// CertificateSigningRequestInterface and KeyUsage are imported from.
func requestCertificate(client certificatesclient.CertificateSigningRequestInterface, csrData []byte, usages []certificates.KeyUsage) (certData []byte, err error) {
	req, err := client.Create(&certificates.CertificateSigningRequest{
		// Username, UID, Groups will be injected by API server.
		TypeMeta:   metav1.TypeMeta{Kind: "CertificateSigningRequest"},
		ObjectMeta: metav1.ObjectMeta{GenerateName: "csr-"},

		Spec: certificates.CertificateSigningRequestSpec{
			Request: csrData,
			Usages:  usages,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create certificate signing request: %v", err)
	}

	// Make a default timeout = 3600s.
	var defaultTimeoutSeconds int64 = 3600
	certWatch, err := client.Watch(metav1.ListOptions{
		Watch:          true,
		TimeoutSeconds: &defaultTimeoutSeconds,
		FieldSelector:  fields.OneTermEqualSelector("metadata.name", req.Name).String(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot watch on the certificate signing request: %v", err)
	}
	defer certWatch.Stop()
	ch := certWatch.ResultChan()

	for {
		event, ok := <-ch
		if !ok {
			break
		}

		if event.Type == watch.Modified || event.Type == watch.Added {
			if event.Object.(*certificates.CertificateSigningRequest).UID != req.UID {
				continue
			}
			status := event.Object.(*certificates.CertificateSigningRequest).Status
			for _, c := range status.Conditions {
				if c.Type == certificates.CertificateDenied {
					return nil, fmt.Errorf("certificate signing request is not approved, reason: %v, message: %v", c.Reason, c.Message)
				}
				if c.Type == certificates.CertificateApproved && status.Certificate != nil {
					return status.Certificate, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("watch channel closed")
}
