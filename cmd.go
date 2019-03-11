package main

import (
	"fmt"
	"strconv"
)

func handleCmd(cmdQue []string) {
	initial, cmdQue := cmdQue[0], cmdQue[1:]
	switch initial {
	case "accounts":
		handleAccounts(cmdQue)
	case "space":
		handleSpace(cmdQue)
	case "file":
		handleFile(cmdQue)
	default:
		printHelpMsg()
	}
}

func handleDownload(cmdQue []string) {
	if len(cmdQue) != 1 {
		fmt.Println("Usage: download {filename}")
	}
	downloadFile(cmdQue[0])
}

func handleDelete(cmdQue []string) {
	if len(cmdQue) != 1 {
		fmt.Println("Usage: delete {filename}")
	}
	deleteFileSt(cmdQue[0])
}

func handleUpload(cmdQue []string) {
	if len(cmdQue) != 2 {
		fmt.Println("Usage: upload {path to file} {chunk size in MB}")
		return
	}
	size, err := strconv.ParseFloat(cmdQue[1], 2)
	checkErr(err)

	uploadBigFile(cmdQue[0], int64(size*1024*1024))
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
	if len(cmdQue) == 0 {
		fmt.Println("Files stored:")
		getAllFileStInfo()
		return
	}
	fileCmd, cmdQue := cmdQue[0], cmdQue[1:]
	switch fileCmd {
	case "upload":
		handleUpload(cmdQue)
	case "get":
		handleDownload(cmdQue)
	case "delete":
		handleDelete(cmdQue)
	default:
		printHelpMsg()
	}
}

func printHelpMsg() {
	fmt.Println("Grive Usage:")
	fmt.Println("$ grive file [upload/get/delete]")
	fmt.Println("$ grive accounts")
	fmt.Println("$ grive space")

}
