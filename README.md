# DESCRIPTION: an authentication service using Go-chi and Jwt.

- THis application is build in Go, Go-chi and jwt + Reactjs for the client side.
- Jwt tokens are being stored in httpOnly cookies after being encrypted with AES-GCM.
- When the access token is expired, it gets automatically refereshed by the middleware if the refersh token is still not expired.

# INSTRUCTIONS:

- add .env file in `server` directory. ( add SECRET_JWT_KEY ="a 32 chars. example: 2a6F9d$%l1@#3s&u7x9A1C3E5G7I9K1M985s7m9U" )
- To run the application: make all
- To run only the frontend: make frontend
- To run only the backend: make backend

# SCREENSHOTS:

Login page:
![Login page](./ui.png)
Profile page:
![Profile page](./ui-2.png)
