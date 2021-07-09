# leakless

Run sub-process and make sure to kill it when the parent process exits.
The way how it works is to output a standalone executable file to guard the subprocess and check parent TCP connection with a UUID.
So that it works consistently on Linux, Mac, and Windows.

If you don't trust the executable, you can build it yourself from the source code by running `go generate` at the root of this repo, then use the [replace](https://golang.org/ref/mod#go-mod-file-replace) to use your own module. Usually, it won't be a concern, all the executables are committed by this [Github Action](https://github.com/ysmood/leakless/actions?query=workflow%3ARelease), the Action will print the hash of the commit, you can compare it with the repo.

Not using the PID is because after a process exits, a newly created process may have the same PID.

## How to Use

See the [examples](example_test.go).
