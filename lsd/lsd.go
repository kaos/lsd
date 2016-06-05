package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"strings"
)

var backends = map[string]NewBackendFactory{}
var lsd_info = Rows{}
var opts struct {
	//flags     *flag.FlagSet
	host      *string
	port      *int
	file      *string
	daemonize *bool
	debug     *bool
	backend   *string
}

func init() {
	//flag := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	//opts.flags = flag
	opts.host = flag.String("addr", "127.0.0.1", "listen address")
	opts.port = flag.Int("port", 0, "listen port")
	opts.file = flag.String("file", "", "listen on unix socket file")
	opts.daemonize = flag.Bool("d", false, "daemonize process")
	opts.debug = flag.Bool("debug", false, "enable debug logging")
	opts.backend = flag.String("backend", "", "backend system")

	addInfo(
		"name", "livestatusd",
		"version", "0.1",
		"url", "https://github.com/.../livestatusd",
	)
}

func usage(list_backends bool, msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	if list_backends {
		i, keys := 0, make([]string, len(backends))
		for k := range backends {
			keys[i] = k
			i++
		}
		fmt.Fprintf(os.Stderr, "available backends: %s\n", strings.Join(keys, ", "))
	}
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()
	/*err := opts.flags.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		return
	}
	fmt.Printf("parse: %#v\nflags: %#v\n", err, opts.flags)*/

	if len(*opts.backend) == 0 {
		usage(false, "missing required flag: -backend")
	}

	backendFactory := backends[*opts.backend]
	if backendFactory == nil {
		usage(true, "unknown backend: %s", *opts.backend)
	}

	if *opts.daemonize {
		fmt.Println("daemonize NYI... exiting")
		os.Exit(127)
	}

	if *opts.debug {
		log.SetLevel(log.DebugLevel)
	}

	start(backendFactory())
}

func start(backend Backend) {
	server := make(chan interface{})
	servers := 0

	if len(*opts.file) > 0 {
		servers++
		go MainLoop("unix", *opts.file, server)
	}

	if *opts.port > 0 {
		servers++
		go MainLoop(
			"tcp",
			net.JoinHostPort(
				*opts.host,
				strconv.Itoa(*opts.port)),
			server,
		)
	}

	for servers > 0 {
		switch cmd := (<-server).(type) {
		case Client:
			go ServeClient(cmd, backend)
		case nil:
			servers--
		}
	}

	log.Info("no listening sockets, exiting")
}

func MainLoop(netw, laddr string, cmd chan interface{}) {
	defer func() {
		log.Info("closing")
		cmd <- nil
	}()

	server, err := net.Listen(netw, laddr)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.WithFields(
		log.Fields{
			"server": server.Addr(),
		})

	logger.Info("livestatusd: listening for connections")
	for {
		conn, err := server.Accept()
		if err != nil {
			logger.Fatal(err)
		}

		cmd <- NewSocketClient(conn, logger)
	}
}

func RegisterBackend(name string, factory NewBackendFactory) {
	backends[name] = factory
}

func lsdInfo() Rows {
	return lsd_info
}

func addInfo(info ...string) {
	lsd_info = append(lsd_info, parseInfo(info)...)
}

func parseInfo(info []string) (out Rows) {
	if len(info)%2 > 0 {
		info = append(info, "")
	}
	obj := Object{}
	for i, v := range info {
		switch i % 2 {
		case 0:
			obj["key"] = v
		case 1:
			obj["val"] = v
			out = append(out, obj)
			obj = Object{}
		}
	}

	return out
}
