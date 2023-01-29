# Short url generator

![Lines of code](https://img.shields.io/tokei/lines/github/gudimz/urlshortener)
![License](https://img.shields.io/github/license/gudimz/urlShortener)
![Go version](https://img.shields.io/github/go-mod/go-version/gudimz/urlshortener)
![Last commit](https://img.shields.io/github/last-commit/gudimz/urlshortener)

This is simple project for understanding how work with postgresql, docker in golang. It generates for an url a short url or uses 
the custom short url sent in the request.
## REST API
###Method POST
This method creates new short url.
- path: `api/v1/create/`
- body:
```yml
{
    "url":"https://example.com",
    "short_url":"example"
}
```
if you need to generate a short url - sent only the url
- return body:
```yml
{
  "message": "http://localhost:8080/example"
}
```
- code: `200`
###Method GET
This method redirects to the original url or returns all short url information.
- redirect:
  - path: `short_url`
  - code: `301`
- short url info:
    - path: `api/v1/short_url`
    - return body:
```yml
{
  "short_url": "example",
  "origin_url": "https://example.com",
  "visits": 1,
  "date_created": "2023-01-29T17:33:42.901111Z",
  "date_updated": "2023-01-29T19:22:15.431546Z"
}
```
- code: `200`
###Method DELETE
This method delete short url info from database.
- path: `api/v1/delete/short_url`
- code: `204`
##Usage
```shell
#build docker image
make build

#run the application
make run

# test application into container
make test

#stop the application
make stop
```