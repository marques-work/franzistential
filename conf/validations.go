package conf

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func validateNetworkListenerConfig(label, service string) error {
	last := strings.LastIndex(service, ":")

	if 0 == last {
		return fmt.Errorf("`--in:%s` is missing the ip address", label)
	}

	if last < 0 || last == len(service) {
		return fmt.Errorf("`--in:%s` is missing the port number", label)
	}

	ip := service[:last]
	port := service[last+1:]

	if nil == net.ParseIP(ip) {
		return fmt.Errorf("`--in:%s` has an invalid IP address: `%s`", label, ip)
	}

	if !isOnlyNumbers(port) {
		return fmt.Errorf("`--in:%s` has an invalid port number: `%s`", label, port)
	}

	if portNum, err := strconv.ParseInt(port, 10, 64); err == nil {
		if portNum == 0 {
			return fmt.Errorf("`--in:%s` has an invalid port number: `%s`; not going to dynamically assign port when specifying port == 0", label, port)
		}

		if portNum > 65535 || portNum < 1 {
			return fmt.Errorf("`--in:%s` has an invalid port number: `%s`; ports must be between 1 and 65535, inclusively", label, port)
		}
	} else {
		return fmt.Errorf("`--in:%s` has an invalid port number: `%s`", label, port)
	}

	return nil
}

func isOnlyNumbers(subj string) bool {
	return strings.IndexFunc(subj, nonDigit) == -1
}

func nonDigit(c rune) bool {
	return c < '0' || c > '9'
}
