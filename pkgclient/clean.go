package pkgclient

import "os/exec"

func ExecClean() error {
	//cmd := exec.Command("bash", "-c", "sleep 5 && rm go_reversh_client")
	cmd := exec.Command("bash", "-c", "sleep 5 && rm doc.txt")

	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
