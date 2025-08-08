# same - File integrity checker

## Usage

```bash
git clone https://github.com/rohitaryal/same.git
cd same
go build -o same ./cmd/same
./same --help
```

## Help

```bash
same: File integrity checker
  -b    Initiates backup mode of a directory (default true)
  -c    Initiates checkup mode using a backup file
  -dir string
        Directory that needs to be backed up (default ".")
  -file string
        Path to save backup file (default "backup1754619205")
  -hash string
        Hash to use for integrity check (default "MD5")
  -v    Make the operation verbose

Author: @rohitaryal :)
```
