####################################################################
# FRONTEND BUILD
#FROM alpine:3.10 AS ffmpeg-build

#RUN apk update
#RUN #apk add -u --no-cache ffmpeg

FROM node:14.16.1-alpine3.10 AS fronend-build

ARG REACT_APP_BACKEND_URL=''
# ARG REACT_APP_ASSET_URL=''

# YARN REQUIRES GIT BINARY
RUN apk update
RUN apk add git
RUN npm install -g typescript

# INSTALL PACKAGES
RUN mkdir -p /app
WORKDIR /app/frontend

# COPY SOURCE CODE
COPY ./frontend .

RUN yarn

# Overwrite files in node_modules
COPY ./overwrite/Canvas.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/
COPY ./overwrite/ThreeDEngine.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/threeDEngine/

# BUILD
RUN yarn build
####################################################################

# MAIN BUILD

# https://github.com/strapi/strapi-docker/blob/master/examples/custom/Dockerfile
FROM strapi/base as base

WORKDIR /my-path

COPY ./backend/package.json ./
COPY ./backend/yarn.lock ./

RUN yarn install

COPY ./backend .

ENV NODE_ENV production

RUN yarn build

COPY --from=fronend-build /app/frontend/build ./public

# copy ffmpeg from frontend build to backend
#COPY --from=ffmpeg-build /usr/bin/ffmpeg /usr/bin/ffmpeg

COPY ./scripts ./scripts

RUN chmod +x scripts/*.sh

EXPOSE 1337

CMD scripts/start.sh
