package spider

import (
	"net/http"
	"sync"
	"github.com/alezh/novel/system/http/request"
	"github.com/PuerkitoBio/goquery"
	"github.com/alezh/novel/system/data"
	"time"
	"unsafe"
	"bytes"
	"mime"
	"strings"
	"io"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"github.com/henrylee2cn/pholcus/common/util"
)

type Context struct {
	spider   *Spider           // 规则
	Request  *request.Request  // 原始请求
	Response *http.Response    // 响应流，其中URL拷贝自*request.Request
	text     []byte            // 下载内容Body的字节流格式
	dom      *goquery.Document // 下载内容Body为html时，可转换为Dom的对象
	items    []data.Data   // 存放以文本形式输出的结果数据
	err      error             // 错误标记
	sync.Mutex
}

var (
	contextPool = &sync.Pool{
		New: func() interface{} {
			return &Context{
				items: []data.Data{},
			}
		},
	}
)

func GetContext(sp *Spider, req *request.Request) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.spider = sp
	ctx.Request = req
	return ctx
}

func PutContext(ctx *Context) {
	if ctx.Response != nil {
		ctx.Response.Body.Close() // too many open files bug remove
		ctx.Response = nil
	}
	ctx.items = ctx.items[:0]
	ctx.spider = nil
	ctx.Request = nil
	ctx.text = nil
	ctx.dom = nil
	ctx.err = nil
	contextPool.Put(ctx)
}

func (self *Context) SetResponse(resp *http.Response) *Context {
	self.Response = resp
	return self
}

// 标记下载错误。
func (self *Context) SetError(err error) {
	self.err = err
}

//**************************************** Set与Exec类公开方法 *******************************************\\
// 生成并添加请求至队列。
// Request.Url与Request.Rule必须设置。
// Request.Spider无需手动设置(由系统自动设置)。
// Request.EnableCookie在Spider字段中统一设置，规则请求中指定的无效。
// 以下字段有默认值，可不设置:
// Request.Method默认为GET方法;
// Request.DialTimeout默认为常量request.DefaultDialTimeout，小于0时不限制等待响应时长;
// Request.ConnTimeout默认为常量request.DefaultConnTimeout，小于0时不限制下载超时;
// Request.TryTimes默认为常量request.DefaultTryTimes，小于0时不限制失败重载次数;
// Request.RedirectTimes默认不限制重定向次数，小于0时可禁止重定向跳转;
// Request.RetryPause默认为常量request.DefaultRetryPause;
// Request.DownloaderID指定下载器ID，0为默认的Surf高并发下载器，功能完备，1为PhantomJS下载器，特点破防力强，速度慢，低并发。
// 默认自动补填Referer。
func (self *Context) AddQueue(req *request.Request) *Context {
	// 若已主动终止任务，则崩溃爬虫协程
	self.spider.tryPanic()

	err := req.
		SetSpiderName(self.spider.GetName()).
		SetEnableCookie(self.spider.GetEnableCookie()).
		Prepare()

	if err != nil {
		//logs.Log.Error(err.Error())
		return self
	}

	// 自动设置Referer
	if req.GetReferer() == "" && self.Response != nil {
		req.SetReferer(self.GetUrl())
	}

	self.spider.RequestPush(req)
	return self
}

// 在请求中保存临时数据。
func (self *Context) SetTemp(key string, value interface{}) *Context {
	self.Request.SetTemp(key, value)
	return self
}

func (self *Context) SetUrl(url string) *Context {
	self.Request.Url = url
	return self
}

func (self *Context) SetReferer(referer string) *Context {
	self.Request.Header.Set("Referer", referer)
	return self
}
// 调用指定Rule下辅助函数AidFunc()。
// 用ruleName指定匹配的AidFunc，为空时默认当前规则。
func (self *Context) Aid(aid map[string]interface{}, ruleName ...string) interface{} {
	// 若已主动终止任务，则崩溃爬虫协程
	self.spider.tryPanic()

	_, rule, found := self.getRule(ruleName...)
	if !found {
		if len(ruleName) > 0 {
			//logs.Log.Error("调用蜘蛛 %s 不存在的规则: %s", self.spider.GetName(), ruleName[0])
		} else {
			//logs.Log.Error("调用蜘蛛 %s 的Aid()时未指定的规则名", self.spider.GetName())
		}
		return nil
	}
	if rule.AidFunc == nil {
		//logs.Log.Error("蜘蛛 %s 的规则 %s 未定义AidFunc", self.spider.GetName(), ruleName[0])
		return nil
	}
	return rule.AidFunc(self, aid)
}

// 解析响应流。
// 用ruleName指定匹配的ParseFunc字段，为空时默认调用Root()。
func (self *Context) Parse(ruleName ...string) *Context {
	// 若已主动终止任务，则崩溃爬虫协程
	self.spider.tryPanic()

	_ruleName, rule, found := self.getRule(ruleName...)
	if self.Response != nil {
		self.Request.SetRuleName(_ruleName)
	}
	if !found {
		self.spider.RuleTree.Root(self)
		return self
	}
	if rule.ParseFunc == nil {
		//logs.Log.Error("蜘蛛 %s 的规则 %s 未定义ParseFunc", self.spider.GetName(), ruleName[0])
		return self
	}
	rule.ParseFunc(self)
	return self
}

// 设置自定义配置。
func (self *Context) SetKeyin(keyin string) *Context {
	self.spider.SetKeyin(keyin)
	return self
}

// 设置采集上限。
func (self *Context) SetLimit(max int) *Context {
	self.spider.SetLimit(int64(max))
	return self
}

// 自定义暂停区间(随机: Pausetime/2 ~ Pausetime*2)，优先级高于外部传参。
// 当且仅当runtime[0]为true时可覆盖现有值。
//func (self *Context) SetPausetime(pause int64, runtime ...bool) *Context {
//	self.spider.SetPausetime(pause, runtime...)
//	return self
//}

// 设置定时器，
// @id为定时器唯一标识，
// @bell==nil时为倒计时器，此时@tol为睡眠时长，
// @bell!=nil时为闹铃，此时@tol用于指定醒来时刻（从now起遇到的第tol个bell）。
func (self *Context) SetTimer(id string, tol time.Duration, bell *Bell) bool {
	return self.spider.SetTimer(id, tol, bell)
}

// 启动定时器，并获取定时器是否可以继续使用。
func (self *Context) RunTimer(id string) bool {
	return self.spider.RunTimer(id)
}

// 重置下载的文本内容，
func (self *Context) ResetText(body string) *Context {
	x := (*[2]uintptr)(unsafe.Pointer(&body))
	h := [3]uintptr{x[0], x[1], x[1]}
	self.text = *(*[]byte)(unsafe.Pointer(&h))
	self.dom = nil
	return self
}

//**************************************** 私有方法 *******************************************\\

// 获取规则。
func (self *Context) getRule(ruleName ...string) (name string, rule *Rule, found bool) {
	if len(ruleName) == 0 {
		if self.Response == nil {
			return
		}
		name = self.GetRuleName()
	} else {
		name = ruleName[0]
	}
	rule, found = self.spider.GetRule(name)
	return
}

// GetHtmlParser returns goquery object binded to target crawl result.
func (self *Context) initDom() *goquery.Document {
	if self.text == nil {
		self.initText()
	}
	var err error
	self.dom, err = goquery.NewDocumentFromReader(bytes.NewReader(self.text))
	if err != nil {
		panic(err.Error())
	}
	return self.dom
}

// GetBodyStr returns plain string crawled.
func (self *Context) initText() {
	var err error

	// 采用surf内核下载时，尝试自动转码
	if self.Request.DownloaderID == request.SURF_ID {
		var contentType, pageEncode string
		// 优先从响应头读取编码类型
		contentType = self.Response.Header.Get("Content-Type")
		if _, params, err := mime.ParseMediaType(contentType); err == nil {
			if cs, ok := params["charset"]; ok {
				pageEncode = strings.ToLower(strings.TrimSpace(cs))
			}
		}
		// 响应头未指定编码类型时，从请求头读取
		if len(pageEncode) == 0 {
			contentType = self.Request.Header.Get("Content-Type")
			if _, params, err := mime.ParseMediaType(contentType); err == nil {
				if cs, ok := params["charset"]; ok {
					pageEncode = strings.ToLower(strings.TrimSpace(cs))
				}
			}
		}

		switch pageEncode {
		// 不做转码处理
		case "utf8", "utf-8", "unicode-1-1-utf-8":
		default:
			// 指定了编码类型，但不是utf8时，自动转码为utf8
			// get converter to utf-8
			// Charset auto determine. Use golang.org/x/net/html/charset. Get response body and change it to utf-8
			var destReader io.Reader

			if len(pageEncode) == 0 {
				destReader, err = charset.NewReader(self.Response.Body, "")
			} else {
				destReader, err = charset.NewReaderLabel(pageEncode, self.Response.Body)
			}

			if err == nil {
				self.text, err = ioutil.ReadAll(destReader)
				if err == nil {
					self.Response.Body.Close()
					return
				} else {
					//logs.Log.Warning(" *     [convert][%v]: %v (ignore transcoding)\n", self.GetUrl(), err)
				}
			} else {
				//logs.Log.Warning(" *     [convert][%v]: %v (ignore transcoding)\n", self.GetUrl(), err)
			}
		}
	}

	// 不做转码处理
	self.text, err = ioutil.ReadAll(self.Response.Body)
	self.Response.Body.Close()
	if err != nil {
		panic(err.Error())
		return
	}
}

//**************************************** Get 类公开方法 *******************************************\\

// 获取下载错误。
func (self *Context) GetError() error {
	// 若已主动终止任务，则崩溃爬虫协程
	self.spider.tryPanic()
	return self.err
}

// 获取日志接口实例。
//func (*Context) Log() logs.Logs {
//	return logs.Log
//}

// 获取蜘蛛名称。
func (self *Context) GetSpider() *Spider {
	return self.spider
}

// 获取响应流。
func (self *Context) GetResponse() *http.Response {
	return self.Response
}

// 获取响应状态码。
func (self *Context) GetStatusCode() int {
	return self.Response.StatusCode
}

// 获取原始请求。
func (self *Context) GetRequest() *request.Request {
	return self.Request
}

// 获得一个原始请求的副本。
func (self *Context) CopyRequest() *request.Request {
	return self.Request.Copy()
}

// 获取结果字段名列表。
func (self *Context) GetItemFields(ruleName ...string) []string {
	_, rule, found := self.getRule(ruleName...)
	if !found {
		//logs.Log.Error("蜘蛛 %s 调用GetItemFields()时，指定的规则名不存在！", self.spider.GetName())
		return nil
	}
	return self.spider.GetItemFields(rule)
}

// 由索引下标获取结果字段名，不存在时获取空字符串，
// 若ruleName为空，默认为当前规则。
func (self *Context) GetItemField(index int, ruleName ...string) (field string) {
	_, rule, found := self.getRule(ruleName...)
	if !found {
		//logs.Log.Error("蜘蛛 %s 调用GetItemField()时，指定的规则名不存在！", self.spider.GetName())
		return
	}
	return self.spider.GetItemField(rule, index)
}

// 由结果字段名获取索引下标，不存在时索引为-1，
// 若ruleName为空，默认为当前规则。
func (self *Context) GetItemFieldIndex(field string, ruleName ...string) (index int) {
	_, rule, found := self.getRule(ruleName...)
	if !found {
		//logs.Log.Error("蜘蛛 %s 调用GetItemField()时，指定的规则名不存在！", self.spider.GetName())
		return
	}
	return self.spider.GetItemFieldIndex(rule, field)
}

func (self *Context) PullItems() (ds []data.Data) {
	self.Lock()
	ds = self.items
	self.items = []data.Data{}
	self.Unlock()
	return
}

// 获取自定义配置。
func (self *Context) GetKeyin() string {
	return self.spider.GetKeyin()
}

// 获取采集上限。
func (self *Context) GetLimit() int {
	return int(self.spider.GetLimit())
}

// 获取蜘蛛名。
func (self *Context) GetName() string {
	return self.spider.GetName()
}

// 获取规则树。
func (self *Context) GetRules() map[string]*Rule {
	return self.spider.GetRules()
}

// 获取指定规则。
func (self *Context) GetRule(ruleName string) (*Rule, bool) {
	return self.spider.GetRule(ruleName)
}

// 获取当前规则名。
func (self *Context) GetRuleName() string {
	return self.Request.GetRuleName()
}

// 获取请求中临时缓存数据
// defaultValue 不能为 interface{}(nil)
func (self *Context) GetTemp(key string, defaultValue interface{}) interface{} {
	return self.Request.GetTemp(key, defaultValue)
}

// 获取请求中全部缓存数据
func (self *Context) GetTemps() request.Temp {
	return self.Request.GetTemps()
}

// 获得一个请求的缓存数据副本。
func (self *Context) CopyTemps() request.Temp {
	temps := make(request.Temp)
	for k, v := range self.Request.GetTemps() {
		temps[k] = v
	}
	return temps
}

// 从原始请求获取Url，从而保证请求前后的Url完全相等，且中文未被编码。
func (self *Context) GetUrl() string {
	return self.Request.Url
}

func (self *Context) GetMethod() string {
	return self.Request.GetMethod()
}

func (self *Context) GetHost() string {
	return self.Response.Request.URL.Host
}

// 获取响应头信息。
func (self *Context) GetHeader() http.Header {
	return self.Response.Header
}

// 获取请求头信息。
func (self *Context) GetRequestHeader() http.Header {
	return self.Response.Request.Header
}

func (self *Context) GetReferer() string {
	return self.Response.Request.Header.Get("Referer")
}

// 获取响应的Cookie。
func (self *Context) GetCookie() string {
	return self.Response.Header.Get("Set-Cookie")
}

// GetHtmlParser returns goquery object binded to target crawl result.
func (self *Context) GetDom() *goquery.Document {
	if self.dom == nil {
		self.initDom()
	}
	return self.dom
}

// GetBodyStr returns plain string crawled.
func (self *Context) GetText() string {
	if self.text == nil {
		self.initText()
	}
	return util.Bytes2String(self.text)
}