/*
Copyright 2023 The Kubernetes Authors.

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

package portforward

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/util/httpstream/wsstream"

	"k8s.io/klog/v2"

	"golang.org/x/net/websocket"

	"k8s.io/apimachinery/pkg/util/runtime"
	netutils "k8s.io/utils/net"
)

// PortForwardProtocolV1Name is the subprotocol used for port forwarding.
// TODO move to API machinery and re-unify with kubelet/server/portfoward
const (
	PortForwardProtocolV1Name        = "portforward.k8s.io"
	PortForwardWebsocketProtocolName = "v4." + wsstream.ChannelWebSocketProtocol
)

var ErrLostConnectionToPod = errors.New("lost connection to pod")

// PortForwarder knows how to listen for local connections and forward them to
// a remote pod via an upgraded HTTP request.
type PortForwarder struct {
	addresses []listenAddress
	ports     []ForwardedPort
	stopChan  <-chan struct{}

	dialer        DialFunc
	closed        bool
	connections   []*websocket.Conn
	listeners     []io.Closer
	Ready         chan struct{}
	requestIDLock sync.Mutex
	requestID     int
	out           io.Writer
	errOut        io.Writer
}

// DialFunc will establish a websocket connection using the requested protocols for the provided port or return
// an error.
type DialFunc func(id int, port uint16, subprotocols ...string) (*websocket.Conn, error)

// ForwardedPort contains a Local:Remote port pairing.
type ForwardedPort struct {
	Local  uint16
	Remote uint16
}

/*
valid port specifications:

5000
- forwards from localhost:5000 to pod:5000

8888:5000
- forwards from localhost:8888 to pod:5000

0:5000
:5000
  - selects a random available local port,
    forwards from localhost:<random port> to pod:5000
*/
func parsePorts(ports []string) ([]ForwardedPort, error) {
	var forwards []ForwardedPort
	for _, portString := range ports {
		parts := strings.Split(portString, ":")
		var localString, remoteString string
		if len(parts) == 1 {
			localString = parts[0]
			remoteString = parts[0]
		} else if len(parts) == 2 {
			localString = parts[0]
			if localString == "" {
				// support :5000
				localString = "0"
			}
			remoteString = parts[1]
		} else {
			return nil, fmt.Errorf("invalid port format '%s'", portString)
		}

		localPort, err := strconv.ParseUint(localString, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("error parsing local port '%s': %s", localString, err)
		}

		remotePort, err := strconv.ParseUint(remoteString, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("error parsing remote port '%s': %s", remoteString, err)
		}
		if remotePort == 0 {
			return nil, fmt.Errorf("remote port must be > 0")
		}

		forwards = append(forwards, ForwardedPort{uint16(localPort), uint16(remotePort)})
	}

	return forwards, nil
}

type listenAddress struct {
	address     string
	protocol    string
	failureMode string
}

func parseAddresses(addressesToParse []string) ([]listenAddress, error) {
	var addresses []listenAddress
	parsed := make(map[string]listenAddress)
	for _, address := range addressesToParse {
		if address == "localhost" {
			if _, exists := parsed["127.0.0.1"]; !exists {
				ip := listenAddress{address: "127.0.0.1", protocol: "tcp4", failureMode: "all"}
				parsed[ip.address] = ip
			}
			if _, exists := parsed["::1"]; !exists {
				ip := listenAddress{address: "::1", protocol: "tcp6", failureMode: "all"}
				parsed[ip.address] = ip
			}
		} else if netutils.ParseIPSloppy(address).To4() != nil {
			parsed[address] = listenAddress{address: address, protocol: "tcp4", failureMode: "any"}
		} else if netutils.ParseIPSloppy(address) != nil {
			parsed[address] = listenAddress{address: address, protocol: "tcp6", failureMode: "any"}
		} else {
			return nil, fmt.Errorf("%s is not a valid IP", address)
		}
	}
	addresses = make([]listenAddress, len(parsed))
	id := 0
	for _, v := range parsed {
		addresses[id] = v
		id++
	}
	// Sort addresses before returning to get a stable order
	sort.Slice(addresses, func(i, j int) bool { return addresses[i].address < addresses[j].address })

	return addresses, nil
}

// New creates a new PortForwarder with localhost listen addresses.
func New(dialer DialFunc, ports []string, stopChan <-chan struct{}, readyChan chan struct{}, out, errOut io.Writer) (*PortForwarder, error) {
	return NewOnAddresses(dialer, []string{"localhost"}, ports, stopChan, readyChan, out, errOut)
}

// NewOnAddresses creates a new PortForwarder with custom listen addresses.
func NewOnAddresses(dialer DialFunc, addresses []string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}, out, errOut io.Writer) (*PortForwarder, error) {
	if len(addresses) == 0 {
		return nil, errors.New("you must specify at least 1 address")
	}
	parsedAddresses, err := parseAddresses(addresses)
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, errors.New("you must specify at least 1 port")
	}
	parsedPorts, err := parsePorts(ports)
	if err != nil {
		return nil, err
	}
	return &PortForwarder{
		dialer:    dialer,
		addresses: parsedAddresses,
		ports:     parsedPorts,
		stopChan:  stopChan,
		Ready:     readyChan,
		out:       out,
		errOut:    errOut,
	}, nil
}

// ForwardPorts formats and executes a port forwarding request. The connection will remain
// open until stopChan is closed.
func (pf *PortForwarder) ForwardPorts() error {
	defer pf.Close()

	return pf.forward()
}

// forward dials the remote host specific in req, upgrades the request, starts
// listeners for each port specified in ports, and forwards local connections
// to the remote host via streams.
func (pf *PortForwarder) forward() error {
	var err error

	listenSuccess := false
	for i := range pf.ports {
		port := &pf.ports[i]
		err = pf.listenOnPort(port)
		switch {
		case err == nil:
			listenSuccess = true
		default:
			if pf.errOut != nil {
				fmt.Fprintf(pf.errOut, "Unable to listen on port %d: %v\n", port.Local, err)
			}
		}
	}

	if !listenSuccess {
		return fmt.Errorf("unable to listen on any of the requested ports: %v", pf.ports)
	}

	if pf.Ready != nil {
		close(pf.Ready)
	}

	// wait for interrupt or conn closure
	select {
	case <-pf.stopChan:
		// TODO listen server chan via ping/pong and close connection when it fails
		// and return ErrLostConnectionToPod
		return nil
	}
}

// listenOnPort delegates listener creation and waits for connections on requested bind addresses.
// An error is raised based on address groups (default and localhost) and their failure modes
func (pf *PortForwarder) listenOnPort(port *ForwardedPort) error {
	var errors []error
	failCounters := make(map[string]int, 2)
	successCounters := make(map[string]int, 2)
	for _, addr := range pf.addresses {
		err := pf.listenOnPortAndAddress(port, addr.protocol, addr.address)
		if err != nil {
			errors = append(errors, err)
			failCounters[addr.failureMode]++
		} else {
			successCounters[addr.failureMode]++
		}
	}
	if successCounters["all"] == 0 && failCounters["all"] > 0 {
		return fmt.Errorf("%s: %v", "Listeners failed to create with the following errors", errors)
	}
	if failCounters["any"] > 0 {
		return fmt.Errorf("%s: %v", "Listeners failed to create with the following errors", errors)
	}
	return nil
}

// listenOnPortAndAddress delegates listener creation and waits for new connections
// in the background f
func (pf *PortForwarder) listenOnPortAndAddress(port *ForwardedPort, protocol string, address string) error {
	listener, err := pf.getListener(protocol, address, port)
	if err != nil {
		return err
	}
	pf.listeners = append(pf.listeners, listener)
	go pf.waitForConnection(listener, *port)
	return nil
}

// getListener creates a listener on the interface targeted by the given hostname on the given port with
// the given protocol. protocol is in net.Listen style which basically admits values like tcp, tcp4, tcp6
func (pf *PortForwarder) getListener(protocol string, hostname string, port *ForwardedPort) (net.Listener, error) {
	listener, err := net.Listen(protocol, net.JoinHostPort(hostname, strconv.Itoa(int(port.Local))))
	if err != nil {
		return nil, fmt.Errorf("unable to create listener: Error %s", err)
	}
	listenerAddress := listener.Addr().String()
	host, localPort, _ := net.SplitHostPort(listenerAddress)
	localPortUInt, err := strconv.ParseUint(localPort, 10, 16)

	if err != nil {
		fmt.Fprintf(pf.out, "Failed to forward from %s:%d -> %d\n", hostname, localPortUInt, port.Remote)
		return nil, fmt.Errorf("error parsing local port: %s from %s (%s)", err, listenerAddress, host)
	}
	port.Local = uint16(localPortUInt)
	if pf.out != nil {
		fmt.Fprintf(pf.out, "Forwarding from %s -> %d\n", net.JoinHostPort(hostname, strconv.Itoa(int(localPortUInt))), port.Remote)
	}

	return listener, nil
}

// waitForConnection waits for new connections to listener and handles them in
// the background.
func (pf *PortForwarder) waitForConnection(listener net.Listener, port ForwardedPort) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			// TODO consider using something like https://github.com/hydrogen18/stoppableListener?
			if !strings.Contains(strings.ToLower(err.Error()), "use of closed network connection") {
				runtime.HandleError(fmt.Errorf("error accepting connection on port %d: %v", port.Local, err))
			}
			return
		}

		go func() {
			defer conn.Close()
			var requestID int
			if started, id := pf.startConnection(); !started {
				runtime.HandleError(fmt.Errorf("too many connections, can't open port %d -> %d", port.Local, port.Remote))
				return
			} else {
				requestID = id
			}
			klog.Infof("Starting connection %v", port)
			if err := pf.handleConnection(conn, port, requestID); err != nil {
				runtime.HandleError(err)
			} else {
				klog.Infof("Finished connection %v", port)
			}

		}()
	}
}

// handleConnection copies data between the local connection and the stream to
// the remote server.
func (pf *PortForwarder) handleConnection(conn net.Conn, port ForwardedPort, id int) error {
	if pf.out != nil {
		fmt.Fprintf(pf.out, "Handling connection for %d\n", port.Local)
	}

	klog.Infof("Dialing remote port %d with protocol %s", port.Remote, PortForwardWebsocketProtocolName)

	ws, err := pf.dialer(id, port.Remote, PortForwardWebsocketProtocolName)
	if err != nil {
		return fmt.Errorf("error upgrading connection: %s", err)
	}

	pf.addConnection(ws)

	defer func() {
		klog.Infof("About to close connection for port %d -> %d", port.Local, port.Remote)
		ws.Close()
		pf.removeConnection(ws)
		klog.Infof("Closed connection for %d -> %d", port.Local, port.Remote)
	}()

	wsConn := wsstream.NewConn(map[string]wsstream.ChannelProtocolConfig{
		PortForwardWebsocketProtocolName: {Binary: true, Channels: []wsstream.ChannelType{wsstream.ReadWriteChannel, wsstream.ReadChannel}},
	})

	actualProtocol, streams, err := wsConn.OpenChannels(ws)
	if err != nil {
		return fmt.Errorf("error setting up channels for %d -> %d: %v", port.Local, port.Remote, err)
	}
	if actualProtocol != PortForwardWebsocketProtocolName {
		return fmt.Errorf("server did not use our supported protocol %s: %s", PortForwardWebsocketProtocolName, actualProtocol)
	}

	klog.Infof("Got %d streams for communication", len(streams))

	errorChan := make(chan error)
	go func() {
		defer close(errorChan)
		if err := checkPort(streams[1], port.Remote); err != nil {
			errorChan <- fmt.Errorf("error establishing error channel: %v", err)
			return
		}
		klog.Infof("Got expected prefix from stream 1")

		message, err := io.ReadAll(streams[1])
		switch {
		case err != nil:
			errorChan <- fmt.Errorf("error reading from error stream for port %d -> %d: %v", port.Local, port.Remote, err)
		case len(message) > 0:
			errorChan <- fmt.Errorf("an error occurred forwarding %d -> %d: %v", port.Local, port.Remote, string(message))
		}
		klog.Infof("Reading error channel port %d -> %d finished", port.Local, port.Remote)
	}()

	localError := make(chan struct{})
	remoteDone := make(chan struct{})

	go func() {
		defer close(remoteDone)

		if err := checkPort(streams[0], port.Remote); err != nil {
			runtime.HandleError(fmt.Errorf("error establishing stream: %v", err))
			return
		}
		klog.Infof("Got expected prefix from stream 0")
		// Copy from the remote side to the local port.
		if _, err := io.Copy(conn, streams[0]); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			runtime.HandleError(fmt.Errorf("error copying from remote stream to local connection: %v", err))
		}

		// inform the select below that the remote copy is done
		klog.Infof("Reading remote port %d -> %d finished", port.Local, port.Remote)
	}()

	go func() {
		defer streams[0].Close()

		// Copy from the local port to the remote side.
		if _, err := io.Copy(streams[0], conn); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			runtime.HandleError(fmt.Errorf("error copying from local connection to remote stream: %v", err))
			// break out of the select below without waiting for the other copy to finish
			close(localError)
		}
		klog.Infof("Writing local port %d -> %d finished", port.Local, port.Remote)
	}()

	// wait for either a local->remote error or for copying from remote->local to finish
	select {
	case <-remoteDone:
	case <-localError:
	}

	wsConn.Close()

	// always expect something on errorChan (it may be nil)
	err = <-errorChan
	if err != nil {
		runtime.HandleError(err)
	}

	return nil
}

// Close stops all listeners of PortForwarder.
func (pf *PortForwarder) Close() {
	pf.requestIDLock.Lock()
	defer pf.requestIDLock.Unlock()

	pf.closed = true

	// stop all listeners
	for _, l := range pf.listeners {
		if err := l.Close(); err != nil {
			runtime.HandleError(fmt.Errorf("error closing listener: %v", err))
		}
	}

	pf.listeners = nil

	for _, c := range pf.connections {
		if err := c.Close(); err != nil {
			runtime.HandleError(fmt.Errorf("error closing connection: %v", err))
		}
	}

	pf.connections = nil
}

// GetPorts will return the ports that were forwarded; this can be used to
// retrieve the locally-bound port in cases where the input was port 0. This
// function will signal an error if the Ready channel is nil or if the
// listeners are not ready yet; this function will succeed after the Ready
// channel has been closed.
func (pf *PortForwarder) GetPorts() ([]ForwardedPort, error) {
	if pf.Ready == nil {
		return nil, fmt.Errorf("no Ready channel provided")
	}
	select {
	case <-pf.Ready:
		return pf.ports, nil
	default:
		return nil, fmt.Errorf("listeners not ready")
	}
}

func (pf *PortForwarder) startConnection() (bool, int) {
	pf.requestIDLock.Lock()
	defer pf.requestIDLock.Unlock()
	if pf.closed {
		return false, 0
	}

	id := pf.requestID
	pf.requestID++
	return true, id
}

func (pf *PortForwarder) addConnection(conn *websocket.Conn) {
	pf.requestIDLock.Lock()
	defer pf.requestIDLock.Unlock()
	if pf.closed {
		conn.Close()
		return
	}
	pf.connections = append(pf.connections, conn)
}

func (pf *PortForwarder) removeConnection(conn *websocket.Conn) {
	pf.requestIDLock.Lock()
	defer pf.requestIDLock.Unlock()
	for i, c := range pf.connections {
		if c == conn {
			pf.connections = append(pf.connections[:i], pf.connections[i+1:]...)
			break
		}
	}
}

// checkPort verifies that the first two bytes from the stream have the expected header.
func checkPort(r io.Reader, expected uint16) error {
	var data [2]byte
	if _, err := io.ReadAtLeast(r, data[:], 2); err != nil {
		return err
	}
	if actual := binary.LittleEndian.Uint16(data[:]); actual != expected {
		return fmt.Errorf("expected to receive port %d as header, but got %d", expected, actual)
	}
	return nil
}
