package file

import "github.com/shrewx/ginx"

// Router 文件模块路由分组
var Router = ginx.NewRouter(ginx.Group("files"))

func init() {
	// 注册文件相关接口
	// Router.Register(&UploadFile{})    // POST /files
	// Router.Register(&DownloadFile{})  // GET /files/:id
	// Router.Register(&DeleteFile{})    // DELETE /files/:id
}