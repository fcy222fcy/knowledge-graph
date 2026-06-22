## SQl通用语法
1. SQl语句可以单行或者多行书写,以分号结尾
2. SQl语句可以使用空格/缩进来增强语句的可读性
3. MySQL数据库的SQl语句不区分大小写,关键字建议使用大写
4. 注释:
```SQl
单行注释: -- 注释内容 或者 # 注释内容
多行注释: /* 注释内容 * /
```
## SQl语句分类
| 分类  | 全称                         | 说明                          |
| --- | -------------------------- | --------------------------- |
| DDL | Data Definition Language   | 数据定义语言,用来定义数据库对象(数据库,表,字段)  |
| DML | Data Manipulation Language | 数据操作语言,用来对数据库表中的数据进行增删改     |
| DQL | Data Query Language        | 数据查询语言,用来查询数据库中表的记录         |
| DCL | Data Control Language      | 数据控制语言,用来创建数据库用户,控制数据库的访问权限 |
## DDl
数据库安装代码
```SQL
CREATE table emp (
 id int comment '编号---排序用',
 workno varchar(10) comment '员工工号',
 name varchar(10) comment'员工姓名',
 gender char(1) comment '性别',
 age tinyint UNSIGNED COMMENT '正常人年龄',
 idcard char(18) comment '身份证号',
 entrydate date comment '入职时间'
) comment '员工信息表';

INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (1, '00001', '柳岩666', '女', 20, '123456789012345678', '北京', '2000-01-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (2, '00002', '张无忌', '男', 18, '123456789012345670', '北京', '2005-09-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (3, '00003', '韦一笑', '男', 38, '123456789712345670', '上海', '2005-08-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (4, '00004', '赵敏', '女', 18, '123456757123845670', '北京', '2009-12-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (5, '00005', '小昭', '女', 16, '123456769012345678', '上海', '2007-07-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (6, '00006', '杨逍', '男', 28, '12345678931234567X', '北京', '2006-01-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (7, '00007', '范瑶', '男', 40, '123456789212345670', '北京', '2005-05-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (8, '00008', '黛绮丝', '女', 38, '123456157123645670', '天津', '2015-05-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (9, '00009', '范凉凉', '女', 45, '123156789012345678', '北京', '2010-04-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (10, '00010', '陈友谅', '男', 53, '123456789012345670', '上海', '2011-01-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (11, '00011', '张士诚', '男', 55, '123567897123465670', '江苏', '2015-05-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (12, '00012', '常遇春', '男', 32, '123446757152345670', '北京', '2004-02-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (13, '00013', '张三丰', '男', 88, '123656789012345678', '江苏', '2020-11-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (14, '00014', '灭绝', '女', 65, '123456719012345670', '西安', '2019-05-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (15, '00015', '胡青牛', '男', 70, '12345674971234567X', '西安', '2018-04-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (16, '00016', '周芷若', '女', 18, null, '北京', '2012-06-01');


```
### 数据库操作
1. 查询所有数据库
```SQl
show databases;
```
2. 查询当前数据库
```SQl
select database();
```
3. 创建数据库
```SQl
create database [if not exists ] 数据库名 [ default charset 字符集 ] [ collate 排序规则 ];

# 创建一个itcast数据库,使用数据库默认的字符集
create database itcast;

# 在同一个数据库服务器中,不能创建两个名称相同的数据库,否则将会报错
# 通过 if not exists 参数来解决这个问题,如果数据库不存在,则创建该数据库,如果存在,则不创建
create database if not exists itcast;

# 创建一个数据库,并使用指定字符集
create database itheima default charset uft8mb4;
```
4. 删除数据库
```SQl
drop database [ if exists ] 数据库名;

# 如果删除一个不存在的数据库,将会报错.此时,可以添加上参数 if exists ,如果数据库存在,再执行删除,否则不执行删除
drop database if exists test;
```
5. 切换数据库
```SQl
use 数据库名 ;

# 例如
use itcast;
```

### 表操作
#### 表操作-查询创建
1. 查询当前数据库的所有表
```SQL

show tables;

# 例如 切换到sys,并查看系统数据库中所有表的结构
use sys;
show tables;

```
2. 查看指定表的结构(表的字段,字段的类型,是否可以为NULL,是否存在默认值等信息)
```SQL
desc/describe 表名;

# 表的结构包括:表的字段,字段的类型,是否可以为NULL,是否存在默认值等信息
```
3. 查询指定表的建表语句
```SQL
show create table 表名;

# 通过这条指令,主要是用来查看建表语句的,而有部分参数我们在创建表的时候,并未指定也会查询到,因为这部分是数据库的默认值,比如存储引擎,字符集等;
```
4. 创建表的结构
```SQL
CREATE TABLE 表名(
	字段1 字段1类型 [COMMENT '字段1注释' ],
	字段2 字段2类型 [COMMENT '字段2注释' ],
	字段3 字段3类型 [COMMENT '字段3注释' ],
	...
	字段n 字段n类型 [COMMENT '字段n注释' ]


)[COMMENT 表注释];

```
> 注意 : 最后一个字段后面没有逗号

比如,创建一个表 tb_user,对应的结构如下,那么建表语句为:

| id  | name | age | gender |
| --- | ---- | --- | ------ |
| 1   | 令狐冲  | 28  | 男      |
| 2   | 风清扬  | 68  | 男      |
| 3   | 东方不败 | 32  | 男      |
```SQL
CREATE table tb_user(
	id int comment '编号',
	name varchar(50) comment '姓名',
	age int '年龄',
	gender varchar(1) '性别'
) comment '用户表';

# varchar(50) 表示最大长度为50个字符;varchar表示可变长度字符串

```
5. 表的复制
```SQL
# 仅仅复制表的结构
CREATE TABLE copy1 LIKE author;

# 复制表的结构和数据
CREATE TABLE copy2 SELECT * FROM author;

# 只复制部分数据
CREATE TABLE copy3 SELECT id,au_name FROM author WHERE nation ='中国';

# 仅仅复制某些字段
CREATE TABLE copy4 SELECT id,au_name FROM author WHERE 1=2;

```
#### 表操作-数据类型
MySQL中的数据类型主要分为三类:数值类型,字符串类型,日期时间类型
1. 数值类型

| 类型          | 大小     | 有符号(SIGNED)范围                                        | 无符号(UNSIGNED)范围                                       | 描述         |
| ----------- | ------ | ---------------------------------------------------- | ----------------------------------------------------- | ---------- |
| TINYINT     | 1byte  | (-128,127)                                           | (0,255)                                               | 小整数值       |
| SMALLINT    | 2bytes | (-32768,32767)                                       | (0,65535)                                             | 大整数值       |
| MEDIUMINT   | 3bytes | (-8388608,8388607)                                   | (0,16777215)                                          | 大整数值       |
| INT/INTEGER | 4bytes | (-2147483648,2147483647)                             | (0,4294967295)                                        | 大整数值       |
| BIGINT      | 8bytes | (-2~63,2~63-1)                                       | (0,2~64-1)                                            | 极大整数值      |
| FLOAT       | 4bytes | (-3.402823466 E+38,3.402823466351 E+38)              | 0和(1.175494351 E-38,3.402823466 E+38)                 | 单精度浮点数值    |
| DOUBLE      | 8bytes | (-1.7976931348623157 E+308,1.7976931348623157 E+308) | 0和(2.2250738585072014 E-308,1.7976931348623157 E+308) | 双精度浮点数值    |
| DECIMAL     |        | 依赖于M(精度)和D(标度)的值                                     | 依赖于M(精度)和D(标度)的值                                      | 小数值(精确定点数) |
```SQL
# 比如
1). 年龄字段--不会出现负数,而且人的年龄不会太大
age tinyint unsigned

2). 分数--总分100分,最多出现一位小数
score double(4,1)
```
2. 字符串类型
![[Pasted image 20251030201429.png]]
char和varchar都可以表述字符串,char是定长字符串,指定长度多长,就占用多少个字符,和字段值的长度无关.而varchar是变长字符串,指定的长度为最大占用长度;相对来说,char的性能会更高一些
```SQL
# 如
1). 用户名 username --->长度不变,最长不超过50
username varchar(50)

2). 性别 gender --->存储值,不是男,就是女
gener char(1)

3). 手机号 phone ---> 固定长度为11
phone char(11)
```
3. 时间和日期类型
![[Pasted image 20251030202147.png]]
```SQL

如
1). 生日字段 brithday
birthday date

2). 创建时间 createtime
createtime datetime
```

#### 表操作-修改
1. 添加字段
```SQL
ALTER TABLE 表名 ADD 字段名 类型(长度) [ COMMENT 注释 ] [ 约束 ];

# 为emp表增加一个新的字段"昵称"为nickname,类型为 varchar(20)
ALTER TABLE emp ADD nickname varchar(20) COMMENT '昵称';
```
2. 修改数据类型或者约束
```SQL
ALTER TABLE 表名 MODIFY 字段名 新的数据类型(长度);

# 修改表时添加约束
ALTER TABLE 表名 MODIFY COLUMN 字段名 字段类型 新约束 ;

# 修改表时添加表级约束
ALTER TABLE 表名 ADD (CONSTRAINT 约束名) 约束类型(字段名);
```
3. 修改==字段名和字段类型==
```SQL
ALTER TABLE 表名 CHANGE 旧字段名 新字段名 类型(长度) [ COMMENT 注释 ] [ 约束 ];

# 将emp表的nickname字段修改为username,类型为varchar(30)
ALTER TABLE emp CHANGE nickname username varchar(30) COMMENT '昵称';

```
4. 删除字段
```SQL
ALTER TABLE 表名 DROP 字段名;

# 将emp表中的字段username删除;
ALTER TABLE emp DROP username;
```
5. 修改==表名==
```SQL
ALTER TABLE 表名 RENAME TO 新表名;

# 将emp表的表名修改为 employee
ALTER TABLE emp RENAME TO employee;
```

#### 表操作-删除
1. 删除表
```SQL
DROP TABLE [ IF EXISTS ] 表名;

# 如果tb_user表存在,则删除tb_user表
DROP TABLE IF EXISTS tb_user;
```
2. 删除指定表,并重新创建表
```SQL
TRUNCATE TABLE 表名;

# 删除表的时候,表中的数据也都会被删除
```

## DML
### 添加数据(INSERT)
1. 给指定字段添加数据
```SQL
INSERT INTO 表名(字段名1,字段名2,...) VALUES(值1,值2,...);

# 例如:给employee表所有的字段添加数据;
insert into employee(id,workno,name,gener,age,idcard,entrydate)
values(1,'1','Itcast','男',10,'123456789012345678','2000-01-01');

```
> 插入数据成功之后 ,查看数据库的数据

方式一:在图形化界面中,双击表名,即可查看表的数据
方式二:`select * from employee;`

2. 给全部字段添加数据
```SQL
INSERT INTO 表名 VALUES (值1,值2,...);

# 插入数据到employee表,具体的SQL如下
insert into employee values(2，'2'，'张无忌'，'男',18，'123456789012345670'，'2005-01-01');
```
3. 批量添加数据
```SQL

INSERT INTO 表名(字段名1,字段名2,...) VALUES(值1,值2,...),(值1,值2,...),(值1,值2,...);

INSERT INTO 表名 VALUES(值1,值2,...),(值1,值2,...),(值1,值2...);

```
>注意:
>插入数据时,指定的字段顺序需要与值的顺序是一一对应的
>字符串和日期型数据应该包含在引号中
>插入的数据大小,应该在字段的规定范围内

### 修改数据(UPDATE)
```SQL
UPDATE 表名 SET 字段名1 = 值1,字段名2 = 值2,....[ WHERE 条件];

#注意:修改语句的条件可以有,也可以没有,如果没有条件,则会修改整张表的所有数据
# 修改id为1的数据,将name修改为itheima
update employee set name = 'itheima' where id = 1 ;

# 修改id为1的数据,将name修改为小昭,gender修改为女
update employee set name = '小昭',gender ='女' where id =1;

# 将所有员工入职日期修改为 2008-01-01
update employee set entrydate='2008-01-01';
```
### 删除数据(DELETE)
```SQL
DELETE FROM 表名 [ WHERE 条件 ];

# 删除gender 为女的员工
delete from employee where gender='女';

# 删除所有员工
delete from employee ;

# truncate table 表名
truncate table boys 

# 区别
1.delete 可以加where条件,但是trancate不能添加
2.truncate删除,效率更高一点
3.加入要删除的表中有自增长列,
如果delete删除之后,再插入数据,自增长列的值从断点开始
如果truncate删除后,再插入数据,自增长列的值从1开始
4.truncate删除没有返回值,delete删除有返回值
5.truncate删除不能回滚,delete删除可以回滚

```
> 注意事项:
> 	DELETE语句的条件可以有,也可以没有,如果没有条件,则会删除整张表的所有数据
> 	DELETE语句语句不能删除某一个字段的值(可以使用UPDATE，将该字段值置为NULL即可)
> 	当进行删除全部数据操作时,datagrip会提示我们,询问是否确认删除,我们直接点击Execute即可

## DQL
### 基本语法
```SQL
SELECT 
	字段列表
FROM
	表名列表
WHERE
	条件列表
GROUP BY
	分组字段列表
HAVING
	分组后条件列表
ORDER BY
	排序字段列表
LIMIT
	分页参数
```
### 数据准备
```SQL
drop table if exists employee;
create table emp(
id int comment '编号',
workno varchar(10) comment '工号',
name varchar(10) comment '姓名',
gender char(1) comment '性别',
age tinyint unsigned comment '年龄',
idcard char(18) comment '身份证号',
workaddress varchar(50) comment '工作地址',
entrydate date comment '入职时间'
)comment '员工表';

INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (1, '00001', '柳岩666', '女', 20, '123456789012345678', '北京', '2000-01-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (2, '00002', '张无忌', '男', 18, '123456789012345670', '北京', '2005-09-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (3, '00003', '韦一笑', '男', 38, '123456789712345670', '上海', '2005-08-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (4, '00004', '赵敏', '女', 18, '123456757123845670', '北京', '2009-12-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (5, '00005', '小昭', '女', 16, '123456769012345678', '上海', '2007-07-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (6, '00006', '杨逍', '男', 28, '12345678931234567X', '北京', '2006-01-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (7, '00007', '范瑶', '男', 40, '123456789212345670', '北京', '2005-05-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (8, '00008', '黛绮丝', '女', 38, '123456157123645670', '天津', '2015-05-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (9, '00009', '范凉凉', '女', 45, '123156789012345678', '北京', '2010-04-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (10, '00010', '陈友谅', '男', 53, '123456789012345670', '上海', '2011-01-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (11, '00011', '张士诚', '男', 55, '123567897123465670', '江苏', '2015-05-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (12, '00012', '常遇春', '男', 32, '123446757152345670', '北京', '2004-02-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (13, '00013', '张三丰', '男', 88, '123656789012345678', '江苏', '2020-11-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (14, '00014', '灭绝', '女', 65, '123456719012345670', '西安', '2019-05-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (15, '00015', '胡青牛', '男', 70, '12345674971234567X', '西安', '2018-04-01');
INSERT INTO emp (id, workno, name, gender, age, idcard, workaddress, entrydate) VALUES (16, '00016', '周芷若', '女', 18, null, '北京', '2012-06-01');

```
### 基本查询
1. 查询多个字段
```SQL
SELECT 字段1,字段2,字段3 ... FROM 表名;

SELECT * FROM 表名
# *代表查询所有字段,在实际开发中尽量少用(不直观,影响效率)
```
2. 查询字段并设置别名
```SQL
SELECT 字段1 [ AS 别名1 ],字段2 [ AS 别名 ] ... FROM 表名;

SELECT 字段1 [别名1],字段2 [别名2] ... FROM 表名;
```
3. 去除重复记录
```SQL
SELECT DISTINCT 字段列表 FROM 表名;
```
### 条件查询(where)
1. 语法
```SQL
SELECT 字段列表 FROM 表名 WHERE 条件列表;
```
2. 条件
常用的比较运算符如下:

| 比较运算符              | 功能                     |
| ------------------ | ---------------------- |
| >                  | 大于                     |
| >=                 | 大于等于                   |
| <                  | 小于                     |
| <=                 | 小于等于                   |
| =                  | 等于                     |
| <>或!=              | 不等于                    |
| BETWEEN ...AND ... | 在某个范围之内(含最小,最大值)       |
| IN(...)            | 在in之后的列表中的值,多选一        |
| LIKE 占位符           | 模糊匹配(_匹配单个字段,%匹配任意个字符) |
| IS NULL            | 是NULL                  |
常用的逻辑运算符如下:

| 逻辑运算符   | 功能             |
| ------- | -------------- |
| ADD或&&  | 并且(多个条件同时成立)   |
| OR或\|\| | 或者(多个条件任意一个成立) |
| NOT或!   | 非,不是           |


### 聚合函数(count,max,min,avg,sum)
将一列数据作为一个整体,进行纵向计算
1. 常见的聚合函数

| 函数    | 功能   |
| ----- | ---- |
| count | 统计数量 |
| max   | 最大值  |
| min   | 最小值  |
| avg   | 平均值  |
| sum   | 求和   |
2. 语法
```SQL
SELECT 聚合函数(字段列表) FROM 表名;

# NULL值是不参与所有聚合函数运算的

# 统计企业员工数量
select count(*) from emp;---统计总记录数
select count(idcard) from emp;---统计idcard不为null的字段数

# 统计该企业员工的平均年龄
select avg(age) from emp;

# 统计最大年龄
select max(age) from emp;

# 统计最小
select min(age) from emp;

# 统计西安地区员工的年龄之和
select sum(age) from emp where workaddress='西安';
```

>注意区分`count(*),count(1),count(字段)

### 分组查询(group by)

1. 语法
```SQL

SELECT 字段列表 FROM 表名 [ WHERE 条件 ] GROUP BY 分组字段名 [ HAVING 分组后过滤条件 ];

# 根据性别分组,统计男性员工和 女性员工的数量
select gender, count(*) from emp group by gender;

# 根据性别分组,统计男性员工和女性员工的平均年龄
select gender, avg(age) from emp group by gender;

# 查询年龄小于45的员工,并根据工作地址分组,获取员工数量大于等于3的工作地址
select workaddress,coount(*) addresss_count from emp where age< 45 group by workaddress having address_count >= 3;
```
2. `where 和 having的区别`
	1. 执行时机不同:where 是分组之前进行过滤,不满足where条件,不参与分组;而having是分组之后对结果进行过滤
	2. 判断条件不同:where 不能对聚合函数进行判断,而having可以

> 分组之后,查询的字段一般为聚合函数和分组字段,查询其他字段无任何意义
> 执行顺序: where > 聚合函数 > having
> 支持多字段分组,具体语法为:group by columnA,columnB
### 排序查询(order by)
1. 语法
```SQL
SELECT 字段列表 FROM 表名 ORDER BY 字段1 排序方式1,字段2 排序方式2;

# 根据年龄对公司的员工进行升序排序
select * from emp order by age asc;
select * from emp order by age; 

# 根据入职时间,对员工进行降序排序
select * from emp order by entrydate desc;

# 根据年龄对公司的员工进行升序排序,年龄相同,再按照入职时间进行降序排序
select * from emp order by age asc,entrydate desc;
```

>ASC:升序(默认值)
>DESC:降序
>注意:
>如果是升序,可以不指定排序方式ASC;
>如果是多字段排序,当第一个字段值相同时,才会根据第二个字段进行排序;
### 分页查询(limit)
1. 语法
```SQL

# 语法
SELECT 字段列表 FROM 表名 LIMIT 起始索引,查询记录数;

# 案例

# 查询第1页员工数据,每页展示10条记录
select * from emp limit 0,10;
select * from emp limit 10;

# 查询第2页员工数据,每页展示10条记录---->(页码-1)*页展示记录数
select * from limit 10,10; 

```

> 注意事项:
> 起始索引从0开始,起始索引=(查询页码-1)\*每页显示记录数
> 分页查询是数据库的方言,不同的数据库有不同的实现,MySQL中是LIMIT
> 如果查询的是第一页数据,起始索引可以省略,简写为 limit  10.

### DQL语句的执行顺序

| 编写顺序     | 执行顺序     |
| -------- | -------- |
| SELECT   | FROM     |
| FROM     | WHERE    |
| WHERE    | GROUP BY |
| GROUP BY | HAVING   |
| HAVING   | SELECT   |
| ORDER BY | ORDER BY |
| LIMIT    | LIMIT    |

## DCL
### 管理用户
1. 查询用户
```SQL
select * from mysql.user;
```
![[PixPin_2025-11-05_16-50-00.png]]
> Host代表当前用户访问的主机,如果为localhost,仅代表只能够在当前本机,是==不可以远程访问的==;user代表的是该数据库的用户名;在mysql中需要通过Host和User来唯一标识一个用户

2. 创建用户
```SQL
CREATE USER '用户名'@'主机名' IDENTIFIED BY '密码';
```
3. 修改用户密码
```SQL
ALTER USER '用户名'@'主机名' IDENTIFIED WITH mysql_native_password BY '新密码';
```
4. 删除用户
```SQL
DROP USER '用户名'@'主机名';
```

- 在mysql中需要通过`用户名@主机名`的方式来唯一标识一个用户
- 主机名可以用`%`通配
5. 案例
```SQL
# 创建用户itcast,只能够在当前主机localhost访问,密码123456;
create user 'itcast'@'localhost' identified by '123456';

# 创建用户heima,可以在任意主机上访问数据库,密码123456
create user 'heima'@'%' identified by '123456';

# 修改用户heima的访问密码为1234
alter user 'heima'@'%' identified with mysql_native_password by '1234';

# 删除itcast@localhost用户
drop user 'itcast'@'localhost';
```
### 权限控制
mysql中定义了很多种权限

| 权限                  | 说明         |
| ------------------- | ---------- |
| ALL, ALL PRIVILEGES | 所有权限       |
| SELECT              | 查询数据       |
| INSERT              | 插入数据       |
| UPDATE              | 修改数据       |
| DELETE              | 删除数据       |
| ALTER               | 修改表        |
| DROP                | 删除数据库/表/视图 |
| CREATE              | 创建数据库/表    |
```SQL
# 查询权限
SHOW GRANTS FOR '用户名'@'主机名';

# 授予权限
GRANT 权限列表 ON 数据库名.表名 TO '用户名'@'主机名';

# 撤销权限
REVOKE 权限列表 ON 数据库名.表名 FROM '用户名'@'主机名';



# 查询 'heima'@'%' 用户的权限
show grants for 'heima'@'%';

# 授予 'heima'@'%' 用户itcast数据库所有表的所有操作权限
grant all on itcast.* to 'heima'@'%';

# 撤销 'heima'@'%' 用户的itcast数据库的所有权限
revoke all on itcast.* from 'heima'@'%' ;

```
> 注意事项:
> 多个权限之间,使用逗号分隔
> 授权时,数据库名和表名可以使用 * 进行通配,代表所有;








































