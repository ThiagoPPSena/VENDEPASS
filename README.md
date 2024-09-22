# VENDEPASS

docker-compose up -d --build
docker-compose exec client /bin/bash

## ANDAMENTO

Fizemos dois arquivos GO LANG, um cliente e um servidor para se comunicarem via protocolo TCP, somente enviando uma mensagem de COMPRA. Estamos criando métodos agora baseados na ideia de serviço HTTP, como GET, PUSH, PUT, ... No entanto, estamos criando métodos próprios, por enquanto BUY (para compra da passagem) e GET (para requisições de visualização de trechos).
