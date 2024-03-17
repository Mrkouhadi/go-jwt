# DESCRIPTION:

- I am using Go-chi + jwt + Redis.
- I am storing jwt tokens in httpOnly cookies after encrypting them with AES-GCM.
- I am using Redis to cache the credentials for the purpose of validation in middleware.go
