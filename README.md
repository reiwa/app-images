## Go

```
$GOROOT: not set
$GOPATH: /Users/uufish/go
```

```
$ go get google.golang.org/appengine
```

## Deploy

```
$ gcloud config set project <project-id>
$ gcloud app deploy --stop-previous-version
```

## How to use


POST: https://rmzmhcxihubl.appspot.com/

BODY:

```
{
	"bucketName": "<project-id>.appspot.com",
	"filePath": "<file-path>"
}
```

For example:

```
{
	"bucketName": "rmzmhcxihubl.appspot.com",
	"filePath": "default/skSsqe8C4gb1QNfsfN16"
}
```

Response:

```
{
	"status": 200,
	"data": "https://lh3.googleusercontent.com/abcdefg"
}
```
