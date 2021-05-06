/* based on:
 *    https://github.com/shaj13/go-guardian/blob/master/_examples/kubernetes/mock.go
 *    https://golangr.com/golang-http-server/
*/
package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "strings"
)

const (
  agentJWT    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
  serviceJWT  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTYiLCJuYW1lIjoic3lzdGVtOnNlcnZpY2U6YWNjb3VudCIsImlhdCI6MTUxNjIzOTAyMn0.4pHu9y6vJvtOnLhpz7M3Znnvcdpm7GCiHPCPYzyxps8"

  authenticatedUser = `{
  "metadata":{
    "creationTimestamp":null
  },
  "spec":{},
  "status":{
    "authenticated":true,
    "user":{
      "username":"system:serviceaccount:curl_agent",
      "uid":"1"
    }
  }
}
`

  unauthenticatedUser = `{
  "metadata":{
    "creationTimestamp":null
  },
  "spec":{},
  "status":{
    "authenticated":false,
  }
}
`
)

func main() {

  tls := flag.Bool("tls", false, "Enable TLS listener")
  key := flag.String("key", "ssl_key.pem", "TLS private key")
  cert := flag.String("cert", "ssl_cert.pem", "TLS certificate")
  addr := flag.String("addr", "localhost", "Bind address")
  port := flag.Int("port", 8080, "Listener port (8080 | 8443)")
  flag.Parse()

  if (*tls && *port == 8080) { *port = 8443 }
  bind_addr :=  fmt.Sprintf("%s:%d", *addr, *port)

  // Print init info
  log.Printf("JWT service account for service: %s \n", serviceJWT)
  log.Printf("JWT service account for agent: %s \n", agentJWT)
  if (*tls) { log.Printf("TLS enabled") }
  log.Printf("Listen address: %s", bind_addr)

  // Set routing rules
  http.HandleFunc("/apis/authentication.k8s.io/v1/tokenreviews", TokenReview)
  http.HandleFunc("/", Nope)

  //Use the default DefaultServeMux; start plain or tls server
  var err error
  if !(*tls) {
    err = http.ListenAndServe(bind_addr, nil)
  } else {
    err = http.ListenAndServeTLS(bind_addr, *cert, *key, nil)
  }

  if err != nil {
    log.Fatal(err)
  }

}

func TokenReview(w http.ResponseWriter, r *http.Request) {
  body, _ := ioutil.ReadAll(r.Body)
  if strings.Contains(string(body), agentJWT) {
    w.WriteHeader(200)
    w.Write([]byte(authenticatedUser))
    return
  }

  w.WriteHeader(401)
  w.Write([]byte(unauthenticatedUser))
  return
}

func Nope(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(418)
  w.Write([]byte("Nope...\n"))
  w.Write([]byte("Try '/apis/authentication.k8s.io/v1/tokenreviews'\n"))
  return
}
