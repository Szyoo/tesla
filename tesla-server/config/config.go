package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Region    string // 部署区域：cn / jp / na / eu（见 REGION 环境变量）
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Tesla     TeslaConfig
	JWT       JWTConfig
	Map       MapConfig
	AI        AIConfig
	Telemetry TelemetryConfig
}

// regionEndpoints 保存各区域的 Tesla 端点默认值。
// 注意：Tesla Fleet API 只有 3 个区域 —— 北美/亚太(na)、欧洲(eu)、中国(cn)。
// 日本属于亚太(APAC)，走北美区 na 端点，不存在 *.tesla.jp 域名。
type regionEndpoints struct {
	AuthBase    string // OAuth 授权/令牌主机（cn 用 auth.tesla.cn，其余用 auth.tesla.com）
	FleetAPIURL string // Fleet API 基址，同时用作 audience
	PairingBase string // 虚拟钥匙配对页主机（tesla.cn 或 tesla.com）
}

var teslaRegions = map[string]regionEndpoints{
	"cn": {
		AuthBase:    "https://auth.tesla.cn",
		FleetAPIURL: "https://fleet-api.prd.cn.vn.cloud.tesla.cn",
		PairingBase: "https://tesla.cn",
	},
	"jp": { // 日本：APAC，复用 na 端点
		AuthBase:    "https://auth.tesla.com",
		FleetAPIURL: "https://fleet-api.prd.na.vn.cloud.tesla.com",
		PairingBase: "https://tesla.com",
	},
	"na": {
		AuthBase:    "https://auth.tesla.com",
		FleetAPIURL: "https://fleet-api.prd.na.vn.cloud.tesla.com",
		PairingBase: "https://tesla.com",
	},
	"eu": {
		AuthBase:    "https://auth.tesla.com",
		FleetAPIURL: "https://fleet-api.prd.eu.vn.cloud.tesla.com",
		PairingBase: "https://tesla.com",
	},
}

type TelemetryConfig struct {
	Enabled         bool
	ListenAddr      string
	Hostname        string
	PrivateKey      string
	PublicKeyFile   string
	TLSCertFile     string
	TLSKeyFile      string
	CACertFile      string
	UseDefaultEngCA bool
	CertSync        CertSyncConfig
}

type CertSyncConfig struct {
	Enabled    bool
	VCPKeysDir string
	CertSrcDir string
	CertsDir   string
}

type AIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type TeslaConfig struct {
	ClientID            string
	ClientSecret        string
	RedirectURI         string
	AuthURL             string
	TokenURL            string
	FleetAPIURL         string
	Audience            string
	VCPURL              string
	FrontendCallbackURL string
	PartnerDomain       string
	PairingBase         string // 虚拟钥匙配对页主机，按区域而定（如 https://tesla.com）
}

type JWTConfig struct {
	Secret           string
	ExpiresIn        int
	RefreshExpiresIn int
}

type MapConfig struct {
	TencentKey string
}

func Load() *Config {
	// REGION 决定 Tesla 端点的默认值（cn/jp/na/eu）。默认 jp（本 fork 为日本版）。
	// 显式设置的 TESLA_*_URL 环境变量始终优先于区域默认值，便于自定义/代理场景。
	region := strings.ToLower(getEnv("REGION", "jp"))
	ep, ok := teslaRegions[region]
	if !ok {
		ep = teslaRegions["jp"]
		region = "jp"
	}

	teslaTokenURL := getEnv("TESLA_TOKEN_URL", ep.AuthBase+"/oauth2/v3/token")
	teslaAuthURL := getEnv("TESLA_AUTH_URL", ep.AuthBase+"/oauth2/v3/authorize")
	teslaFleetAPIURL := getEnv("TESLA_FLEET_API_URL", ep.FleetAPIURL)
	teslaAudience := getEnv("TESLA_AUDIENCE", ep.FleetAPIURL)
	teslaRedirectURI := getEnv("TESLA_REDIRECT_URI", "http://localhost:8080/api/tesla/callback")

	teslaTokenURL = ensureHTTPS(teslaTokenURL)
	teslaAuthURL = ensureHTTPS(teslaAuthURL)
	teslaFleetAPIURL = ensureHTTPS(teslaFleetAPIURL)
	teslaAudience = ensureHTTPS(teslaAudience)
	teslaRedirectURI = ensureHTTPS(teslaRedirectURI)

	return &Config{
		Region: region,
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "tesla_platform"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Tesla: TeslaConfig{
			ClientID:            getEnv("TESLA_CLIENT_ID", ""),
			ClientSecret:        getEnv("TESLA_CLIENT_SECRET", ""),
			RedirectURI:         teslaRedirectURI,
			AuthURL:             teslaAuthURL,
			TokenURL:            teslaTokenURL,
			FleetAPIURL:         teslaFleetAPIURL,
			Audience:            teslaAudience,
			VCPURL:              getEnv("TESLA_VCP_URL", ""),
			FrontendCallbackURL: getEnv("TESLA_FRONTEND_CALLBACK_URL", ""),
			PartnerDomain:       getEnv("TESLA_PARTNER_DOMAIN", ""),
			PairingBase:         ep.PairingBase,
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "your-secret-key"),
			ExpiresIn:        getEnvAsInt("JWT_EXPIRES_IN", 604800),      // 7天
			RefreshExpiresIn: getEnvAsInt("JWT_REFRESH_EXPIRES_IN", 2592000), // 30天
		},
		Map: MapConfig{
			TencentKey: getEnv("TENCENT_MAP_KEY", ""),
		},
		AI: AIConfig{
			APIKey:  getEnv("AI_API_KEY", ""),
			Model:   getEnv("AI_MODEL", "glm-4-flash"),
			BaseURL: getEnv("AI_BASE_URL", "https://open.bigmodel.cn/api/paas/v4"),
		},
		Telemetry: TelemetryConfig{
			Enabled:         getEnvAsBool("TELEMETRY_ENABLED", false),
			ListenAddr:      getEnv("TELEMETRY_LISTEN_ADDR", ":8443"),
			Hostname:        getEnv("TELEMETRY_HOSTNAME", ""),
			PrivateKey:      getEnv("TELEMETRY_PRIVATE_KEY", ""),
			PublicKeyFile:   getEnv("TELEMETRY_PUBLIC_KEY", ""),
			TLSCertFile:     getEnv("TELEMETRY_TLS_CERT", ""),
			TLSKeyFile:      getEnv("TELEMETRY_TLS_KEY", ""),
			CACertFile:      getEnv("TELEMETRY_CA_CERT", ""),
			UseDefaultEngCA: getEnvAsBool("TELEMETRY_USE_ENG_CA", false),
			CertSync: CertSyncConfig{
				Enabled:    getEnvAsBool("CERT_SYNC_ENABLED", true),
				VCPKeysDir: getEnv("CERT_SYNC_VCP_KEYS_DIR", "/opt/tesla-vcp/keys"),
				CertSrcDir: getEnv("CERT_SYNC_SRC_DIR", ""),
				CertsDir:   getEnv("CERT_SYNC_DST_DIR", ""),
			},
		},
	}
}

func ensureHTTPS(url string) string {
	if len(url) > 7 && url[:7] == "http://" {
		// 本地调试地址豁免：Tesla 允许 localhost 使用 http，
		// 强制转 https 会导致 redirect_uri 与门户登记值不符。
		host := url[7:]
		if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
			return url
		}
		return "https://" + host
	}
	return url
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return defaultValue
}
