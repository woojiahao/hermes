FROM node:19-alpine

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY package.json ./
COPY yarn.lock ./

RUN yarn install --registry=https://registry.npmjs.org --network-timeout 1000000
RUN yarn global add react-scripts
RUN yarn global add serve

COPY . ./

RUN yarn build

CMD ["serve", "-s", "build/"]
