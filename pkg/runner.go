package pkg

import (
    "github.com/yhy0/Ouroboros/conf"
    "github.com/yhy0/Ouroboros/pkg/process"
    "github.com/yhy0/logging"
    "os/user"
    "time"
)

/**
   @author yhy
   @since 2024/1/20
   @desc //TODO
**/

func Run() {
    if !conf.AllProcesses {
        current, err := user.Current()
        if err != nil {
            logging.Logger.Error(err)
            return
        }
        conf.CurrentUser = current.Username
        logging.Logger.Infof("CurrentUser is %s", conf.CurrentUser)
    }
    
    process.GetUserProcesses()
    
    Monitor()
}

// Monitor 每 5s 监控进程
func Monitor() {
    for {
        for pid := range process.MonitorPid {
            if !process.IsProcessRunning(pid) {
                logging.Logger.Warnf("Process '[%s]:%s' not running. Restarting", pid, process.MonitorPid[pid].CommandArr)
                process.RestartProcess(process.MonitorPid[pid])
            }
        }
        time.Sleep(time.Duration(conf.Interval) * time.Second)
    }
}
