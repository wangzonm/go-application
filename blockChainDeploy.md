# 区块链部署文档

## 链部署

### 节点1部署
```sh
# 下载节点基础镜像
docker pull buildpack-deps:jessie-curl

# 节点目录
mkdir private_bsc && cd private_bsc

# 创建并编写Dockerfile文件 
touch Dockerfile && vim Dockerfile  
# FROM buildpack-deps:jessie-curl
# RUN mkdir /data
# RUN mkdir -p /usr/local/bin/
# RUN cd /usr/local/bin
# COPY geth /usr/local/bin 
# RUN chmod +x /usr/local/bin/geth
# EXPOSE 50777 30303
# WORKDIR /data
# ENTRYPOINT ["/data/start.sh"]

# Dockerfile 与 geth 放到同目录
# 下载geth 
wget https://github.com/bnb-chain/bsc/releases/download/v1.0.7/geth_linux

# 节点数据挂载目录
mkdir data && cd data 

# 创建并编写创世文件
touch genesis.json && vim genesis.json
# {
#   "config": {
#     "chainId": 6668,
#     "homesteadBlock": 0,
#     "eip150Block": 0,
#     "eip155Block": 0,
#     "eip158Block": 0,
#     "byzantiumBlock": 0,
#     "constantinopleBlock": 0,
#     "petersburgBlock": 0,
#     "istanbulBlock": 0
#   },
#   "nonce": "0x0",
#   "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
#   "difficulty": "0x400",
#   "coinbase": "0xe7C58c28C8802c581Ec6bA40329504Cd4f36a32E",
#   "timestamp": "0x0",
#   "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
#   "gasLimit": "0x47b760",
#   "alloc": {
#     "0xe7C58c28C8802c581Ec6bA40329504Cd4f36a32E":{
#       "balance": "100000000000000000000000000"
#     },
#     "0x07820c5687843a7ee0c89ee960f9ba559851c9ff":{
#       "balance": "100000000000000000000000000"
#     },
#     "0x582A9e054757bEe0b3c80bAA1f4520edAe6dB361":{
#       "balance": "100000000000000000000000000"
#     }
#   },
#   "number":"0x0",
#   "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000"
# }

# 新建并编写Docker容器启动脚本
touch run.sh && vi /private_bsc/data/run.sh
# #!/bin/bash
# docker run -itd --restart=unless-stopped -v /etc/localtime:/etc/localtime --name node_bsc \
#     -v $(pwd):/data \
#     -v $(pwd)/../ethereum:/root/.ethereum \ 
#     -v $(pwd)/../ethash:/root/.ethash \
#     -p 30303:30303 -p 50777:50777 private_bsc:v1.0.7

# 新建并编写geth执行脚本
touch start.sh && vi /private_bsc/data/start.sh
# #!/bin/bash
# set -e

# # Init
# echo ""
# echo "Init geth"
# geth --nousb --datadir /data/node init /data/genesis.json
# sleep 3

# # Start geth
# echo ""
# echo "Start geth"
# geth --datadir /data/node --networkid 6668 --nousb --nodiscover --rpc --rpcapi eth,net,web3,miner,txpool --rpcaddr "0.0.0.0" --rpcport "50777" --syncmode "full" --gcmode "archive"  --allow-insecure-unlock &
# sleep 10

# while true; do
#     sleep 1000000000
# done


# 生成节点镜像
# 需要在与Dockerfile同一级目录下执行此命令
docker build . -t private_bsc:v1.0.7 

# 启动节点容器
sh /private_bsc/data/run.sh
```

### 节点2部署
```sh
# 节点目录
mkdir private_bsc1 && cd private_bsc1

# 编写Dockerfile文件 
touch Dockerfile && vim Dockerfile
# FROM buildpack-deps:jessie-curl
# RUN mkdir /data
# RUN mkdir -p /usr/local/bin/
# RUN cd /usr/local/bin
# COPY geth /usr/local/bin 
# RUN chmod +x /usr/local/bin/geth
# EXPOSE 51777 30304
# WORKDIR /data
# ENTRYPOINT ["/data/start.sh"]

# Dockerfile 与 geth 放到同目录
# 下载geth 
wget https://github.com/bnb-chain/bsc/releases/download/v1.0.7/geth_linux

# 节点数据挂载目录
$ mkdir data && cd data 

# 创建编写创世文件
touch genesis.json && vim genesis.json
# {
#   "config": {
#     "chainId": 6668,
#     "homesteadBlock": 0,
#     "eip150Block": 0,
#     "eip155Block": 0,
#     "eip158Block": 0,
#     "byzantiumBlock": 0,
#     "constantinopleBlock": 0,
#     "petersburgBlock": 0,
#     "istanbulBlock": 0
#   },
#   "nonce": "0x0",
#   "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
#   "difficulty": "0x400",
#   "coinbase": "0xe7C58c28C8802c581Ec6bA40329504Cd4f36a32E",
#   "timestamp": "0x0",
#   "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
#   "gasLimit": "0x47b760",
#   "alloc": {
#     "0xe7C58c28C8802c581Ec6bA40329504Cd4f36a32E":{
#       "balance": "100000000000000000000000000"
#     },
#     "0x07820c5687843a7ee0c89ee960f9ba559851c9ff":{
#       "balance": "100000000000000000000000000"
#     },
#     "0x582A9e054757bEe0b3c80bAA1f4520edAe6dB361":{
#       "balance": "100000000000000000000000000"
#     }
#   },
#   "number":"0x0",
#   "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000"
# }

# 新建并编写Docker容器启动脚本
touch run.sh && vi /private_bsc1/data/run.sh
# #!/bin/bash
# docker run -itd --restart=unless-stopped -v /etc/localtime:/etc/localtime --name node_bsc1 -v $(pwd):/data -p 30304:30304 -p 51777:51777 private_bsc:v1.0.7

# 新建并编写geth执行脚本
touch start.sh && vi /private_bsc1/data/start.sh
# #!/bin/bash
# set -e

# # Init
# echo ""
# echo "Init geth"
# geth --nousb --datadir /data/node init /data/genesis.json
# sleep 3

# # Start geth
# echo ""
# echo "Start geth"
# geth --datadir /data/node --networkid 6668 --nousb --nodiscover --port 30304 --rpc --rpcapi eth,net,web3,miner,txpool --rpcaddr "0.0.0.0" --rpcport "51777" --syncmode "full" --gcmode "archive"  --allow-insecure-unlock &
# sleep 10

# while true; do
#     sleep 1000000000
# done

# 生成节点镜像
docker build . -t private_bsc1:v1.0.7

# 启动节点容器
sh /private_bsc1/data/run.sh
```

### 节点3部署
```sh
# 节点目录
mkdir private_bsc2 && cd private_bsc2

# 新建并编写Dockerfile文件 
touch Dockerfile && vim Dockerfile
# FROM buildpack-deps:jessie-curl
# RUN mkdir /data
# RUN mkdir -p /usr/local/bin/
# RUN cd /usr/local/bin
# COPY geth /usr/local/bin 
# RUN chmod +x /usr/local/bin/geth
# EXPOSE 52777 30305
# WORKDIR /data
# ENTRYPOINT ["/data/start.sh"]

# Dockerfile 与 geth 放到同目录
# 下载geth 
wget https://github.com/bnb-chain/bsc/releases/download/v1.0.7/geth_linux

# 节点数据挂载目录
mkdir data && cd data 

# 创建编写创世文件
touch genesis.json && vim genesis.json
# {
#   "config": {
#     "chainId": 6668,
#     "homesteadBlock": 0,
#     "eip150Block": 0,
#     "eip155Block": 0,
#     "eip158Block": 0,
#     "byzantiumBlock": 0,
#     "constantinopleBlock": 0,
#     "petersburgBlock": 0,
#     "istanbulBlock": 0
#   },
#   "nonce": "0x0",
#   "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
#   "difficulty": "0x400",
#   "coinbase": "0xe7C58c28C8802c581Ec6bA40329504Cd4f36a32E",
#   "timestamp": "0x0",
#   "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
#   "gasLimit": "0x47b760",
#   "alloc": {
#     "0xe7C58c28C8802c581Ec6bA40329504Cd4f36a32E":{
#       "balance": "100000000000000000000000000"
#     },
#     "0x07820c5687843a7ee0c89ee960f9ba559851c9ff":{
#       "balance": "100000000000000000000000000"
#     },
#     "0x582A9e054757bEe0b3c80bAA1f4520edAe6dB361":{
#       "balance": "100000000000000000000000000"
#     }
#   },
#   "number":"0x0",
#   "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000"
# }

# 新建编写Docker容器启动脚本
touch run.sh && vi /private_bsc2/data/run.sh
# #!/bin/bash
# docker run -itd --restart=unless-stopped -v /etc/localtime:/etc/localtime --name node_bsc2 -v $(pwd):/data -p 30305:30305 -p 52777:52777 private_bsc:v1.0.7

# 新建编写geth执行脚本
$ touch start.sh && vi /private_bsc2/data/start.sh
# #!/bin/bash
# set -e

# # Init
# echo ""
# echo "Init geth"
# geth --nousb --datadir /data/node init /data/genesis.json
# sleep 3

# # Start geth
# echo ""
# echo "Start geth"
# geth --datadir /data/node --networkid 6668 --nousb --nodiscover --port 30305 --rpc --rpcapi eth,net,web3,miner,txpool --rpcaddr "0.0.0.0" --rpcport "52777" --syncmode "full" --gcmode "archive"  --allow-insecure-unlock &
# sleep 10

# while true; do
#     sleep 1000000000
# done

# 生成节点镜像
docker build . -t private_bsc2:v1.0.7

# 启动节点容器
sh /private_bsc2/data/run.sh
```

## 节点连接
```sh
# 进入节点1容器
docker exec -it node_bsc bash

# 进入节点1geth客户端
geth attach geth.ipc

# 查看节点信息
> admin.nodeInfo.enode
# out
"enode://948430af04c2715ed1289928571ef0bddaf79476a4b7598f506445d6eb4a9d479016f7f8e5430fd77bca67f6ba86b19ccbd9549d08e9addaf63b40665f57029f@172.17.0.2:30303?discport=0"

# 连接节点2
> admin.addPeer("enode://8726a1f8b70f6be2576da946caa4ac0a291d18c13e98077021344c6c881b60ab056510a856203eee6ab41a342c235873f1e5540182536fb15457769166f7b12a@172.17.0.3:30304?discport=0")

# 连接节点3
> admin.addPeer("enode://7ba3ea6d79ed82314d4acf758745d6115e4b5358d8fd0d7388c6b90926b5b7626cba63cdb182d5d25930ee9040563cd0f0c94d2c2b28d19af7db1f243fd5be3a@172.17.0.4:30305?discport=0")

# 查看节点连接情况
> admin.peers

# 进入节点2容器
docker exec -it node_bsc1 bash

# 进入节点2geth客户端
geth attach geth.ipc

# 连接节点3
> admin.addPeer("enode://7ba3ea6d79ed82314d4acf758745d6115e4b5358d8fd0d7388c6b90926b5b7626cba63cdb182d5d25930ee9040563cd0f0c94d2c2b28d19af7db1f243fd5be3a@172.17.0.4:30305?discport=0")

# 查看节点2连接情况
> admin.peers
```

## 挖矿
```sh
# 进入节点容器
docker exec -it node_bsc bash

# 进入节点geth客户端
geth attach geth.ipc

# 设置矿工地址
> miner.setEtherbase(eth.accounts[0])

# 解锁矿工地址
> personal.unlockAccount(eth.coinbase,"passwd",0)

# 开启挖矿
> miner.start(1)

# 停止挖矿
> miner.stop()
```
