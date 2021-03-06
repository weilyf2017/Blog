1.查看当前数据库的事务隔离级别：

```
show variables like 'tx_isolation';
```

2.通过检查InnoDB_row_lock 状态变量分析系统上的行锁的争夺情况 

```
mysql> show status like 'innodb_row_lock%';
+-------------------------------+-------+
| Variable_name                 | Value |
+-------------------------------+-------+
| Innodb_row_lock_current_waits | 0     |
| Innodb_row_lock_time          | 0     |
| Innodb_row_lock_time_avg      | 0     |
| Innodb_row_lock_time_max      | 0     |
| Innodb_row_lock_waits         | 0     |
+-------------------------------+-------+
innodb_row_lock_current_waits: 当前正在等待锁定的数量
innodb_row_lock_time: 从系统启动到现在锁定总时间长度；非常重要的参数，
innodb_row_lock_time_avg: 每次等待所花平均时间；非常重要的参数，
innodb_row_lock_time_max: 从系统启动到现在等待最常的一次所花的时间；
innodb_row_lock_waits: 系统启动后到现在总共等待的次数；非常重要的参数。直接决定优化的方向和策略。
```

3.可以通过检查table_locks_waited 和 table_locks_immediate 状态变量分析系统上的表锁定：show status like 'table_locks%';

```
mysql> show status like 'table_locks%';
+----------------------------+-------+
| Variable_name              | Value |
+----------------------------+-------+
| Table_locks_immediate      | 104   |
| Table_locks_waited         | 0     |
+----------------------------+-------+
table_locks_immediate: 表示立即释放表锁数。
table_locks_waited: 表示需要等待的表锁数。此值越高则说明存在着越严重的表级锁争用情况。
```

4. 查看一些超时时间的设置 
```
show variables like "%timeout%";
+-------------------------------+----------+
| Variable_name                 | Value    |
+-------------------------------+----------+
| connect_timeout               | 30       |
| delayed_insert_timeout        | 300      |
| have_statement_timeout        | YES      |
| innodb_flush_log_at_timeout   | 1        |
| innodb_lock_wait_timeout      | 50       |
| innodb_rollback_on_timeout    | OFF      |
| interactive_timeout           | 7200     |
| lock_wait_timeout             | 31536000 |
| net_read_timeout              | 30       |
| net_write_timeout             | 60       |
| rds_trx_changes_idle_timeout  | 0        |
| rds_trx_idle_timeout          | 0        |
| rds_trx_readonly_idle_timeout | 0        |
| rpl_semi_sync_master_timeout  | 1000     |
| rpl_stop_slave_timeout        | 31536000 |
| slave_net_timeout             | 60       |
| wait_timeout                  | 7200     |

```
5.查看最大连接数  

```
show variables like 'max_connections';
```

6.

```
mysql> show status like 'thread%';
+——————-+——-+
| Variable_name | Value |
+——————-+——-+
| Threads_cached | 0 | <—当前被缓存的空闲线程的数量
| Threads_connected | 1 | <—正在使用（处于连接状态）的线程
| Threads_created | 1498 | <—服务启动以来，创建了多少个线程
| Threads_running | 1 | <—正在忙的线程（正在查询数据，传输数据等等操作）
+——————-+——-+
```

7.关于数据库

```
show databases; 显示所有数据库

use databasename; 选择数据库

show tables; 显示表

flush privileges //刷新数据库

use dbname； 打开数据库：

DELETE FROM table:删除表数据  truncate table;//不保留表记录值

DROP TABLE table:删除数据表

show create table from  :显示表结构详细信息

```

8.查看所有数据库容量大小

```
select 
table_schema as '数据库',
sum(table_rows) as '记录数',
sum(truncate(data_length/1024/1024, 2)) as '数据容量(MB)',
sum(truncate(index_length/1024/1024, 2)) as '索引容量(MB)'
from information_schema.tables
group by table_schema
order by sum(data_length) desc, sum(index_length) desc;
```

9.查看所有数据库各表容量大小

```
select 
table_schema as '数据库',
table_name as '表名',
table_rows as '记录数',
truncate(data_length/1024/1024, 2) as '数据容量(MB)',
truncate(index_length/1024/1024, 2) as '索引容量(MB)'
from information_schema.tables
order by data_length desc, index_length desc;
```



10.查看指定数据库容量大小

例：查看mysql库容量大小

```
select 
table_schema as '数据库',
sum(table_rows) as '记录数',
sum(truncate(data_length/1024/1024, 2)) as '数据容量(MB)',
sum(truncate(index_length/1024/1024, 2)) as '索引容量(MB)'
from information_schema.tables
where table_schema='mysql';
```




11.查看指定数据库各表容量大小

例：查看mysql库各表容量大小

```
select 
table_schema as '数据库',
table_name as '表名',
table_rows as '记录数',
truncate(data_length/1024/1024, 2) as '数据容量(MB)',
truncate(index_length/1024/1024, 2) as '索引容量(MB)'
from information_schema.tables
where table_schema='mysql'
order by data_length desc, index_length desc;
```
