FROM node:alpine

COPY . /usr/src/app
WORKDIR /usr/src/app

RUN yarn

CMD yarn start