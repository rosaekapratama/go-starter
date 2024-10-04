package sshclient

import (
	"context"
	"errors"
	"github.com/rosaekapratama/go-starter/log"
	"golang.org/x/crypto/ssh"
	"time"
)

var (
	errSSHClientNotFound = errors.New("SSH client not found")

	Manager IManager
)

func (m *managerImpl) initClient(ctx context.Context, clientId string, config ssh.Config) (client *ssh.Client, err error) {
	return nil, nil
}

func (m *managerImpl) GetClient(ctx context.Context, clientId string) (client *ssh.Client, err error) {
	if v, exists := m.clientMap[clientId]; exists {
		if v.client == nil {
			v.mu.Lock()
			defer v.mu.Unlock()

			err = v.connect(ctx)
			if err != nil {
				log.Errorf(ctx, err, "error on connect() clientId=%s", clientId)
				return
			}
		}
		client = v.client
		return
	}

	err = errSSHClientNotFound
	log.Errorf(ctx, err, "SSH client not found, clientId=%s", clientId)
	return
}

// connect establishes the SSH connection
func (c *SSHClient) connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error
	c.client, err = ssh.Dial("tcp", c.address, c.config)
	if err != nil {
		log.Errorf(ctx, err, "error on ssh.Dial(), clientId=%s, address=%s", c.id, c.address)
		return err
	}

	// Reset idle timer
	c.resetIdleTimer(ctx)
	return nil
}

// resetIdleTimer resets the idle timeout
func (c *SSHClient) resetIdleTimer(_ context.Context) {
	if c.idleTimer != nil {
		c.idleTimer.Stop()
	}
	c.idleTimer = time.AfterFunc(c.idleTimeout, c.close)
}

// close closes the SSH connection
func (c *SSHClient) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
}
