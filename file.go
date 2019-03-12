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
			TotalSize: Total size of the file
			AverageChunkSize: Average chunk size, not including the tailing chunk (equal or smaller than average)
		}
*/

// File .
type File struct {
	TotalSize    int64
	AvgChunkSize int64
	Chunks       []Chunk
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
	// Get the size of the file in local system
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
		TotalSize:    getSize(path),
		AvgChunkSize: size,
		Chunks:       chunks,
	}

	// Fetch all Google Accounts & Construct a round queue
	srvs := getAllAccounts(ACCCONFIG)
	pos := 0

	for {
		// Fetch the next chunk and upload
		content := make([]byte, size)
		_, err := f.Read(content)
		content = bytes.TrimRight(content, "\x00")
		fmt.Print("Content is: ")
		fmt.Println(content)
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
		f, err := createFileCloud(srv, checksum, bytes.NewReader(content))
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
	_ = ioutil.WriteFile(CLOUDDIR+fName, fOut, 0644)
}

func _getFileSt(fName string) File {
	content, err := ioutil.ReadFile(CLOUDDIR + fName)
	checkErr(err)

	fileSt := File{}
	err = json.Unmarshal(content, &fileSt)
	checkErr(err)

	return fileSt
}

func getAllFileStInfo() {
	fs, err := ioutil.ReadDir(CLOUDDIR)
	checkErr(err)

	for _, f := range fs {
		fName := f.Name()
		fileSt := _getFileSt(fName)

		fmt.Println("-> "+fName+":\n", " Chunk Size:", fileSt.AvgChunkSize, "  Total Size:", fileSt.TotalSize)
	}
}

func deleteFileSt(fName string) {
	fmt.Println("Remove", fName)
	// Remove all chunks that a File struct points to
	fileSt := _getFileSt(fName)

	srvMapper := makeEmailSrvMapper()

	for _, chunk := range fileSt.Chunks {
		fmt.Println("Deleting", chunk.FileID, "from", chunk.Email)
		err := deleteFileCloud(srvMapper[chunk.Email], chunk.FileID)
		checkErr(err)
	}

	// Remove the config file at last
	err := os.Remove(CLOUDDIR + fName)
	checkErr(err)
}

func makeEmailSrvMapper() map[string]*drive.Service {
	// Construct an email:serviceVar mapping
	fmt.Println("Searching for chunks in accounts...")
	srvMapper := make(map[string]*drive.Service)
	for _, srv := range getAllAccounts(ACCCONFIG) {
		srvMapper[getUserInfo(srv).EmailAddress] = srv
	}
	return srvMapper
}

func downloadFile(fName string) {
	// Download chunks from a given file name
	fSt := _getFileSt(fName)
	srvMapper := makeEmailSrvMapper()

	// Fetch data from each chunk
	for i, c := range fSt.Chunks {
		fmt.Println("Downloading chunk", i+1, "out of", len(fSt.Chunks))
		b := downloadFileCloud(srvMapper[c.Email], c.FileID)

		// Append chunk to the file
		f, err := os.OpenFile(LOCALDIR+fName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkErr(err)
		defer f.Close()
		_, err = f.Write(b)
		checkErr(err)
	}
}
