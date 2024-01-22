package process

import (
    "fmt"
    "github.com/yhy0/Ouroboros/conf"
    "github.com/yhy0/logging"
    "os/exec"
    "strconv"
    "strings"
    "sync"
    "time"
)

/**
   @author yhy
   @since 2024/1/19
   @desc //TODO
**/

var Processes = make(map[string]*Info)

var lock sync.Mutex

// MonitorPid 监控中的 pid
var MonitorPid = make(map[string]*Info)

func GetUserProcesses() {
    Processes = make(map[string]*Info)
    cmd := exec.Command("ps", "aux")
    out, err := cmd.Output()
    if err != nil {
        logging.Logger.Error(err)
        return
    }
    
    lines := strings.Split(string(out), "\n")
    
    for _, line := range lines[1:] {
        if line == "" {
            continue
        }
        info := strings.Fields(line)
        
        if conf.CurrentUser != "" && info[0] != conf.CurrentUser {
            continue
        }
        
        process := &Info{
            User:       info[0],
            Pid:        info[1],
            Cpu:        info[2],
            Mem:        info[3],
            VsZ:        info[4],
            Rss:        info[5],
            Tty:        info[6],
            Stat:       info[7],
            Start:      info[8],
            Time:       info[9],
            Command:    info[10],
            CommandArr: strings.Join(info[10:], " "),
        }
        
        if strings.HasPrefix(process.Command, "./") || strings.Contains(process.CommandArr, "go run ") || strings.Contains(process.CommandArr, "python3 ") || strings.Contains(process.CommandArr, "python ") || strings.Contains(process.CommandArr, "java ") {
            process.AbsolutePath = GetProcessExePath(process.Pid)
        }
        
        if MonitorPid[info[1]] != nil {
            lock.Lock()
            process.Counter = MonitorPid[info[1]].Counter
            MonitorPid[info[1]] = process
            lock.Unlock()
            continue
        }
        
        lock.Lock()
        Processes[process.Pid] = process
        lock.Unlock()
    }
    return
}

// GetProcessExePath 根据 pid 获取进程的绝对路径
func GetProcessExePath(pid string) string {
    out, err := exec.Command("lsof", "-p", pid, "-F", "n").Output()
    if err != nil {
        logging.Logger.Errorf("Error: [%s]:%v", pid, err)
        return ""
    }
    
    lines := strings.Split(string(out), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "n") && !strings.Contains(line, " (") {
            return line[1:]
        }
    }
    logging.Logger.Errorf("Error: [%s]:%v", pid, "executable path not found")
    return ""
}

// IsProcessRunning 通过进程状态判断进程是否在运行
func IsProcessRunning(pid string) bool {
    cmd := exec.Command("ps", "-p", fmt.Sprintf("%s", pid), "-o", "stat=")
    output, err := cmd.Output()
    if err != nil {
        return false
    }
    
    if strings.Contains(string(output), "S") || strings.Contains(string(output), "R") {
        return true
    }
    
    return false
}

func RestartProcess(process *Info) {
    command := process.CommandArr
    
    cmd := exec.Command("sh", "-c", command)
    
    // 设置命令的工作目录
    if process.AbsolutePath != "" {
        cmd.Dir = process.AbsolutePath
    }
    
    err := cmd.Start()
    if err != nil {
        logging.Logger.Errorf("Error restarting process: [%s]:%v", process.Command, err)
        return
    }
    
    // 获取重新启动的进程的 PID
    newPid := strconv.Itoa(cmd.Process.Pid)
    
    time.Sleep(1 * time.Second)
    
    lock.Lock()
    
    // 删除旧的进程信息
    delete(MonitorPid, process.Pid)
    delete(Processes, process.Pid)
    
    // 更新为新的 PID
    process.Pid = newPid
    process.Counter += 1
    // 将新的进程信息添加回 map 中
    MonitorPid[newPid] = process
    lock.Unlock()
    
    // 重新获取一遍进程信息
    GetUserProcesses()
}
