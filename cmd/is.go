package cmd

import (
    "github.com/yhy0/InfiniteSnake/pkg"
    "github.com/yhy0/InfiniteSnake/pkg/web"
)

/**
   @author yhy
   @since 2024/1/19
   @desc //TODO
**/

func RunApp() {
    pkg.ParseOptions()
    go pkg.Run()
    web.Init()
}
