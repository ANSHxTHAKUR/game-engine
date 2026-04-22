 Run it the same way 

 #Terminal 1 
 go run ./cmd/server/

 #Terminal 2 

 go run ./cmd/simulator/ --users=1000

 # check metrics anytime
 curl http://localhost:8080/metrics
