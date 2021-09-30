package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "Library"), "Preferences")
	return filepath.Join(desktopConfig, "myhosts")
}

func SaveHosts(content, pwd string) (err error) {
	cmd := exec.Command("sudo", "-S", "bash", "-c", fmt.Sprintf(`echo "%s" > %s`, content, HostsFile()))
	cmd.Stdin = bytes.NewReader([]byte(pwd))
	_, err = cmd.Output()
	return
}

func HostsFile() string {
	return "/etc/hosts"
}
