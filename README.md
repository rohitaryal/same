# same - File integrity checker

## Usage

```bash
git clone https://github.com/rohitaryal/same.git
cd same
go build -o same ./cmd/same
./same --help
```

## How does it look?

<p style="text-align: center">
<img src="./assets/display.gif" style="max-width: 500px;border-radius:1rem;">
</p>

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
