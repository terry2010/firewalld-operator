# firewalld-operator
golang firewalld operator , punch a hole in firewalld 

#how to

```bash
go build main.go

./main --port=8080 --hole=3306

```

#run
```bash
curl http://localhost:8080/firewall/update/id123id
```

