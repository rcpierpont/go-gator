## Overview

Gator is an RSS feed aggregator CLI tool written in Go. The user can enter feeds by a name/url pair and the feed is saved to the db. 

By starting the aggregator function, posts are retrieved from the feeds that user has added. 

Basic information in the form of post title, publish date, and url can be accessed by using the browse commmand. This allows the user to easily navigate to that post in a browser.

Basic user management is also loosely implemented. If logged in as a user, posts from feeds added to that user will automatically be followed. This is useful if feeds need to be grouped into categories (for example "gaming" or "world news")


## Setup Info

Dependencies:
- macOS or linux
- postgres installed
- go installed
- goose installed (for running migrations to set up database, located in gator/sql/schema)
- config file in home directory named ".gatorconfig.json"
    - file should contain a single json object with two key/value pairs, representing connection string for postgres db and current logged in user (starting as empty string)
    - this can be copy/pasted into an empty file as a template: {"db_url":"postgres://username:password@localhost:5432/gator?sslmode=disable","current_user_name":"ryan"}



## Example usage

Gator is intended to be run in two separate terminal windows. 
1. First window runs the aggregator function to collect posts from feeds (see first usage example for 'agg' command) 
2. Second window retrieves posts from what is aggregated by querying the database.


Aggregate posts from all feeds by sending http request every 1 minute (specified by a duration string, e.g. 1h is 1 hour):

```
gator agg 1m
```

Register new user named "firstname":

```
gator register firstname
```


Log in to already-existing user "firstname":

```
gator login firstname
```


Add feed for bootdev blog:

```
gator addfeed "boot.dev blog" https://www.boot.dev/blog/index.xml
```


Follow feed associated with url https://www.boot.dev/blog/index.xml (requires feed already in database):

```
gator follow https://www.boot.dev/blog/index.xml
```


Unfollow feed associated with url https://www.boot.dev/blog/index.xml:

```
gator unfollow https://www.boot.dev/blog/index.xml
```


Browse 5 posts, sorted by most recent first, from all followed feeds:

```
gator browse 5
```