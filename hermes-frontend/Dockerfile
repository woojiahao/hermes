FROM node:19-alpine

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY package.json ./
COPY yarn.lock ./
RUN yarn install --silent
RUN yarn global add react-scripts --silent
RUN yarn global add serve

COPY . ./

CMD ["serve", "-s" "build"]
