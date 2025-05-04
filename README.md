# Dumb Text Storage Server

This is a simple Go-based HTTP server that allows clients to:
- Generate unique files
- Store text under that file
- Retrieve the stored text
- Automatically delete files after they become unused

---

##  Features

- `GET /generate` — Generates a new unique file and returns its key.
- `POST /{key}` — Appends text to the file associated with the key.
- `GET /{key}` — Reads and returns the file contents with lines in reverse order.

---

##  Requirements

- Go 1.18 or later

---

## How to Use

### 1. **Build and Run the Server**

```bash
  go build -o dumb-server main.go
  ./dumb-server
```
### 2. **Generate a Key**
```bash
  // a GET request to /generate
  curl http://localhost:6969/generate
  // example response: d3f1a9c0-1d53-4269-bc80-a77d7b994c8a
```

### 3. **Add Text**
```bash
  // POST to /{id}
  curl -X POST http://localhost:6969/d3f1a9c0-1d53-4269-bc80-a77d7b994c8a -d "First line of text"
  curl -X POST http://localhost:6969/d3f1a9c0-1d53-4269-bc80-a77d7b994c8a -d "Second line"
```

### 4. **Share URL to share the file**

```bash
  // GET /{id}
  curl http://localhost:6969/d3f1a9c0-1d53-4269-bc80-a77d7b994c8a
/* output:
@14:02:11 - :
Second line



@14:01:55 - :
First line of text
*/
```

## Note

- Text is stored as actual files written to disk
- If you intend to expose this server to the internet, it is recommended to handle rate limiting first

