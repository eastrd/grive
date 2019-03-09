package main

import (
	"fmt"
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	srv := retrieveAccount("e0t3rx")

	for _, f := range listAllFiles(srv) {
		fmt.Println(f.Name, f.Id, f.Size)
	}

	// Create a sample file
	// f, err := os.Open("a.png")
	// checkErr(err)
	// defer f.Close()

	// file, err := createFile(srv, "a.png", f, "root")
	// checkErr(err)

	// fmt.Println(file.Name, file.Id)
	// err = deleteFile(srv, "1_yeT7cLmbtx5PgggFE5rQepEtQW-d-QQ")
	// checkErr(err)

	// downloadFile(srv, "1bLObXT3D3ZQgMcjLH1TftvZ1v7WRcE-_", "aaa.png")
	// fmt.Println(getUsageQuota(srv))
}

/*
	$grive accounts
	$grive space
	$grive file
*/
