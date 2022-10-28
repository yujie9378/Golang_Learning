package ginlearning

import (
	"context"
	"fmt"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestGinDay01_1(t *testing.T) {
	r := gin.Default()
	/*r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")

	})*/
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")

}
func TestGinDay02_1(t *testing.T) {
	r := gin.Default()
	// 此规则能够匹配/user/john这种格式，但不能匹配/user/ 或 /user这种格式
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	// 但是，这个规则既能匹配/user/john/格式也能匹配/user/john/send这种格式
	// 如果没有其他路由器匹配/user/john，它将重定向到/user/john/
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + "is" + action
		c.String(http.StatusOK, message)
	})
	r.Run(":8081")

}

func TestGinDay02_2(t *testing.T) {
	r := gin.Default()
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")
		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)

	})
	r.Run(":8082")
}

func TestGinDay02_3(t *testing.T) {
	r := gin.Default()
	r.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")
		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	r.Run(":8083")
}

func TestGinDay02_4(t *testing.T) {
	r := gin.Default()
	r.POST("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		fmt.Printf("id:%s;page:%s;name:%s;message:%s", id, page, name, message)

	})
	r.Run(":8084")

}

func TestGinDay02_5(t *testing.T) {
	//http://localhost:8085/upload -F "file=/Users/bytedance/a.txt" -H "Content-Type:multipart/form-data"
	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename)
		c.SaveUploadedFile(file, "/Users/bytedance/work")
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	r.Run(":8085")
}

func TestGinDay02_6(t *testing.T) {
	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]
		for _, file := range files {
			log.Println(file.Filename)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded", len(files)))

	})
	r.Run(":8086")
}

func TestGinDay02_7(t *testing.T) {
	r := gin.New() //变成new后就没有日志
	v1 := r.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}
	v2 := r.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}
	r.Run(":8087")
}

func readEndpoint(context *gin.Context) {
	context.String(http.StatusOK, "readEndpoint")
}

func submitEndpoint(context *gin.Context) {
	context.String(http.StatusOK, "submitEndpoint")
}

func loginEndpoint(context *gin.Context) {
	context.String(http.StatusOK, "loginendpoint")
}

func TestGinDay02_8(t *testing.T) {
	r := gin.New()
	//全局中间件
	r.Use(gin.Logger())
	// Recovery 中间件会 recover 任何 panic。如果有 panic 的话，会写入 500。
	r.Use(gin.Recovery())
	//自定义中间件
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}
	r.Run(":8088")
}
func TestGinDay03_1(t *testing.T) {
	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")

	})
	r.Run(":8080")
}
func TestGinDay03_2(t *testing.T) {
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \" %s %s %s %d %s \" %s\" %s\"\n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC1123),
			params.Method,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)

	}))
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.Run(":8081")
}

type Login struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `json:"password" form:"password" xml:"password" binding:"required"`
}

func TestGinDay03_3(t *testing.T) {
	r := gin.Default()
	r.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})

	})
	r.POST("/loginXML", func(c *gin.Context) {
		var xml Login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	r.POST("/loginForm", func(c *gin.Context) {
		var form Login
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if form.User != "manu" || form.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	r.Run(":8082")

}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)

	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

type Booking struct {
	CheckIn  time.Time `form:"check_in" json:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" json:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func TestGinDay03_4(t *testing.T) {
	//gtfield=CheckIn代表字段值要大于CheckIn字段的值。
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}
	r.POST("/bookable", getBookable)
	//?check_in=2022-12-13&check_out=2022-12-14
	r.Run(":8084")

}

/*
<?xml version="1.0" encoding="UTF-8"?>
<root>

	<name>yj</name>
	<address>345</address>
	<birthday>2000-09-09</birthday>

</root>
*/
type Person struct {
	Name     string    `form:"name" xml:"name" json:"name"`
	Address  string    `form:"address" xml:"address" json:"address"`
	Birthday time.Time `form:"birthday" xml:"birthday" json:"birthday" time_format:"2006-01-02" `
}

func TestGinDay03_5(t *testing.T) {
	r := gin.Default()
	r.Any("/testing", startPage)
	r.Run(":8085")
}
func TestGinDay03_6(t *testing.T) {
	//若要get请求得用这种方式，普通方式不行curl -X GET localhost:8086/testing --data '{"name":"JJ", "address":"xyz"}' -H "Content-Type:application/json"
	r := gin.Default()
	//<?xml version="1.0" encoding="UTF-8"?>
	r.Any("/testing", startPage1)
	r.Run(":8086")
}

type Person1 struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func TestGinDay04_1(t *testing.T) {
	r := gin.Default()
	r.GET("/:name/:id", func(c *gin.Context) {
		var person Person1
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			fmt.Println("err", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": person.Name, "uuid": person.ID})
	})
	r.Run(":8080")
}

type myForm struct {
	Colors []string `form:"colors[]"`
}

func TestGinDay04_2(t *testing.T) {
	r := gin.Default()
	r.POST("/form", func(c *gin.Context) {
		var myform myForm
		c.ShouldBind(&myform)
		c.JSON(200, gin.H{"colors": myform.Colors})

	})
	r.Run(":8081")
}
func TestGinDay04_3(t *testing.T) {
	r := gin.Default()
	r.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})
	r.GET("/moreJSON", func(c *gin.Context) {
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		c.JSON(http.StatusOK, msg)
	})
	r.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})

	})
	r.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})
	r.GET("/someProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		c.ProtoBuf(http.StatusOK, data)
	})
	r.GET("/secureJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}
		c.SecureJSON(http.StatusOK, names)
	})
	r.GET("/JSONP", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}
		c.JSONP(http.StatusOK, data)
	})
	r.GET("/AJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}
		c.AsciiJSON(http.StatusOK, data)

	})
	//通常情况下，JSON会将特殊的HTML字符替换为对应的unicode字符，比如<替换为\u003c，如果想原样输出html，则使用PureJSON
	r.GET("/PureJSON", func(c *gin.Context) {
		c.PureJSON(200, gin.H{
			"html": "<b>Hello,world</b>",
		})
	})
	r.Run(":8082")
}
func TestGinDay04_4(t *testing.T) {
	r := gin.Default()
	r.Static("/a", "./a.html")
	r.Run(":8083")

}
func TestGinDay04_5(t *testing.T) {
	r := gin.Default()
	response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
	r.GET("/someDataFromReader", func(c *gin.Context) {
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		extraHeaders := map[string]string{
			"Content-Disposition": `attachment;filename="gopher.png"`,
		}
		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

	})
	r.Run(":8084")

}

func TestGinDay05_1(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "Main website"})
	})
	r.Run(":8080")
}
func TestGinDay05_2(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	r.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})
	r.Run(":8081")

}
func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d%02d", year, month, day)
}
func TestGinDay06_1(t *testing.T) {
	r := gin.Default()
	r.Delims("{[{", "}]}")
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLFiles("./templates/raw.tmpl")
	r.GET("/raw", func(c *gin.Context) {
		c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
			"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
		})
	})
	r.Run(":8080")
}
func TestGinDay07_1(t *testing.T) {
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
	})
	r.GET("/test1", func(c *gin.Context) {
		c.Request.URL.Path = "/test2"
		r.HandleContext(c)
	})
	r.GET("/test2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	r.Run(":8080")

}
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "12345")
		c.Next()
		latency := time.Since(t)
		log.Print(latency)
		status := c.Writer.Status()
		log.Println(status)

	}
}

func TestGinDay07_2(t *testing.T) {
	r := gin.New()
	r.Use(Logger())
	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)
		log.Println(example)
	})
	r.Run(":8081")

}

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func TestGinDay07_3(t *testing.T) {
	r := gin.Default()
	v := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))
	v.GET("/secrets", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		}
	})
	r.Run(":8082")
}
func TestGinDay07_4(t *testing.T) {
	r := gin.Default()
	r.GET("/long_async", func(c *gin.Context) {
		cCp := c.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			log.Println("Done! in path", cCp.Request.URL)
		}()
	})
	r.GET("/long_sync", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		log.Println("Done! in path", c.Request.URL.Path)
	})
	r.Run(":8083")
}
func TestGinDay07_5(t *testing.T) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	log.Fatal(autotls.Run(r, "yj.com", "yjx.com"))
}

var (
	g errgroup.Group
)

func router01() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			})
	})
	return r

}
func router02() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":  http.StatusOK,
			"error": "Welcome server 02",
		})
	})
	return r
}

func TestGinDay07_6(t *testing.T) {
	server01 := &http.Server{
		Addr:         ":8085",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server02 := &http.Server{
		Addr:         ":8084",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	g.Go(func() error {
		return server01.ListenAndServe()
	})
	g.Go(func() error {
		return server02.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
func TestGinDay07_7(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(15 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	h := &http.Server{
		Addr:    ":8086",
		Handler: r,
	}
	go func() {
		if err := h.ListenAndServe(); err != nil && err != http.ErrServerClosed {

			log.Fatalf("listen : %s \n", err)
			//	time.Sleep(10 * time.Second)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.Shutdown(ctx); err != nil { //链接仍在处理，而h已经等了5秒了，这个时候会返回错误且直接关掉
		log.Fatal("Server Shutdown:", err)
	}
	//没有处理的连接，则直接关掉，关掉连接是无报错的，
	log.Println("Server exiting")

}

type StructA struct {
	FieldA string `form:"field_a"`
}
type StructB struct {
	NestedStruct StructA
	FieldB       string `form:"field_b"`
}
type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}
type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
	}
	FieldD string `form:"field_d"`
}

func GetDataB(c *gin.Context) {
	var b StructB
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStruct,
		"b": b.FieldB,
	})
	fmt.Println("header", c.Request.Header)

}
func GetDataC(c *gin.Context) {
	var b StructC
	c.Bind(&b)
	c.JSON(http.StatusOK, gin.H{
		"x": b.NestedStructPointer,
		"d": b.FieldC,
	})

}
func GetDataD(c *gin.Context) {
	var b StructD
	c.Bind(&b)
	c.JSON(200, gin.H{
		"x": b.NestedAnonyStruct,
		"d": b.FieldD,
	})

}

func TestGinDay07_8(t *testing.T) {
	r := gin.Default()
	r.GET("/getb", GetDataB)
	r.GET("/getc", GetDataC)
	r.GET("/getd", GetDataD)

	r.Run(":8087")
}
func TestGinDay08_1(t *testing.T) {
	r := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	r.POST("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, "foo")
	})
	r.GET("/bar", func(c *gin.Context) {
		c.JSON(http.StatusOK, "bar")
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	// Listen and Server in http://0.0.0.0:8080
	r.Run(":8080")
}
func TestGinDay08_2(t *testing.T) {
	r := gin.Default()
	r.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			cookie = "NOTSET"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}
		c.JSON(http.StatusOK, gin.H{"cookie": cookie})
	})
	r.Run(":8081")
}
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}
func TestGin08_3(t *testing.T) {
	r := setupRouter()
	r.Run(":8082")
}
func TestPingRoute(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
func startPage1(c *gin.Context) {
	var person Person
	if err := c.ShouldBind(&person); err != nil {
		fmt.Println(err)

	}
	fmt.Println("====only bind by query string ---")
	fmt.Println(person.Name)
	fmt.Println(person.Address)
	fmt.Println(person.Birthday)
	c.String(200, "Success")
}

func startPage(c *gin.Context) {
	var person Person
	if c.ShouldBindQuery(&person) == nil {
		fmt.Println("====only bind by query string ---")
		fmt.Println(person.Name)
		fmt.Println(person.Address)

	}
	c.String(200, "Success")
}

func getBookable(context *gin.Context) {
	var b Booking

	if err := context.ShouldBindWith(&b, binding.Query); err != nil {
		fmt.Println("haha")
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	}

}
func analyticsEndpoint(context *gin.Context) {

	context.String(http.StatusOK, "analyticsEndpoint")
}

func AuthRequired() gin.HandlerFunc {

	return func(c *gin.Context) {
		fmt.Println("请求打到了路由authorized")

	}
}

func MyBenchLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		fmt.Printf("有请求打到了/benchmark中,时间是%s,状态码是%d", time.Since(t), context.Writer.Status())
		fmt.Println()

	}

}

func benchEndpoint(context *gin.Context) {
	context.String(http.StatusOK, "benchEndpoint")
}
