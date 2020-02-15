# doxie-receiver
A tool for automatically downloading jpegs from the Doxie Go scanner.

## Behaviour
The tool connects to the IP of Doxie on port 80 when it connects to the local wifi, and downloads all files.
By default files will be deleted from the unit after they are downloaded.

Files are deposited to "./images" by default.

## Basic execution
Turn on your doxie, and identify the IP address - How to do this is out of the scope of this document.
Run the tool as follows:
```./doxie-receiver -ip 1.2.3.4```
This will download and delete all jpegs.

## Switches
* -daemon: Continue looking for the given IP, and begin downloading when Doxie returns a successful status page
* -ip: Specify an IP address for the Doxie
* -no-delete: Leave downloaded files on the Doxie
* -output: Specify an output path for downloaded files
* -port: Provide a different port for doxie - Default 80

## Support
This tool is provided without any knowledge or support of the creators of Doxier, and is an independent project built on their published API.

## License
This software is provided under the terms of the [Apache License, Version 2.0]("https://opensource.org/licenses/Apache-2.0") license