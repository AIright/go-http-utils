package go_http_utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type readinessInfo struct {
	IP              string `json:"ip,omitempty"`
	Host            string `json:"host,omitempty"`
	OS              string `json:"os,omitempty"`
	Language        string `json:"language,omitempty"`
	LanguageVersion string `json:"languageVersion,omitempty"`
	GitCommit       string `json:"gitCommit,omitempty"`
}

func readinessProbe() http.Handler {
	inf := readinessInfo{
		IP:   execute("hostname", "-I"),
		Host: execute("uname", "-n"),
		OS: strings.Join([]string{
			execute("uname", "-s"),
			execute("uname", "-r"),
			execute("uname", "-v"),
			execute("uname", "-m"),
		}, "  "),
		Language:        "go",
		LanguageVersion: runtime.Version(),
		GitCommit:       os.Getenv("GIT_COMMIT"),
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		body, _ := json.MarshalIndent(&inf, "", "  ")
		_, _ = w.Write(body)
	})
}

func execute(name string, args ...string) string {
	var out bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(out.String())
}
