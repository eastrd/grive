package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
)

/*
	filename:
		{
			Chunks: {
				checksum,
				FileID,
				size,
				email,
			},
			Size:
		}
*/

// File .
type File struct {
	TotalSize int64
	ChunkSize int64
	Chunks    []Chunk
}

// Chunk .
type Chunk struct {
	Checksum string
	FileID   string
	Email    string
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func getSize(path string) int64 {
	f, err := os.Open(path)
	checkErr(err)

	fileInfo, err := f.Stat()
	checkErr(err)

	return fileInfo.Size()
}

func uploadBigFile(path string, size int64) {
	// Split the given file into chunks and upload each chunk onto gDrives
	// Saves the chunk configurations as JSON file
	f, err := os.Open(path)
	checkErr(err)

	chunks := make([]Chunk, 0)
	fileSt := File{
		TotalSize: getSize(path),
		ChunkSize: size,
		Chunks:    chunks,
	}

	// Fetch all Google Accounts & Construct a round queue
	srvs := getAllAccounts(ACCCONFIG)
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
		fmt.Println("Chunk detected:", len(content)/1024, "KB, Checksum:", checksum)

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
		fileSt.Chunks = append(fileSt.Chunks, Chunk{
			Checksum: checksum,
			FileID:   f.Id,
			Email:    getUserInfo(srv).EmailAddress,
		})
	}
	fOut, err := json.MarshalIndent(fileSt, "", " ")
	checkErr(err)

	// Extract filename from path
	s := strings.Split(path, "/")
	fName := s[len(s)-1]
	// Store to storage/
	_ = ioutil.WriteFile(STORAGEDIR+fName, fOut, 0644)
}

func _getFileSt(fName string) File {
	content, err := ioutil.ReadFile(STORAGEDIR + fName)
	checkErr(err)

	fileSt := File{}
	err = json.Unmarshal(content, &fileSt)
	checkErr(err)

	return fileSt
}

func getAllFileStInfo() {
	fs, err := ioutil.ReadDir(STORAGEDIR)
	checkErr(err)

	for _, f := range fs {
		fName := f.Name()
		fileSt := _getFileSt(fName)

		fmt.Println("-> "+fName+":\n", " Chunk Size:", fileSt.ChunkSize, "  Total Size:", fileSt.TotalSize)
		for _, chunk := range fileSt.Chunks {
			fmt.Println("["+chunk.Email+"]: ", " Checksum:", chunk.Checksum, " FileID:", chunk.FileID)
		}
	}
}

func deleteFileSt(fName string) {
	fmt.Println("Remove", fName)
	// Remove all chunks that a File struct points to
	fileSt := _getFileSt(fName)

	// Construct an email:serviceVar mapping
	fmt.Println("Searching for chunks in accounts...")
	srvMapper := make(map[string]*drive.Service)
	for _, srv := range getAllAccounts(ACCCONFIG) {
		srvMapper[getUserInfo(srv).EmailAddress] = srv
	}

	for _, chunk := range fileSt.Chunks {
		fmt.Println("Deleting", chunk.FileID, "from", chunk.Email)
		err := deleteFile(srvMapper[chunk.Email], chunk.FileID)
		checkErr(err)
	}
}
