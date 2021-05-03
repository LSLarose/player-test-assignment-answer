# player-test-assignment-answer
This repository was created to answer an internship's technical interview assignment.

# rationale
This program was written in go, as I understood that the internship would use golang and I had never used it before.

Upon further investigation, it was found to offer the researched caracteristics out of the box, that is:
- it offered automatic documentation for developpers with the 'godoc' tool
- it offered out of the box unit-testing capabilities (tutorial followed: https://golangdocs.com/golang-unit-testing)
- golang comes preinstalled on docker containers, which is where it would be expected to be deployed (see: https://www.docker.com/blog/docker-golang/)

This confirmed my choice and I started doing the tutorials, slowly customizing the code to provide the intended functionnalities

Assumptions:
 - This tool's final environment would be the cloud

How to "use" this tool
//TODO

To get developper documentation for this project:
[godoc -http=:6060]

This wil compile and execute a Go Documentation Server, viewable in a browser at localhost:6060.

It offers documentation for all modules and packages used in this project

The code from originating from this particular project can be found documented under the section "Third party/exemple.com"
e.g. http://localhost:6060/pkg/example.com/greetings/

To run this project on your own machine:
[cd ./main]
[go run .]

This will start a server on your machine accessible at *IP*.
It will then answer normally to the documented API Calls
