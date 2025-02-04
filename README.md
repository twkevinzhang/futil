# futil

**futil** is a command-line utility written in Golang for basic file operations, such as counting lines and calculating file checksums. The tool supports reading from both files and standard input, and it can calculate checksums using multiple algorithms (md5, sha1, and sha256). In addition, it handles error cases such as non-existent files, directories, or binary files (for line count).

## Acknowledgements and Honest Disclosure

This project was developed as part of an interview assignment.
A significant portion of the code was generated with the assistance of Github Copilot, but I performed thorough code reviews and handled the CI/CD integration myself.
Leveraging AI tools to enhance work efficiency without compromising transparency and integrity is a core principle of this work.

---

## Project Design

The project is organized into several key parts:

- **`cmd/main.go`**  
  The main entry point for the application. It uses [github.com/urfave/cli/v2](https://github.com/urfave/cli/tree/v2) to implement a CLI with the following subcommands:
    - **`linecount`**: Reads file content (or standard input) and counts the number of lines. If the file is detected as binary, it returns an error.
    - **`checksum`**: Computes the checksum of a file using one of three supported algorithms (md5, sha1, sha256). The user must specify exactly one algorithm.
    - **`version`**: Displays the version information of the application.

- **`internal/utils/utils.go`**  
  Contains common helper functions, such as opening files and detecting binary content.

- **`main_test.go`**  
  Contains unit tests for the `lineCountCmd` and `checksumCmd` functions, covering scenarios described in the project requirements.

### Design Considerations

- **Modular Structure**: Common functions (like file opening and binary detection) are placed in the `utils` folder so that the main file focuses solely on CLI logic.
- **Robust Error Handling**: The tool returns clear error messages when files do not exist, when a directory is passed instead of a file, or when a binary file is used for line counting.
- **Flag Parsing**: The [urfave/cli/v2](https://github.com/urfave/cli/tree/v2) package is used for command-line argument parsing, which automatically generates help and version information.

---

## Third-Party Libraries

- [**github.com/urfave/cli/v2**](https://github.com/urfave/cli/tree/v2)  
  Provides a simple way to build command-line applications, including support for subcommands, flag parsing, and automatic help/version generation.

---

## Building and Running

### Prerequisites

1. Ensure that [Go](https://golang.org/dl/) is installed on your system.
2. Clone the repository and navigate to the project directory:
   ```bash
   git clone https://github.com/twkevinzhang/futil.git
   cd futil
   ```

### Building
Use Go modules to build the project:
```bash
go mod tidy
go build -o futil ./cmd
```

### Usage Examples

#### 1. Count Lines in a File
Prepare an input file:
```bash
cat <<EOF > myfile.txt
how do
you
turn this
on
EOF
```

Count the number of lines:
```bash
$ ./futil linecount -f myfile.txt
4

$ ./futil linecount --file myfile.txt
4
```

#### 2. Compute File Checksum
Calculate checksums using different algorithms:
```bash
$ ./futil checksum -f myfile.txt --md5
a8c5d553ed101646036a811772ffbdd8

$ ./futil checksum -f myfile.txt --sha1
a656582ca3143a5f48718f4a15e7df018d286521

$ ./futil checksum -f myfile.txt --sha256
495a3496cfd90e68a53b5e3ff4f9833b431fe996298f5a28228240ee2a25c09d
```

#### 3. Read from Standard Input
```bash
$ cat myfile.txt | ./futil linecount -f -
4

$ cat myfile.txt | ./futil checksum -f - --sha256
495a3496cfd90e68a53b5e3ff4f9833b431fe996298f5a28228240ee2a25c09d
```

#### 4. Display Version and Help
```bash
$ ./futil version
futil v0.0.1

$ ./futil help
File Utility

Usage:
  futil [global options] command [command options] [arguments...]

Commands:
   linecount    Print line count of file
   checksum     Print checksum of file
   version      Show the version info
   help, h      Shows a list of commands or help for one command

Global Options:
   --help, -h  show help
   
$ ./futil help linecount
NAME:
   futil linecount - Print line count of file

USAGE:
   futil linecount [flags]

OPTIONS:
   --file value, -f value  the input file
   --help, -h              show help
```

## Unit Testing
Unit tests have been provided to ensure the correctness of the `lineCountCmd` and `checksumCmd` functions. To run the tests, execute:
```bash
go test -v ./cmd
```
