# gocqlastra

This is a wrapped of gocql to use with DataStax Astra. For more information, please go to see [gocql](https://github.com/gocql/gocql) document.

## Usage
```go
    cluster, _ := gocqlastra.NewCluster("<<PATH/TO/>>secure-connect.zip")
    cluster.Authenticator = &gocql.PasswordAuthenticator{
        Username: "<<CLIENT ID>>",
        Password: "<<CLIENT SECRET>>",
    }
    session, _ := cluster.CreateSession()
    defer session.Close()
```

## Credits
I have taken most of the code logic from the [cassandra-driver](https://github.com/datastax/nodejs-driver/blob/master/lib/datastax/cloud/index.js) project written in Javascript and rewrite it using Golang.
And this is a wrapped of [gocql](https://github.com/gocql/gocql). I will give credit to gocql. You don't need to use my package, but it can save your development time.
