## 一、Through NET/HTTP to achieve HTTPS two-way authentication

`openss`Related Information Links：   
 `openssl`：`https://www.openssl.org/docs/`

**Generate CA certificates, server-side, client-side certificates, and private keys from OpenSS**   
1、Generate CA Certificates  
   `openssl genrsa -out ca.key 2048`  
   `openssl req -x509 -new -nodes -key ca.key -subj "//CN=localhost" -days 1 -out ca.crt`

2、Generate the server certificate and private key  
   `openssl genrsa -out server.key 2048`  
   `openssl req -new -key server.key -subj "//CN=localhost" -out server.csr`  
   `openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 1`

3、Generate the client certificate and private key  
   `openssl genrsa -out client.key 2048`  
   `openssl req -new -key client.key -subj "//CN=localhost" -out client.csr`  
   `openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 1`
   
   
## 二、HTTPS bi-directional authentication is implemented using Let's Encrypt

`letsEncrypt`Related Information Links：   
 letsEncrypt：`https://letsencrypt.org/`
 
**Generate Let's Encrypt certificates with official requirements**  
1、Start by configuring the site link for HTTPS one-way validation based on Let's Encrypt certificates  
   `wget https://dl.eff.org/certbot-auto`  
   `chmod u+x certbot-auto`   
   `./certbot-auto certonly --standalone -m *@search.cn --agree-tos -d *.search.cn`  
2、Modify the configuration of NGINX  
`server {`   
  `listen 80;`   
  `listen 443 ssl;`  
  `server_name abc.yubangweb.com;`    
  `root /var/web;`  
  `ssl_certificate /etc/letsencrypt/live/abc.yubangweb.com/fullchain.cer;`   
  `ssl_certificate_key /etc/letsencrypt/live/abc.yubangweb.com/privkey.key; `   
  `ssl_trusted_certificate /etc/letsencrypt/live/abc.yubangweb.com/ca.cer; `  
  `location / { }`  
`}`  
3、自签客户端证书  
**# generate primary key**  
`openssl genrsa -des3 -out ca.key 4096`  
**# or if you don't want a password:**      
`openssl genrsa -out ca.key 4096`   
**# generate a cert**  
`openssl req -new -x509 -days 365 -key ca.key -out ca.crt`  

You will need to fill in some information when generating the ca.crt file, so you can fill in any information you want：  
$ `openssl req -new -x509 -days 365 -key ca.key -out ca.crt You are about to be asked to enter information that will be incorporated into your certificate request.What you are about to enter is what is called a Distinguished Name or a DN.There are quite a few fields but you can leave some blank For some fields there will be a default value,If you enter '.', the field will be left blank.`    
`Country Name (2 letter code) [AU]:CN`  
`State or Province Name (full name) [Some-State]:ShenZheng`  
`Locality Name (eg, city) []:ShenZheng`  
`Organization Name (eg, company) [Internet Widgits Pty Ltd]:`  
`Organizational Unit Name (eg, section) []:ShenZheng`  
`Common Name (e.g. server FQDN or YOUR name) []:ShenZheng`  
`Email Address []:1447560092@qq.com`  

`The CA.CRT created here is to be placed on the server side, and then the certificate installed on the client side is created`  

**Create the client certificate**  
Create a certificate request：  
`openssl genrsa -out client.key 4096`  
`openssl req -new -key client.key -out client.csr`  
`You can just type it in as you did up here`  
**The CA certificate is then used for issuance：**  
`openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt`  
`Issue the client certificate client.crt`

**Install the client certificate**  
`The client. CRT generated in the previous step cannot be installed directly and needs to be converted to PKCS format：`  
`openssl pkcs12 -export -clcerts -in client.crt -inkey client.key -out client.p12`  
`Import the resulting client certificate into the browser`    

**Deploy to the server**  
`The server only needs a single ca.crt file`  
`scp sa.crt /nginx/conf/ssl/`  
**Add the following configuration to Nginx:**  
`server {`  
    `...`  
    `ssl_client_certificate /nginx/conf/ssl/ca.crt;`  
    `ssl_verify_client optional;`  
    `...`  
`}`   
**The ssl_verify_client value can be optional | on two：**  
1、When on is selected, client authentication is enforced and cannot be accessed if it fails；  
2、Authentication is optional when optional and can be determined by the $SSL_CLIENT_VERIFY variable。  
`For example, if we want the root directory of the website to be accessible arbitrarily, but the /admin path requires authentication to access it, we can configure it this way：`  
`server {`  
    `...`  
    `ssl_client_certificate /nginx/conf/ssl/ca.crt;`  
    `ssl_verify_client optional;`  
    `...`  
    `location /admin {if ($ssl_client_verify != SUCCESS) {return 401;}proxy_pass http://localhost:1001;}`    
`}`  

**complete**  
`When the website is opened and the machine is equipped with the corresponding client certificate, a prompt box requesting certificate authentication will appear.`
