FROM node:22 AS build

WORKDIR /app

COPY package*.json ./
COPY postcss.config.js ./
COPY svelte.config.js ./
COPY tailwind.config.js ./
COPY tsconfig.json ./
COPY vite.config.ts ./
COPY ./src ./src
COPY ./static ./static

RUN npm install
RUN npm run build
RUN npm prune --production

FROM node:22 as run

WORKDIR /app
COPY --from=build /.svelte-kit/build ./build
COPY --from=build /app/package.json ./package.json

FROM 