package environment

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Env struct {
	AppName    string
	AppVersion string
	LogLevel   string
	AdminName  string

	RestListenAddr      string
	WebsocketListenAddr string

	DSN string

	BypassAuth bool
}

func Get() (Env, error) {
	var err error

	if err = load(); err != nil {
		return Env{}, err
	}

	var appName string
	if os.Getenv("APP_NAME") == "" {
		return Env{}, fmt.Errorf("app name is required")
	} else {
		appName = os.Getenv("APP_NAME")
	}

	var appVersion string
	if os.Getenv("APP_VERSION") == "" {
		return Env{}, fmt.Errorf("app version is required")
	} else {
		appVersion = os.Getenv("APP_VERSION")
	}

	var logLevel string
	if os.Getenv("LOG_LEVEL") == "" {
		logLevel = "INFO"
	} else {
		logLevel = os.Getenv("LOG_LEVEL")
	}

	var adminName string
	if os.Getenv("ADMIN_NAME") == "" {
		adminName = "admin"
	} else {
		adminName = os.Getenv("ADMIN_NAME")
	}

	var restListenAddr string
	if os.Getenv("REST_LISTEN_ADDR") == "" {
		restListenAddr = "localhost:8090"
	} else {
		restListenAddr = os.Getenv("REST_LISTEN_ADDR")
	}

	var websocketListenAddr string
	if os.Getenv("WEBSOCKET_LISTEN_ADDR") == "" {
		websocketListenAddr = "localhost:8091"
	} else {
		websocketListenAddr = os.Getenv("WEBSOCKET_LISTEN_ADDR")
	}

	var dsn string
	if os.Getenv("DSN") == "" {
		return Env{}, fmt.Errorf("dsn is required")
	} else {
		dsn = os.Getenv("DSN")
	}

	var bypassAuth bool
	if os.Getenv("BYPASS_AUTH") == "" {
		bypassAuth = false
	} else {
		bypassAuth, err = strconv.ParseBool(os.Getenv("BYPASS_AUTH"))
		if err != nil {
			return Env{}, err
		}
	}

	return Env{
		AppName:             appName,
		AppVersion:          appVersion,
		LogLevel:            logLevel,
		AdminName:           adminName,
		RestListenAddr:      restListenAddr,
		WebsocketListenAddr: websocketListenAddr,
		DSN:                 dsn,
		BypassAuth:          bypassAuth,
	}, nil
}

var (
	envLoaded = false
	mu        sync.Mutex
)

func load() error {
	mu.Lock()
	defer mu.Unlock()
	if envLoaded {
		return nil
	}

	//定位到根目录的.env文件
	_, f, _, _ := runtime.Caller(0) // 当前执行的文件，即此文件environment/env.go
	basepath := filepath.Dir(f)
	envFile := path.Join(basepath, "../.env")
	_ = godotenv.Load(envFile)

	//从os读取环境变量，忽略error
	_ = godotenv.Load()
	envLoaded = true
	return nil
}
