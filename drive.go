package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func listAllFilesCloud(service *drive.Service) []*drive.File {
	allFiles := make([]*drive.File, 0)
	r, err := service.Files.List().Fields("files(id, name, size)").Do()
	checkErr(err)

	if len(r.Files) > 0 {
		for _, f := range r.Files {
			if f.Size > 0 {
				// Google Docs / Slides files are free storage, hence ignore them. (Not our concern)
				allFiles = append(allFiles, f)
			}
		}
	}
	return allFiles
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func retrieveAccount(name string) *drive.Service {
	b, err := ioutil.ReadFile(ACCDIR + name + ".json")
	checkErr(err)

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	checkErr(err)

	client := _getClient(config, name)

	srv, err := drive.New(client)
	checkErr(err)

	return srv
}

// Retrieve a token, saves the token, then returns the generated client.
func _getClient(config *oauth2.Config, name string) *http.Client {
	// The file {username}_token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := ACCDIR + name + "_token.json"
	tok, err := _tokenFromFile(tokFile)
	if err != nil {
		tok = _getTokenFromWeb(config)
		_saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func _getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Go to link and enter auth code:", authURL)

	var authCode string
	_, err := fmt.Scan(&authCode)
	checkErr(err)

	tok, err := config.Exchange(context.TODO(), authCode)
	checkErr(err)

	return tok
}

// Retrieves a token from a local file.
func _tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func _saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	checkErr(err)

	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func createFileCloud(service *drive.Service, name string, content io.Reader) (*drive.File, error) {
	f := &drive.File{
		MimeType: "application/x-grivefile",
		Name:     name,
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func createDirCloud(service *drive.Service, name string, parentID string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

func deleteFileCloud(service *drive.Service, fileID string) error {
	err := service.Files.Delete(fileID).Do()
	return err
}

func downloadFileCloud(service *drive.Service, fileID string) []byte {
	resp, err := service.Files.Get(fileID).Download()
	checkErr(err)
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	return content
}

func getAllAccounts(configPath string) []*drive.Service {
	fb, err := ioutil.ReadFile(ACCCONFIG)
	checkErr(err)
	names := strings.Split(string(fb), "\r\n")

	srvs := make([]*drive.Service, 0)
	for _, name := range names {
		srvs = append(srvs, retrieveAccount(name))
	}
	return srvs
}

func getUserInfo(service *drive.Service) *drive.User {
	about, err := service.About.Get().Fields("user").Do()
	checkErr(err)

	// fmt.Println(about.User.EmailAddress)
	// fmt.Println(about.User.DisplayName)
	return about.User
}

func getUsageQuota(service *drive.Service) *drive.AboutStorageQuota {
	about, err := service.About.Get().Fields("storageQuota").Do()
	checkErr(err)

	// fmt.Println("Limit:", about.StorageQuota.Limit/1024/1024, "MB")
	// fmt.Println("Total Usage Across all services (Gmail, Photos, Drive):", about.StorageQuota.Usage/1024/1024, "MB")
	// fmt.Println("UsageInDrive:", about.StorageQuota.UsageInDrive/1024/1024, "MB")
	// fmt.Println("UsageInDriveTrash:", about.StorageQuota.UsageInDriveTrash/1024/1024, "MB")
	return about.StorageQuota
}

// srv := retrieveAccount("e0t3rx")

// for _, f := range listAllFiles(srv) {
// 	fmt.Println(f.Name, f.Id, f.Size)
// }

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
