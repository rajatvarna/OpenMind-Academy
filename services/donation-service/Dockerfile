# Use an official Node.js runtime as a parent image.
FROM node:18-alpine

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm ci --only=production

COPY . .

EXPOSE 3007

CMD [ "node", "index.js" ]
