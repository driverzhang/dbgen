// 配置, 如果有新增配置,在这里完善配置的结构
package config

import (
	"advertisement/pkg/app"
	"advertisement/pkg/authlogin"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	elastic "github.com/olivere/elastic/v7/config"
)

type (
	application struct {
		Name          string `toml:"name"`
		Domain        string `toml:"domain"`
		Addr          string `toml:"addr"`
		PasswordToken string `toml:"password_token"`
		JwtToken      string `toml:"jwt-token"`
		CertFile      string `toml:"cert_file"`
		KeyFile       string `toml:"key_file"`
	}
	master struct {
		Addr     string `toml:"addr"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		DbName   string `toml:"dbname"`
		MaxIdle  int    `toml:"max_idle"`
		MaxOpen  int    `toml:"max_open"`
	}
	slave struct {
		Addr     string `toml:"addr"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		DbName   string `toml:"dbname"`
		MaxIdle  int    `toml:"max_idle"`
		MaxOpen  int    `toml:"max_open"`
	}
	database struct {
		Master master  `toml:"master"`
		Slaves []slave `toml:"slave"`
	}
	mongo struct {
		Url             string `toml:"url"`
		Database        string `toml:"database"`
		MaxConnIdleTime int    `toml:"max_conn_idle_time"`
		MaxPoolSize     int    `toml:"max_pool_size"`
		Username        string `toml:"username"`
		Password        string `toml:"password"`
	}
	redis struct {
		Addr         string `toml:"addr"`
		Password     string `toml:"password"`
		Db           int    `toml:"db"`
		PoolSize     int    `toml:"pool_size"`
		MinIdleConns int    `toml:"min_idle_conns"`
	}
	sessions struct {
		Key          string `toml:"key"`
		Name         string `toml:"name"`
		Domain       string `toml:"domain"`
		Addr         string `toml:"addr"`
		Password     string `toml:"password"`
		Db           int    `toml:"db"`
		PoolSize     int    `toml:"pool_size"`
		MinIdleConns int    `toml:"min_idle_conns"`
	}
	elasticsearch struct {
		URL         string `toml:"url"`
		Index       string `toml:"index"`
		Username    string `toml:"username"`
		Password    string `toml:"password"`
		Shards      int    `toml:"shards"`
		Replicas    int    `toml:"replicas"`
		Sniff       bool   `toml:"sniff"`
		HealthCheck bool   `toml:"health"`
		InfoLog     string `toml:"info_log"`
		ErrorLog    string `toml:"error_log"`
		TraceLog    string `toml:"trace_log"`
	}
)

// 配置顶级结构
type config struct {
	Application   application   `toml:"application"`
	Database      database      `toml:"database"`
	Mongo         mongo         `toml:"mongo"`
	Redis         redis         `toml:"redis"`
	Sessions      sessions      `toml:"sessions"`
	ElasticSearch elasticsearch `toml:"elastic"`
	Log           logConfig     `toml:"log"`
	FilePath      filePath      `toml:"filepath"`
	Sms           sms           `toml:"sms"`
}

// 三方支付配置信息
type (
	payment struct {
		Alipay alipay `toml:"alipay"`
		Wechat wechat `toml:"wechat"`
	}
	alipay struct {
		AppId              string `toml:"appid"`
		AlipayRsaPublicKey string `toml:"alipay_rsa_public_key"`
		RsaPrivateKey      string `toml:"rsa_private_key"`
		NotifyUrl          string `toml:"notify_url"`
		ReturnUrl          string `toml:"return_url"`
		Product            bool   `toml:"product"`
	}
	wechat struct {
		AppId     string `toml:"appid"`
		MchId     string `toml:"mch_id"`
		NotifyUrl string `toml:"notify_url"`
		SignKey   string `toml:"sign_key"`
	}
)

// 日志配置信息
type logConfig struct {
	EsMode bool   `toml:"esmode"`
	Index  string `toml:"index"`
	Dir    string `toml:"dir"`
}

// photo文件配置信息
type filePath struct {
	PhotoDir string `toml:"photo_dir"`
}

// 短信验证配置信息
type sms struct {
	ApiKey string `toml:"api_key"`
	SmsApi string `toml:"sms_api"`
	User   string `toml:"user"`
	Passwd string `toml:"password"`
}

// 友盟推送
type upush struct {
	Android struct {
		AppKey        string `toml:"app_key"`
		MessageSecret string `toml:"message_secret"`
		AppSecret     string `toml:"app_secret"`
	} `toml:"android"`
	IOS struct {
		AppKey        string `toml:"app_key"`
		MessageSecret string `toml:"message_secret"`
		AppSecret     string `toml:"app_secret"`
	} `toml:"ios"`
}

//
type bizConfig struct {
	ExpireTime         int64  `toml:"expire_time"`           // 新手任务有效期
	SignStart          string `toml:"sign_start"`            // 签到开始时间
	SignEnd            string `toml:"sign_end"`              // 签到结束时间
	DivideStart        string `toml:"divide_start"`          // 瓜分大奖开始时间
	DivideEnd          string `toml:"divide_end"`            // 瓜分大奖结束时间
	InviteFirstAward   int    `toml:"invite_first_award"`    // 首次邀请用户所获取的奖励
	InviteAward        int    `toml:"invite_award"`          // 邀请用户所获得的奖励
	InviteTaskFinshNum int    `toml:"invite_task_finsh_num"` // 用户完成多少次任务算是有效邀请
	WealthRankingNum   int    `toml:"wealth_ranking_show"`   // 财富排行显示条数
	WealthRankingMax   int    `toml:"wwalth_ranking_max"`    // 财富排行最大排名条数
	CronPushNew        struct {
		NewbieCheckDay int `toml:"newbie_check_day"` // 注册日起15日内有没有完成的新手任务
		NewbieSendHz   int `toml:"newbie_send_hz"`   // 注册日起每7日一次
		NewbieSendNum  int `toml:"newbie_send_num"`  //持续1周

	} `toml:"cron_push_new"`
	JQexchange int `toml:"jq_exchange"` // 荐令转化为金钱汇率
	Oss        struct {
		EndPoint        string `toml:"end_point"` // ENDPOINT
		AccessKey       string `toml:"access_key"`
		AccessKeySecret string `toml:"access_key_secret"`
		BucketName      string `toml:"bucket_name"`
		Thumbnail       string `toml:"thumbnail"`
		ImageSize       int    `toml:"image_size"`
		OssAddr         string `toml:"oss_addr"`
		VideoSnapshot   string `toml:"video_snapshot"`
	} `toml:"oss"`
}

var Config config
var Payment payment
var Upush upush
var BizConfig bizConfig
