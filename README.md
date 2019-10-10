# 🏢 Jail
This is memory/persistent jail engine. Sometimes you want to put some data into jail for some reason. Now you can do it simple with Jail. Its a db with expiration mode. Every item has a expire period.

### Install
```go get github.com/prestone/jail```

### Example
```go
//lets create a jail engine
//if you write a filepath it will a persistent store
//or you can just memory based store it

//memory based
j := jail.New()

//persistent
j := jail.New("jail.db")

//so lets put some into db
//you can use any ID, jail time duration
//and you can write a reason as optional
j.Put("127.0.0.1", time.Hour * 10, "reason is ban")
j.Put("127.0.0.2", time.Second)
j.Put(42, time.Minute)
j.Put([]byte("user1"), time.Minute, "ban")


//check the id
j.Check("127.0.0.1") //true... in jail now

//get reason
j.Reason("127.0.0.1") //string "reason is ban"

//delete
k.Delete("127.0.0.1")
```

### Thanks
I hope it help you.
