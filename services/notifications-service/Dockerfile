# Use an official Node.js runtime as a parent image.
FROM node:18-alpine

# Set the working directory in the container.
WORKDIR /usr/src/app

# Copy package.json and package-lock.json (if available)
COPY package*.json ./

# Install app dependencies using 'npm ci' for reproducible builds
RUN npm ci --only=production

# Bundle app source
COPY . .

# This service is a background worker and does not expose any ports.

# Define the command to run your app
CMD [ "node", "main.js" ]
