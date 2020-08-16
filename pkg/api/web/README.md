# pkg/api/web

The `web` package contains assets and views for server-side rendering on the primary server. All assets are packaged with the binary using [go.rice](https://github.com/GeertJohan/go.rice), and the html is generated using standard go templating.

## views (templates)

We use [goview](https://github.com/foolin/goview) to extend standard go templating. This makes it easier to create layouts and partials (includes).

## signup (request access)

The sign-up process is used to request access to a sandpiper primary server. This is accomplished through a standard html form served from the root (/) endpoint.

## login (gain access)

The login screen is the front-end for checking credentials and returning a bearer token (jwt stored in a a cookie) for subsequent use. Access allows showing your subscriptions and downloading slices.

## download (retrieve grains)

todo: document this.
