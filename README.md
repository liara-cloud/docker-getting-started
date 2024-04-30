## Nginx conf for reverse proxy
### Installation (Liara)
- create a [docker app](https://console.liara.ir/apps/create) on Liara
- connect your domain to the created app
```
git clone https://github.com/liara-cloud/docker-getting-started.git
```
```
cd docker-getting-started
```
```
git checkout nginx
```
```
liara deploy --port 80
```
