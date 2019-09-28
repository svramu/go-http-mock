# go-http-mock


**Usage**

Default run, just logs all requests.  
```go-http-mock.exe```  

For Callbacks, have a [conf.yaml](https://github.com/svramu/go-http-mock/blob/master/conf.yaml) at same place as the executable.  
```go-http-mock.exe```  
  
**Development**

Get the dependencies. If in golang 1.3, any folder is fine.  
```go get gopkg.in/yaml.v2```  
  
Build for your os  
```go build main.go```  

Append console out to a log file, if needed.  
```go run main.go >> history.log```  
  
or, create a new file for every run:  
```go run main.go > history.log```  
