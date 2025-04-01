# Sprint 3

## Backend

### Integration of Wails Framework

#### Overview
In Sprint 3, Integrating the [Wails framework](https://github.com/wailsapp/wails) has transformed Sector into a cross-platform desktop application, delivering a native experience on Windows, macOS, and Linux. By embedding the React frontend within the Go backend, Wails enables the creation of a single executable file, simplifying distribution and installation. This integration facilitates direct function calls between the frontend and backend, eliminating the need for HTTP requests and enhancing communication efficiency. However, in order to maintain the web interface's functionality, the team has decided that both the web and desktop versions will interact with the backend via HTTP requests for simplicity of development. Additionally, packaging the application as a unified executable streamlines deployment, reduces compatibility issues, and improves user accessibility.

Key implementation steps included:

- Restructuring the application to align with Wails' architecture.
- Configuring the embedding of frontend assets within the Go binary.
- Ensuring seamless communication between frontend and backend components.
- Updating the startup sequence to initialize both backend services and the frontend interface.
- Implementing robust context management to handle the application's lifecycle effectively.

This integration enhances Sector's functionality and user experience by leveraging Wails' capabilities to create a cohesive desktop application.

## Backend Unit Tests Overview

In Sprint 3, the unit tests have been leveled up by scaling test coverage across multiple modules:

### Account Management

- Creation, retrieval, updating, and deletion of user accounts, ensuring full lifecycle support.
- Searching accounts by criteria such as ID, creation date, and username, improving query validation.

### Group Management

- Operations for creating, retrieving, updating, and deleting groups, covering all core functionalities.
- Managing group memberships, including adding and removing members, to verify membership consistency.
- Searching groups based on ID, creation date, name, and membership, enhancing group discoverability.

### Channel Management

- Handling creation, retrieval, updating, and deletion of channels within groups, ensuring robust channel operations.
- Searching channels by ID, creation date, name, and associated group, validating channel relationships.

### Message Management

- Managing creation, retrieval, updating, and deletion of messages within channels, supporting complete message handling.
- Searching messages using criteria like ID, creation date, author, channel, pinned status, and content, ensuring comprehensive message filtering.

These unit tests are implemented using the Go testing framework and are located in the `backend/tests/api/v1/` directory of the repository. To run the tests, navigate to the `backend` directory and execute the following command:
```
go test -p 1 ./...
```

# Backend API Documentation

The following API documentation is **auto-generated** using **Swagger UI** for this project, which is hosted by the server.
A PDF printout of the Swagger UI has been inserted into the repository. Please view it here: [Sprint 3 Swagger UI.pdf](Swagger%20UI.pdf)
