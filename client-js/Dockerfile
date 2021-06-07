FROM node:lts-alpine

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY . .

RUN chmod +x entrypoint.sh

ENTRYPOINT ./entrypoint.sh
