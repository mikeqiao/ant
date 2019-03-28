
RPC 调用

方法结构 
	有返回类型 func( call *CallInfo) 
	无返回类型 func( call *CallInfo)
	
	callback(in interface{}, e error)
	
	参数类型 返回值类型
	fancId  inID  outID
	
	
	call{
		排队等待
		处理中等待结果
	}