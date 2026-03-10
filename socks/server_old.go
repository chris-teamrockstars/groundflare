package socks

import "bufio"
import "context"
import "crypto/tls"
import "errors"
import "fmt"
import "io"
import "log"
import "net"
import "groundflare/socks/bufferpool"
import "groundflare/socks/authenticators"
import socks_errors "groundflare/socks/errors"
import "groundflare/socks/interfaces"
import "groundflare/socks/loggers"
import "groundflare/socks/protocol"
import "groundflare/socks/statute"
import "groundflare/socks/types"

// GPool is used to implement custom goroutine pool default use goroutine
type GPool interface {
	Submit(f func()) error
}

// Server is responsible for accepting connections and handling
// the details of the SOCKS5 protocol
type Server struct {

	// XXX: Done
	authMethods []interfaces.Authenticator
	credentials interfaces.Credentials
	resolver NameResolver
	logger interfaces.Logger

	// rules is provided to enable custom logic around permitting
	// various commands. If not provided, NewPermitAll is used.
	rules RuleSet

	// rewriter can be used to transparently rewrite addresses.
	// This is invoked before the RuleSet is invoked.
	// Defaults to NoRewrite.
	rewriter AddressRewriter
	// bindIP is used for bind or udp associate
	bindIP net.IP
	// useBindIpResolveAsUdpAddr is used to resolve bindIP as udp address
	// default false, use  &net.UDPAddr{IP: request.LocalAddr.(*net.TCPAddr).IP, Port: 0}
	useBindIpBaseResolveAsUdpAddr bool
	// logger can be used to provide a custom log target.
	// Defaults to io.Discard.
	// Optional function for dialing out.
	// The callback set by dialWithRequest will be called first.
	dial func(ctx context.Context, network, addr string) (net.Conn, error)
	// Optional function for dialing out with the access of request detail.
	dialWithRequest func(ctx context.Context, network, addr string, request *Request) (net.Conn, error)
	// buffer pool
	bufferPool bufferpool.BufPool
	// XXX: OLD goroutine pool
	// XXX: gPool GPool
	// user's handle
	userConnectHandle   func(ctx context.Context, writer io.Writer, request *Request) error
	userBindHandle      func(ctx context.Context, writer io.Writer, request *Request) error
	userAssociateHandle func(ctx context.Context, writer io.Writer, request *Request) error
	// user's middleware
	userConnectMiddlewares   MiddlewareChain
	userBindMiddlewares      MiddlewareChain
	userAssociateMiddlewares MiddlewareChain
}

// ServeConn is used to serve a single connection.
func (sf *Server) ServeConn(conn net.Conn) error {

	var authContext *types.AuthContext

	defer conn.Close() // nolint: errcheck

	bufConn := bufio.NewReader(conn)

	mr, err := statute.ParseMethodRequest(bufConn)
	if err != nil {
		return err
	}
	if mr.Ver != protocol.Version5 {
		return statute.ErrNotSupportVersion
	}

	// Authenticate the connection
	userAddr := ""
	if conn.RemoteAddr() != nil {
		userAddr = conn.RemoteAddr().String()
	}
	authContext, err = sf.authenticate(conn, bufConn, userAddr, mr.Methods)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// The client request detail
	request, err := ParseRequest(bufConn)
	if err != nil {
		if errors.Is(err, statute.ErrUnrecognizedAddrType) {
			if err := SendReply(conn, statute.RepAddrTypeNotSupported, nil); err != nil {
				return fmt.Errorf("failed to send reply %w", err)
			}
		}
		return fmt.Errorf("failed to read destination address, %w", err)
	}

	if request.Request.Command != protocol.CommandConnect && // nolint: staticcheck
		request.Request.Command != protocol.CommandBind && // nolint: staticcheck
		request.Request.Command != protocol.CommandAssociate { // nolint: staticcheck
		if err := SendReply(conn, statute.RepCommandNotSupported, nil); err != nil {
			return fmt.Errorf("failed to send reply, %v", err)
		}
		return fmt.Errorf("unrecognized command[%d]", request.Request.Command) // nolint: staticcheck
	}

	request.AuthContext = authContext
	request.LocalAddr = conn.LocalAddr()
	request.RemoteAddr = conn.RemoteAddr()
	// Process the client request
	return sf.handleRequest(conn, request)
}

// authenticate is used to handle connection authentication
func (sf *Server) authenticate(conn io.Writer, bufConn io.Reader,
	userAddr string, methods []byte) (*types.AuthContext, error) {
	// Select a usable method
	for _, auth := range sf.authMethods {
		for _, method := range methods {
			if auth.GetCode() == method {
				return auth.Authenticate(bufConn, conn, userAddr)
			}
		}
	}
	// No usable method found
	conn.Write([]byte{protocol.Version5, protocol.MethodNoAcceptableMethods}) //nolint: errcheck
	return nil, socks_errors.AuthUnsupportedMethod
}

