# Grive
##### A Golang Google Drive Parallelization Tool

### What Is It?
- Combine multiple Google Drives into one continuous storage space.
- Allows you to bind an infinite amount of Google Drives altogether.

### What Does It Do?
- It splits a large file into small chunks and upload them on the Google Drives.
- The chunks are stored across all available Google Drives with no obvious sequences.
- The chunks can be encrypted(TODO) and the sequence is recorded locally.
- The Google Drives can be configured to be a RAID(TODO) in case you lose access to one of your Google Drives. 

### Why Do I Need It?
- If an attacker gains access to one of your Google Drives, they won't be able to access the complete data, as the chunks are obfuscated and mixed across Drives.
- Infinite Backups...
- Large sensitive files to be stored safely on cloud.

### Usage:
- `$ grive file`: Show a list of files in Google Drives that's being uploaded by Grive.
- `$ grive file upload [filename] [chunk size in MB]`: Upload a file onto Google Drives.
- `$ grive file get [filename]`: Download a file from Google Drives.
- `$ grive file delete [filename]`: Remove a file from Google Drives.
- `$ grive accounts`: Shows the accounts details that is connected to Grive.
- `$ grive space`: Show Google Drives space usage.

### Setup Tutorial:
1. Install `Golang` environment, and install dependencies for Google Drives: `go get golang.org/x/oauth2` and `go get google.golang.org/api/drive/v3`.
2. Create two folders: `uploaded` and `downloaded` in this repo after `git clone`. `uploaded` directory will be used to store uploaded files' configurations, and `downloaded` directory contains all files that are downloaded by Grive.
3. Goto https://developers.google.com/drive/api/v3/quickstart/go and click `ENABLE THE DRIVE API`, then `DOWNLOAD CLIENT CONFIGURATION` and rename the `credential.json` into any JSON file you like. (In my case, I just renamed it to be my {AccountName}.json)
4. Place the JSON file into `accounts/` folder.
5. Add the name of your JSON file without the suffix into `accounts/accounts.txt` in a newline. (In my case, just {AccountName})
6. Execute `$ grive accounts` and Grive will ask you to click the link and paste the token in the terminal, do so.
7. You are good to go :D 


### Build Grive
- Easy, just run `$ go build -o grive *.go`



### TODO
- Utilize goroutines to increase network I/O efficiencies
- Craft a frontend & backend for Grive (use websocket)