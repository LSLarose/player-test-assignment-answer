If this was being exposed online, you wouldn't be able to use self-signed certificates and would insert
certificates signed by a valid authority in here.
For testing purposes, such as running the tool locally, one can do so.
These commands has been taken from https://medium.com/rungo/secure-https-servers-in-go-a783008b36da

This program expects two particular files to start its HTTPS server:

- "assets/localhost.key"
    which can be generated with:
`openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr`

- "assets/localhost.crt"
    which can be generated with:
`openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt`