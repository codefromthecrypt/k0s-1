package client

import (
	"k0s.io"
	"k0s.io/pkg/hub"
)

type Client interface {
	Config

	RunRedir() error
	RunSocks() error
	RunDoh() error
	Run() error
	MiniRun() error
	ListAgents() ([]hub.AgentInfo, error)
}

type Config interface {
	GetHost() string
	GetPort() string
	GetAddr() string
	GetScheme() string
	GetSchemeWS() string

	GetRedir() string
	GetSocks() string
	GetDoh() string
	GetCacheCredentials() bool
	GetCredentials() map[string]KeyStore
	GetConfigLocation() string

	GetVersion() k0s.Version
	GetInsecure() bool
	GetRecord() bool
}

type KeyStore map[string]string
