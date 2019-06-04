
# growler

A cli web browser/crawler written in Go
## Getting Started

These instructions will get you a copy of the project up and running on your local machine for various purposes.

### Prerequisites

In order to be able to install and run the application you will need the following programs in your machine.

```
go version >= go1.12.1

```

### Installing

You can install the app by running the following command in the root folder of the project.

```
go install -v .
```

### Running

Considering you installed the app by running the previous step

```
#The following command will launch the app regularly
growler

#The next command will run the app with extra logging capabilities
growler debug
```

To verify the extra logging capabilities, the log's output is being done to the `info.log` file.
```
#checking the log
cat info.log

#getting the last 10 lines
tail info.log

#checking for a certain error/output
cat info.log | grep -i 'tab-1-start'
```

### Testing

Writing the command `help` will print out a more detailed explanation of how commands work.
```
[GRWLR] > help

 	 Command usage help

get <url> 					 - Opens a new tab and gives it the job of requesting the given url.

list 						 - Lists all open tabs, their current or last job, and their current status.

read <image_filename> 				 - Reads an image and tries to decode it into a txt with RGB values.

download <url_to_download> <output_filename> 	 - Requests the URL and saves the response or downloadable file to the output file.

stop <tab id> 					 - Stops the thread represented by the input id.
							 If a thread is working it will wait until it has finished to stop it.

job <url> 					 - Adds the given url to the pool of jobs, it will be executed once a tab is available.

exit 						 - Stops all tabs and exits the browser.


```

## Built With

* [golang](https://golang.org) - The programming language


