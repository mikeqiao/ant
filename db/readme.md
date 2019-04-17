
data
1 get data
	{
		get <--redis 
		no{
			get <--DB
			yes{
				set -->redis
				return data
			}
			no{
				return nil
			}
		}
		yes{
			return data
		}
	}
	
2 set data
	{
		set-->DB
		yes{
			set-->redis
			yes{
				return true
			}
			no{
				return error
			}
		}
		no{
			return error
		}
		
	}
	

3 DB 
	{
		都是表 操作需要指定：
		{
			表，
			key 
		}	
	}
4 redis 
	{
		存储的结构不同 操作需要指定： 
		{
			结构，
			表，
			key
		}
	}
5 struct 
	{
		name 表名
		{
			key: 	value,
			key1: 	value1,
			...
		}
	}
	
	
// 增 删 改 查  where  value 

1 Select    where [] 

2 del      where []

3 update   where [] value[]

4 insert   value[]

内存修改时候 是变化量  根据变化量 判断哪些字段修改 

或者是记录改变的字段值

包装一个壳子
｛
	set（  struct 变化量）
	{
		记录 改变的字段
	}
	update
	{
		变化的字段
	}
	del
	｛
		key
	｝
	insert
	{
		key
	}
	select
	{
		key
	}
｝
