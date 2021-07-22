package oidb_0x5eb

// http://oidb2.server.com/metronic/html/protocolfile/protocolList.php?appid=1515
const (
	RetCodeOk             = uint32(0)
	RetCodeInvalidPackage = uint32(0x1)  // 包格式非法
	RetCodeAddressing     = uint32(0x40) // 寻址后端失败
	RetCodeAccess         = uint32(0x41) // 访问后端失败
	RetCodeInvalidUin     = uint32(0x42) // 号码无效
	RetCodeFieldForbidden = uint32(0x43) // 字段没有申请权限，请参考【调用逻辑】进行申请
	RetCodeInvalidSkey    = uint32(0x97) // skey无效
)
