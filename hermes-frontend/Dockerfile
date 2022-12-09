FROM node:19-alpine

WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json ./
COPY yarn.lock ./
RUN pwd
RUN cd /app
RUN yarn install --silent
RUN yarn global add react-scripts --silent
COPY . ./
EXPOSE 3000
CMD ["yarn", "start"]
