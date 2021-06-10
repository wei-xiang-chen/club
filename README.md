# club
此專案是結合了club house及woo talk的部分功能而成的匿名群聊系統。

展示連結 : (http://104.214.48.227:8090/login)

API文件 : (https://docs.google.com/document/d/1KVG8n_bbyeBAtAurPq8pybsug6FQT6-Ne9WmeuwKnpg/edit?usp=sharing)

## 使用語言及技術

語言框架 : GO-Gin 

核心套件 : 

* Gorm : 與資料庫溝通

* Gorilla Websocket : 建立聊天室的通訊服務及使用者連線狀態的控制

* rs/cors : 網頁跨域連線處理

# 環境建置
## Server 建置
至Azure租用 Ubuntu 18.04的VM

將未來的服務及資料庫等，需對外開放的 Port打開。
至 Azure portal 找到該VM，設定 => 網路 =>新增輸入連接阜規則，至此設定

以下操作皆須連上VM做操作，可透過ssh至VM

1. **安裝 Docker**

`sudo apt-get update`

`sudo apt-get install docker.io`

  - 將使用者加入Docker
  
  `sudo usermod -aG docker “user name”`
  
  - 查看 Docker版本
  
  `docker version`
  
2. **安裝 Postgres**

`docker run -p 5432:5432 --name “postgres-db” -e POSTGRES_PASSWORD=”11111111” -d postgres`

3. **安裝 docker-compose**

`sudo curl -L "https://github.com/docker/compose/releases/download/1.10.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose`

  - 給予 docker-compose權限
  
  `sudo chmod +x /usr/local/bin/docker-compose`
  
  - 查看 docker-compose 版本
  
  `docker-compose --version`
  
4. **安裝Docker Registry 私庫**

建立在5000port

`docker run -d -p 5000:5000 -v /home/xiang/storage:/var/lib/registry --name registry registry:2`
  - 需要修改 client 的 Docker 設定
  
  進入設定檔
  
  `vi /etc/docker/daemon.json`
  
  貼上以下內容 : 
  
  ```
  {
    "live-restore": true,
    "group": "dockerroot",
    "insecure-registries": ["ip:5000"] 
   }
   ````
   
5. **安裝 Drone**

Drone是可以幫我們執行CI/CD的工作

可透過docker-compose去建立drone server 及drone runner

建立一份 docker-compose.yaml，內容如下 : 

```
version: '2'
 
services:
 
  drone-server:
    image: drone/drone:1
    container_name: drone-server
    ports:
      - 80:80
      - 443:443
      
    volumes:
      - /var/lib/drone:/data
      - /var/run/docker.sock:/var/run/docker.sock
    restart: always
    environment:
      - DRONE_SERVER_HOST=你架設那台機器的ip或doman
      - DRONE_SERVER_PROTO=[http/https]
      - DRONE_RPC_SECRET=隨便打，跟drone runner 溝通用的
 
      # GitHub Config
      - DRONE_GITHUB_SERVER=https://github.com
      - DRONE_GITHUB_CLIENT_ID=到 github -> settings -> Developer settings -> OAuthApps 申請
      - DRONE_GITHUB_CLIENT_SECRET=到 github -> settings -> Developer settings -> OAuthApps 申請
 
      - DRONE_LOGS_PRETTY=true
      - DRONE_LOGS_COLOR=true
      - DRONE_AGENTS_ENABLED=true
      - DRONE_TLS_AUTOCERT=true
      - DRONE_USER_FILTER=lnikell
      - DRONE_USER_CREATE=username:wei-xiang-chen,admin:true
      - DRONE_SERVER_NAME=drone
 
  # runner for docker version
  drone-runner:
    image: drone/drone-runner-docker:1
    container_name: drone-runner
    restart: always
    ports:
      - 3000:3000
    depends_on:
      - drone-server
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - DRONE_RPC_HOST=drone-server
      - DRONE_RPC_PROTO=[http/https]
      - DRONE_RPC_SECRET=跟drone一樣
      - DRONE_RUNNER_CAPACITY=3
      - DRONE_RUNNER_NAME=drone-runner
```
在VM上執行

` docker-compose -f docker-compose.yaml up -d`

因為該服務對外開放是80及443 port，可在瀏覽器輸入VM的IP位置或Domain及可連到drone的GUI，接下來會導到github去，需要同意權限即可完成

# CI/CD


我們會透過drone將我們的程式build出執行檔再包成docker image上傳至我們的Docker Registry。當你在你git的history上某個節點上下了V*的tag會執行CI/CD流程

## CI

Build出執行檔

建立一份 Dockerfile，內容如下 : 

```
FROM golang
RUN mkdir -p /club
WORKDIR /club
COPY . .
RUN go mod download
RUN go build -o club
ENTRYPOINT ["./club"]
```

## CD

將執行檔打包成docker image上傳至我們的Docker Registry

建立一份 .drone.yaml，內容如下 : 

```
kind: pipeline
type: docker
name: build

steps:
- name: build
  image: golang:1.12
  commands:
  - go build
  when:
    ref:
    - refs/tags/latest
    - refs/tags/v*

- name: build-images-push-nexus
  image: plugins/docker
  settings: 
    username: 
      from_secret: docker_username #到drone上該專案的Repository裡的settings設定
    password: 
      from_secret: docker_password #到drone上該專案的Repository裡的settings設定
    repo: 104.214.48.227:5000/wei-xiang-chen/club
    registry: 104.214.48.227:5000
    insecure: true
    auto_tag: true
  when:
    ref:
    - refs/tags/latest
    - refs/tags/v*
```

因為drone預設是抓 .drone.yml來執行，所以可以直接改副檔名，或是到drone上該專案的Repository -> settings -> Configuration修改讀取的檔名
