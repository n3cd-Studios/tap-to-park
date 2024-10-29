FROM node:22 AS build

WORKDIR /app

COPY package*.json ./
COPY postcss.config.js ./
COPY svelte.config.js ./
COPY tailwind.config.js ./
COPY tsconfig.json ./
COPY vite.build.config.ts ./vite.config.ts
COPY ./src ./src
COPY ./static ./static

RUN npm ci
RUN npx vite build

RUN rm -rf src/ static/

EXPOSE 3000

CMD [ "npx", "vite", "preview", "--port", "3000" ]
