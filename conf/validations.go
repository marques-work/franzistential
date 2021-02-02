package conf

import (
	"net"
	"strconv"
	"strings"

	"github.com/marques-work/franzistential/logging"
)

func validateNetworkListenerConfig(label, service string) {
	last := strings.LastIndex(service, ":")

	if 0 == last {
		logging.Die("-%sListen is missing the ip address", label)
	}

	if last < 0 || last == len(service) {
		logging.Die("-%sListen is missing the port number", label)
	}

	ip := service[:last]
	port := service[last+1:]

	if nil == net.ParseIP(ip) {
		logging.Die("-%sListen has an invalid IP address: `%s`", label, ip)
	}

	if !isOnlyNumbers(port) {
		logging.Die("-%sListen has an invalid port number: `%s`", label, port)
	}

	if portNum, err := strconv.ParseInt(port, 10, 64); err == nil {
		if portNum == 0 {
			logging.Die("-%sListen has an invalid port number: `%s`; not going to dynamically assign port when specifying port == 0", label, port)
		}

		if portNum > 65535 || portNum < 1 {
			logging.Die("-%sListen has an invalid port number: `%s`; ports must be between 1 and 65535, inclusively", label, port)
		}
	} else {
		logging.Die("-%sListen has an invalid port number: `%s`", label, port)
	}
}

func isOnlyNumbers(subj string) bool {
	return strings.IndexFunc(subj, nonDigit) == -1
}

func nonDigit(c rune) bool {
	return c < '0' || c > '9'
}
