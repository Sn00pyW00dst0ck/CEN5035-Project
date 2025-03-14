# Sector

**Team Members:** 
- Chase Hap - Frontend Engineer
- Abhignan Sai Arcot - Frontend Engineer
- Gabriel Aldous - Backend Engineer
- Sai Neha Ratakonda - Backend Engineer

## Project Description

Numerous communication platforms currently exist to allow groups to communicate effectively and rapidly. Being that they often rely upon centralized servers to store and transmit user messages, they require great funding. Many of these platforms decide to gain this funding by selling user data/analytics. Our solution ***Sector*** proposes the implementation of p2p technologies to provide users with an end-to-end encrypted solution in which no user data is stored anywhere except on user devices. With this solution, users will be able to create group chats in which all participating members contain a complete record of all messages.

## Development Setup

Utilize the `generate.sh` script to generate the frontend and backend data models/route representations.

### Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

### Building

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

## Notes

- The http server listens on port `3000`.