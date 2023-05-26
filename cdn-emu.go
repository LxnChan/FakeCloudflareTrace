package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	/* "strings" */

	"github.com/oschwald/geoip2-golang"
)

func main() {
	// 创建HTTP服务器
	http.HandleFunc("/cdn-cgi/trace", func(w http.ResponseWriter, r *http.Request) {
		// 获取请求的IP地址
		ip := getIPAddress(r)

		// 获取IP地址对应的位置信息
		location := getLocation(ip)

		// 获取一个随机数指代数据中心
		fl := generateRandomString(5)

		// 获取时间戳
		ts := generateTimestamp()

		// 获取UserAgent
		uag := getUserAgent(r)

		// 生成IATA标准代码
		colo := generateAirportCode()

		// 获取当前http协议版本
		http := getProtocol(r)

		// 获取TLS版本
		/* tls := getTLSVersion(r) */

		// 返回信息
		fmt.Fprintf(w, "fl=%s\nh=lxnchan.cn\nip=%s\nts=%s\nvisit_scheme=https\nuag=%s\ncolo=%s\nsliver=none\nhttp=%s\nloc=%s\ntls=null\nsni=plaintext\nwarp=off\ngateway=off\nrbi=off\nkex=X25519", fl, ip, ts, uag, colo, http, location)
	})

	// 启动服务器，监听端口8080
	fmt.Println("Server is starting, listening 8080...")
	http.ListenAndServe(":8080", nil)
}

func getIPAddress(r *http.Request) string {
	// 从请求头中获取"CF-Connecting-IP"字段的值
	cfConnectingIP := r.Header.Get("CF-Connecting-IP")

	// 如果该字段的值为空，则从RemoteAddr中获取IP地址
	if cfConnectingIP == "" {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		return ip
	}

	return cfConnectingIP
}

func getLocation(ip string) string {
	// 加载IP地理位置数据库（请根据实际情况修改数据库文件路径）
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		fmt.Println("无法加载IP地理位置数据库：", err)
		return "null"
	}
	defer db.Close()

	// 查询IP地址对应的位置信息
	ipObj := net.ParseIP(ip)
	record, err := db.City(ipObj)
	if err != nil {
		fmt.Println("无法查询IP地址对应的位置信息：", err)
		return "null"
	}

	// 格式化位置信息
	location := fmt.Sprintf("%s, %s, %s", record.City.Names["en"], record.Subdivisions[0].Names["en"], record.Country.Names["en"])

	return location
}

func generateRandomString(length int) string {
	// 定义包含可用字符的字符串
	charSet := "abcdefghijklmnopqrstuvwxyz0123456789"

	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 生成随机字符串
	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomString[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(randomString)
}

func generateTimestamp() string {
	// 获取当前的 Unix 时间戳（秒）
	seconds := time.Now().Unix()

	// 将秒转换为字符串
	secondsStr := strconv.FormatInt(seconds, 10)

	// 获取当前的毫秒部分
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)

	// 将毫秒转换为字符串
	millisecondsStr := strconv.FormatInt(milliseconds, 10)

	// 构建时间戳字符串，使用小数点分隔秒和毫秒部分，并取毫秒部分的前三位
	timestamp := secondsStr + "." + millisecondsStr[:3]

	return timestamp
}

func getUserAgent(r *http.Request) string {
	// 获取请求头中的User-Agent字段值
	userAgent := r.Header.Get("User-Agent")

	return userAgent
}

func generateAirportCode() string {
	// 定义可用的字母集合
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 生成随机的 IATA 机场代码
	airportCode := ""
	for i := 0; i < 3; i++ {
		randomIndex := rand.Intn(len(alphabet))
		airportCode += string(alphabet[randomIndex])
	}

	return airportCode
}

func getProtocol(r *http.Request) string {
	// 获取当前 http 协议版本
	protocol := r.Proto
	return protocol
}

/* func getTLSVersion(r *http.Request) string {
	// 获取 TLS 版本
	tlsVersion := r.TLS.Version.String()
	return tlsVersion
} */
