# Sprint 3

## Backend

### Integration of Wails Framework

#### Overview
In Sprint 3, Integrating the [Wails framework](https://github.com/wailsapp/wails) has transformed Sector into a cross-platform desktop application, delivering a native experience on Windows, macOS, and Linux. By embedding the React frontend within the Go backend, Wails enables the creation of a single executable file, simplifying distribution and installation. This integration facilitates direct function calls between the frontend and backend, eliminating the need for HTTP requests and enhancing communication efficiency. Additionally, packaging the application as a unified executable streamlines deployment, reduces compatibility issues, and improves user accessibility

Key implementation steps included:

- Restructuring the application to align with Wails' architecture.
- Configuring the embedding of frontend assets within the Go binary.
- Ensuring seamless communication between frontend and backend components.
- Updating the startup sequence to initialize both backend services and the frontend interface.
- Implementing robust context management to handle the application's lifecycle effectively.

This integration enhances Sector's functionality and user experience by leveraging Wails' capabilities to create a cohesive desktop application.

### Unit Testing Improvements

During Sprint 3, we focused on enhancing the backend's unit testing framework to ensure code reliability and maintainability. Our efforts encompassed:

1. **Expanded Test Coverage**: We developed comprehensive tests for all API endpoints, emphasizing Channel and Message operations to ensure thorough validation.

2. **Refined Test Fixtures**: By improving setup and teardown procedures, we established consistent and isolated test environments, reducing dependencies and potential conflicts.

3. **Edge Case Validation**: We introduced tests targeting boundary conditions and error scenarios, verifying the API's robustness under various circumstances.

4. **Enhanced Mocking Techniques**: Utilizing tools like [Testify](https://github.com/stretchr/testify), we improved our mock database implementations, closely simulating production environments for more accurate testing.

5. **Integration Testing**: We conducted tests to validate interactions between system components, ensuring that operations such as message deletion appropriately update related channel data.

These improvements bolster our confidence in the backend's stability and performance, facilitating ongoing development and scalability.

## Backend Unit Tests Overview

In Sprint 3, we enhanced and expanded unit tests across various modules:

### Account Management

- Creation, retrieval, updating, and deletion of user accounts.
- Searching accounts by criteria such as ID, creation date, and username.

### Group Management

- Operations for creating, retrieving, updating, and deleting groups.
- Managing group memberships, including adding and removing members.
- Searching groups based on ID, creation date, name, and membership.

### Channel Management

- Handling creation, retrieval, updating, and deletion of channels within groups.
- Searching channels by ID, creation date, name, and associated group.

### Message Management

- Managing creation, retrieval, updating, and deletion of messages within channels.
- Searching messages using criteria like ID, creation date, author, channel, pinned status, and content.

These unit tests, implemented with Go's [testing package](https://pkg.go.dev/testing) and [Testify](https://github.com/stretchr/testify), are located in the `backend/tests/api/v1/` directory. To execute the tests, navigate to the `backend` directory and run:

```bash
go test -p 1 ./...

The -p 1 flag ensures sequential test execution, preventing conflicts in the test database.

By fortifying our testing framework, we aim to maintain high code quality and reliability as Sector continues to evolve.
