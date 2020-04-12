# GamesDB Mirror API

This is intended to be an easy-to-run solution for mirroring the data served by TheGamesDB.

## Getting started

```sh
make
./out/fetch # NOTE: must be run at least once before running the server
./out/server

# test that it works with
#     curl 'http://localhost:5000/Games/ByIds?ids=1'
```

## Running on heroku

One of the primary goals of this project is to be able to run on a free heroku dyno. To that end, things should be pretty drop-in, you just need to:

1. Create a heroku app
2. Run `make heroku-all APP={name-of-heroku-app}`

I like to then run `heroku logs -t --app={name-of-heroku-app}` to make sure the deploy was successful. You should now be able to hit `https://{name-of-heroku-app}.herokuapp.com/Games/ById?id=1` to test.

## Architecture

### Fetcher

Run with

```sh
make fetch
```

The fetcher does the work of pulling the latest GamesDB data dump and massaging/caching it into a form that the server expects. The fetcher is meant to be run often, as it has logic in it that prevents it from doing unnecessary work if no updates are available.

### Server

Run with

```sh
make serve
```

The server is a barebones API server that loads the work previously cached by the fetcher into memory and then performs queries against it in reponse to HTTP requests. You can provide a custom bind port by setting the `PORT` environment variable.
