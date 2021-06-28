# dianthus-server

## Docker
  
### Version

```bash
> docker --version
Docker version 19.03.12, build 48a66213fe

> docker-compose --version
docker-compose version 1.27.2, build 18f557f9
```
  
### Build
  
```bash
> docker-compose up --build -d
```
  
### Run
  
```bash
> docker-compose exec app go run main.go
```
  
### Usage
  
```bash
> curl -u user:pass "localhost:8080/v1/roman?target=kasumi" | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  6121    0  6121    0     0   229k      0 --:--:-- --:--:-- --:--:--  229k
[
  {
    "raw": "暑い",
    "roman": "atsui",
    "vowels": "aui"
  },
  {
    "raw": "略儀",
    "roman": "ryakugi",
    "vowels": "aui"
  },
  {
    "raw": "休み",
    "roman": "yasumi",
    "vowels": "aui"
  },
  ...
]
```