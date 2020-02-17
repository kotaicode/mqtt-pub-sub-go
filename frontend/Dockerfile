FROM node:12.16-alpine as builder

ARG NODE_ENV=production
ENV NODE_ENV ${NODE_ENV}

ARG REACT_APP_BACKEND_URL=http://localhost:8080
ENV REACT_APP_BACKEND_URL ${REACT_APP_BACKEND_URL}

RUN mkdir -p /app
WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY . .
RUN yarn build

FROM nginx:1.17.8-alpine

#add curl for healthcheck
RUN apk update && apk add --no-cache curl

COPY --from=builder /app/build/ /usr/share/nginx/html
EXPOSE 80
