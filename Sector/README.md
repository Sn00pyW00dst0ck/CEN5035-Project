# README

## About

This is the official Wails React template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

# Backend

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

- The http server listens on port `3000`.
