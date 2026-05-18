# SQL 查询优化指南

## 1. 基础优化原则

### 避免 SELECT *

```sql
-- 不推荐
SELECT * FROM users WHERE id = 1;

-- 推荐
SELECT id, name, email FROM users WHERE id = 1;
```

### 使用 LIMIT

```sql
-- 获取前 10 条记录
SELECT id, name FROM users LIMIT 10;

-- 分页查询
SELECT id, name FROM users ORDER BY created_at DESC LIMIT 20 OFFSET 40;
```

## 2. 索引优化

### 创建索引

```sql
-- 单列索引
CREATE INDEX idx_users_email ON users(email);

-- 复合索引
CREATE INDEX idx_orders_user_date ON orders(user_id, created_at DESC);

-- 唯一索引
CREATE UNIQUE INDEX idx_users_email ON users(email);
```

### 索引使用原则

- 索引WHERE子句中经常使用的列
- 索引ORDER BY中的列
- 避免在索引列上使用函数
- 复合索引考虑列的顺序

```sql
-- 复合索引顺序示例
-- 索引 (a, b, c) 支持:
WHERE a = 1              -- 使用索引
WHERE a = 1 AND b = 2   -- 使用索引
WHERE a = 1 AND b = 2 AND c = 3  -- 使用索引

-- 不支持:
WHERE b = 2              -- 不使用索引
WHERE b = 2 AND c = 3   -- 不使用索引
```

## 3. 查询优化

### 使用 EXPLAIN

```sql
EXPLAIN SELECT * FROM users WHERE email = 'test@example.com';
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';
```

### 避免全表扫描

```sql
-- 不推荐：使用函数导致索引失效
SELECT * FROM users WHERE YEAR(created_at) = 2024;

-- 推荐：使用范围查询
SELECT * FROM users WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01';
```

### 使用覆盖索引

```sql
-- 创建覆盖索引
CREATE INDEX idx_users_email_name ON users(email, name);

-- 查询只需索引，不需要回表
SELECT name, email FROM users WHERE email = 'test@example.com';
```

## 4. JOIN 优化

### 小表驱动大表

```sql
-- 使用小表作为驱动表
SELECT * FROM small_table s
INNER JOIN large_table l ON s.id = l.small_id;

-- 性能差异巨大
```

### 避免JOIN过多表

```sql
-- 不推荐：JOIN 太多表
SELECT * FROM a
JOIN b ON a.id = b.a_id
JOIN c ON b.id = c.b_id
JOIN d ON c.id = d.c_id
JOIN e ON d.id = e.d_id;
```

### 使用子查询代替 JOIN

```sql
-- 子查询方式
SELECT * FROM users
WHERE id IN (SELECT user_id FROM orders WHERE total > 1000);

-- 可能比 JOIN 更高效
```

## 5. 分页优化

### 传统分页问题

```sql
-- 性能随页数增加而下降
SELECT * FROM users ORDER BY id LIMIT 1000000, 20;
```

### 优化方案

```sql
-- 使用主键分页
SELECT * FROM users
WHERE id > 1000000
ORDER BY id
LIMIT 20;

-- 基于上一页最后一条记录
SELECT * FROM users
WHERE id > :last_id
ORDER BY id
LIMIT 20;
```

### 游标分页

```sql
SELECT * FROM users
WHERE created_at > :cursor_created_at
  OR (created_at = :cursor_created_at AND id > :cursor_id)
ORDER BY created_at ASC, id ASC
LIMIT 20;
```

## 6. 写入优化

### 批量插入

```sql
-- 不推荐：逐条插入
INSERT INTO users (name, email) VALUES ('A', 'a@test.com');
INSERT INTO users (name, email) VALUES ('B', 'b@test.com');

-- 推荐：批量插入
INSERT INTO users (name, email) VALUES
('A', 'a@test.com'),
('B', 'b@test.com'),
('C', 'c@test.com');
```

### 使用事务

```sql
START TRANSACTION;
INSERT INTO orders (user_id, total) VALUES (1, 100);
INSERT INTO order_items (order_id, product_id, qty) VALUES (LAST_INSERT_ID(), 1, 2);
COMMIT;
```

### 异步写入

```sql
-- 写入临时表
INSERT INTO orders_temp (user_id, total)
VALUES (1, 100), (2, 200);

-- 定期合并到主表
INSERT INTO orders (user_id, total)
SELECT user_id, total FROM orders_temp;
```

## 7. 数据类型优化

### 选择合适的数据类型

```sql
-- 不推荐
CREATE TABLE orders (
    id BIGINT,
    status VARCHAR(10),
    created_at DATETIME
);

-- 推荐
CREATE TABLE orders (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    status ENUM('pending', 'processing', 'completed'),
    created_at DATETIME(3)
);
```

### 使用 NOT NULL

```sql
-- 指定 NOT NULL
ALTER TABLE users MODIFY name VARCHAR(100) NOT NULL;
ALTER TABLE users MODIFY email VARCHAR(255) NOT NULL;
```

## 8. 监控和分析

### 慢查询日志

```sql
-- 查看慢查询
SHOW VARIABLES LIKE 'slow_query_log';
SHOW VARIABLES LIKE 'long_query_time';

-- 分析慢查询
mysqldumpslow -s t /var/log/mysql/slow.log;
```

### 使用性能 Schema

```sql
-- 查看语句统计
SELECT * FROM events_statements_summary ORDER BY avg_timer_exec DESC LIMIT 10;
```

## 9. 常见模式

### 预计算

```sql
-- 创建统计表
CREATE TABLE order_stats (
    date DATE NOT NULL,
    total_orders INT DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    PRIMARY KEY (date)
);

-- 定时更新
INSERT INTO order_stats (date, total_orders, total_amount)
SELECT DATE(created_at), COUNT(*), SUM(total)
FROM orders
WHERE created_at >= CURRENT_DATE
ON DUPLICATE KEY UPDATE
    total_orders = total_orders + VALUES(total_orders),
    total_amount = total_amount + VALUES(total_amount);
```

### 分区表

```sql
CREATE TABLE logs (
    id BIGINT AUTO_INCREMENT,
    created_at DATETIME NOT NULL,
    message TEXT,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (TO_DAYS(created_at)) (
    PARTITION p202401 VALUES LESS THAN (TO_DAYS('2024-02-01')),
    PARTITION p202402 VALUES LESS THAN (TO_DAYS('2024-03-01')),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);
```

## 10. 注意事项

### 避免 OR

```sql
-- 不推荐
SELECT * FROM users WHERE name = 'A' OR email = 'A@test.com';

-- 使用 UNION
SELECT * FROM users WHERE name = 'A'
UNION
SELECT * FROM users WHERE email = 'A@test.com';
```

### 避免 LIKE 开头

```sql
-- 不推荐：无法使用索引
SELECT * FROM users WHERE name LIKE '%john';

-- 推荐：可以使用索引
SELECT * FROM users WHERE name LIKE 'john%';
```

### 谨慎使用 DISTINCT

```sql
-- 性能较低
SELECT DISTINCT email FROM users;
```