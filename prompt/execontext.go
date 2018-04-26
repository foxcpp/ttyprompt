package prompt

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"time"
)

var startupTime time.Time

func init() {
	startupTime = time.Now()
}

func procname(pid int) (string, error) {
	path, err := os.Readlink("/proc/" + strconv.Itoa(pid) + "/exe")
	if err != nil {
		return "", err
	}

	return path, nil
}

type ExeContext struct {
	User        string
	Group       string
	ParentPID   string
	ParentPath  string
	StartupTime time.Time
}

func CaptureExeCtx() (*ExeContext, error) {
	ctx := ExeContext{}

	current, err := user.Current()
	if err != nil {
		return nil, err
	}
	if usr, err := user.LookupId(current.Uid); err != nil {
		ctx.User = current.Uid
	} else {
		ctx.User = usr.Username
	}
	if grp, err := user.LookupGroupId(current.Gid); err != nil {
		ctx.Group = current.Gid
	} else {
		ctx.Group = grp.Name
	}

	ctx.StartupTime = startupTime

	ppid := os.Getppid()
	ctx.ParentPID = strconv.Itoa(ppid)

	if procpath, err := procname(ppid); err != nil {
		ctx.ParentPath = "-"
	} else {
		ctx.ParentPath = procpath
	}

	return &ctx, nil
}

func (ctx ExeContext) String() string {
	return fmt.Sprintf("Running as %s:%s, started by PID %s (%s) on %s",
		ctx.User, ctx.Group, ctx.ParentPID, ctx.ParentPath, ctx.StartupTime.Format(time.UnixDate))
}
