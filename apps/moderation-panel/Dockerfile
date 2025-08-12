# --- Build Stage ---
# Use a Node.js image to build the React application
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package.json and install dependencies
COPY package.json ./
RUN npm install

# Copy the rest of the application source code
COPY . .

# Build the application, creating a 'build' folder with static assets
RUN npm run build

# --- Serve Stage ---
# Use a lightweight Nginx image to serve the static files
FROM nginx:1.23-alpine

# Copy the custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy the built static assets from the builder stage
COPY --from=builder /app/build /usr/share/nginx/html

# Expose port 80 for the Nginx server
EXPOSE 80

# The default Nginx command will start the server
CMD ["nginx", "-g", "daemon off;"]
