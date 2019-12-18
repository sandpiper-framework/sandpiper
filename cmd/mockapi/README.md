# **mockapi**

mockapi creates an http server that mimics responses to a subset of the sandpiper api for client testing purposes (without a database and before the sandpiper server is completed). We anticipate this utility having a short lifespan.

### **Getting Started**

Copy to a new folder called mockapi and then execute these commands.

```
go mod download
go build
```

### **Running the server**

Simply run the server.

```
.\mockapi
```
Then use a browser to see the registered routes:

```
http://localhost:3030/v1/routes
```
You should see something like the following:

```
// 20191217163813
// http://localhost:3030/v1/routes

[
  {
    "method": "GET",
    "path": "/v1/login",
    "name": "main.login"
  },
  {
    "method": "GET",
    "path": "/v1/slices",
    "name": "main.getMySlices"
  },
  {
    "method": "POST",
    "path": "/v1/slices/:id",
    "name": "main.postObject"
  },
  {
    "method": "GET",
    "path": "/v1/routes",
    "name": "main.listRoutes"
  }
]
```

### **Using the server**

HTTP commands can be sent to the mock api using Postman, curl, our Sandpiper CLI utilities or any other means.

**For example:**

```
curl --data "{'token':'eyJhbGsI...', 'oid':'GQBFRvf112p33Q', \
  'type':'aces-file', 'payload':'--base64 data--'}" \
  -H "Content-Type:application/json" \
  http://localhost:3030/v1/slices/aap-brakes
```

This sends the POST request to add an 'aces-file' object to the 'aap-brakes' slice. 

### **Stopping the server**

```
Ctrl-C
```