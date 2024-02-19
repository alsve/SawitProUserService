package handler

import (
	"crypto/rsa"

	"github.com/alsve/SawitProUserService/repository"
)

type Server struct {
	Repository                   repository.RepositoryInterface
	defaultJWTExpirationInSecond int
	jwtRSAPrivateKey             *rsa.PrivateKey
}

type NewServerOptions struct {
	Repository                   repository.RepositoryInterface
	DefaultJWTExpirationInSecond int
	JWTRSAPrivateKey             *rsa.PrivateKey
}

// NewServer creates a new instance of Server.
func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository:                   opts.Repository,
		defaultJWTExpirationInSecond: opts.DefaultJWTExpirationInSecond,
		jwtRSAPrivateKey:             opts.JWTRSAPrivateKey,
	}
}
