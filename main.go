package main

import (
  "context"
  "encoding/json"
  "fmt"
  "google.golang.org/appengine"
  "google.golang.org/appengine/blobstore"
  "google.golang.org/appengine/image"
  "io"
  "net/http"
  "strconv"
)

type Body struct {
  BucketName string `json:"BucketName"`
  FilePath   string `json:"FilePath"`
}

type Response struct {
  Status int    `json:"status"`
  Data   string `json:"data"`
}

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  switch r.Method {
  case "POST":
    postHandler(w, r)
    break
  case "DELETE":
    deleteHandler(w, r)
    break
  default:
    http.NotFound(w, r)
  }
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  if r.Header.Get("Content-Type") != "application/json" {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  length, err := strconv.Atoi(r.Header.Get("Content-Length"))

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  bodyStr := make([]byte, length)

  length, err = r.Body.Read(bodyStr)

  if err != nil && err != io.EOF {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  var body Body

  err = json.Unmarshal(bodyStr[:length], &body)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  ctx := appengine.NewContext(r)

  bucketName := body.BucketName
  filePath := body.FilePath

  servingURL, err := serveURL(ctx, bucketName, filePath)

  if err != nil {
    fmt.Fprintf(w, "[error!] %s", err)
    return
  }

  response := Response{http.StatusOK, servingURL}

  responseStr, err := json.Marshal(response)

  w.Header().Set("Content-Type", "application/json")

  w.Write(responseStr)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
  if r.Header.Get("Content-Type") != "application/json" {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  length, err := strconv.Atoi(r.Header.Get("Content-Length"))
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  bodyStr := make([]byte, length)

  length, err = r.Body.Read(bodyStr)

  if err != nil && err != io.EOF {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  var body Body

  err = json.Unmarshal(bodyStr[:length], &body)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  ctx := appengine.NewContext(r)

  bucketName := body.BucketName
  filePath := body.FilePath

  err = deleteURL(ctx, bucketName, filePath)

  if err != nil {
    fmt.Fprintf(w, "[error!] %s", err)
    return
  }

  response := Response{http.StatusOK, ""}

  responseStr, err := json.Marshal(response)

  w.Header().Set("Content-Type", "application/json")

  w.Write(responseStr)
}

func serveURL(ctx context.Context, Bucket, Path string) (string, error) {
  gsURL := fmt.Sprintf("/gs/%s/%s", Bucket, Path)

  blobKey, err := blobstore.BlobKeyForFile(ctx, gsURL)

  if err != nil {
    return "", err
  }

  servingURLOpts := &image.ServingURLOptions{Secure: true}

  url, err := image.ServingURL(ctx, blobKey, servingURLOpts)

  if err != nil {
    return "", err
  }

  return url.String(), nil
}

func deleteURL(ctx context.Context, Bucket, Path string) error {
  gsURL := fmt.Sprintf("/gs/%s/%s", Bucket, Path)

  blobKey, err := blobstore.BlobKeyForFile(ctx, gsURL)

  if err != nil {
    return err
  }

  err = image.DeleteServingURL(ctx, blobKey)

  if err != nil {
    return err
  }

  return nil
}
