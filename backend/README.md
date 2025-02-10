# Sector Backend

## How to use

Follow the `kubo` installation instructions for your operating system. Kubo version 0.27.0 must be used. We highly recommend utilizing the [ipfs-update](https://docs.ipfs.tech/how-to/ipfs-updater/) tool to install this version of kubo. 

> [!Warning]
> Ensure that kubo version 0.27.0 is utilized, otherwise you may need to install migration tooling to run the databse.

If you haven't used IPFS so far, initialize the IPFS repository using the following command:

`ipfs init`

If you had used IPFS an already have an IPFS repository in place, either (re)move it from ~/.ipfs or make sure to export IPFS_PATH before running the ipfs init command, e.g.:

```
export IPFS_PATH=~/.ipfs-sector
ipfs init
```

Follow instructions in main `README.md` to generate the backend server files (if needed).
Then run the following commands to build the backend.

- Download the Go modules: `go mod download`
- Start the server: `go run main.go`
- Open the browser or any api test program to `http://127.0.0.1:3000`
- Start coding!

## Notes

- The http server listens on port `3000`, unless overridden with a command line flag. 
