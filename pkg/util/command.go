/**
 * Package util
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/18 11:43
 */

package util

import (
	"bytes"
	"os/exec"
	"syscall"
	"time"
)

func RunCommand(name string, arg []string, timeout time.Duration) ([]byte, []byte, error) {
	var (
		stdout = bytes.NewBuffer(nil)
		stderr = bytes.NewBuffer(nil)
	)

	cmd := exec.Command(name, arg...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	time.AfterFunc(timeout, func() {
		if cmd.Process != nil {
			//The minus sign in -cmd.Process.Pid is not a mathematical operation on PID,
			//but a special use of Go's syscall.Kill to indicate that it applies to the entire process group.
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	})

	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}
