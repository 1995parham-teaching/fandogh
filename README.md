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

new home:

```sh
curl 127.0.0.1:1378/api/homes -X POST  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyIiwiZXhwIjoxNjI1Mzc5MDE3LCJqdGkiOiI4NDRhMzQ4Yy03OGVjLTRlNTctODJhZi03YjU3NTNmNjk5ZjciLCJpYXQiOjE2MjUzNzU0MTcsImlzcyI6ImZhbmRvZ2giLCJuYmYiOjE2MjUzNzU0MTcsInN1YiI6InBhcmhhbS5hbHZhbmlAZ21haWwuY29tIn0.EZUWQ-sLP1ClA0vtK6vZEcQ4qf3ZaBm9VpFV6smEwUc' -F 'title=sweet' -F 'location=italy' -F 'description=a place to live' -F 'peoples=3' -F 'room=good' -F 'bed=single' -F 'rooms=4' -F'bathrooms=1' -F'contract=good' -F'price=100' -F'security_deposit=1000'
```

```json
{
  "ID": "ObjectID(\"60e148279e3aeaa3fe6075bf\")",
  "Owner": "parham.alvani@gmail.com",
  "Title": "sweet",
  "Location": "italy",
  "Description": "a place to live",
  "Peoples": 3,
  "Room": "good",
  "Bed": 1,
  "Rooms": 4,
  "Bathrooms": 1,
  "Smoking": false,
  "Guest": false,
  "Pet": false,
  "BillsIncluded": false,
  "Contract": "good",
  "SecurityDeposit": 1000,
  "Photos": {},
  "Price": 100
}
```
