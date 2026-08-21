package modules

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func RunLateral(params map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	if runtime.GOOS == "windows" {
		res["message"] = "Lateral movement on Windows requires WMI/WinRM – stub."
		return res
	}
	home, _ := os.UserHomeDir()
	sshPath := filepath.Join(home, ".ssh")
	files, _ := ioutil.ReadDir(sshPath)
	var keys []string
	for _, f := range files {
		if !f.IsDir() {
			keys = append(keys, f.Name())
		}
	}
	res["ssh_keys"] = keys
	if target, ok := params["target_ip"].(string); ok {
		res["target"] = target
		res["bruteforce"] = "SSH brute-force stub (use keys found)"
	}
	return res
}