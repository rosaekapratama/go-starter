package sshclient

import (
	"context"
	"golang.org/x/crypto/ssh"
	"sync"
	"time"
)

type IManager interface {
	initClient(ctx context.Context, clientId string, config ssh.Config) (client *ssh.Client, err error)
	GetClient(ctx context.Context, clientId string) (client *ssh.Client, err error)
}

type managerImpl struct {
	clientMap map[string]*SSHClient
}

type SSHClient struct {
	id          string
	address     string
	client      *ssh.Client
	config      *ssh.ClientConfig
	idleTimeout time.Duration
	idleTimer   *time.Timer
	mu          sync.Mutex
}
