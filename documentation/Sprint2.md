# Sprint 2

## Video

https://www.youtube.com/watch?v=ji-p_NHyHpk

## General Notes

This sprint the backend ran into multiple complicated blockers which slowed development significantly. Mocking the backend database and server was much more involved than first predicted due to needing to mock an underlying IPFS instance, while at the same time development of the backend api routes became much slower due to intricate connections between data that have to be maintained when deleting or updating various database entries. 

Due to this, the integration of front and backend is much less developed than we would like, but the frontend and backend are able to communicate and have been since sprint 1. 

Future work involves refactoring the backend slightly to make further development easier, adding authentication, and encryption. 

## Frontend

## Frontend Unit Tests

List of Frontend Unit Tests for Sprint 2 :

- **User Badge:** Examines creation and rendering of user badge components.
- **Server List:** Examines creation, rendering, and filtering of complex server list component and integration with search filtering.

Our unit tests are constructing using the Cypress testing framework and are each located in the directory of their respective components.

## Frontend Work Completed

- **User Badge:** Designed and implemented flexible user badge to be used throughout the interface. Now interactable, with a placeholder alert which can be replaced later with a more-detailed user information dropdown. 
- **Search Component:** Standardized search component developed to be used throughout the interface. Allows for text input with a modular placeholder text and label.
- **Dynamic Server List:** Server list now upgraded to fetch from an array, which can be pulled from the backend. Also implements search component to allow for dynamic filtering of server results in real time.
- **Active Server Top Bar:** Added a menu bar at the top of the active server which includes a server badge for the active server, search component for message filtering, and a toggle button for a member list.
- **Dynamic Member List:** Added toggle functionality to the member list, allowing a user to toggle its visibility.
## Backend Unit Tests
- **Interface Wrapping Issues Addressed:** Patched css issues relating to component movement outside of expected bounds. Previously, some components would render improperly in the interface, especially if using abnormal window aspect ratios.

List of Backend Unit tests for Sprint 2 : 

- **Account Management:** Tests for creating, updating, retrieving, and deleting user accounts.
- **Group Management:** Tests for creating, updating, retrieving, and deleting groups.
- **Channel Management:** Tests for creating, updating, retrieving, and deleting channels.
- **Message Management:** Tests for creating, updating, retrieving, and deleting messages.

These unit tests are implemented using the Go testing framework and are located in the `backend/tests/api/v1/` directory of the repository. To run the tests, navigate to the `backend` directory and execute the following command:
```
go test -p 1 ./...
```


# Backend API Documentation

The following API documentation is **auto-generated** using **Swagger UI** for this project, which is hosted by the server.
A PDF printout of the Swagger UI has been inserted into the repository. Please view it here: [Swagger UI.pdf](Swagger%20UI.pdf)
