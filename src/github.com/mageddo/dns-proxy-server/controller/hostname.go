package controller

import (
	"net/http"
	"encoding/json"
	"github.com/mageddo/dns-proxy-server/events/local"
	"golang.org/x/net/context"
	"github.com/mageddo/log"
	"fmt"
)

func init(){
	Get("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		c := local.GetConfiguration(ctx)
		envName := req.URL.Query().Get("env")
		env, _ := c.GetEnv(envName)
		if env == nil {
			BadRequest(res, fmt.Sprintf("Env %s not found", envName))
			return
		}
		json.NewEncoder(res).Encode(env)
	})

	Get("/hostname/find/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		c := local.GetConfiguration(ctx)
		env := req.URL.Query().Get("env")
		hostname := req.URL.Query().Get("hostname")
		hostnames, err := c.FindHostnameByNameAndEnv(ctx, env, hostname)
		if err != nil {
			BadRequest(res, fmt.Sprintf(err.Error()))
			return
		}
		json.NewEncoder(res).Encode(hostnames)
	})

	Post("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.GetLogger(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/, status=begin, action=create-hostname")
		var hostname local.HostnameVo
		json.NewDecoder(req.Body).Decode(&hostname)
		logger.Infof("m=/hostname/, status=parsed-host, host=%+v", hostname)
		err := local.GetConfiguration(ctx).AddHostname(ctx, hostname.Env, hostname)
		if err != nil {
			logger.Infof("m=/hostname/, status=error, action=create-hostname, err=%+v", err)
			BadRequest(res, err.Error())
			return
		}
		logger.Infof("m=/hostname/, status=success, action=create-hostname")
	})

	Put("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.GetLogger(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/, status=begin, action=update-hostname")
		var hostname local.HostnameVo
		json.NewDecoder(req.Body).Decode(&hostname)
		logger.Infof("m=/hostname/, status=parsed-host, host=%+v", hostname)
		err := local.GetConfiguration(ctx).UpdateHostname(ctx, hostname.Env, hostname)
		if err != nil {
			logger.Infof("m=/hostname/, status=error, action=update-hostname, err=%+v", err)
			BadRequest(res, err.Error())
			return
		}
		logger.Infof("m=/hostname/, status=success, action=update-hostname")
	})

	Delete("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.GetLogger(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/, status=begin, action=delete-hostname")
		var hostname local.HostnameVo
		json.NewDecoder(req.Body).Decode(&hostname)
		logger.Infof("m=/hostname/, status=parsed-host, action=delete-hostname, host=%+v", hostname)
		err := local.GetConfiguration(ctx).RemoveHostnameByEnvAndHostname(ctx, hostname.Env, hostname.Hostname)
		if err != nil {
			logger.Infof("m=/hostname/, status=error, action=delete-hostname, err=%+v", err)
			BadRequest(res, err.Error())
			return
		}
		logger.Infof("m=/hostname/, status=success, action=delete-hostname")
	})
}