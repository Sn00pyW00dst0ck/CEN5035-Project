[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![project_license][license-shield]][license-url]

[contributors-shield]: https://img.shields.io/github/contributors/Sn00pyW00dst0ck/CEN5035-Project.svg?style=for-the-badge
[contributors-url]: https://github.com/Sn00pyW00dst0ck/CEN5035-Project/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/Sn00pyW00dst0ck/CEN5035-Project.svg?style=for-the-badge
[forks-url]: https://github.com/Sn00pyW00dst0ck/CEN5035-Project/network/members
[stars-shield]: https://img.shields.io/github/stars/Sn00pyW00dst0ck/CEN5035-Project.svg?style=for-the-badge
[stars-url]: https://github.com/Sn00pyW00dst0ck/CEN5035-Project/stargazers
[issues-shield]: https://img.shields.io/github/issues/Sn00pyW00dst0ck/CEN5035-Project.svg?style=for-the-badge
[issues-url]: https://github.com/Sn00pyW00dst0ck/CEN5035-Project/issues
[license-shield]: https://img.shields.io/github/license/Sn00pyW00dst0ck/CEN5035-Project.svg?style=for-the-badge
[license-url]: https://github.com/Sn00pyW00dst0ck/CEN5035-Project/blob/master/LICENSE.txt

# Sector

Numerous communication platforms currently exist to allow groups to communicate effectively and rapidly. Being that they often rely upon centralized servers to store and transmit user messages, they require great funding. Many of these platforms decide to gain this funding by selling user data/analytics. Our solution ***Sector*** proposes the implementation of p2p technologies to provide users with a solution in which no user data is stored anywhere except on user devices. With this solution, users will be able to create group chats in which all members contain a complete record of all messages.

## Key Features

* Decentralized architecture: No central server storing your messages or metadata.
* Group chat functionality: Create and manage group conversations with multiple users.
* Channel organization: Structure conversations by topic within groups.
* Server handling: Create new servers or join existing servers.
* Cross-platform support: Available as a desktop application for Windows, macOS, and Linux.
* User profiles: Customize your profile with username and profile picture.
* Message pinning: Pin important messages for easy reference.

## Usage

Clone this repository, then to build a redistributable, production mode package, utilize `./generate.sh` then `wails build` within your terminal. You will also need to ensure that Kubo version 0.27.0 is installed on your system!

Then, you can run the executable file/bundle which was generated and the application should open in a desktop window.

## Technology Stack

<div align="center">
<img src="https://github-readme-tech-stack.vercel.app/api/cards?title=Sector+Tech+Stack&align=center&titleAlign=center&lineCount=2&line1=react%2Creact%2C61DAFB%3Bgo%2Cgo%2C00ADD8%3Bvite%2Cvite%2C646CFF%3Bipfs%2Cipfs%2C65C2CB%3B&line2=cypress%2Ccypress%2C69D3A7%3Bgithub%2Cgithub%2C181717%3Bgit%2Cgit%2CF05032%3Bvitest%2Cvitest%2C6E9F18%3B" alt="Sector Tech Stack" />
</div>

## Development Setup

1.  Follow the `kubo` installation instructions for your operating system. Kubo version 0.27.0 must be used!

    If you haven't used IPFS so far, initialize the IPFS repository     using the following command:
    ```
    ipfs init
    ```

    If you had used IPFS an already have an IPFS repository in  place, either (re)move it from ~/.ipfs or make sure to export    IPFS_PATH before running the ipfs init command, e.g.:
    ```
    export IPFS_PATH=~/.ipfs-sector
    ipfs init
    ```

2. Follow the [wails](https://wails.io/) installation instructions for your operating system! 

> [!Warning]
> Ensure that Kubo version 0.27.0 is utilized, otherwise you may need to install migration tooling to run the databse.

### Frontend Setup

Install Node.js and npm from [the website](https://nodejs.org/en).

Navigate to the frontend directory:
```
cd frontend
```

Install the following dependencies for development:
```
npm install
npm install react react-dom
npm install @mui/material @emotion/react @emotion/styled @mui/icons-material
npm install -D vite @vitejs/plugin-react
npm install -D vitest jsdom @testing-library/react @testing-library/jes
npm install -D cypress @cypress/react
```

### Frontend Development

Navigate to the frontend directory:
```
cd frontend
```
Start the development server:
```
npm run dev
```
Open your browser and navigate to [http://localhost:5173/
](http://localhost:5173/).

For running tests during development:
```
npm run test
```

### Backend Setup

Ensure Kubo v0.27.0 is installed and that Go version 1.23.4 or later is installed on your system. 

Install the dependencies listed within the `go.mod` file and then run the project.

### Backend Development

Utilize `wails dev` to run the project. The http server listens on port `3000`. The swagger ui page may be reached from [http://localhost:3000/v1/swagger-ui/](http://localhost:3000/v1/swagger-ui/).

Unit testing can be performed using the following command from the root of the project:
```
go test -v ./...
```

### Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

### Building

To build a redistributable, production mode package, use `wails build`.

## Key Milestones

### Sprint 1

#### Frontend
- Designed wireframes for the main and login pages.
- Chose MUI component library for consistent and scalable design.
- Developed a modular React template with reusable components (e.g., user badges, server selection window, message window).
- Established a responsive desktop layout that adapts across different screen sizes.

#### Backend
- Evaluated and selected technology stack: Go and OrbitDB.
- Chose oapi-codegen for generating REST-like architecture.
- Set up the development process using openapi-typescript for frontend model generation.
- Designed the initial database schema for common objects.

---

### Sprint 2

#### Frontend
- Implemented dynamic server list with search functionality.
- Created a flexible user badge component.
- Added active server top bar with server badge and search functionality.
- Developed toggle functionality for member list visibility.
- Resolved CSS wrapping issues for various window aspect ratios.
- Added unit tests for key components using Cypress.

#### Backend
- Developed core API endpoints for account, group, channel, and message management.
- Implemented CRUD operations for major entities.
- Created comprehensive unit tests for all backend modules.
- Generated auto-documentation with Swagger UI.
- Encountered blockers with database mocking and complex data relationships.

---

### Sprint 3

#### Frontend
- Integrated backend endpoints for login, group fetching, and channel creation.
- Implemented the main user interface with server selection and channel display.
- Added user profile management functionality.
- Developed the login page with username/password forms.
- Created a complete test suite with unit tests and Cypress E2E tests.

#### Backend
- Integrated the Wails framework for cross-platform desktop support.
- Restructured the application to align with Wails architecture.
- Enhanced unit test coverage, including advanced search functionality.
- Improved lifecycle management for the desktop application.
- Maintained compatibility with the web interface alongside the desktop version.

---

### Sprint 4

#### Frontend
- Added a registration page for new user creation.
- Retained all functionality from Sprint 3.
- Continued integration with backend endpoints.
- Enhanced user profile management capabilities.

#### Backend
- Implemented JWT authentication for secure token management.
- Developed challenge-response authentication using public key cryptography.
- Optimized database queries for improved performance.
- Improved API error handling and response standardization.
- Created a comprehensive test suite for authentication, account, group, channel, and message management.
- Updated API documentation with detailed endpoint descriptions.

## Project Flow

```mermaid
graph TD
    %% Main Application Flow
    A[User] -->|Opens Application| B[Login Screen]
    B -->|Authentication| C[Main Interface]

    %% Authentication Flow
    subgraph Authentication
        style Authentication fill:#f8bbd0,stroke:#333,stroke-width:2px
        B1[Enter Username] -->|Request Challenge| B2[Backend: Generate Challenge]
        B2 -->|Return Challenge| B3[Client: Sign Challenge with Private Key]
        B3 -->|Submit Signature| B4[Backend: Verify Signature]
        B4 -->|Valid| B5[Generate JWT Token]
        B5 -->|Return Token| B6[Store Token for Requests]
        B4 -->|Invalid| B7[Show Error]
    end

    %% Main Interface Components
    subgraph MainInterface
        style MainInterface fill:#e0f7fa,stroke:#333,stroke-width:2px
        C1[Server List] -->|Select Server| C2[Channel List]
        C2 -->|Select Channel| C3[Message Display]
        C3 -->|Send Message| C4[Message Encryption]
        C4 -->|Store & Sync| C5[IPFS/OrbitDB]
        C1 -->|Create New Server| C6[Server Creation]
        C2 -->|Create New Channel| C7[Channel Creation]
        C3 -->|Receive Message| C8[Message Decryption]
    end

    %% Backend Components
    subgraph Backend
        style Backend fill:#fff3e0,stroke:#333,stroke-width:2px
        D1[HTTP API Server] -->|Process Requests| D2[Authentication Middleware]
        D2 -->|Validate JWT| D3[Request Handler]
        D3 -->|Database Operations| D4[OrbitDB Interface]
        D4 -->|Distributed Storage| D5[IPFS Node]
        D5 -->|P2P Communication| D6[Network Interface]
        D6 -->|Connect to Peers| D7[Peer Discovery]
    end

    %% Data Flow
    subgraph DataFlow
        style DataFlow fill:#d1c4e9,stroke:#333,stroke-width:2px
        E1[User Input] -->|React Components| E2[Frontend State]
        E2 -->|API Requests| E3[Backend API]
        E3 -->|Process| E4[Database Operations]
        E4 -->|Store| E5[Local Data]
        E5 -->|Sync| E6[Peer Data]
        E6 -->|Update| E7[UI Components]
    end

    %% Encryption Flow
    subgraph Encryption
        style Encryption fill:#ffecb3,stroke:#333,stroke-width:2px
        F1[Plaintext Message] -->|RSA Encryption| F2[Encrypted Message]
        F2 -->|Transmit| F3[Receive Encrypted Message]
        F3 -->|RSA Decryption| F4[Plaintext Message]
        F5[Generate Key Pair] -->|Store Private Key| F6[Local Storage]
        F5 -->|Share Public Key| F7[Public Key Database]
    end

    %% Connection Between Major Components
    B -->|Success| C
    C1 --- E2
    C2 --- E2
    C3 --- E2
    E3 --- D1
    E5 --- D4
    E6 --- D6
    C4 --- F2
    C8 --- F4
