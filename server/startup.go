package server

import (
	"flag"
	"fmt"
	"httpstore/log4g"
	"net"
	"os"
	"runtime"
	"strconv"
)

var (
	validPortRangeMin = 1024
	validPortRangeMax = 65535
	host              = "localhost"
)

// ValidatePort Set port globally if arg value is valid, otherwise exit
func ValidatePort() int {
	portArg := flag.String("port", "none", "server port")
	flag.Parse()

	p, err := strconv.Atoi(*portArg)
	if err != nil || p < validPortRangeMin || p > validPortRangeMax {
		os.Exit(ExitStatus(-1))
	}
	return p
}

// ExitOnErrors Check for Port Binding and successful start of application, otherwise exit
func ExitOnErrors(port int, err error) {
	switch t := err.(type) {
	case *net.OpError:
		log4g.Error.Printf("Error Binding port %d - %s\n", port, t)
		os.Exit(ExitStatus(-2))
	default:
		log4g.Error.Println("Error Starting Server - ", err)
	}
}

// ExitStatus Make negative exit status codes positive for windows
func ExitStatus(i int) int {
	hostOS := runtime.GOOS
	if hostOS == "windows" && i < 0 {
		i *= -1
	}
	return i
}

func ApiHost(port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
