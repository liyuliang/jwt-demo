package main

import (
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
)

const (
    SecretKey = "I have login"
)


func main() {
    router := gin.Default()
    router.GET("/login", loginHandler)
    router.Use(authMiddleware)
    router.GET("/getData", getData)
    router.Run(":2323")
}

//验证token中间件
func authMiddleware(ctx *gin.Context) {
    //从cookie中获取token
    if tokenStr, err := ctx.Cookie("token"); err == nil {
        //获取验证之后的结果
        token, err := parseToken(tokenStr)
        if err != nil {
            ctx.JSON(200, "token verify error")
        }
        //如果验证结果是false直接返回token错误我 如果成功则继续下一个handler
        if token.Valid {
            ctx.Next()
        } else {
            ctx.JSON(200, "token verify error")
            ctx.Abort()
        }
    } else {
        ctx.JSON(200, "no token")
        ctx.Abort()
    }
}

func getData(ctx *gin.Context) {
    ctx.JSON(200, "data")
}

func loginHandler(ctx *gin.Context) {
    user := ctx.Query("user")
    pwd := ctx.Query("pwd")

    if user == "peter" && pwd == "pwd" {
        token := CreateToken(user, pwd)
        //ctx.Header("Authorization", token)
        ctx.SetCookie("token", token, 10, "/", "localhost", false, true)
        ctx.JSON(200, "ok")
    } else {
        ctx.JSON(200, "user is not exit")
    }
}

func parseToken(s string) (*jwt.Token, error) {
    fn := func(token *jwt.Token) (interface{}, error) {
        return []byte(SecretKey), nil
    }
    return jwt.Parse(s, fn)
}

//创建token
func CreateToken(user, pwd string) string {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := make(jwt.MapClaims)
    claims["user"] = user
    // 这边的pwd 不应该放到claims 荷载中不应该有机密的数据
    claims["pwd"] = pwd
    token.Claims = claims
    if tokenString, err := token.SignedString([]byte(SecretKey)); err == nil {
        return tokenString
    } else {
        return ""
    }
}
