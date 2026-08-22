# Snippetbox

Web application which lets people paste and share snippets of text.
It is my walkthrough of Alex Edwards [book](https://lets-go.alexedwards.net/) "Let's Go".

Even though I followed the book last bits of work are mine. I have completed
guided exercises and added a few new features myself.

### Tools used:
 - Mostly Go's standard library (net/http, html/template)
 - MySQL
 - external Go libraries for purposes like session management (see go.mod file) 

### Guided exercises:
 - add an 'About' page to the application
 - add a debug mode
 - test the snippetCreate handler
 - add an 'Account' page to the application
 - redirect user appropriately after login
 - implement a 'Change Password' feature

### My additions:
* only logged users can see who is a creator of a particular snippet
* users can see snippets their added (one-to-many relationship)
* users can manage their snippets (edit, delete)

### Major concepts covered:
 - structuring web app
 - centralized logging and error handing
 - middleware
 - session management
 - authentication and authorization
 - form validation
 - testing

### Reflection:
It was my first practical exposure to the world of backend development.
I did it for learning purposes and I think the book served me well.
I happened to pick Go because I wanted to see a bigger picture without framework
abstracting new concepts away from me.
