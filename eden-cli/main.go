package main

//import m "gitlab.lainuoniao.cn/rhinobird/backend/eden-cli.git/eden-cli"

func main() {
	rootCmd := createCmd()
	_ = rootCmd.Execute()
}
