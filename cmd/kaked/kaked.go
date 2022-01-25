package kaked

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Joe-Degs/kake"
	"google.golang.org/grpc"
)

// Daemon represents a kake server daemon
type Daemon struct {
	// srv represents the server that the daemon will be running as a service
	// which is a grpc service
	srv *kake.Server

	// out is where all logs are spewed to.
	out io.Writer

	logging        *log.Logger
	shutting, done chan bool
	shutCloser     chan io.Closer
}

// return a default daemon for testing purposes
func DefaultDaemon() *Daemon {
	return &Daemon{
		srv:     kake.DefaultServer(),
		out:     os.Stdout,
		logging: setupLogging(os.Stdout),
	}
}

func (d *Daemon) handleSignals() {
	d.shutting, d.done = make(chan bool), make(chan bool)
	d.shutCloser = make(chan io.Closer)
	signalc := make(chan os.Signal, 1)

	signal.Notify(signalc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-signalc
		ss, ok := sig.(syscall.Signal)
		if !ok {
			d.logging.Fatal("signal is not a posix signal")
		}
		switch ss {
		case syscall.SIGHUP:
			// restart the process
			d.logging.Printf("recieved '%s': restarting process", sig)
		case syscall.SIGTERM, syscall.SIGINT:
			// shutdown process gracefully
			d.logging.Printf("recieved '%s': shutting down", sig)
			d.shutting <- true
			donec := make(chan bool)
			go func() {
				c := <-d.shutCloser
				if err := c.Close(); err != nil {
					exitf("error shuting down: %v", err)
				}
				donec <- true
			}()
			select {
			case <-donec:
				log.Println("shutdown")
				d.done <- true
				close(d.done)
				os.Exit(0)
			case <-time.After(3 * time.Second):
				exitf("exiting uncleanly: timeout")
			}
		default:
			log.Fatal("recieved unknown signal")
		}
	}
}

func exitf(format string, v ...interface{}) {
	log.SetOutput(os.Stderr)
	log.Fatalf(format, v...)
}

// register grpc service and start listening for requests
func (d *Daemon) Serve(lstn net.Listener) {
	gsrv := grpc.NewServer()

	d.srv.SelfRegisterService(gsrv)

	// start the grpc server
	go func() {
		if err := gsrv.Serve(lstn); err != nil {
			d.logging.Fatalf("error starting server: %v", err)
		}
	}()

	// wait for shutdown signal then proceed to terminate resources
	<-d.shutting
	gsrv.GracefulStop()
	d.shutCloser <- lstn
}

// init is a blocking function does these things
// 1. handle config parsing
// 2. start the server and graceful shutdown
// 3. wait for service to shutdown
func (d *Daemon) Init(cfg *kake.Config) error {
	if cfg != nil {
		d.srv.Config = cfg
	}
	d.logging.Println("daemon is starting...")

	d.logging.Printf("starting tcp listener on %s", d.srvAddr())
	lstn, err := net.ListenTCP("tcp", d.srvAddr())
	if err != nil {
		return err
	}

	go d.handleSignals()
	go d.Serve(lstn)

	// service is started wait till shutdown
	<-d.done
	return nil
}

func (d Daemon) srvAddr() *net.TCPAddr {
	return d.srv.Config.Addr
}

type logWriter struct {
	io.Writer
	timeFormat, prefix string
}

func (l *logWriter) Write(b []byte) (int, error) {
	t := time.Now().Format(l.timeFormat)
	return l.Writer.Write(append([]byte(t), b...))
}

func setupLogging(w io.Writer) *log.Logger {
	lw := &logWriter{w, "2006-01-02 15:04:05 ", "[kaked] "}
	return log.New(lw, lw.prefix, 0)
}
