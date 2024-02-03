package service

import "os/exec"

func ReloadZone() (string, error) {
	cmd := exec.Command("rndc", "reload chn")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

func DnsServiceStatus() (string, error) {
	cmd := exec.Command("systemctl", "status", "named")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

func RestartDnsService() (string, error) {
	cmd := exec.Command("systemctl", "restart", "named")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
