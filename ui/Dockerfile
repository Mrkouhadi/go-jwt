# Stage 1: Build the Next.js app
FROM node:latest as build

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code to the container
COPY . .

# Build the Next.js app
RUN npm run build

# Stage 2: Serve the Next.js app using a lightweight Node.js image
FROM node:latest

# Set the working directory in the container
WORKDIR /app

# Copy the built Next.js app from the previous stage
COPY --from=build /app/.next ./.next
COPY package.json package-lock.json ./

# Install only production dependencies
RUN npm install --only=production

# Expose the port that the Next.js app runs on
EXPOSE 3000

# Run the Next.js app
CMD ["npm", "start"]
