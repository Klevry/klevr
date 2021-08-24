## Build stage section from here
FROM node as builder

RUN mkdir /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

#ENV REACT_APP_API_URL="http://klevr-manager:8090"
## This Variable will be replace when the runtime by this -> Entrypoint:  sed -i  's#%%KLEVR_API_SERVER_IP_MARKUP%%#192.168.0.1:8090#g' /app/build/static/js/*.js
ENV REACT_APP_API_URL="http://%%KLEVR_API_SERVER_IP_MARKUP%%" 

COPY ./console/package.json /usr/src/app/package.json
RUN npm install --silent

COPY ./console/ /usr/src/app/
RUN npm run build


## Runtime stage section from here
FROM nginx:latest

RUN mkdir /app
WORKDIR /app
RUN mkdir ./build
COPY --from=builder /usr/src/app/build ./build
COPY ./Dockerfiles/console/entrypoint.sh /app/entrypoint.sh

## Setup the Nginx
RUN rm /etc/nginx/conf.d/default.conf
COPY ./Dockerfiles/console/nginx.conf /etc/nginx/conf.d/

EXPOSE 80
ENTRYPOINT /app/entrypoint.sh
