# IDARU - Batch URL Manager

IDARU is a command-line tool for efficiently managing URLs in batch. It allows you to work with multiple URLs at once, whether it's for validation, manipulation, or other URL-related tasks.

## Features

- **URL Validation:** Check the syntax of a list of URLs and identify any issues.
- **URL Manipulation:** Perform bulk operations on URLs, like adding or removing parameters.
- **URL Export:** Save results and collections to files for future use.

## Installation

Follow these steps to install IDARU:

1. **Clone the Repository:**
   ```sh
   git clone https://github.com/siriil/idaru.git
   cd idaru
   go build
   mv idaru /usr/bin/

## Usage

Example Input File:
```
▶ cat urls.txt 
https://domain.tld/path?test=debug&file=index.html
https://domain.tld/path?file=index.html&test=debug&id=1
https://domain.tld/pathfile?file=index.html&test=debug
https://domain2.tld/parent/path?file=index.html&test=debug
https://domain2.tld/path
```

### Set All Keys in Query String Values

```
▶ cat urls.txt | idaru -s "*=new"
https://domain.tld/path?test=new&file=new
https://domain.tld/path?file=new&test=new&id=new
https://domain.tld/pathfile?file=new&test=new
https://domain2.tld/parent/path?file=new&test=new
https://domain2.tld/path
```

### Set Specific Keys in Query String Values

```
▶ cat urls.txt | idaru -s "test=true" -s "id=2"
https://domain.tld/path?test=true&file=index.html
https://domain.tld/path?file=index.html&test=true&id=2
https://domain.tld/pathfile?file=index.html&test=true
https://domain2.tld/parent/path?file=index.html&test=true
https://domain2.tld/path
```

### Append All Keys in Query String Values

```
▶ cat urls.txt | idaru -a "*=add"
https://domain.tld/path?test=debugadd&file=index.htmladd
https://domain.tld/path?file=index.htmladd&test=debugadd&id=1add
https://domain.tld/pathfile?file=index.htmladd&test=debugadd
https://domain2.tld/parent/path?file=index.htmladd&test=debugadd
https://domain2.tld/path
```

### Append Specific Keys in Query String Values

```
▶ cat urls.txt | idaru -a "id=add"
https://domain.tld/path?test=debug&file=index.html
https://domain.tld/path?file=index.html&test=debug&id=1add
https://domain.tld/pathfile?file=index.html&test=debug
https://domain2.tld/parent/path?file=index.html&test=debug
https://domain2.tld/path
```

### Merge All Keys in One Path

```
▶ cat urls.txt | idaru -m
https://domain.tld/path?file=index.html&test=debug&id=1
https://domain.tld/pathfile?file=index.html&test=debug
https://domain2.tld/parent/path?file=index.html&test=debug
```

### Filter URLs with Parameters 

```
▶ cat urls.txt | idaru -fP -m
https://domain.tld/path?file=index.html&test=debug&id=1
https://domain.tld/pathfile?file=index.html&test=debug
https://domain2.tld/parent/path?file=index.html&test=debug
```


