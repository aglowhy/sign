package guilinlife

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aglowhy/sign/internal/app/config"
	"github.com/aglowhy/sign/pkg/logger"
	"github.com/aglowhy/sign/pkg/notice"
	"github.com/benbjohnson/phantomjs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	SignInPage = "http://bbs.guilinlife.com/plugin.php?id=ljdaka:newmobile&newmobile=1"
)

type Form struct {
	Action         string
	FormHash       string
	Referer        string
	FastLoginField string
	CookieTime     string
	uAsec          string
	Cookies        []*http.Cookie
}

type Client struct {
	Client   *http.Client
	Integral int
	Ch       chan string
	Config   *config.GuilinlifeConf
	Form
}

func (c *Client) BeforeLogin() error {
	// Start the process once.
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		return err
	}
	defer phantomjs.DefaultProcess.Close()

	// Create a web page.
	// IMPORTANT: Always make sure you close your pages!
	page, err := phantomjs.CreateWebPage()
	if err != nil {
		return err
	}
	defer page.Close()

	wps, _ := page.Settings()
	wps.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/12.0 Mobile/15A372 Safari/604.1"
	if err := page.SetSettings(wps); err != nil {
		return err
	}

	// Open a URL.
	if err := page.Open("http://bbs.guilinlife.com/member.php?mod=logging&action=login&mobile=2"); err != nil {
		return err
	}

	v1, err := page.EvaluateJavaScript(`
		function (){
			return document.getElementById("loginform").action;
		}
	`)
	if err != nil {
		return err
	}
	action := v1.(string)

	v2, err := page.EvaluateJavaScript(`
		function (){
			return document.getElementsByName("formhash")[0].value
		}
	`)
	if err != nil {
		return err
	}
	formhash := v2.(string)

	v3, err := page.EvaluateJavaScript(`
		function (){
			return document.getElementsByName("referer")[0].value
		}
	`)
	if err != nil {
		return err
	}
	referer := v3.(string)

	v4, err := page.EvaluateJavaScript(`
		function (){
			return document.getElementsByName("fastloginfield")[0].value
		}
	`)
	if err != nil {
		return err
	}
	fastloginfield := v4.(string)

	v5, err := page.EvaluateJavaScript(`
		function (){
			return document.getElementsByName("cookietime")[0].value
		}
	`)
	if err != nil {
		return err
	}
	cookietime := v5.(string)

	v6, err := page.Evaluate(`
		function (){
			var pattern = RegExp("u_asec" + "=.[^;]*")
			var matched = document.cookie.match(pattern)
			if (matched) {
				var cookie = matched[0].split('=')
				return cookie[1]
			}
			return ''
		}
	`)
	if err != nil {
		return err
	}
	uasec := v6.(string)

	cookies, _ := page.Cookies()

	c.Form = Form{
		Action:         action,
		FormHash:       formhash,
		Referer:        referer,
		FastLoginField: fastloginfield,
		CookieTime:     cookietime,
		Cookies:        cookies,
		uAsec:          uasec,
	}

	return nil
}

func (c *Client) Login() error {

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	//userName, err := Utf8ToGbk([]byte(c.User.Username))
	//if err != nil {
	//	fmt.Println(err)
	//}

	data := url.Values{
		"formhash":       {c.Form.FormHash},
		"referer":        {c.Form.Referer},
		"fastloginfield": {c.Form.FastLoginField},
		"cookietime":     {c.Form.CookieTime},
		"username":       {c.Config.Username},
		"password":       {c.Config.Password},
		"questionid":     {"0"},
		"answer":         {""},
	}
	body := strings.NewReader(data.Encode())

	loginUrl := c.Form.Action + "&handlekey=loginform&inajax=1&u_atype=2&u_asec=" + c.Form.uAsec

	u, _ := url.Parse(loginUrl)
	jar.SetCookies(u, c.Form.Cookies)

	c.Client.Jar = jar
	req, err := http.NewRequest("POST", loginUrl, body)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/xml, text/xml, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/12.0 Mobile/15A372 Safari/604.1")
	req.Header.Set("Host", "bbs.guilinlife.com")
	req.Header.Set("Referer", "http://bbs.guilinlife.com/portal.php")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, _ := c.Client.Do(req)
	defer resp.Body.Close()

	//utf8Reader := transform.NewReader(resp.Body,
	//	simplifiedchinese.GBK.NewDecoder())
	//htmlBytes, err := ioutil.ReadAll(utf8Reader)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%s\n", htmlBytes)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	loginResult, _ := doc.Html()
	logger.Infof(nil, "登录结果: %s", loginResult)

	if !strings.Contains(loginResult, "现在将转入登录前页面") {
		notice.ServerChanNotice(&notice.Message{
			Title: "桂林人论坛登录失败",
			Desc:  loginResult,
		})

		return errors.New("登陆异常")
	}

	return nil
}

func (c *Client) Sign() error {
	req, err := http.NewRequest("GET", SignInPage, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:71.0) Gecko/20100101 Firefox/71.0")
	req.Header.Set("Host", "bbs.guilinlife.com")
	req.Header.Set("Referer", "http://bbs.guilinlife.com/ljdaka-ranklist.html")

	resp, _ := c.Client.Do(req)
	defer resp.Body.Close()

	//for _, cookie := range resp.Cookies() {
	//	fmt.Println("Found a cookie named:", cookie.Name)
	//}

	utf8Reader := transform.NewReader(resp.Body,
		simplifiedchinese.GBK.NewDecoder())
	htmlBytes, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		return err
	}

	//doc, err := goquery.NewDocumentFromReader(resp.Body)
	//if err != nil {
	//	fmt.Printf("Parsing HTML error:%s\r\n", err)
	//	os.Exit(1)
	//}
	//
	//signResult, _ := doc.Html()
	//fmt.Println("签到结果: " + signResult)

	signResult := string(htmlBytes)
	logger.Infof(nil, "签到结果: %s", signResult)

	if strings.Contains(signResult, "您已打卡") {
		notice.ServerChanNotice(&notice.Message{
			Title: "桂林人论坛重复签到",
			Desc:  signResult,
		})

		return nil
	}

	notice.ServerChanNotice(&notice.Message{
		Title: "桂林人论坛签到成功",
		Desc:  signResult,
	})

	return nil
}

func Exec(conf *config.GuilinlifeConf) error {
	client := Client{
		Client:   &http.Client{},
		Integral: 0,
		Ch:       make(chan string, 1),
		Config:   conf,
	}

	err := client.BeforeLogin()
	if err != nil {
		return err
	}

	err = client.Login()
	if err != nil {
		return err
	}

	err = client.Sign()
	if err != nil {
		return err
	}

	return nil
}
