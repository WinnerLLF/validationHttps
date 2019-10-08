## 一、通过net/http实现https双向验证

`openss`相关资料链接：   
 `openssl`：`https://www.openssl.org/docs/`

**通过openss生成CA证书、服务端、客户端的证书和私钥**   
1、生成CA证书  
   `openssl genrsa -out ca.key 2048`  
   `openssl req -x509 -new -nodes -key ca.key -subj "//CN=localhost" -days 1 -out ca.crt`

2、生成服务端证书和私钥  
   `openssl genrsa -out server.key 2048`  
   `openssl req -new -key server.key -subj "//CN=localhost" -out server.csr`  
   `openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 1`

3、生成客户端证书和私钥  
   `openssl genrsa -out client.key 2048`  
   `openssl req -new -key client.key -subj "//CN=localhost" -out client.csr`  
   `openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 1`
   
   
## 二、通过Let’s Encrypt实现https双向验证 

`letsEncrypt`相关资料链接：   
 letsEncrypt：`https://letsencrypt.org/`
 
**通过官方要求生成Let’s Encrypt证书**  
1、先基于Let’s Encrypt证书配置https单向校验的网站链接  
   `wget https://dl.eff.org/certbot-auto`  
   `chmod u+x certbot-auto`   
   `./certbot-auto certonly --standalone -m abc@yubangweb.com --agree-tos -d abc.yubangweb.com`  
2、修改nginx的配置  
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

生成 ca.crt 文件的时候会需要填写一些信息，随便填写都行：  
$ `openssl req -new -x509 -days 365 -key ca.key -out ca.crt You are about to be asked to enter information that will be incorporated into your certificate request.What you are about to enter is what is called a Distinguished Name or a DN.There are quite a few fields but you can leave some blank For some fields there will be a default value,If you enter '.', the field will be left blank.`    
`Country Name (2 letter code) [AU]:CN`  
`State or Province Name (full name) [Some-State]:ShenZheng`  
`Locality Name (eg, city) []:ShenZheng`  
`Organization Name (eg, company) [Internet Widgits Pty Ltd]:`  
`Organizational Unit Name (eg, section) []:ShenZheng`  
`Common Name (e.g. server FQDN or YOUR name) []:ShenZheng`  
`Email Address []:1447560092@qq.com`  

`这里创建的 ca.crt 是会放到服务端的，接下来创建安装在客户端的证书`  

**创建客户端证书**  
创建证书请求：  
`openssl genrsa -out client.key 4096`  
`openssl req -new -key client.key -out client.csr`  
`这里和上面一样随便输入即可`  
**然后使用 ca 证书进行签发：**  
`openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt`  
`签发好客户端证书 client.crt`

**安装客户端证书**  
`上一步产生的 client.crt 并不能直接安装，需要转成 PKCS 格式：`  
`openssl pkcs12 -export -clcerts -in client.crt -inkey client.key -out client.p12`  
`将产生的客户端证书导入到浏览器`    

**部署到服务器**  
`服务端只需要一个 ca.crt 文件即可`  
`scp sa.crt /nginx/conf/ssl/`  
**Nginx 中加入如下的配置：**  
`server {`  
    `...`  
    `ssl_client_certificate /nginx/conf/ssl/ca.crt;`  
    `ssl_verify_client optional;`  
    `...`  
`}`   
**其中 ssl_verify_client 的值可以是 optional | on 者两个：**  
1、当选择 on 时会强制进行客户端认证，失败无法访问；  
2、当选择 optional 的时候，认证是可选的，是否认证成功可以从 $ssl_client_verify 变量得知。  
`比如，我们想要网站根目录是任意访问的，但是 /admin 路径下是需要认证才能访问的，就可以这么配置：`  
`server {`  
    `...`  
    `ssl_client_certificate /nginx/conf/ssl/ca.crt;`  
    `ssl_verify_client optional;`  
    `...`  
    `location /admin {if ($ssl_client_verify != SUCCESS) {return 401;}proxy_pass http://localhost:1001;}`    
`}`  

**完成**  
`当打开网站且本机装有对应的客户端证书时，就会出现请求证书认证的提示框.`