# DESCRIPTION: an authentication service using Go-chi and Jwt.

- THis application is build in Go, Go-chi and jwt + Reactjs for the client side.
- Jwt tokens are being stored in httpOnly cookies after being encrypted with AES-GCM.
- When the access token is expired, it gets automatically refereshed by the middleware if the refersh token is still not expired.

# INSTRUCTIONS:

- To run the application: make all
- To run only the frontend: make frontend
- To run only the backend: make backend
