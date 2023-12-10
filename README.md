## Usage
Add package to project

```console
go get github.com/sylvesterogoe/calmshutdown
```

## Example Usage for Start()
Run function in the background as a goroutine

```go
server := http.Server{}
shutdownError := make(chan error)

go calmshutdown.Start(shutdownError, &server, 5*time.Second)

err := server.ListenAndServe()
	if err != nil{
		//handle error
	}

err = <-shutdownError
if err != nil {
	// shutdown failed
	//handle failure
}
```

## Example Usage for AwaitGoroutinesAndStart()
Run function in the background as a goroutine
This function is suitable when there several background processes which are tracked

```go
//application wide waitgroup incremented for every goroutine added
	var globalWG sync.WaitGroup

	server := http.Server{}
	shutdownError := make(chan error)

	go calmshutdown.AwaitGoroutinesAndStart(shutdownError, &server, 5*time.Second, &globalWG)

	err := server.ListenAndServe()
	if err != nil{
		//handle error
	}

	err = <-shutdownError
	if err != nil {
		// shutdown failed
		// handle failure
	}
```
