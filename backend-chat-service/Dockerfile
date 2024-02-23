FROM node:16.13.2-alpine
WORKDIR /usr/src/socket-io
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 5300
CMD ["npm", "start"]