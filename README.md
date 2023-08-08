### DESCRIPTION

- A guideline of how to use JWT with golang using this external [pkg](https://github.com/golang-jwt/jwt)


### a little explanaition
- when a user successfully logs in our app, he/she would receive a JWT Token that looks like this: 
`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTE1MDgwOTQsInVzZXJJRCI6ImhhcmRjb2RlZC1JRCJ9.9dOEZP0nR793uJCrPKWMLm0PAt1ilfWkjHwy7ol3fig`
- That Toke is divided into 3 parts by a `.` : 1)The header, 2) The payload, 3) The signature.
- NOTE: The header and payload are not encrypted; they are simply base64 encoded. This means that anyone can decode them using a base64 decoder.
