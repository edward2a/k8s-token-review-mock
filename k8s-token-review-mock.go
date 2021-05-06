/* based on:
 *    https://github.com/shaj13/go-guardian/blob/master/_examples/kubernetes/mock.go
 *    https://golangr.com/golang-http-server/
*/
package main

import (
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
  // Print token info
  log.Printf("JWT service account for service: %s \n", serviceJWT)
  log.Printf("JWT service account for agent: %s \n", agentJWT)

  // Set routing rules
  http.HandleFunc("/apis/authentication.k8s.io/v1/tokenreviews", TokenReview)
  http.HandleFunc("/", Nope)

  //Use the default DefaultServeMux.
  err := http.ListenAndServe(":8080", nil)
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
  return
}
