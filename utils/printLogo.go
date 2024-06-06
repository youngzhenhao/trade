package utils

import "fmt"

func PrintAsciiLogo() {
	PrintAsciiShadowLogo()
	fmt.Println(GetTimeSuffixString())
}

func PrintAsciiShadowLogo() {
	fmt.Println("")
	fmt.Println("   ████████╗██████╗  █████╗ ██████╗ ███████╗")
	fmt.Println("   ╚══██╔══╝██╔══██╗██╔══██╗██╔══██╗██╔════╝")
	fmt.Println("      ██║   ██████╔╝███████║██║  ██║█████╗  ")
	fmt.Println("      ██║   ██╔══██╗██╔══██║██║  ██║██╔══╝  ")
	fmt.Println("      ██║   ██║  ██║██║  ██║██████╔╝███████╗")
	fmt.Println("      ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚══════╝")
	fmt.Println("")
}
