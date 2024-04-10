# FileTree-API
FileTree-API is a high-performance, scalable server application built with Go that provides a RESTful API for recursively traversing and listing files within a specified directory structure. Utilizing the concurrency power of Go's goroutines and optimized file system traversal strategies, this API is designed to efficiently handle large-scale file operations across multi-threaded environments. Packaged with Docker, the FileTree-API is easy to deploy and integrate into any infrastructure, offering a robust solution for applications requiring rapid access to file metadata and structure. Perfect for cloud storage services, file management systems, or any application dealing with extensive file directories.

## Features
- **Secure Encryption**: Uses AES encryption and signature verification to secure directory path.
- **Easy Integration**: Built with the Go, allowing for straightforward integration into any Go application.
- **Optimized Performance**:
    - **Concurrent File Traversal**: Utilizes Go's goroutines to traverse directories concurrently.
    - **Efficient File System Operations**: Implements optimized file system traversal strategies for high-performance file operations.
    - **Scalable Architecture**: Designed to handle large-scale file operations across multi-threaded environments.
    - **Optimize JSON Marshalling**: Uses `jsoniter` for faster JSON marshalling.
    - **Support WebSockets**: Supports WebSocket for real-time file system monitoring.
- **Dockerized**: Packaged with Docker for easy deployment and integration into any infrastructure.


## Getting Started
1. Build the Go application
```bash
go build -o filetree-api ./cmd/server
```

2. Create the necessary environment variables for secure encryption:
    Generate `FILETREE_SECRET_KEY`:
    ```sh
    echo FILETREE_SECRET_KEY=$(xxd -g 2 -l 32 -p /dev/random | tr -d '\n')
    ```

    Generate `FILETREE_SECRET_SALT`:
    ```sh
    echo FILETREE_SECRET_SALT=$(xxd -g 2 -l 32 -p /dev/random | tr -d '\n')
    ```

    After generating these keys, make sure to set them as environment variables in your development environment or include them in your deployment configuration.

3. Set the environment variables
```bash
export FILETREE_PORT=8080
export FILETREE_SECRET_KEY=your-secret-key
export FILETREE_SECRET_SALT=your-secret-salt
```

4. Run the Go application
```bash
./filetree-api
```

## API Usage
Make a GET request to the service with a signature and an encrypted folder path to retrieve the file tree structure of the specified directory.  
The signature is generated using the `FILETREE_SECRET_KEY` and `FILETREE_SECRET_SALT` environment variables.  
The encrypted folder path is generated using the `FILETREE_SECRET_KEY` and `FILETREE_SECRET_SALT` environment variables.
```
http://your-domain.com:your-port/<signature>/enc/<encrypted_folder_path>
```

## Projects Using FileTree-API
Several projects are built on top of or with FileTree-API to extend its capabilities and offer more features. Here's a list of such projects:

- [PHP-FileTree](https://github.com/carry0987/PHP-FileTree): A PHP script for generating signed and encrypted URLs with FileTree-API, using AES-256-GCM and HMAC-SHA256.

We encourage the community to build more projects leveraging FileTree-API's powerful image processing capabilities. If you have a project that uses FileTree-API, feel free to open a pull request to add it to this list!

## Contributing
We welcome all forms of contributions, whether it be submitting issues, writing documentation, or sending pull requests.

## License
This project is licensed under the [MIT](LICENSE) License.
