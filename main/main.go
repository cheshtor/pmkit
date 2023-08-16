package main

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	"pmkit/internal/app/controller"
	"pmkit/internal/pkg"
	"pmkit/internal/pkg/database"
	"reflect"
	"sync"
)

// ControllerMethod 路由模型
type ControllerMethod struct {
	// HTTP 请求方法
	MethodType string
	// 请求全路径
	FullPath string
	// 执行路由需要具备的权限
	Permissions []string
	// 路由的业务逻辑
	Handler func(*fiber.Ctx) error
}

// HTTP GET 路由
var getRoutes = make(map[string]*ControllerMethod)

// HTTP POST 路由
var postRoutes = make(map[string]*ControllerMethod)

// 配置文件加载锁
var configLoadLock sync.Mutex

// Configs 配置信息
var Configs *viper.Viper

// 从 Controller 中解析出函数名称与函数实例的对应关系
func analyzeController(controller controller.Controller) map[string]reflect.Value {
	m := make(map[string]reflect.Value)
	t := reflect.TypeOf(controller)
	v := reflect.ValueOf(controller)
	for i := 0; i < t.NumMethod(); i++ {
		typeMethod := t.Method(i)
		methodName := typeMethod.Name
		method := v.MethodByName(methodName)
		m[methodName] = method
	}
	return m
}

// 从 Controller 的函数中解析出路由信息
func analyzeRoute(methods map[string]reflect.Value, commonPrefix string) {
	restController := methods["RestController"]
	call := restController.Call([]reflect.Value{})
	prefix := call[0].Interface().(string)
	delete(methods, "RestController")

	for _, method := range methods {
		call := method.Call([]reflect.Value{})
		methodType := call[0].Interface().(string)
		path := call[1].Interface().(string)
		permissions := call[2].Interface().([]string)
		handler := call[3].Interface().(func(ctx *fiber.Ctx) error)
		fullPath := commonPrefix + prefix + path
		controllerMethod := &ControllerMethod{
			MethodType:  methodType,
			FullPath:    fullPath,
			Permissions: permissions,
			Handler:     handler,
		}
		switch methodType {
		case fiber.MethodGet:
			_, ok := getRoutes[fullPath]
			if ok {
				panic("路由地址重复。" + fullPath)
			}
			getRoutes[fullPath] = controllerMethod
		case fiber.MethodPost:
			_, ok := postRoutes[fullPath]
			if ok {
				panic("路由地址重复。" + fullPath)
			}
			postRoutes[fullPath] = controllerMethod
		default:
			panic("不支持的路由请求方式" + fullPath)
		}
	}
}

func bootstrap() {
	app := fiber.New(fiber.Config{
		// 异常处理
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Errorf("系统发生异常。%s", err.Error())
			// fiber 异常
			var e *fiber.Error
			if errors.As(err, &e) {
				return c.Status(fiber.StatusOK).JSON(pkg.Exception(e.Message, e.Code))
			}
			// 其他异常
			return c.Status(fiber.StatusOK).JSON(pkg.Exception(err.Error(), -1))
		},
	})
	app.Use(fiberRecover.New())

	// 从请求头中获取 JWT Token 并解析用户 ID
	app.Use(func(c *fiber.Ctx) error {
		err := pkg.SetCurrentUserId(c)
		if err != nil {
			return err
		}
		return c.Next()
	})

	api := app.Group("/api")

	methods := analyzeController(&controller.UserController{})
	analyzeRoute(methods, "/api")
	methods = analyzeController(&controller.ProjectController{})
	analyzeRoute(methods, "/api")

	api.Add(fiber.MethodGet, "/*", func(c *fiber.Ctx) error {
		fullPath := pkg.FixRequestURL(c.OriginalURL())
		controllerMethod := getRoutes[fullPath]
		if controllerMethod == nil {
			return errors.New("未知的请求地址")
		}
		return controllerMethod.Handler(c)
	})

	api.Add(fiber.MethodPost, "/*", func(c *fiber.Ctx) error {
		fullPath := pkg.FixRequestURL(c.OriginalURL())
		controllerMethod := postRoutes[fullPath]
		if controllerMethod == nil {
			return errors.New("未知的请求地址")
		}
		return controllerMethod.Handler(c)
	})

	app.Listen(":8080")

}

func readConfig() *viper.Viper {
	configLoadLock.Lock()
	defer configLoadLock.Unlock()
	if Configs == nil {
		config := viper.New()
		config.SetConfigFile("./configs/application.yaml")
		err := config.ReadInConfig()
		if err != nil {
			fmt.Printf("%s", err.Error())
			panic("系统配置文件异常")
		}
		Configs = config
	}
	return Configs
}

func main() {
	config := readConfig()

	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=true&loc=Local", config.GetString("datasource.username"), config.GetString("datasource.password"), config.GetString("datasource.url"), config.GetString("datasource.database"))
	driverName := config.GetString("datasource.driver-name")
	maxOpenConns := config.GetInt("datasource.max-open-connections")
	maxIdleConns := config.GetInt("datasource.max-idle-connections")

	database.NewDBConfig(dsn, driverName, maxOpenConns, maxIdleConns)
	db := database.GetDBInstance()
	defer db.Close()
	bootstrap()
}
