package main

import "fmt"

func handleCmd(cmdQue []string) {
	initial, cmdQue := cmdQue[0], cmdQue[1:]
	switch initial {
	case "accounts":
		handleAccounts(cmdQue)
	case "space":
		handleSpace(cmdQue)
	case "file":
		handleFile(cmdQue)
	case "upload":
		handleUpload(cmdQue)
	default:
		printHelpMsg()
	}
}

func handleUpload(cmdQue []string) {

}

func handleAccounts(cmdQue []string) {
	srvs := getAllAccounts(ACCCONFIG)
	fmt.Println("Connected Accounts:", len(srvs))
	for _, srv := range srvs {
		user := getUserInfo(srv)
		fmt.Println("["+user.EmailAddress+"] ", user.DisplayName)
	}
}

func handleSpace(cmdQue []string) {
	srvs := getAllAccounts(ACCCONFIG)
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

func handleFile(cmdQue []string) {
	fmt.Println("File stuff")
}

func printHelpMsg() {
	fmt.Println("Invalid option")
}
