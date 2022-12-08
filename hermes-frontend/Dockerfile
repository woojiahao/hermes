FROM node:19-alpine

WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json ./
COPY yarn.lock ./
RUN ls
RUN yarn install --silent
RUN yarn global add react-scripts --silent
COPY . ./
CMD ["yarn", "start"]
