package scan

import (
	"net"
	"strconv"
	"time"
)

type PortState struct {
	Port int
	Open state
}

type state bool

func (s state) String() string {
	if s {
		return "Open"
	}
	return "Closed"
}

func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}
	address := net.JoinHostPort(host, strconv.Itoa(port))
	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return p
	}
	scanConn.Close()
	p.Open = true
	return p
}

type Result struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

func Run(hl *HostsList, ports []int) []Result {
	results := make([]Result, 0, len(hl.Hosts))
	for _, host := range hl.Hosts {
		result := Result{
			Host: host,
		}
		//	_, err := net.LookupAddr(host)
		//	if err != nil {
		//		log.Println("host not found,", host)
		//		result.NotFound = true
		//		results = append(results, result)
		//		continue
		//	}
		result.NotFound = false
		portStates := make([]PortState, 0, len(ports))
		for _, port := range ports {
			portState := scanPort(host, port)
			portStates = append(portStates, portState)
		}
		result.PortStates = portStates
		results = append(results, result)
	}
	return results
}
