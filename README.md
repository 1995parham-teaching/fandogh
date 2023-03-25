<h1 align="center"> Fandogh 🌰 </h1>

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/1995parham-teaching/fandogh/ci.yaml?label=ci&logo=github&style=flat-square&branch=main)
[![Codecov](https://img.shields.io/codecov/c/gh/1995parham-teaching/fandogh?logo=codecov&style=flat-square)](https://codecov.io/gh/1995parham-teaching/fandogh)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/1995parham-teaching/fandogh)](https://pkg.go.dev/github.com/1995parham-teaching/fandogh)

## APIs

Register new user with JSON request as follows:

```bash
curl 127.0.0.1:1378/register -X POST -d '{ "email": "parham.alvani@gmail.com", "name": "Parham Alvani", "password": "123456" }' -H 'Content-Type: application/json'
```

```json
{
  "Email": "parham.alvani@gmail.com",
  "Password": "123456",
  "Name": "Parham Alvani"
}
```

Login into system and getting the token:

```bash
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

Creating new home requires using the POST request with form data because it contains images:

```bash
curl -vvv 127.0.0.1:1378/api/homes -X POST  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyIiwiZXhwIjoxNjI1Mzc5MDE3LCJqdGkiOiI4NDRhMzQ4Yy03OGVjLTRlNTctODJhZi03YjU3NTNmNjk5ZjciLCJpYXQiOjE2MjUzNzU0MTcsImlzcyI6ImZhbmRvZ2giLCJuYmYiOjE2MjUzNzU0MTcsInN1YiI6InBhcmhhbS5hbHZhbmlAZ21haWwuY29tIn0.EZUWQ-sLP1ClA0vtK6vZEcQ4qf3ZaBm9VpFV6smEwUc' -F 'title=sweet' -F 'location=italy' -F 'description=a place to live' -F 'peoples=3' -F 'room=good' -F 'bed=single' -F 'rooms=4' -F'bathrooms=1' -F'contract=good' -F'price=100' -F'security_deposit=1000' -F'photos=1,2' -F'1=@1.png' -F'2=@2.png'
```

```json
{
  "ID": "60e1535541e125c415973cd2",
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
  "Photos": {
    "1": "60e1535541e125c415973cd2_1",
    "2": "60e1535541e125c415973cd2_2"
  },
  "Price": 100
}
```
