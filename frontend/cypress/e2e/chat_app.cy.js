// cypress/e2e/chatApp.cy.js
describe('Chat Application', () => {
  beforeEach(() => {
    // Create a custom Cypress command for login
    Cypress.Commands.add('login', () => {
      // Try to detect if we're already logged in
      cy.get('body').then(($body) => {
        // If we don't see the login button, we're already logged in
        if (!$body.find('button:contains("Sign In")').length) {
          return;
        }
        // Otherwise login
        cy.get('input').first().type('testuser')
        cy.get('input[type="password"]').type('password123')
        cy.get('button').contains('Sign In').click()
        // Verify login succeeded
        cy.contains('Login').should('not.exist')
      })
    })
    // Visit the app before each test
    cy.visit('http://localhost:5173')
  })
  
  it('should display login screen on initial load', () => {
    // Verify login form is displayed
    cy.contains('Login').should('be.visible')
    cy.get('input[type="password"]').should('be.visible')
    cy.get('button').contains('Sign In').should('be.visible')
  })
  
  it('should require both username and password to login', () => {
    // Try to login with empty fields
    cy.get('button').contains('Sign In').click()
    cy.contains('Login').should('be.visible') // Should still be on login page
    
    // Try with only username
    cy.get('input').first().type('testuser')
    cy.get('button').contains('Sign In').click()
    cy.contains('Login').should('be.visible') // Should still be on login page
    
    // Try with only password
    cy.get('input').first().clear()
    cy.get('input[type="password"]').type('password123')
    cy.get('button').contains('Sign In').click()
    cy.contains('Login').should('be.visible') // Should still be on login page
  })
  
  it('should log in with valid credentials and show the main screen', () => {
    // Enter any non-empty credentials and login
    cy.get('input').first().type('testuser')
    cy.get('input[type="password"]').type('password123')
    cy.get('button').contains('Sign In').click()
    
    // Verify we're on the main screen - we should no longer see the Login header
    cy.contains('Login').should('not.exist')
    // And we should see server badges
    cy.get('button').contains('test1').should('exist')
    cy.get('button').contains('test2').should('exist')
  })
  
  context('After logging in', () => {
    beforeEach(() => {
      // Use our custom login command
      cy.login()
    })
    
    it('should display server list with multiple servers', () => {
      // Check that multiple servers are displayed
      cy.contains('test1').should('be.visible')
      cy.contains('test2').should('be.visible')
    })
    
    
    it('should select a server and display its channels', () => {
      // Click on a server
      cy.contains('test1').click()
      
      // Verify the server is selected and channels are displayed
      cy.contains('test1 Channels').should('be.visible')
      cy.contains('General').should('be.visible')
      cy.contains('Gaming').should('be.visible')
      cy.contains('Music').should('be.visible')
    })
    
    it('should allow adding a new channel to a server', () => {
      // Select a server
      cy.contains('test1').click()
      
      // Click the add channel button
      cy.contains('+ Add Channel').click()
      
      // Enter a new channel name
      cy.get('input[placeholder="Channel name"]').type('New Test Channel')
      
      // Submit the form
      cy.contains('Add').click()
      
      // Verify the new channel appears in the list
      cy.contains('New Test Channel').should('be.visible')
    })
    
    it('should switch between channels and display correct messages', () => {
      // Select a server
      cy.contains('test1').click()
      
      // Click on the channel button
      cy.get('button').contains('General').first().click()
      
      // Select the Gaming channel from the menu
      cy.contains('Gaming').click()
      
      // Verify we see the appropriate message for Gaming
      cy.contains('Anyone up for a game?').should('be.visible')
      
      // Click on the channel button again
      cy.get('button').contains('Gaming').first().click()
      
      // Select the General channel
      cy.contains('General').click()
      
      // Verify we see the appropriate messages for General
      cy.contains('Welcome to the server!').should('be.visible')
      cy.contains('Hey everyone!').should('be.visible')
    })
    
    it('should send a new message in a channel', () => {
      // Select a server
      cy.contains('test1').click()
      
      // Click on the channel button and select General
      cy.get('button').contains('General').first().click()
      cy.contains('General').click()
      
      // Type a message and send it
      cy.get('input[placeholder="Text Message"]').type('Hello, this is a test message')
      cy.contains('Send').click()
      
      // Verify that the input is cleared after sending
      cy.get('input[placeholder="Text Message"]').should('have.value', '')
    })
    
    
    it('should open the user profile edit modal', () => {
      // Click the "Edit" button on the user profile
      cy.contains('Edit').click()
      
      // Verify the modal is open
      cy.contains('Edit Profile').should('be.visible')
      
      // Change the username
      cy.get('input').first().clear().type('UpdatedUsername')
      
      // Save changes
      cy.contains('Save').click()
      
      // Verify the modal is closed
      cy.contains('Edit Profile').should('not.exist')
      
      // Verify the username has been updated
      cy.contains('UpdatedUsername').should('be.visible')
    })
  })
})