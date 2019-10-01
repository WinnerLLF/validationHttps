## 通过net/http实现https双向验证

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
   
