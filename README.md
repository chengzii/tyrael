@2013-12-16
@love 橙子
@haddis_123@163.com
说明: 
	本工具是一个简单的linux下进程cpu/mem动态展示平台
	
    bin/ 目录  dyrael 为可执行程序.
	conf/ 目录  public.conf为平台配置文件,包含以下内容:
		[system]
		process = walk   	//侦听进程名称，多个用‘,’分割
		port = 5436      	//侦听进程端口，多个用‘,’分割     在存在同名process时，可配合port进行区分
		aspect = cpu,mem 	//侦听类型，单选或两者都选
		intervaltime = 15 	//写数据库的轮询时间 
		totaltime = 259200 	//写入数据库的保存时间
		issave = 0 			//是否写入数据库 @目前还不完整支持数据库存取,敬请等待

		[db]					//目前可写入的数据仅支持 redis
		ip = 127.0.0.1			//连接数据库的IP
		port = 6379				//连接数据库的PORT
		db = 0					//连接数据库的NAME
		password = chengzi		//连接数据库的AUTH
		maxpoolsize = 500000	//连接数据库的最大数据池

		[http]
		port = 12345			//服务启动后的侦听端口
	src/ 目录	可执行程序的原始代码（go语言）
	log/ 目录   执行过程记录的日志
	tpl/ 目录   html文件所在目录
	
××注意：
	平台目前仅支持linux环境，centos6.2上完成测试，其它os上可能存在某种预期外结果；
	启动后，可访问  
		$IP:12345/show 查看进程的cpu/mem的当前使用状态的页面，每15秒刷新一次；
		$IP:12345/line 查看一段时间内的进程的cpu/mem的使用率，数据间隔等于intervaltime参数值(注:数据是存储在内存中,进程结束后即全部清除,请慎重!)；
	进程的cpu/mem信息为top命令返回的及时结果，请自行区别与ps结果不同；
	平台结果展示基于google chart api实现，服务所在机器需要联网条件，由于某种“不可抗拒的原因”，展现可能存在一定的延迟或失败；

Todo:
	1	多服务结果展示
	2	页面美化
	3	连续时间的性能展示   -- 完成60%
	4	支持更多数据存储方式 
	5   性能优化			 
	