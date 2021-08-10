package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerApp(t *testing.T) {
	port := pickRandomUnusedPort()

	app, ctx, cancel := prepServerApp(t, func(rc RunCommand) RunCommand {
		rc.API.Port = port
		return rc
	})

	go func() { _ = app.run(ctx) }()
	waitForHTTPServerStart(port)

	// send ping
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", port))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "pong", string(body))

	cancel()
	app.Wait()
}

func TestServerApp_Shutdown(t *testing.T) {
	app, ctx, cancel := prepServerApp(t, func(rc RunCommand) RunCommand {
		rc.API.Port = pickRandomUnusedPort()
		return rc
	})
	time.AfterFunc(100*time.Millisecond, func() {
		cancel()
	})
	st := time.Now()
	err := app.run(ctx)
	assert.NoError(t, err)
	assert.True(t, time.Since(st).Seconds() < 1, "should take about 100msec")
	app.Wait()
}

func TestServerApp_WithSSL(t *testing.T) {
	cmd := RunCommand{}
	sslPort := pickRandomUnusedPort()
	logger := log.New(os.Stdout, "NODE-TEST ", 0)
	cmd.SetCommon(CommonOpts{Logger: logger, NodeURL: fmt.Sprintf("https://localhost:%d", sslPort)})

	// prepare options
	p := flags.NewParser(&cmd, flags.Default)
	port := pickRandomUnusedPort()
	_, err := p.ParseArgs([]string{"--api.port=" + strconv.Itoa(port),
		"--ssl.type=static", "--ssl.cert=fixtures/cert.pem", "--ssl.key=fixtures/key.pem",
		"--ssl.port=" + strconv.Itoa(sslPort)})
	require.NoError(t, err)

	// create app
	app, err := cmd.newServerApp()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = app.run(ctx) }()
	waitForHTTPSServerStart(sslPort)

	client := http.Client{
		// prevent http redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},

		// allow self-signed certificate
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// check http to https redirect response
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/ping", port))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 307, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("https://localhost:%d/ping", sslPort), resp.Header.Get("Location"))

	// check https server
	resp, err = client.Get(fmt.Sprintf("https://localhost:%d/ping", sslPort))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "pong", string(body))

	cancel()
	app.Wait()
}

func Test_ACMEEmail(t *testing.T) {
	cmd := RunCommand{}
	logger := log.New(os.Stdout, "NODE-TEST ", 0)
	cmd.SetCommon(CommonOpts{NodeURL: "https://obs.node.obada.io:443", Logger: logger})
	p := flags.NewParser(&cmd, flags.Default)
	args := []string{"--ssl.type=auto"}
	_, err := p.ParseArgs(args)
	require.NoError(t, err)
	conf, err := cmd.makeSSLConfig()
	require.NoError(t, err)
	assert.Equal(t, "techops@obada.io", conf.ACMEEmail)

	cmd = RunCommand{}
	cmd.SetCommon(CommonOpts{NodeURL: "https://obs.node.obada.io:443", Logger: logger})
	p = flags.NewParser(&cmd, flags.Default)
	args = []string{"--ssl.type=auto", "--ssl.acme-email=john.doe@gmail.com"}
	_, err = p.ParseArgs(args)
	require.NoError(t, err)
	conf, err = cmd.makeSSLConfig()
	require.NoError(t, err)
	assert.Equal(t, "john.doe@gmail.com", conf.ACMEEmail)
}

func pickRandomUnusedPort() (port int) {
	for i := 0; i < 10; i++ {
		port = 50000 + int(rand.Int31n(10000))
		if ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
			_ = ln.Close()
			break
		}
	}
	return port
}

func prepServerApp(t *testing.T, fn func(o RunCommand) RunCommand) (*serverApp, context.Context, context.CancelFunc) {
	cmd := RunCommand{}
	logger := log.New(os.Stdout, "NODE-TEST ", 0)
	cmd.SetCommon(CommonOpts{
		Logger: logger,
		NodeURL: fmt.Sprintf("https://localhost:%d", cmd.API.Port),
	})

	// prepare options
	p := flags.NewParser(&cmd, flags.Default)
	_, err := p.ParseArgs([]string{})

	require.NoError(t, err)

	cmd = fn(cmd)

	return createAppFromCmd(t, cmd)
}

func createAppFromCmd(t *testing.T, cmd RunCommand) (*serverApp, context.Context, context.CancelFunc) {
	app, err := cmd.newServerApp()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	rand.Seed(time.Now().UnixNano())
	return app, ctx, cancel
}

func waitForHTTPServerStart(port int) {
	// wait for up to 3 seconds for server to start before returning it
	client := http.Client{Timeout: time.Second}
	for i := 0; i < 300; i++ {
		time.Sleep(time.Millisecond * 10)
		if resp, err := client.Get(fmt.Sprintf("http://localhost:%d", port)); err == nil {
			_ = resp.Body.Close()
			return
		}
	}
}

func waitForHTTPSServerStart(port int) {
	// wait for up to 3 seconds for HTTPS server to start
	for i := 0; i < 300; i++ {
		time.Sleep(time.Millisecond * 10)
		conn, _ := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Millisecond*10)
		if conn != nil {
			_ = conn.Close()
			break
		}
	}
}
