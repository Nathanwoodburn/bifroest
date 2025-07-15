package environment

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/engity-com/bifroest/pkg/crypto"
	"github.com/engity-com/bifroest/pkg/errors"
	bnet "github.com/engity-com/bifroest/pkg/net"
	"github.com/engity-com/bifroest/pkg/session"
	"github.com/engity-com/bifroest/pkg/configuration"
)

type remote struct {
	repository *RemoteRepository
	session    session.Session
}

func (this *remote) SessionId() session.Id {
	return this.session.Id()
}

func (this *remote) PublicKey() crypto.PublicKey {
	return nil
}

func (this *remote) Banner(req Request) (io.ReadCloser, error) {
	b, err := this.repository.conf.Banner.Render(req)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(strings.NewReader(b)), nil
}

func (this *remote) Run(t Task) (int, error) {
	fail := func(err error) (int, error) {
		return -1, err
	}

	// Resolve connection details
	host, err := this.repository.conf.Host.Render(t)
	if err != nil {
		return fail(fmt.Errorf("cannot render host: %w", err))
	}

	port, err := this.repository.conf.Port.Render(t)
	if err != nil {
		return fail(fmt.Errorf("cannot render port: %w", err))
	}

	user, err := this.repository.conf.User.Render(t)
	if err != nil {
		return fail(fmt.Errorf("cannot render user: %w", err))
	}

	password, err := this.repository.conf.Password.Render(t)
	if err != nil {
		return fail(fmt.Errorf("cannot render password: %w", err))
	}

	// Create SSH client config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key verification
	}

	// Connect to remote server
	addr := net.JoinHostPort(host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fail(fmt.Errorf("cannot connect to remote server: %w", err))
	}
	defer client.Close()

	// Create session
	sshSession, err := client.NewSession()
	if err != nil {
		return fail(fmt.Errorf("cannot create SSH session: %w", err))
	}
	defer sshSession.Close()

	// Set up PTY if requested
	if ptyReq, _, isPty := t.SshSession().Pty(); isPty {
		if err := sshSession.RequestPty(ptyReq.Term, int(ptyReq.Window.Height), int(ptyReq.Window.Width), ssh.TerminalModes{}); err != nil {
			return fail(fmt.Errorf("cannot request PTY: %w", err))
		}
	}

	// Connect pipes
	sshSession.Stdin = t.SshSession()
	sshSession.Stdout = t.SshSession()
	sshSession.Stderr = t.SshSession().Stderr()

	// Execute command or shell
	var cmd string
	switch t.TaskType() {
	case TaskTypeShell:
		if rawCmd := t.SshSession().RawCommand(); len(rawCmd) > 0 {
			cmd = rawCmd
		}
	case TaskTypeSftp:
		cmd = "sftp-server"
	default:
		return fail(fmt.Errorf("unsupported task type: %v", t.TaskType()))
	}

	// Run the command
	var runErr error
	if cmd != "" {
		runErr = sshSession.Run(cmd)
	} else {
		runErr = sshSession.Shell()
		if runErr == nil {
			// Wait for the shell session to complete
			runErr = sshSession.Wait()
		}
	}

	if runErr != nil {
		if exitErr, ok := runErr.(*ssh.ExitError); ok {
			return exitErr.ExitStatus(), nil
		}
		return fail(fmt.Errorf("SSH session failed: %w", runErr))
	}

	return 0, nil
}

func (this *remote) IsPortForwardingAllowed(dest bnet.HostPort) (bool, error) {
	// Use a dummy request context for rendering
	req := &dummyRequest{session: this.session}
	allowed, err := this.repository.conf.PortForwardingAllowed.Render(req)
	if err != nil {
		return false, err
	}
	return allowed, nil
}

func (this *remote) NewDestinationConnection(ctx context.Context, dest bnet.HostPort) (io.ReadWriteCloser, error) {
	return nil, errors.Newf(errors.Permission, "port forwarding not implemented for remote environment")
}

func (this *remote) Dispose(context.Context) (bool, error) {
	return false, nil
}

func (this *remote) Close() error {
	return nil
}

// dummyRequest is a minimal implementation for template rendering
type dummyRequest struct {
	session session.Session
}

func (this *dummyRequest) Authorization() interface{} {
	return &dummyAuth{session: this.session}
}

func (this *dummyRequest) Connection() interface{} {
	return nil
}

func (this *dummyRequest) Context() context.Context {
	return context.Background()
}

func (this *dummyRequest) StartPreparation(string, string, PreparationProgressAttributes) (PreparationProgress, error) {
	return nil, nil
}

type dummyAuth struct {
	session session.Session
}

func (this *dummyAuth) Flow() configuration.FlowName {
	return this.session.Flow()
}

func (this *dummyAuth) FindSession() session.Session {
	return this.session
}

func (this *dummyAuth) EnvVars() map[string]string {
	return nil
}