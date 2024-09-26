# tap-to-park

A parking application created with Go and Svelte.

## Installation

First you need to install [Docker](https://www.docker.com/) to get a database setup.

Once Docker is installed, you can install a Postgres docker container: `docker run --name ttp-db -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres`

Now, copy the `.env.defaults` file to a `.env` file and fill out the connection string.

To install all of the Node packages, run `npm install`

## Getting Started

Run the container with `docker start ttp-db`

After that, you can:
- Run production: `npm start`
- Bring up the frontend: `npm run frontend`
- Bring up the backend: `npm run backend`

## Important Links

- [Google Drive](https://drive.google.com/drive/u/0/folders/1sLmxW9ZR5giioXCH832F8guUtJmecOBy)
- [Project Board](https://github.com/orgs/n3cd-Studios/projects/1/)
- [UI Wireframes](https://www.figma.com/design/bSKJXLx0NjJ8s6dcdf4cEO/Untitled?node-id=38-183&node-type=canvas&t=232RmGiC1CwHbwfI-0)
- [Color Palette](https://coolors.co/a0d8e3-f0f4f9-021427-76be37-d12335)