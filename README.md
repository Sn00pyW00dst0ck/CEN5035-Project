# Sector

A peer-to-peer end-to-end encrypted communication application

**Team Members:** 
- Chase Hap - Frontend Engineer
- Abhignan Sai Arcot - Frontend Engineer
- Gabriel Aldous - Backend Engineer
- Sai Neha Ratakonda - Backend Engineer

## Project Description

Numerous communication platforms currently exist to allow groups to communicate effectively and rapidly. Being that they often rely upon centralized servers to store and transmit user messages, they require great funding. Many of these platforms decide to gain this funding by selling user data/analytics. Our solution ***Sector*** proposes the implementation of p2p technologies to provide users with an end-to-end encrypted solution in which no user data is stored anywhere except on user devices. With this solution, users will be able to create group chats in which all participating members contain a complete record of all messages.

## About

Sector is a cross-platform desktop application built with Go and React that tackles the privacy concerns of modern messaging platforms. By leveraging peer-to-peer technologies and end-to-end encryption, Sector ensures that your conversations remain private and are stored only on the devices of the participants. The application provides a user-friendly interface for creating group chats, adding channels within those groups, and communicating securely with other users.

## Key Features

End-to-end encryption: All messages are encrypted, ensuring only the intended recipients can read them

Decentralized architecture: No central server storing your messages or metadata

Group chat functionality: Create and manage group conversations with multiple users

Channel organization: Structure conversations by topic within groups

Server handling: Create new servers or join existing servers.

Cross-platform support: Available as a desktop application for Windows, macOS, and Linux

User profiles: Customize your profile with username, status, and profile picture

Message pinning: Pin important messages for easy reference

## Development Setup

Utilize the `generate.sh` script to generate the frontend and backend data models/route representations.


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
## Techonology Stack

Backend: Go (Golang) for server-side development

Frontend: React with Material UI for building the user interface

Database: IPFS (InterPlanetary File System) with OrbitDB for distributed data storage

Framework: Wails for cross-platform desktop application development

Testing: Go testing framework, Vitest for React components, Cypress for E2E testing

API Documentation: Swagger UI

Version Control: Git/GitHub

## Installation and Setup

# Frontend Setup

Install Node.js and npm:

For Linux/Mac

curl -fsSL https://deb.nodesource.com/setup_16.x | sudo -E bash -

sudo apt-get install -y nodejs

For Windows: Download from https://nodejs.org/ and run installer

Navigate to the frontend directory:

cd frontend

Install dependencies:

npm install

Install Vite for development and building:

npm install -D vite @vitejs/plugin-react

Install React and related dependencies:

Install other dependencies:

npm install

npm install react react-dom

Install Vitest for unit testing:

npm install -D vitest jsdom @testing-library/react @testing-library/jes

Install Cypress for end-to-end and component testing:

npm install -D cypress @cypress/react

Install Material UI components:

npm install @mui/material @emotion/react @emotion/styled @mui/icons-materi

# Frontend Development

Navigate to the frontend directory:

cd frontend

Start the development server:

npm run dev

Open your browser and navigate to:

http://localhost:5173/

### Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

### Building

To build a redistributable, production mode package, use `wails build`.

## Notes

- The http server listens on port `3000`.
