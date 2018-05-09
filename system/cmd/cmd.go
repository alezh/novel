package cmd

import (
	"github.com/kataras/golog"
	"os"
	"os/exec"
)

func installVersion() {
	golog.Infof("Downloading...\n")
	repo := "github.com/kataras/iris/..."
	cmd := exec.Command("go", "get", "-u", "-v", repo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		golog.Warnf("unexpected message while trying to go get,\nif you edited the original source code then you've to remove the whole $GOPATH/src/github.com/kataras folder and execute `go get -u github.com/kataras/iris/...` manually\n%v", err)
		return
	}

	golog.Infof("Update process finished.\nManual rebuild and restart is required to apply the changes...\n")
	return
}