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

const (
	validPortRangeMin = 1024
	validPortRangeMax = 65535
	host              = "localhost"
)

// ValidateFlags Set port globally if arg value is valid, otherwise exit
func ValidateFlags() (int, int) {
	portArg := flag.String("port", "none", "server port")
	depthArg := flag.String("depth", "10", "key store depth")
	flag.Parse()

	port, err := strconv.Atoi(*portArg)
	if err != nil || port < validPortRangeMin || port > validPortRangeMax {
		os.Exit(ExitStatus(-1))
	}

	depth, err := strconv.Atoi(*depthArg)
	if err != nil || depth < 0 {
		os.Exit(ExitStatus(-1))
	}
	return port, depth
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
