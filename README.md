# web3

task3
#####mysql安装在了本地docker里
🔄 安装本地MySQL
docker --version
docker run --name mysql-local -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=testdb -p 3306:3306 -d mysql:8.0
docker ps
🔄 管理MySQL容器
# 启动容器
docker start mysql-local
# 停止容器
docker stop mysql-local
# 查看状态
docker ps -a
# 查看日志
docker logs mysql-local