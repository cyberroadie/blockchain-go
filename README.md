# blockchain-go
Simple blockchain implementation in Go


Running:
  go run main.go
  
Testing
```
  cd blockchain
  go test
```

New block:

    http://localhost:8080/block

Add transaction:
```
curl -X POST -H "Content-Type: application/json" -d '{
"sender": "d4ee26eee15148ee92c6cd394edd974e",
"recipient": "someone-other-address",
"amount": 5
}' "http://localhost:8080/transactions/new"
```

Get whole block chain:

    http://localhost:8080/chain
