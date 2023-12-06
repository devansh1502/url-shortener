# url-shortener

The url shortner service consists of two components
1. Http REST
2. MongoDB/In Memory Backend


## Usage
In order to use the service using mongodb, follow below instructions:
1. You need to have mongodb installed on your local.
2. After mongodb is installed, open your terminal and go to project path and enter `go run main.go`.This will start the server at localhost:8080
3. Once the server is up and running, you can try shortening the urls.
4. This will return a shortened url. Please save the shortened URl for your reference later on.
5. You can also try redirecting to the orginal path by giving the shortened url in the path with the redirect keyword.
6. Once you have tried shortening the URls, you can also get the metric count of top three most shortened URLs.

| Urls| Result|
|----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| localhost:8080/short/https:/www.youtube.com/ | {"short_url":"youtube.com/46O6pjZf"}                                                                                    |
| localhost:8080/redirect/youtube.com/46O6pjZf | Redirects to the Original URL                                                                                           |
| localhost:8080/metrics                       | [{"domain": "youtube.com","counter": 3},{"domain": "cricbuzz.com","counter": 2},{"domain": "mongodb.com","counter": 2}] |
|

## Note:
By default the service is using mongodb. In order to test the service with in memory backend, please make respective changes in main.go