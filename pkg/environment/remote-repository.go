package environment

import (
	"context"
	"fmt"

	glssh "github.com/gliderlabs/ssh"

	"github.com/engity-com/bifroest/pkg/alternatives"
	"github.com/engity-com/bifroest/pkg/configuration"
	"github.com/engity-com/bifroest/pkg/errors"
	"github.com/engity-com/bifroest/pkg/imp"
	"github.com/engity-com/bifroest/pkg/session"
)

var (
	_ = RegisterRepository(NewRemoteRepository)
)

func NewRemoteRepository(_ context.Context, flow configuration.FlowName, conf *configuration.EnvironmentRemote, _ alternatives.Provider, _ imp.Imp) (*RemoteRepository, error) {
	if conf == nil {
		return nil, fmt.Errorf("nil configuration")
	}

	return &RemoteRepository{
		flow: flow,
		conf: conf,
	}, nil
}

type RemoteRepository struct {
	flow configuration.FlowName
	conf *configuration.EnvironmentRemote
}

func (this *RemoteRepository) WillBeAccepted(ctx Context) (ok bool, err error) {
	fail := func(err error) (bool, error) {
		return false, err
	}

	if ok, err = this.conf.LoginAllowed.Render(ctx); err != nil {
		return fail(fmt.Errorf("cannot evaluate if user is allowed to login or not: %w", err))
	}

	return ok, nil
}

func (this *RemoteRepository) DoesSupportPty(Context, glssh.Pty) (bool, error) {
	return true, nil
}

func (this *RemoteRepository) Ensure(req Request) (Environment, error) {
	fail := func(err error) (Environment, error) {
		return nil, err
	}

	if ok, err := this.WillBeAccepted(req); err != nil {
		return fail(err)
	} else if !ok {
		return fail(ErrNotAcceptable)
	}

	sess := req.Authorization().FindSession()
	if sess == nil {
		return nil, errors.System.Newf("authorization without session")
	}

	return this.FindBySession(req.Context(), sess, nil)
}

func (this *RemoteRepository) FindBySession(_ context.Context, sess session.Session, _ *FindOpts) (Environment, error) {
	return &remote{this, sess}, nil
}

func (this *RemoteRepository) Close() error {
	return nil
}

func (this *RemoteRepository) Cleanup(context.Context, *CleanupOpts) error {
	return nil
}
