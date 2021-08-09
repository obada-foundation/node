package cmd

import (
	"context"
	"github.com/obada-foundation/node/rest"
	"os"
	"os/signal"
	"syscall"
)

type RunCommand struct {
	API RestAPIGroup `group:"api" namespace:"api" env-namespace:"API"`
	SSL SSLGroup `group:"ssl" namespace:"ssl" env-namespace:"SSL"`

	CommonOpts
}

// RestAPIGroup defines options group for REST API params
type RestAPIGroup struct {
	Port    int    `long:"port" default:"3000" description:"port"`
	Address string `long:"address" default:"127.0.0.1" description:"listening address"`
}

// SSLGroup defines options group for server ssl params
type SSLGroup struct {
	Type         string `long:"type" env:"TYPE" description:"ssl (auto) support" choice:"none" choice:"static" choice:"auto" default:"none"` //nolint
	Port         int    `long:"port" env:"PORT" description:"port number for https server" default:"3443"`
	Cert         string `long:"cert" env:"CERT" description:"path to cert.pem file"`
	Key          string `long:"key" env:"KEY" description:"path to key.pem file"`
	ACMELocation string `long:"acme-location" env:"ACME_LOCATION" description:"dir where certificates will be stored by autocert manager" default:"./var/acme"`
	ACMEEmail    string `long:"acme-email" env:"ACME_EMAIL" description:"admin email for certificate notifications"`
}

type serverApp struct {
	*RunCommand

	restSrv    *rest.Rest
	terminated chan struct{}
}

// Execute is the entry point for "run" command, called by flag parser
func (rc *RunCommand) Execute(_ []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		rc.Logger.Printf("interrupt signal")
		cancel()
	}()

	app, err := rc.newServerApp()
	if err != nil {
		return err
	}

	if err = app.run(ctx); err != nil {
		return err
	}

	return nil
}

func (rc *RunCommand) newServerApp() (*serverApp, error) {
	srv := &rest.Rest{
		Version: rc.Version,
		Logger: rc.Logger,
	}

	return &serverApp{
		RunCommand: rc,

		restSrv:    srv,
		terminated: make(chan struct{}),
	}, nil
}

func (a *serverApp) run(ctx context.Context) error {
	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		a.Logger.Print("server :: shutdown initiated")
		a.restSrv.Shutdown()
	}()

	a.restSrv.Run(a.API.Address, a.API.Port)

	close(a.terminated)

	return nil
}

// Wait for application completion (termination)
func (a *serverApp) Wait() {
	<-a.terminated
}

