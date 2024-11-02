package main

import main2 "gitlab.lainuoniao.cn/rhinobird/backend/eden-cli.git/eden-cli"

func main() {
	rootCmd := main2.createCmd()
	_ = rootCmd.Execute()
}
