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
    "path": "/login",
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

You could also use a standard web browser for the GET commands. For example:

`http://localhost:3030/v1/slices`

should display a list of your slices:

```
{
  "slices": [
    {
      "slice-id": "08efdf90-a815-4cf7-b71c-008e5fd31cce",
      "slice-name": "AAP-Brakes",
      "slice-hash": "cf23df2207d99a74fbe169e3eba035e633b65d94",
      "metadata": {
        "pcdb-version": "2019-09-27",
        "vcdb-version": "2019-09-27"
      },
      "count": 2919
    },
    {
      "slice-id": "cb4b768b-6d6b-4965-a29a-9052a80dbbbb",
      "slice-name": "AAP-Wipers",
      "slice-hash": "1a804c61e1a70ab37b912792ee846de7378c4a36",
      "metadata": {
        "pcdb-version": "2019-09-27",
        "vcdb-version": "2019-09-27"
      },
      "count": 2342
    }
  ],
  "count": 2
}
```

### **Stopping the server**

```
Ctrl-C
```