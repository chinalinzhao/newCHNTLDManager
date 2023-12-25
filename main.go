package main

import (
	_ "newCHNTLDManager/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"newCHNTLDManager/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
