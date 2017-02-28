digitalocean-ssh-cleanup
========================

A simple programm to delete ssh public keys via the api.


## Example

```bash
  go run main.go -token=mytoken -keyname=swarm
```

Replace `mytoken` with an valid api token

This will delete all keys that contains the keyname `swarm`


## ToDo 
- [] Handle Ratelimiting
- [] Make a Container

