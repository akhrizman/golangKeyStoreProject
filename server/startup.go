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
		//fmt.Println("Missing or Invalid Port")
		os.Exit(ExitStatus(-1))
	}
	return p
}

// ExitOnErrors Check for Port Binding and successful start of application, otherwise exit
func ExitOnErrors(port int, err error) {
	switch t := err.(type) {
	case *net.OpError:
		fmt.Printf("Error Binding port %d - %s\n", port, t)
		log4g.Error.Printf("Error Binding port %d - %s\n", port, t)
		os.Exit(ExitStatus(-2))
	default:
		fmt.Println("Error Starting Server - ", err)
		log4g.Error.Println("Error Starting Server - ", err)
	}
}

// ExitStatus Make negative exit status codes positive for windows
func ExitStatus(i int) int {
	hostOS := runtime.GOOS
	//fmt.Println("Detected OS:", hostOS)
	if hostOS == "windows" && i < 0 {
		i *= -1
	}
	return i
}

func ApiHost(port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
