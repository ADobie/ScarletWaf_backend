# ScarletWaf_backend
The backend of scarlet

### 数据库

使用数据库Mysql和Redis
- Mysql负责存储用户数据
- Redis负责存储配置和规则

### TODO

- 配置文件读取（现在还有BUG总之没实现
- 各类基本CURD API
- 日志处理。waf_core会产生日志到一个 per user的List中，后端读取、存储。也要考虑编排多用户的文件路径、考虑如何返回给前端
