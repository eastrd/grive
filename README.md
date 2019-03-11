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
- `$ grive accounts`: Shows the accounts details that is connected to Grive.
- `$ grive upload [filename] [chunk size in MB]`: Upload a file onto Google Drives.
- `$ grive get [filename]`: Download a file from Google Drives.
- `$ grive delete [filename]`: Remove a file from Google Drives.