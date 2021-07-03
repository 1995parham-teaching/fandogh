# Fandogh :chestnut:

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/1995parham/fandogh/ci?label=ci&logo=github&style=flat-square)
[![Codecov](https://img.shields.io/codecov/c/gh/1995parham/fandogh?logo=codecov&style=flat-square)](https://codecov.io/gh/1995parham/fandogh)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/1995parham/fandogh)](https://pkg.go.dev/github.com/1995parham/fandogh)

## APIs

register

```sh
curl 127.0.0.1:1378/register -X POST -d '{ "email": "parham.alvani@gmail.com", "name": "Parham Alvani", "password": "123456" }' -H 'Content-Type: application/json'
```

```json
{
  "Email": "parham.alvani@gmail.com",
  "Password": "123456",
  "Name": "Parham Alvani"
}
```

login:

```sh
curl 127.0.0.1:1378/login -X POST -d '{ "email": "parham.alvani@gmail.com", "password": "123456" }' -H 'Content-Type: application/json'
```

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyIiwiZXhwIjoxNjI1MzU4MDQ1LCJqdGkiOiI4NTYxYzA4NC1kYzAxLTQ0ZmEtODEyZS05ZjNhZDJlNDcxNTAiLCJpYXQiOjE2MjUzNTQ0NDUsImlzcyI6ImZhbmRvZ2giLCJuYmYiOjE2MjUzNTQ0NDUsInN1YiI6InBhcmhhbS5hbHZhbmlAZ21haWwuY29tIn0.hUiEGqQxCSTQOFDPBypKkdI85q7TxSGENY6IwA2QR7E",
  "Email": "parham.alvani@gmail.com",
  "Password": "123456",
  "Name": "Parham Alvani"
}
```
