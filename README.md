 
docker run -it --name rsshub-app --network host --env-file .env rsshub-app bash

сначало нужно запустить докер как обычно up --build 
потом нужно удалить контейнер rsshub-app
потом вот это 
docker run -it --name rsshub-app --network host --env-file .env rsshub-app bash
(Docker-контейнер с интерактивной оболочкой (bash) что бы можно было писать команды ./rsshub )
а то у нас rsshub-app останавливается сразу же после того как создается и мы не можем выполнить ничего а есть его билдить за пределами докера мы не можем подключиться к докеру вот такая ирония 
