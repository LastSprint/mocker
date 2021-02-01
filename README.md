<p align="center">
  <img src="logo.png">
</p>


![Actions](https://github.com/LastSprint/mocker/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/LastSprint/mocker/branch/master/graph/badge.svg)](https://codecov.io/gh/LastSprint/mocker)
[![codebeat badge](https://codebeat.co/badges/2c22d06b-0f69-44af-9b42-00c6cf0cc3e4)](https://codebeat.co/projects/github-com-lastsprint-mocker-master)

# `Mocker` — a Web Server Emulating a Real Backend

# Features

- URL-Query prams matching — selects a mock depending on query parameters given in a query and in a mock.
- JSON-Body prams matching — selects a mock depending on a JSON body in a query and in a mock.
- Caching Proxy — proxies client queries to an actual backend, records the result in a mock and returns it to the client.
- Individual mocks or all mocks except for a selected one can be disabled.
- Response can be delayed for a selected mock.
- Iterative responses: several mocks with a specified URL will be returned one at a time — if no parameters are matched.

**We’re not planning to support mock relation**. Mocks are just files: they are in no way connected to each other or altered by the Mocker itself.

Supporting relation would complicate the service without providing any advantages: if you write something using mocks, you hardly need data relation.

However, if you do need relation, use matching by query or body parameters.

# How it Works

- Mocks are written by users. A mock describes what the Mocker has to return in response to a query.
- `Mocker` reads data when launched or following a `GET /update_models query` (read on to learn how to make it automated)
- When the `Mocker` receives a query it finds a relevant mock and returns it to a client.

The way the mocker works is pretty simple, but it’s much more complex if you look under the hood (:

## Mocks

Mocks are Json files:
```
 {
    "isDisabled": bool,
    "isOnly": bool,
    "isExcludedFromIteration": bool,
    "url": string,
    "method": string,
    "statusCode": int,
    "responseDelay": int,
    "response": object,
    "request": object
    "responseHeaders": object
    "requestHeaders": {
        "key": "value",
         .....
     }
 }
```

This literally means the following:

If a query with `URL = url` and `Method = method is received`, return `response` with a code `statusCode`

### `url`

The following types can be used:

#### `/path/to/endpoint` 
A simple URL. In response to a query a service will compare strings one character at a time.

#### `/path/to/endpoint/{number}`

A URL with a path pattern. A mock with such a URL will react to any query compliant with the set template.

For example:

```
/path/to/endpoint/1 --> OK
/path/to/endpoint/item --> OK
/path/to/endpoint/1/2 --> FALSE
```

#### `/path/to/endpoint/data?param={value}`

A URL with a query pattern. A mock with such a URL will react to a query containing whatever parameters were set. 
Notice that a query with any of the parameters missing will not match a template.

**NOTE**:

A URL must start with a slash (`/`)

### `method`

Names of all HTTP methods should be written in UpperCase (i.e. write `GET` instead of `get`)

### `statusCode`

Any integer, preferably one of the known [HTTP codes](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)

### `response`

This field contains a `Json` to be returned in response to a query.

### `request`

This field contains data needed to search a specific mock. Here we can apply parameterization.

Say, we want to mock an ordering process:

```JSON
{
  "url": "/billing/create",
  "method": "POST",
  "request": {
    "shopId": "123",
    "paymentType": "card",
    "items": []
  },
  "statusCode": 200,
  "response": {
    "orderId": 123
  }
}
```
If an order is made from shop `123` and paid by a card, a query will return the mock given above.

But there’s a catch. Such a mock will only be matched to a query, if an `items` array in the query is empty. 
In other words, we’ll have to create a new mock for each cart (containing different items).

To avoid that we need to update the mock:

```JSON
{
  "url": "/billing/create",
  "method": "POST",
  "request": {
    "shopId": "123",
    "paymentType": "card",
    "items": "{items}"
  },
  "statusCode": 200,
  "response": {
    "orderId": 123
  }
}
```

Now the mock will be received regardless of the `items` value.

Now let’s say we want queries where payment does not equal `card` to return an error.

```JSON
{
  "url": "/billing/create",
  "method": "POST",
  "request": {
    "shopId": "123",
    "paymentType": "{ paymentType != card }",
    "items": "{items}"
  },
  "statusCode": 400,
  "response": {
    "msg": "Current paymentType is unsupported"
  }
}
```

#### Templates

- `{value}` — a template describing a value.
- `{value != | > | < | >= | <= $const$ }` — an expression template with all applicable operators separated by `|`

The operators apply to a limited set of types:

- `!=` for `String`, `Int`, `Dobule`
- `>`, `<`, `>=`, `<=` for `Int`, `Dobule`

What you should remember about expression templates is:

- If you stated a non-existent operation, your mock will match.
- If data type in `request` can not be used in this operator, your mock will not match.
- If `$const$` value does not comply with the data type given in `request`, your mock will not match.

You can write anything you want in the template, but we advise you to copy the name of your variable, because functionality of templates with operators is going to be extended as we go forward.

### `requestHeaders`

If a query contains the headers specified in this field, you mock will match.

A mock will only be considered matched if all match conditions (query, request, headers) are met.

### `responseHeaders`

Contains the list of `key-value` pairs, where `key` is a name of a header, and `value` is a value.

For example, if we want our mocker to return an `X-Example-Header` header with an `example_value` value, we’ll write:

```JSON
"responseHeaders": {
  "X-Example-Header": "example_value"
}
```

### `isDisabled`

This flag is used to switch a mock “off”. If `isDisabled == true`, a mock will not be taken into account.

If the value is `false` or `nil`, all works as usual.

### `isOnly`

This flag switches off all mocks except for the one selected. 
If `IsOnly == true` for  mock, it will be the only mock taken into account. The others will be considered “switched off”.

If a mock has `IsOnly == true` and `IsDisabled == true` at the same time, `IsDisabled` is ignored.

If several mocks have `IsOnly == true` at the same time, the first mock will always be the one you receive. 
No iteration is available (at least for the time being).

Please note that the iteration counter is not zeroed out. 
If the iterator shows the n-th file, it will still be showing it even after the IsOnly is switched off and on again.

### `responseDelay`

We need this field to delay a server response with a specific mock on purpose

In other words, all mocks with != 0 value in this field will be delayed by the time specified. The time should be given in seconds.

Default value: `0`

### `isExcludedFromIteration`

This field is used to exclude a mock from iterative responses.

Say, when you have a mock with a specific body of a query and you want it to be returned only if this exact body is matched.

Default state: `false`
 
## The Mocker is configured via environment variables:

The `Mocker` is configured via environment variables:

- `MOCKER_MOCKS_ROOT_DIR: string` — a pathway to the folder containing all mocks.
- `MOCKER_SERVER_PORT: integer` — a port on which the `Mocker` listens for connections.
- `MOCKER_LOG_PATH` — a pathway to the file where the `Mocker` writes logs. Logs are written in `JSON`.


## How to Install and Get Started

You can find the following files in the root of the repository:
- `docker-compose.yaml` contains all the necessary configurations and is ready to launch - `docker-compose up -d`.
- `Dockerfile` contains configurations needed to launch `Mocker`.
- `FSWatherDockerfile` —  a container listening for changes in the file system (in the folder where mocks are stored) and responding to changes with an automated query `GET /updateModels`.

When launched, compose can return an error like `you try to mount directory to file (or vice versa)`. ЕIf you got this error, create all the necessary files manually.

`.filebrowser_config.json`:
```JSON
{
  "port": 80,
  "baseURL": "",
  "address": "",
  "log": "stdout",
  "database": "/eddb/database.db",
  "root": "/srv"
}
```

And then simply add the files to .git/exclude once you get down to work.

## Roadmap

- Add support of `form-url` for `request` matching

## Contributing

I’d appreciate your bug reports, feature requests, and PRs!
