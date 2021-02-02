package conf

import (
	"flag"

	"github.com/marques-work/franzistential/domain"
	"github.com/marques-work/franzistential/logging"
)

var (
	RawForward   *bool
	Out          *bool
	UnixSocket   *string
	TcpService   *string
	UdpService   *string
	Destinations []domain.Destination
	SendTimeout  *uint64
)

func OptParse() {
	RawForward = flag.Bool("cat", false, "Forward the input data as-is; do not try to parse into syslog-standard formats")
	Out = flag.Bool("stdout", false, "Pass through input to STDOUT")
	UnixSocket = flag.String("unixSocket", "", "Listen on local socket `path`; socket will be created on startup")
	TcpService = flag.String("tcpListen", "", "Listen on TCP/IP `ipaddr:portnum`, e.g. `127.0.0.1:514`; must include IP (v4 or v6) and port")
	UdpService = flag.String("udpListen", "", "Listen on UDP/IP `ipaddr:portnum`, e.g. `127.0.0.1:514`; must include IP (v4 or v6) and port")
	logging.Trace = flag.Bool("debug", false, "Verbose debug output")
	logging.Silent = flag.Bool("quiet", false, "Silence logging output; ``does NOT affect the `-stdout` flag")
	SendTimeout = flag.Uint64("sendTimeoutMs", uint64(20*1000), "Number of `milliseconds` before aborting message to destination (default: 20000)")
	flag.Var(&EventHubFlag{}, "eventHub", "Forward data to the provided Azure Event Hub `connection-uri`; multiple invocations of this flag will multiplex over each connection")

	flag.Parse()
}

func ConfigureAndValidate() {
	OptParse()

	if !*Out && 0 == len(Destinations) {
		logging.Warn("You have not set any destinations and are not printing to STDOUT; this will act like a sink, which might be a configuration error.")
	}

	if *logging.Trace && *logging.Silent {
		logging.Die("Both `-debug` and `-quiet` cannot be simultaneously set; pick one")
	}

	if "" != *TcpService {
		validateNetworkListenerConfig("tcp", *TcpService)
	}
}

func ServerMode() bool {
	return "" != *UnixSocket || "" != *TcpService || "" != *UdpService
}
