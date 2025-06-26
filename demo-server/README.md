# Demo Server

This server is the demo server used in [session 2](/docs/session2.md).

## Setup

### Github Secrets
It uses Github Oauth authentication flow. In order to have it work, you need to create a Github app in your github account [here](https://github.com/settings/developers).

The Github Oauth flow is described on Github [here](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps).

Once you have the CLIENT_ID and the CLIENT_SECRET for your application, create a new file `.env` and copy the content of the `env.example` into it.

### Server Secret
Lastly, you need to generate a new SERVER_KEY. On a linux system, you should be able to use `openssl rand -hex 32`. Paste the resulting value in the `.env` file.

## Run
Now run `go mod tidy`.

And then `go run .`

You should be able to access the home page in your browser: [http://localhost:8090](http://localhost:8090).