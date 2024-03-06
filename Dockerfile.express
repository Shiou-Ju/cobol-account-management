FROM node:latest

# install gnu cobol
RUN apt-get update && apt-get install -y build-essential libgmp-dev libdb-dev libncurses5-dev wget

RUN wget https://sourceforge.net/projects/gnucobol/files/gnucobol/3.2/gnucobol-3.2.tar.xz \
    && tar -xf gnucobol-3.2.tar.xz \
    && cd gnucobol-3.2 \
    && ./configure && make && make install

ENV LD_LIBRARY_PATH=/usr/local/lib

RUN rm -rf gnucobol-3.2.tar.xz gnucobol-3.2

# start building app
WORKDIR /app

COPY ./package.json ./yarn.lock ./
COPY ./backend ./backend

RUN yarn global add typescript @vue/cli

# install backend packages
RUN yarn install

RUN yarn build:backend

# install frontend packages
COPY ./frontend ./frontend
RUN cd frontend && yarn install

# env to prevent vue config linting errors
ENV NODE_ENV=production
RUN cd frontend && yarn build

# copy frontend build to backend for serving
RUN mkdir backend/dist/public/
RUN cp -r frontend/dist/* backend/dist/public/

# copy cobol scripts to backend
RUN cp -r ./backend/cobol ./backend/dist/cobol

# run server
WORKDIR /app/backend

EXPOSE 3000

CMD ["node", "dist/server.js"]
