package handler

import (
	"log"
	"filestore/dao"
	util "filestore/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "#995?"
)

//用户注册
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	//GET请求返回注册页面
	if r.Method == http.MethodGet {
		html, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(html)
		return
	}
	//POST请求注册用户
	if r.Method == http.MethodPost {
		userName := r.FormValue("username")
		password := r.FormValue("password")

		if len(userName) < 3 || len(userName) > 12 ||
			len(password) < 6 || len(password) > 16 {
			w.Write([]byte("invalid parameter"))
			return
		}
		enc_pwd := util.Sha1([]byte(password + pwd_salt))
		ok := dao.UserSignUp(userName, enc_pwd)
		if ok {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte("FAILED"))
		}
	}

}

//用户登录
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("请求用户登录接口")
	userName := r.FormValue("username")
	password := r.FormValue("password")
	encPwd := util.Sha1([]byte(password + pwd_salt))
	//1.校验用户名和密码
	ok := dao.UserSignin(userName, encPwd)
	if !ok {
		w.Write([]byte("SignIn fail"))
		return
	}
	//2.生成token
	token := GenToken(userName)
	ok = dao.UpdateToken(userName, token)
	if !ok {
		w.Write([]byte("SignIn fail"))
		return
	}
	//3.登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.NewRespMsg(0,"OK",struct{
		Location string
		UserName string
		Token string
	}{
		Location: "http://" + r.Host + "/static/view/home.html",
		UserName : userName,
		Token :token,
	})

	w.Write(resp.JSONBytes())

}

//查询用户信息
func UserInfoHandler(w http.ResponseWriter,r *http.Request) {
	log.Println("请求用户信息接口")
	//1.解析请求参数
	userName := r.FormValue("username")
	//token := r.FormValue("token")
	// //2.验证token是否有效
	// ok := IsTokenValid(token)
	// if !ok {
	// 	w.WriteHeader(http.StatusForbidden)
	// } 
	//3.查询用户信息
	user, err := dao.GetUserInfo(userName)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}
	//4.组装并响应用户数据
	resp := util.NewRespMsg(0,"OK",user)
	w.Write(resp.JSONBytes())
}


func GenToken(userName string) string {
	//md5(username+ timestamp + token_salt )+ timestamp[0:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(userName + ts + "_tokensalt"))

	return tokenPrefix + ts[:8]
}

//验证token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	return true
	//判断token是否过期
	//查询数据库token
	//对比两个token
}