# f - A File Manager CLI

`f` is a file manager CLI that's much more intuitive than the built-in commands. It supports operations like copying, moving, and renaming files and directories with ease.

## Installation

### From Source

To install `f` from source, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/cdrice26/f.git
    cd f
    ```

2. Build the project:
    ```sh
    go build -o f
    ```

3. Move the binary to a directory in your PATH:
    ```sh
    mv f /usr/local/bin/
    ```

## Usage

Note that for the copy and move commands, the program will create nonexistent directories for you - no need to use `mkdir`!

### Copy Files and Directories

To copy files or directories, use the `copy` command:
```sh
f copy <source>... <destination>
```
Example:

```sh
f copy file1.txt file2.txt /path/to/destination/
```

The following flags are supported:
- `-o`, `--overwrite` - Overwrite the destination file if it exists.

### Move Files and Directories


To move files or directories, use the move command:
```sh
f move <source>... <destination>
```

Example:

```sh
f move file1.txt file2.txt /path/to/destination/
```

The following flags are supported:
- `-o`, `--overwrite` - Overwrite the destination file if it exists.

### Rename Files and Directories
To rename a file or directory, use the rename command:
```sh
f rename <source> <destination>
```

Example:

```sh
f rename file.txt new_file.txt
```

The following flags are supported:
- `-o`, `--overwrite` - Overwrite the destination file if it exists.

### Delete Files and Directories
To delete a file or directory, use the delete command:
```sh
f delete <source>...
```

Example:

```sh
f delete file.txt
```
The following flags are supported:
- `-f`, `--force` - Force deletion without prompting for confirmation.

### List Files in a Directory
To list files in a directory, use the list command:
```sh
f list [directory]
```
If no directory is provided, the working directory is used.
The following flags are supported:
- `-n`, `--no-directory-sizes` - By default, sizes are calculated for directories, which can be time-consuming. This flag disables it.
- `-t`, `--tree` - Show all subdirectories and files in a tree-style output.
- `-a`, `--hidden` - Include hidden files and directories in the output.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
