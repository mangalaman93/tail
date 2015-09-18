# tail
The package is implemented using `tail` function provided by Linux. It simply reads the stdout of `tail` using `bufio.Scanner` and pass it to the application through a go channel.
