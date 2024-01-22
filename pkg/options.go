package pkg

import (
    "flag"
    "github.com/yhy0/Ouroboros/conf"
    "github.com/yhy0/logging"
)

/**
  @author: yhy
  @since: 2023/7/26
  @desc: //TODO
**/

func ParseOptions() {
    logging.Logger = logging.New(true, "", "Ouroboros", true)
    
    flag.StringVar(&conf.WebPort, "port", "9089", "web report port, (example: 9089)")
    flag.StringVar(&conf.WebUser, "user", "yhy", "web authorized user, (example: yhy)")
    flag.StringVar(&conf.WebPass, "pwd", "", "web authorized pwd")
    flag.BoolVar(&conf.AllProcesses, "all", false, "web authorized pwd")
    flag.IntVar(&conf.Interval, "time", 5, "monitoring interval(default: 5s)")
    flag.Parse()
    
    if conf.WebPass == "" {
        conf.WebPass = RandomString()
    }
    logging.Logger.Infof("web authorized is %s/%s", conf.WebUser, conf.WebPass)
}
