package process

/**
   @author yhy
   @since 2024/1/19
   @desc //TODO
**/

type Info struct {
    User         string `json:"user"`
    Pid          string `json:"pid"`
    Cpu          string `json:"cpu"`
    Mem          string `json:"mem"`
    VsZ          string `json:"vsz"`
    Rss          string `json:"rss"`
    Tty          string `json:"tty"`
    Stat         string `json:"stat"`
    Start        string `json:"start"`
    Time         string `json:"time"`
    AbsolutePath string `json:"absolutePath"` // 该进程的绝对路径，有的程序是通过相对路径启动的，所以需要获取绝对路径
    Command      string `json:"command"`
    CommandArr   string `json:"commandArr"`
    
    Counter int `json:"counter"` // 重启次数
}
