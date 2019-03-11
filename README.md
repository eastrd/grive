# Grive
##### A Golang Google Drive Parallelization Tool

### What Is It?
- Combine multiple Google Drives into one continuous storage space.
- Allows you to bind an infinite amount of Google Drives altogether.

### What Does It Do?
- It splits a large file into small chunks and upload them on the Google Drives.
- The chunks are stored across all available Google Drives with no obvious sequences.
- The chunks can be encrypted(TODO) and their sequences are only recorded on a JSON file locally.
- The Google Drives can be configured to be RAID(TODO). 

### Why Do I Need It?
- If an attacker gains access to one of your Google Drives, they won't be able to access the complete data, as the chunks are obfuscated and mixed across Drives.
- Infinite Backups...
- Large sensitive files to be stored safely on cloud.