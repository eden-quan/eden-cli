package main

func main() {
	rootCmd := createCmd()
	_ = rootCmd.Execute()
}
