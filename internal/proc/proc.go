package proc

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

var ErrProcessNotFound = errors.New("process not found")
var ErrFailedToTerminateProcess = errors.New("failed to terminate process")

func ReadAllProcesses(filter []string, grep string, systemUser string) ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var allProcesses []ProcessInfo
	for _, p := range processes {
		processName, _ := p.Name()
		processExe, _ := p.Exe()
		processUsername, _ := p.Username()
		terminal, _ := p.Terminal()
		parentPid, _ := p.Ppid()
		if len(filter) > 0 && !slices.Contains(filter, strings.ToLower(processName)) {
			continue
		}
		if systemUser != "" && processUsername != systemUser {
			continue
		}
		if grep != "" && !strings.Contains(strings.ToLower(processName), grep) {
			continue
		}

		allProcesses = append(
			allProcesses, ProcessInfo{
				Pid:       p.Pid,
				Name:      processName,
				User:      processUsername,
				Exe:       processExe,
				ParentPid: parentPid,
				Terminal:  terminal,
			},
		)
	}

	slices.SortFunc(
		allProcesses, func(a, b ProcessInfo) int {
			if a.User != b.User {
				return strings.Compare(a.User, b.User)
			}
			if a.Name != b.Name {
				return strings.Compare(a.Name, b.Name)
			}
			return int(a.Pid - b.Pid)
		},
	)
	return allProcesses, nil
}

func TermiateProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("%w pid=%v: %w", ErrProcessNotFound, pid, err)
	}

	err = p.Terminate()
	if err != nil {
		return fmt.Errorf("%w pid=%v: %w", ErrFailedToTerminateProcess, pid, err)
	}

	return nil
}
