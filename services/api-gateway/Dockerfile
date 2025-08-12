# Use an official Node.js runtime as a parent image.
# Using a specific version is good practice for reproducibility.
FROM node:18-alpine

# Set the working directory in the container.
WORKDIR /usr/src/app

# Copy package.json and package-lock.json to leverage Docker cache.
# This step is only re-run if these files change.
COPY package*.json ./

# Install app dependencies.
# Using 'ci' is better for production builds as it uses the lockfile.
RUN npm ci --only=production

# Bundle app source.
COPY . .

# Your app binds to port 8080, so you need to expose it.
EXPOSE 8080

# Define the command to run your app.
CMD [ "node", "index.js" ]
