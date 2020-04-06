# ScarletWaf_backend
The backend of scarlet

### 数据库

使用数据库Mysql和Redis
- Mysql负责存储用户数据
- Redis负责存储配置和规则

### TODO

- 配置文件读取（现在还有BUG总之没实现
- <del>各类基本CURD API</del>
- 日志处理。waf_core会产生日志到一个 per user的List中，后端读取、存储。也要考虑编排多用户的文件路径、考虑如何返回给前端


### ChangeLog


#### 4-6

- 增加表示运行环境的变量：`DEVELOP`；
- 暂时修改URI的`CC_rate`为不可见；
- 如果不给`Flag`，通过是否给`uri_id`来判断时BASE还是CUSTOM的规则；
- 补上`deleteRule`中遗漏的表单验证；
- 信息补全（查库操作）从`handler`移到`Service`中；
- 用户注册增加是否已存在的判断；
- `/user/rule/get` `Endpoint` 对应 `service/ruleService.go` `GetRulePage` 信息补全；

#### 4-X<=4-2

- dgrijalva/jwt-go :引入JWT做session管理
- 原redis库处理zset报错 使用新库 gomodule/redigo
- SESSION中间件设计
- 修了个关于SESSION不能提早Abort的BUG；
- 增加了Rule、Switch相关API；
- 修改接口为Restful变体（因为rule不存在合适的主键

#### Long Time Ago ~

- 写了代码基本框架和简单的CURD API
- gin-gonic/gin ：基本web框架
- go-ozzo/ozzo-validation ：表单验证 
- swaggergo/gin-swagger : 代码文档生成
- spf13/viper : 配置文件读取
- jinzhu/gorm : 数据库操作和orm
- go-redis/redis : redis操作
