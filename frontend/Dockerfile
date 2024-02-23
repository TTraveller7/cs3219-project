FROM node:16.13.2-alpine
WORKDIR /usr/src/app
COPY . .
EXPOSE 3000
RUN npm install
RUN npm run build
RUN npm install -g serve
CMD ["serve", "-s", "build"]