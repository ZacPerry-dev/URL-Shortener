# URL-Shortener
A URL shortening service written in Go.

## Future "why I made this and what does it do" 

## Future Features 

## Future Usage Info 


## TODO
Part 1

- Create a REST API that allows clients to add a URL to the list of URLs that are currently shortened
   - Setup DB to store this information (mapping short to real URLs)
   - Figure out how to generate short key for each url (hashing)
   - Return an HTTP status when the requests succeeds, along with the shortened URL 
   - When the URL already exists, return the existing shorthand url and the same HTTP status
   - API should be idempotent (same endpoint for everything) (return if exists, otherwise add and create and return) 
   - IF invalid request, return an error HTTP status
   - Also check if there is an identical key in the DB whenever you generate a new hash
