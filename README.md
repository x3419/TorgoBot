## TorgoBot
#### PoC Golang malware for educational purposes only.
###### Note: This project is in BETA

### Features
- No C2 required!
    - TorgoBot creates a Tor Hidden Service for encrypted, anonymous communication
    - This means you don't have to worry about port forwarding and don't have to pay for C2 infrastructure
    - You just need the onion ID (e.g. jasf1j2l1ln.onion) of your servers which can be sent as POST data upon initial execution to a staging server
        - Implementation of staging is currently in progress...
- Remote shell
- Execute-assembly
    - Specify any .Net Assembly executable on the client machine to have it execute in memory on the server
        - Fileless, stealthy and extends unlimited functionality 
    - Note: Only x64 servers are compatible by default but x86 can be made compatible very easily 
#
### Commands
- shell
    - Remote interactive shell
- back
    - Returns you to the main menu
- execute-assembly
    - Ex) execute-assembly C:\Hello.exe
        - C:\Hello.exe is the path on the client machine
#    
### Build

go get golang.org/x/crypto/ed25519

go get golang.org/x/net/proxy

go get github.com/x3419/TorgoBot
#

##### Credit
Thanks to the following authors whose projects I modified to make this
- github.com/cretz/bine
- github.com/lesnuages/go-execute-assembly
    - The original package did not correctly capture the .net assembly stdout so I fixed and integrated it into this project
#

##### TODO
- Test interactivity for when .net assemblies are expecting stdin
- Implement secure authentication
- Implement 

