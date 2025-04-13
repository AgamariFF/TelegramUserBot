package internal

import (
	"io"
	"os"
	"os/exec"
	"syscall"
)

func IsBrowserRunning(profilePath string) bool {
	// Проверка через 3 разных метода
	return isFileLocked(profilePath) ||
		checkChromeProcesses(profilePath)
}

func isFileLocked(path string) bool {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return false
	}
	defer file.Close()

	// Пытаемся получить эксклюзивную блокировку через fcntl
	flock := syscall.Flock_t{
		Type:   syscall.F_WRLCK,
		Whence: io.SeekStart,
		Start:  0,
		Len:    0,
	}

	err = syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &flock)
	if err != nil {
		return true // Файл заблокирован
	}

	// Снимаем блокировку
	flock.Type = syscall.F_UNLCK
	syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &flock)
	return false
}

func checkChromeProcesses(profilePath string) bool {
	cmd := exec.Command("pgrep", "-f", "chrome.*"+profilePath)
	output, _ := cmd.CombinedOutput()
	return len(output) > 0
}
