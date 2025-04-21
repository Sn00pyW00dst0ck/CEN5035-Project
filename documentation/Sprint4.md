# Sprint 4

## Frontend

### Integration with Backend Endpoints and Furthered Unit Testing

#### Overview
In Sprint 3, we focused primarily on continued integration with the backend. This includes logging in, fetching user groups, fetching group channels, and creating new group channels. We also addressed our lack of both unit and cypress tests to ensure that future changes will not inadvertently cause component failures. In addition, the following changes to the interface were made:

- Initial login screen rendering
- Successful login with credentials
- Server and Channel Interaction
- Displaying server list after authentication
- Server selection and channel display
- Adding new channels to existing servers
- Switching between channels
- Verifying correct message display for selected channels
- User profile management, can edit user status, Name, and description
- Servers and channels exist in the same state and messages are displayed appropriately for each selected channels

## Frontend Documentation

The frontend has two primary pages: the Login Page and the Main Page.

## Login
The login page consists of two forms: Username and Password. Currently, only the Username box must be filled to permit login as credentials have not been implemented. Upon sign in, a query is sent to the backend to fetch the user's data, before switching to the main page. We also intend to add a Register page in which the user can create a new account.

## Main
This is where the user will spend most of their time. It displays the user's groups, channels, messages, etc.

***Basic Anatomy***  
The main page is broken into two primary components: the Group List and the Active Server.

***Group List***  
Consists of the user's info, a server search bar, and a list of all the servers the client is in. When a server is selected, it becomes the 'active server'.

***Active Server***
Holds the messages of a selected channel in addition to a menu bar and member list.

## Unit Tests

**ActiveServer.test**
- Renders a placeholder when no server is selected
- Renders the main chat interface when a server is selected
- Passes the correct props to MenuBar
- Displays messages correctly
- Forwards channel selection from MenuBar
- Allows sending messages when a server and channel are selected
- Does not send empty messages
- Handles message sending via Enter key
- Does not call onChannelSelect when no channel is selected

**App.test**
- Renders without exceptions

**MainScreen.test**
- Renders without crashing
- Initializes with no selected server and channel
- Updates the selected server and channel when a server is selected
- Updates the selected channel when a channel is selected from the server list
- Updates the selected channel when a channel is selected from the dropdown menu
- Loads messages when a server and channel are selected
- Retains channel selection when changing between servers
- Handles state updates properly when changing channels

**MenuBar.test**
- Renders without crashing
- Renders ServerBadge with correct props
- Renders search component with correct label
- Renders menu button with icon
- Calls setVisible with opposite value when menu button is clicked
- Displays the selected channel name in the channel button
- Displays "select channel" when no channel is selected
- Opens the channel menu when the channel button is clicked
- Selects a channel when a menu item is clicked
- Marks the current selected channel as selected in the menu
- Handles the case when selected Server has no channels
- Handles the case when selectedServer is null

**ServerAndMembers.test**
- Renders without crashing
- Passes the correct props to ActiveServer
- Does not render Members by default (visible is false)
- Renders Members when visible state is true
- Forwards channel selection to parent component
- Handles null server and channel gracefully

**ServerBadge.test**
- Renders without exceptions
- Renders server name
- Renders default server icon
- Renders provided server icon

**ServerList.test**
- Renders without exceptions
- Server list renders all test servers
- Search filters servers returned
- Selecting a server calls onServerSelect and onChannelSelect with default channel
- Clicking on a channel calls onChannelSelect
- Adding a new channel works
- Visual indication is applied to the selected channel

**UserBadge.test**
- Renders without exceptions
- Renders server name
- Renders default server icon
- Renders provided server icon

## Cypress Tests

- Should display login screen on initial load
- Should require both username and password to login
- Should log in with valid credentials and show the main screen
- Should display server list with multiple servers
- Should select a server and display its channels
- Should allow adding a new channel to a server
- Should switch between channels and display correct messages
- Should send a new message in a channel
- Should open the user profile edit modal


# Backend
## **Backend Implementation**

## **Authentication System**
- **JWT Authentication**:  
  Implemented a robust **JSON Web Token (JWT)** authentication system with secure token generation, validation, and management.
  
- **Challenge-Response Authentication**:  
  Developed a secure challenge-response mechanism for user login using **public key cryptography**.
  
- **Session Management**:  
  Created comprehensive session handling with **token expiration** and **refresh capabilities**.

## **Encryption**
- **End-to-End Encryption Support**:  
  Implemented backend support for **end-to-end encrypted messaging**.
  
- **Key Management**:  
  Developed secure storage and retrieval mechanisms for user **public keys** used in message encryption.
  
- **Message Integrity**:  
  Added verification mechanisms to ensure **message integrity** during transmission.

## **Database Enhancements**
- **Optimized Queries**:  
  Improved **database query performance** for faster message and channel retrieval.
  
- **Data Consistency**:  
  Enhanced **data consistency checks** during create, update, and delete operations.
  
- **Relationship Management**:  
  Strengthened **relationship handling** between users, groups, channels, and messages.

## **API Improvements**
- **Standardized Responses**:  
  Implemented **consistent API response formats** across all endpoints.
  
- **Error Handling**:  
  Enhanced **error reporting** and handling for more informative client feedback.
  
- **Documentation**:  
  Updated **API documentation** with comprehensive endpoint descriptions and example usage.

---

# **Backend Unit Tests**

The backend implements extensive testing covering all major components to ensure robustness and functionality.

## **Authentication Tests**
- **JWT Token Generation**:  
  Tests proper generation of JWT tokens with correct claims and expiration.
  
- **Token Validation**:  
  Verifies validation of tokens, including rejection of expired or tampered tokens.
  
- **Challenge Generation**:  
  Tests secure challenge generation for login authentication.
  
- **Signature Verification**:  
  Validates verification of cryptographic signatures during login.
  
- **Authorization Middleware**:  
  Tests middleware that enforces authentication on protected routes.

## **Account Management Tests**
- **Account Creation**:  
  Verifies user account creation with proper validation.
  
- **Account Retrieval**:  
  Tests retrieval of user account information.
  
- **Account Updates**:  
  Validates updating user profile information.
  
- **Account Deletion**:  
  Tests proper account deletion with cascade effects.
  
- **Account Search**:  
  Verifies searching for accounts based on various criteria.

## **Group Management Tests**
- **Group Creation**:  
  Tests creation of new groups with proper validation.
  
- **Group Retrieval**:  
  Verifies retrieval of group information.
  
- **Group Updates**:  
  Tests updating group details.
  
- **Group Deletion**:  
  Validates proper group deletion with cascade effects on channels and messages.
  
- **Member Management**:  
  Tests adding and removing members from groups.
  
- **Group Search**:  
  Verifies searching for groups based on various criteria.

## **Channel Management Tests**
- **Channel Creation**:  
  Tests creation of channels within groups.
  
- **Channel Retrieval**:  
  Verifies retrieval of channel information.
  
- **Channel Updates**:  
  Tests updating channel details.
  
- **Channel Deletion**:  
  Validates proper channel deletion with cascade effects on messages.
  
- **Channel Search**:  
  Verifies searching for channels based on various criteria.

## **Message Management Tests**
- **Message Creation**:  
  Tests creation of messages within channels.
  
- **Message Retrieval**:  
  Verifies retrieval of messages.
  
- **Message Updates**:  
  Tests updating message content.
  
- **Message Deletion**:  
  Validates proper message deletion.
  
- **Message Search**:  
  Verifies searching for messages based on various criteria.
  
- **Message Encryption**:  
  Tests **encryption** and **decryption** of message content.

These unit tests are implemented using the Go testing framework and are located in the `backend/tests/api/v1/` directory of the repository. To run the tests, navigate to the `backend` directory and execute the following command:
```
go test -p 1 ./...
```

# Backend API Documentation

The following API documentation is **auto-generated** using **Swagger UI** for this project, which is hosted by the server.
A PDF printout of the Swagger UI has been inserted into the repository. Please view it here: [Swagger UI Documentation (PDF)](Swagger%20UI%20-%204.pdf)

## Overview of Backend API Documents



