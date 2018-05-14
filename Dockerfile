FROM node:alpine

RUN mkdir /db
VOLUME /db

COPY . /usr/src/app
WORKDIR /usr/src/app

RUN yarn

CMD yarn start