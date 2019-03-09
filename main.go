package main

import (
	"fmt"
	"os"
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func handleCmd(cmd string) {
	switch cmd {
	case "accounts":
		printAccInfo()
	case "space":
		printSpaceInfo()
	case "file":
		printFileInfo()
	default:
		printHelpMsg()
	}
}

func printAccInfo() {
	srvs := getAllAccounts("accounts.txt")
	fmt.Println("Connected Accounts:", len(srvs))
	for _, srv := range srvs {
		user := getUserInfo(srv)
		fmt.Println("["+user.EmailAddress+"] ", user.DisplayName)
	}
}

func printSpaceInfo() {
	srvs := getAllAccounts("accounts.txt")
	totalSpace := int64(0)
	usedSpace := int64(0)

	for _, srv := range srvs {
		quota := getUsageQuota(srv)
		totalSpace += quota.Limit
		usedSpace += quota.Usage
		fmt.Println("["+getUserInfo(srv).EmailAddress+"] ", quota.Usage*100/quota.Limit, "% Used of", quota.Limit/1024/1024/1024, "GB")
	}
	fmt.Println("\nOverall:", usedSpace*100/totalSpace, "% Used  (", (totalSpace-usedSpace)/1024/1024/1024, "GB Available in", totalSpace/1024/1024/1024, "GB )")
}

func printFileInfo() {
	fmt.Println("File stuff")
}

func printHelpMsg() {
	fmt.Println("Invalid option")
}

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		handleCmd(cmd)
	}
}

/*
	$grive accounts
	$grive space
	$grive file
*/
