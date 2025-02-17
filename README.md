# f - A File Manager CLI

`f` is a file manager CLI that's much more intuitive than the built-in commands. It supports operations like copying, moving, and renaming files and directories with ease.

## Installation

### From Source

To install `f` from source, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/f.git
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

### Copy Files and Directories

To copy files or directories, use the `copy` command:
```sh
f copy <source>... <destination>
```
Example:

```sh
f copy file1.txt file2.txt /path/to/destination/
```

### Move Files and Directories


To move files or directories, use the move command:
```sh
f move <source>... <destination>
```

Example:

```sh
f move file1.txt file2.txt /path/to/destination/
```

### Rename Files and Directories
To rename a file or directory, use the rename command:
```sh
f rename <source> <destination>
```

Example:

```sh
f rename file.txt new_file.txt
```

### Delete Files and Directories
To delete a file or directory, use the rename command:
```sh
f delete <source>...
```

Example:

```sh
f delete file.txt
```

## License
This project is licensed under the MIT License. See the LICENSE file for details.