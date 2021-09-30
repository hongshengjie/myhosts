package app

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "AppData"), "Roaming")
	return filepath.Join(desktopConfig, "myhosts")
}

func SaveHosts(content, pwd string) error {
	f, err := os.OpenFile(HostsFile(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	f.Write([]byte(content))
	f.Close()
	return nil
}

func HostsFile() string {
	windir := "C:\\Windows"
	if dir, ok := os.LookupEnv("windir"); ok {
		windir = dir
	}
	return windir + "\\System32\\drivers\\etc\\hosts"
}
