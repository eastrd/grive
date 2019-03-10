package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func getSize(path string) int64 {
	f, err := os.Open(path)
	checkErr(err)

	fileInfo, err := f.Stat()
	checkErr(err)

	return fileInfo.Size()
}

func uploadBigFile(path string, size int64) {
	f, err := os.Open(path)
	checkErr(err)

	// Fetch all Google Accounts & Construct a round queue
	srvs := getAllAccounts(accountConfig)
	pos := 0

	for {
		// Fetch the next chunk and upload
		content := make([]byte, size)
		_, err := f.Read(content)
		if err == io.EOF {
			break
		}

		hash := md5.Sum(content)
		checksum := hex.EncodeToString(hash[:])
		fmt.Println("Chunk detected:", len(content), checksum)

		// Round Robin gDrives
		srv := srvs[pos]

		fmt.Println("Using", getUserInfo(srv).EmailAddress)

		pos++
		if pos == len(srvs) {
			pos = 0
		}

		// Upload
		f, err := createFile(srv, checksum, bytes.NewReader(content))
		checkErr(err)
		fID := f.Id
		fmt.Println(fID)
		// Generate a config for this file
		/*
			filename:
				{
					Chunks: {
						checksum,
						ID,
						size,
						email,
					},
					Size:
				}
		*/

	}
}
