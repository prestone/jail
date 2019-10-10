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
guantanamo := jail.New()

//persistent
guantanamo := jail.New("jail.db")

//so lets put some into db
//you can use any ID, jail time duration
//and you can write a reason as optional
guantanamo.Put("127.0.0.1", time.Hour * 10, "reason is ban")
guantanamo.Put("127.0.0.2", time.Second * 30)
guantanamo.Put(42, time.Minute)
guantanamo.Put([]byte("user1"), time.Minute, "ban")


//check the id
guantanamo.Check("127.0.0.1") //true... in jail now

//get reason
guantanamo.Reason("127.0.0.1") //string "reason is ban"

//delete
guantanamo.Delete("127.0.0.1")

//count
//so how many items in jail
guantanamo.Count() //4
```

### Thanks
I hope it help you.
