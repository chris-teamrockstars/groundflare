package socks

import "bufio"
import "context"
import "crypto/tls"
import "io"
import "log"
import "net"
import "groundflare/socks/authenticators"
import "groundflare/socks/interfaces"
import "groundflare/socks/loggers"
import "groundflare/socks/resolvers"
import "groundflare/socks/types"

type Server struct {

	AuthMethods []interfaces.Authenticator
	BufferPool  types.BufferPool
	Credentials interfaces.Credentials
	Logger      interfaces.Logger
	Resolver    interfaces.Resolver
	Ruleset     interfaces.Ruleset

}

func NewServer(options ...Option) *Server {

	server := &Server{
		AuthMethods: make([]interfaces.Authenticator, 0),
		BufferPool:  types.NewBufferPool(32 * 1024),
		Credentials: nil,
		Logger:      loggers.NewStandard(log.New(io.Discard, "socks5: ", log.LstdFlags)),
		Resolver:    resolvers.NewDNS(),
		Ruleset:     rulesets.NewPermitAll(),
	}

	for _, option := range options {
		option(server)
	}

	if len(server.AuthMethods) == 0 {

		if server.Credentials != nil {
			server.AuthMethods = append(server.AuthMethods, interfaces.Authenticator(&authenticators.UserPass{server.Credentials}))
		} else {
			server.AuthMethods = append(server.AuthMethods, interfaces.Authenticator(&authenticators.NoAuth{}))
		}

	}

	return server

}

func (server *Server) ListenAndServe(network string, address string) error {

	listener, err := net.Listen(network, address)

	if err == nil {
		return server.Serve(listener)
	} else {
		return err
	}

}

func (server *Server) ListenAndServeTLS(network string, address string, config *tls.Config) error {

	listener, err := tls.Listen(network, address, config)

	if err == nil {
		return server.Serve(listener)
	} else {
		return err
	}

}

func (server *Server) Serve(listener net.Listener) error {

	defer listener.Close()

	for {

		connection, err1 := listener.Accept()

		if err1 != nil {
			return err1
		}

		go func() {

			if err2 := server.ServeConnection(connection); err2 != nil {
				server.Logger.Errorf("server: %v", err2)
			}

		}()

	}

}

func (server *Server) ServeConnection(connection net.Conn) error {

	var auth_context *types.AuthContext = nil

	defer connection.Close()

	buffer_reader := bufio.NewReader(connection)



}
