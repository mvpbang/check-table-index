# check-table-index
用于检查oracle表索引存在情况及打印表索引名字

# intent
```
比对表中是否存在对应的索引

索引比对规则：
1、单个字段索
	column_position = 1， 通过
	other:
    	column_position !=1,  不通过  (组合索引的后面位置)
		not exist index，不通过

2、组合字段索引
    all index column, position = 1,  通过
    every one index column map to position,  通过
	other:
		other or exist, 不通过
```


# used
```
./check-table-index.go -h          
  -level string
        log level, value: debug|info (default "info")
注解：
    索引检查是否通过，默认日志级别为info  -level info
    索引检查明细日志(包含检查索引过程日志) -level debug
```

## config.yml
```
# 登陆oracle实例
login:
  address: "1.1.1.1:1521/orcl"  #数据库地址
  username: "xxx"   #用户
  password: "yyy"   #密码

# 检查表是否存在对应索引
check:
  - a:  # 表名
      - b    # 单个字段索引
      - c,d  # 组合索引
  - z:
      - e
      - f,g

注意：多个索引，逗号左右不能有空格
```


## window
- https://www.oracle.com/database/technologies/instant-client/winx64-64-downloads.html
```

C:\green_soft\instantclient_11_2
add path
```