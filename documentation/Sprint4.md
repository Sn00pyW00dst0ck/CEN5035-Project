# Sprint 4

## Frontend

### Integration with Backend Endpoints and Furthered Unit Testing

#### Overview
In Sprint 4, we focused primarily on continued integration with the backend. This includes logging in, fetching user groups, fetching group channels, and creating new group channels. We also addressed our lack of both unit and cypress tests to ensure that future changes will not inadvertently cause component failures. In addition, the following changes to the interface were made:

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

### Authentication Endpoints

#### GET /v1/api/challenge
Retrieves a challenge for user authentication. The challenge is later signed by the client's private key to prove identity.
- **Query Parameters**: `username` (required)
- **Response**: A randomly generated challenge string for the user to sign

#### POST /v1/api/login
Authenticates a user with a signed challenge.
- **Request Body**: `username` and `signature` (the signed challenge)
- **Response**: JWT token for subsequent authenticated requests

### Account Management

#### POST /v1/api/account/
Creates a new user account.
- **Request Body**: Account details including username, public key, and profile picture
- **Response**: Created account information

#### GET /v1/api/account/{id}
Retrieves account information by ID.
- **Path Parameters**: `id` (UUID)
- **Response**: Account details

#### PUT /v1/api/account/{id}
Updates an existing account.
- **Path Parameters**: `id` (UUID)
- **Request Body**: Updated account details
- **Response**: Updated account information

#### DELETE /v1/api/account/{id}
Deletes an account.
- **Path Parameters**: `id` (UUID)
- **Response**: No content on success

#### POST /v1/api/account/search
Searches for accounts based on filter criteria.
- **Request Body**: Filter parameters (username, creation date range, etc.)
- **Response**: List of matching accounts

### Group Management

#### POST /v1/api/group/
Creates a new group (server).
- **Request Body**: Group details including name, description, and member list
- **Response**: Created group information

#### GET /v1/api/group/{groupId}
Retrieves group information by ID.
- **Path Parameters**: `groupId` (UUID)
- **Response**: Group details

#### PUT /v1/api/group/{groupId}
Updates an existing group.
- **Path Parameters**: `groupId` (UUID)
- **Request Body**: Updated group details
- **Response**: Updated group information

#### DELETE /v1/api/group/{groupId}
Deletes a group.
- **Path Parameters**: `groupId` (UUID)
- **Response**: No content on success

#### POST /v1/api/group/search
Searches for groups based on filter criteria.
- **Request Body**: Filter parameters (name, creation date range, members, etc.)
- **Response**: List of matching groups

#### POST /v1/api/group/{groupId}/members/{memberId}
Adds a member to a group.
- **Path Parameters**: `groupId` (UUID), `memberId` (UUID)
- **Response**: Updated group information

#### DELETE /v1/api/group/{groupId}/members/{memberId}
Removes a member from a group.
- **Path Parameters**: `groupId` (UUID), `memberId` (UUID)
- **Response**: No content on success

### Channel Management

#### POST /v1/api/group/{groupId}/channel/
Creates a new channel within a group.
- **Path Parameters**: `groupId` (UUID)
- **Request Body**: Channel details including name and description
- **Response**: Created channel information

#### GET /v1/api/group/{groupId}/channel/{channelId}
Retrieves channel information by ID.
- **Path Parameters**: `groupId` (UUID), `channelId` (UUID)
- **Response**: Channel details

#### PUT /v1/api/group/{groupId}/channel/{channelId}
Updates an existing channel.
- **Path Parameters**: `groupId` (UUID), `channelId` (UUID)
- **Request Body**: Updated channel details
- **Response**: Updated channel information

#### DELETE /v1/api/group/{groupId}/channel/{channelId}
Deletes a channel.
- **Path Parameters**: `groupId` (UUID), `channelId` (UUID)
- **Response**: No content on success

#### POST /v1/api/channel/search
Searches for channels based on filter criteria.
- **Request Body**: Filter parameters (name, group, creation date range, etc.)
- **Response**: List of matching channels

### Message Management

#### POST /v1/api/group/{groupId}/channel/{channelId}/message
Creates a new message in a channel.
- **Path Parameters**: `groupId` (UUID), `channelId` (UUID)
- **Request Body**: Message details including body, author, and pinned status
- **Response**: Created message information

#### GET /v1/api/group/{groupId}/channel/{channelId}/message/{messageId}
Retrieves message information by ID.
- **Path Parameters**: `groupId` (UUID), `channelId` (UUID), `messageId` (UUID)
- **Response**: Message details

#### PUT /v1/api/group/{groupId}/channel/{channelId}/message/{messageId}
Updates an existing message.
- **Path Parameters**: `groupId` (UUID), `channelId` (UUID), `messageId` (UUID)
- **Request Body**: Updated message details
- **Response**: Updated message information

#### DELETE /v1/api/group/{groupId}/channel/{channelId}/message/{messageId}
Deletes a message.
- **Path Parameters**: `groupId` (UUID), `channelId` (UUID), `messageId` (UUID)
- **Response**: No content on success

#### POST /v1/api/message/search
Searches for messages based on filter criteria.
- **Request Body**: Filter parameters (body content, author, channel, pinned status, creation date range, etc.)
- **Response**: List of matching messages

### Miscellaneous

#### GET /v1/api/
Root endpoint returning a welcome message.
- **Response**: Simple welcome message

#### GET /v1/api/health
Health check endpoint to verify API is operational.
- **Response**: 200 OK if API is healthy

### Authentication
All endpoints except `/v1/api/challenge`, `/v1/api/login`, `/v1/api/health`, and `/v1/api/` require authentication via JWT token. The token must be included in the `Authorization` header as `Bearer {token}`.
