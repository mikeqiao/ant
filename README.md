# ant
service for game

没啥新鲜的内容

把 leaf 拆分整改了下，轻量级，不支持各种扩展，不追求什么功能都能做，够用就好，简单快捷，少出错
非常懒，懒得解释

1 目前只能用protobuf 做通信
	msgtype uint16类型  建议从1000起 递增
	
	添加功能 
	 func( call *CallInfo){
		
		call.Args //本服务需要的参数
		call.Data //附加信息
		计算得出结果
		res 
		error
		
		//如果需要异步调用其他服务｛
			cb ：= func(in interface{}, e error){
				call.Setresult(in, e)
			}
			mod.Route(ID， cb, in， data) //调用其他服务， 服务号， 回调函数， 参数， 附加数据
		｝
		//不需要的话
		call.Setresult(res, error)
	}
