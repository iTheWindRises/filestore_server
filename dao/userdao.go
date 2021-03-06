package dao

import (
	mydb "filestore/dao/mysql"
	"fmt"
)

//通用户名和密码注册user
func UserSignUp(userName, passWord string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user(user_name, user_pwd) values(?,?)")

	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(userName, passWord)
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}

	return false
}

func UserSignin(userName, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1")

	if err != nil {
		fmt.Println("Failed to select user,err:" + err.Error())
		return false
	}
	rows, err := stmt.Query(userName)
	if err != nil {
		fmt.Println("Failed to select user,err:" + err.Error())
		return false
	}
	if rows == nil {
		fmt.Println("userName no find:" + userName)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	} else {
		return false
	}
}

//刷新用户登录的token
func UpdateToken(userName, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token(user_name, user_token) values(?,?)")
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(userName, token)
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}
	return true
}

type User struct {
	UserName string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	Status int
}

//用户信息查询
func GetUserInfo(userName string) (*User,error) {
	user := &User{}

	stmt,err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1",
	)
	if err != nil {
		fmt.Println("Failed to select,err:" + err.Error())
		return nil,err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userName).Scan(&user.UserName,&user.SignupAt)
	if err != nil {
		return nil,err
	}
	
	return user,nil
}