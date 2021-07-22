FROM node as builder

RUN mkdir /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

COPY package.json /usr/src/app/package.json
RUN npm install --silent

COPY . /usr/src/app/
RUN npm run build


FROM nginx:latest

RUN mkdir /app
WORKDIR /app
RUN mkdir ./build
COPY --from=builder /usr/src/app/build ./build
RUN rm /etc/nginx/conf.d/default.conf
COPY ./nginx.conf /etc/nginx/conf.d

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]