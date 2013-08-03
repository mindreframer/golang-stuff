/*Routers*/
package hamster

import (
	"net/http/pprof"
)

//add path handlers
func (s *Server) addHandlers() {

	//add pprof path handlers
	s.route.AddRoute("GET", "/debug/pprof", pprof.Index)
	s.route.AddRoute("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	s.route.AddRoute("GET", "/debug/pprof/profile", pprof.Profile)
	s.route.AddRoute("GET", "/debug/pprof/symbol", pprof.Symbol)

	//Route filters
	s.route.FilterPrefixPath("/", s.baseAuth)
	s.route.FilterPrefixPath("/api/v1/developers/", s.developerAuth)
	s.route.FilterPrefixPath("/api/v1/objects/", s.objectAuth)
	s.route.FilterPrefixPath("/api/v1/files/", s.objectAuth)

	/*Developer*/
	s.route.Post("/api/v1/developers/", s.CreateDev)
	//get a developer objectId, email or username
	s.route.Post("/api/v1/developers/login/", s.LoginDev)
	//login
	s.route.Post("/api/v1/developers/logout/", s.LogoutDev)
	//update developer
	s.route.Put("/api/v1/developers/:objectId", s.UpdateDev)
	//queries
	s.route.Get("/api/v1/developers/:objectId", s.QueryDev)
	//delete object
	s.route.Del("/api/v1/developers/:objectId", s.DeleteDev)

	/*App*/
	s.route.Post("/api/v1/developers/:developerId/apps/", s.CreateApp)
	//get an app
	s.route.Get("/api/v1/developers/apps/:objectId", s.QueryApp)
	//queries
	s.route.Get("/api/v1/developers/:developerId/apps/", s.QueryAllApps)
	//update app
	s.route.Put("/api/v1/developers/apps/:objectId", s.UpdateApp)
	//delete app
	s.route.Del("/api/v1/developers/apps/:objectId", s.DeleteApp)

	/*Object*/
	s.route.Post("/api/v1/objects/", s.CreateObjects)
	s.route.Post("/api/v1/objects/:objectName", s.CreateObject)
	//get an object
	s.route.Get("/api/v1/objects/:objectName/:objectId", s.QueryObject)
	//queries
	s.route.Get("/api/v1/objects/:objectName", s.QueryObjects)
	//update object
	s.route.Put("/api/v1/objects/:objectName/:objectId", s.UpdateObject)
	//delete object
	s.route.Del("/api/v1/objects/:objectName/:objectId", s.DeleteObject)

	/*File*/
	s.route.Post("/api/v1/files/:fileName", s.SaveFile)
	s.route.Get("/api/v1/files/:fileName/:fileId", s.GetFile)

}
