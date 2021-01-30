# Status Quode  
  
Tool to loop supplied URLs and get their status codes.

## Why?  
  
Quickly identify specific status codes for discovered endpoints. Tools like ffuf, gobuster (etc) already do this of course but everyone like small unix tools \o/  
  
## Installing  
```
go get github.com/crawl3r/statusquode
```  
  
## Usage  
Standard Run  
```
cat urls.txt | ./statusquode
```
  
Run quiet mode to just receive some easy greppable output
```
cat urls.txt | ./statusquode -q
```  
  
Run in quiet mode and only print results with status codes 200 and 404 to stdout
```
cat urls.txt | ./statusquode -q -s 200,404
```
  
## License  
I'm just a simple skid. Licensing isn't a big issue to me, I post things that I find helpful online in the hope that others can:  
 A) learn from the code  
 B) find use with the code or   
 C) need to just have a laugh at something to make themselves feel better  
  
Either way, if this helped you - cool :)  
