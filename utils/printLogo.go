package utils

import "fmt"

func PrintAsciiLogoAndInfo() {
	PrintAsciiShadowLogo()
	PrintVersion()
	PrintUrl()
	fmt.Println("")
}

func PrintVersion() {
	fmt.Println("   " + GetVersion())
}

func PrintUrl() {
	fmt.Println("   " + GetUrl())
}

func GetUrl() string {
	return "https://bitlong.gitbook.io/api-doc/fairlaunch/fairlaunch-trade-rest-api"
}

func GetVersion() string {
	base := "0.0.1"
	version := base + "-" + GetTimeSuffixString()
	return version
}

func PrintAsciiShadowLogo() {
	fmt.Println("")
	fmt.Println("   ████████╗██████╗  █████╗ ██████╗ ███████╗")
	fmt.Println("   ╚══██╔══╝██╔══██╗██╔══██╗██╔══██╗██╔════╝")
	fmt.Println("      ██║   ██████╔╝███████║██║  ██║█████╗  ")
	fmt.Println("      ██║   ██╔══██╗██╔══██║██║  ██║██╔══╝  ")
	fmt.Println("      ██║   ██║  ██║██║  ██║██████╔╝███████╗")
	fmt.Println("      ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚══════╝")
	//fmt.Println("")
}
