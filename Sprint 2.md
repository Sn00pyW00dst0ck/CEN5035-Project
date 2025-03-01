# Backend API Documentation

## Miscellaneous Endpoints

### Root Endpoint

- **Endpoint:** `GET /v1/api/`
- **Description:** Retrieves the root endpoint.
- **Parameters:** None
- **Responses:**
  - `200 OK`: Success

### Health Check

- **Endpoint:** `GET /v1/api/health`
- **Description:** Checks the health status of the API.
- **Parameters:** None
- **Responses:**
  - `200 OK`: API is healthy

## Account Management Endpoints

### Create or Update an Account

- **Endpoint:** `POST /v1/api/account/`
- **Description:** Creates a new account or updates an existing one.
- **Request Body (application/json):**
  - `id` (string, UUID): Unique identifier for the account.
  - `username` (string): Username of the account holder.
  - `profile_pic` (string): URL to the profile picture.
  - `created_at` (string, date-time): Account creation timestamp.
- **Responses:**
  - `200 OK`: Account creation or update successful.

### Search for Accounts

- **Endpoint:** `POST /v1/api/account/search`
- **Description:** Searches for accounts based on provided criteria.
- **Request Body (application/json):**
  - `ids` (array of UUIDs): List of account IDs to search for.
  - `username` (string): Username to search for.
  - `from` (string, date-time): Start date for the search range.
  - `until` (string, date-time): End date for the search range.
- **Responses:**
  - `200 OK`: Query completed successfully.

### Delete Account by ID

- **Endpoint:** `DELETE /v1/api/account/{id}`
- **Description:** Deletes an account by its unique ID.
- **Parameters:**
  - `id` (string, UUID, path): ID of the account to delete.
- **Responses:**
  - `204 No Content`: Account was deleted.
  - `500 Internal Server Error`: Error occurred during deletion.

### Get Account by ID

- **Endpoint:** `GET /v1/api/account/{id}`
- **Description:** Retrieves account details by its unique ID.
- **Parameters:**
  - `id` (string, UUID, path): ID of the account to retrieve.
- **Responses:**
  - `200 OK`: Account with specified ID retrieved successfully.

## Group Management Endpoints

### Create or Update a Group

- **Endpoint:** `POST /v1/api/group/`
- **Description:** Creates a new group or updates an existing one.
- **Request Body (application/json):**
  - `id` (string, UUID): Unique identifier for the group.
  - `name` (string): Name of the group.
  - `description` (string): Description of the group.
  - `members` (array of UUIDs): List of member account IDs.
  - `channels` (array of strings): List of channels within the group.
  - `created_at` (string, date-time): Group creation timestamp.
- **Responses:**
  - `200 OK`: Group creation or update successful.

### Search for Groups

- **Endpoint:** `POST /v1/api/group/search`
- - **Description:** Searches for groups based on provided criteria.
- **Request Body (application/json):**
  - `id` (array of UUIDs): List of group IDs to search for.
  - `name` (string): Name of the group to search for.
  - `from` (string, date-time): Start date for the search range.
  - `until` (string, date-time): End date for the search range.
- **Responses:**
  - `200 OK`: Query completed successfully.

### Delete Group by ID

- **Endpoint:** `DELETE /v1/api/group/{id}`
- **Description:** Deletes a group by its unique ID.
- **Parameters:**
  - `id` (string, UUID, path): ID of the group to delete.
- **Responses:**
  - `204 No Content`: Group was deleted.

### Get Group by ID

- **Endpoint:** `GET /v1/api/group/{id}`
- **Description:** Retrieves group details by its unique ID.
- **Parameters:**
  - `id` (string, UUID, path): ID of the group to retrieve.
- **Responses:**
  - `200 OK`: Group with specified ID retrieved successfully.

## Channel Management Endpoints

### Create or Update a Channel

- **Endpoint:** `POST /v1/api/channel/`
- **Description:** Creates a new channel or updates an existing one.
- **Request Body (application/json):**
  - `id` (string, UUID): Unique identifier for the channel.
  - `name` (string): Name of the channel.
  - `description` (string): Description of the channel.
  - `messages` (array of message objects): List of messages in the channel.
  - `pinned_messages` (array of message objects): List of pinned messages in the channel.
  - `created_at` (string, date-time): Channel creation timestamp.
- **Responses:**
  - `200 OK`: Channel creation or update successful.

### Search for Channels

- **Endpoint:** `POST /v1/api/channel/search`
- **Description:** Searches for channels based on provided criteria.
- **Request Body (application/json):**
  - `id` (array of UUIDs): List of channel IDs to search for.
  - `name` (string): Name of the channel to search for.
  - `from` (string, date-time): Start date for the search range.
  - `until` (string, date-time): End date for the search range.
- **Responses:**
  - `200 OK`: Query completed successfully.

### Delete Channel by ID

- **Endpoint:** `DELETE /v1/api/channel/{id}`
- **Description:** Deletes a channel by its unique ID.
- **Parameters:**
  - `id` (string, UUID, path): ID of the channel to delete.
- **Responses:**
  - `204 No Content`: Channel was deleted.

### Get Channel by ID

- **Endpoint:** `GET /v1/api/channel/{id}`
- **Description:** Retrieves channel details by its unique ID.
- **Parameters:**
  - `id` (string, UUID, path): ID of the channel to retrieve.
- **Responses:**
  - `200 OK`: Channel with specified ID retrieved successfully.

## Message Management Endpoints

### Create or Update a Message

- **Endpoint:** `POST /v1/api/message/`
- **Description:** Creates a new message or updates an existing one.
- **Request Body (application/json):**
  - `id` (string, UUID): Unique identifier for the message.
  - `author` (UUID): ID of the account that authored the message.
  - `body` (string): Content of the message.
  - `created_at` (string, date-time): Message creation timestamp.
- **Responses:**
 

 
