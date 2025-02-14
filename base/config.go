package base

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ucloud/ucloud-sdk-go/services/uaccount"
	sdk "github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/ucloud/ucloud-sdk-go/ucloud/log"
)

//ConfigFilePath path of config.json
var ConfigFilePath = fmt.Sprintf("%s/%s", GetConfigDir(), "config.json")

//CredentialFilePath path of credential.json
var CredentialFilePath = fmt.Sprintf("%s/%s", GetConfigDir(), "credential.json")

//LocalFileMode file mode of $HOME/ucloud/*
const LocalFileMode os.FileMode = 0600

//DefaultTimeoutSec default timeout for requesting api, 15s
const DefaultTimeoutSec = 15

//DefaultMaxRetryTimes default timeout for requesting api, 15s
const DefaultMaxRetryTimes = 3

//DefaultBaseURL location of api server
const DefaultBaseURL = "https://api.ucloud.cn/"

//DefaultProfile name of default profile
const DefaultProfile = "default"

//Version 版本号
const Version = "0.1.23"

//ConfigIns 配置实例, 程序加载时生成
var ConfigIns = &AggConfig{
	Profile:       DefaultProfile,
	BaseURL:       DefaultBaseURL,
	Timeout:       DefaultTimeoutSec,
	MaxRetryTimes: sdk.Int(DefaultMaxRetryTimes),
}

//AggConfigListIns 配置列表, 进程启动时从本地文件加载
var AggConfigListIns = &AggConfigManager{}

//ClientConfig 创建sdk client参数
var ClientConfig *sdk.Config

//AuthCredential 创建sdk client参数
var AuthCredential *auth.Credential

//BizClient 用于调用业务接口
var BizClient *Client

//Global 全局flag
var Global GlobalFlag

//GlobalFlag 几乎所有接口都需要的参数，例如 region zone projectID
type GlobalFlag struct {
	Debug      bool
	JSON       bool
	Version    bool
	Completion bool
	Config     bool
	Signup     bool
	Profile    string
	PublicKey  string
	PrivateKey string
}

//CLIConfig cli_config element
type CLIConfig struct {
	ProjectID     string `json:"project_id"`
	Region        string `json:"region"`
	Zone          string `json:"zone"`
	BaseURL       string `json:"base_url"`
	Timeout       int    `json:"timeout_sec"`
	Profile       string `json:"profile"`
	Active        bool   `json:"active"` //是否生效
	MaxRetryTimes *int   `json:"max_retry_times"`
}

//CredentialConfig credential element
type CredentialConfig struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Profile    string `json:"profile"`
}

//AggConfig 聚合配置 config+credential
type AggConfig struct {
	Profile       string `json:"profile"`
	Active        bool   `json:"active"`
	ProjectID     string `json:"project_id"`
	Region        string `json:"region"`
	Zone          string `json:"zone"`
	BaseURL       string `json:"base_url"`
	Timeout       int    `json:"timeout_sec"`
	PublicKey     string `json:"public_key"`
	PrivateKey    string `json:"private_key"`
	MaxRetryTimes *int   `json:"max_retry_times"`
}

//ConfigPublicKey 输入公钥
func (p *AggConfig) ConfigPublicKey() error {
	Cxt.Print("Your public-key:")
	_, err := fmt.Scanf("%s\n", &p.PublicKey)
	if err != nil {
		Cxt.Println(err)
		return err
	}
	p.PublicKey = strings.TrimSpace(p.PublicKey)
	AuthCredential.PublicKey = p.PublicKey
	return nil
}

//ConfigPrivateKey 输入私钥
func (p *AggConfig) ConfigPrivateKey() error {
	Cxt.Print("Your private-key:")
	_, err := fmt.Scanf("%s\n", &p.PrivateKey)
	if err != nil {
		Cxt.Println(err)
		return err
	}
	p.PrivateKey = strings.TrimSpace(p.PrivateKey)
	AuthCredential.PrivateKey = p.PrivateKey
	return nil
}

//GetClientConfig 用来生成sdkClient
func (p *AggConfig) GetClientConfig(isDebug bool) *sdk.Config {
	clientConfig := &sdk.Config{
		Region:    p.Region,
		ProjectId: p.ProjectID,
		BaseUrl:   ClientConfig.BaseUrl,
		Timeout:   ClientConfig.Timeout,
		UserAgent: ClientConfig.UserAgent,
		LogLevel:  ClientConfig.LogLevel,
	}
	if isDebug == true {
		clientConfig.LogLevel = log.DebugLevel
	}
	return clientConfig
}

//GetCredential 用来生成SDkClient
func (p *AggConfig) GetCredential() *auth.Credential {
	return &auth.Credential{
		PublicKey:  p.PublicKey,
		PrivateKey: p.PrivateKey,
	}
}

func (p *AggConfig) copyToCLIConfig(target *CLIConfig) {
	target.Profile = p.Profile
	target.BaseURL = p.BaseURL
	target.Timeout = p.Timeout
	target.ProjectID = p.ProjectID
	target.Region = p.Region
	target.Zone = p.Zone
	target.Active = p.Active
	target.MaxRetryTimes = p.MaxRetryTimes
}

func (p *AggConfig) copyToCredentialConfig(target *CredentialConfig) {
	target.Profile = p.Profile
	target.PrivateKey = p.PrivateKey
	target.PublicKey = p.PublicKey
}

//AggConfigManager 配置管理
type AggConfigManager struct {
	activeProfile string
	configs       map[string]*AggConfig
	configFile    *os.File
	credFile      *os.File
}

//NewAggConfigManager create instance
func NewAggConfigManager(cfgFile, credFile *os.File) (*AggConfigManager, error) {
	manager := &AggConfigManager{
		configs:    make(map[string]*AggConfig),
		configFile: cfgFile,
		credFile:   credFile,
	}

	err := manager.Load()
	if err != nil {
		if !os.IsNotExist(err) {
			return manager, err
		}

		aerr := adaptOldConfig()
		if aerr != nil {
			HandleError(fmt.Errorf("adapt to old config failed: %v", aerr))
			return manager, aerr
		}

		err := manager.Load()
		if err != nil {
			HandleError(fmt.Errorf("retry to load cli config failed: %v", err))
			return manager, err
		}
	}
	return manager, nil
}

//Append config to list, override if already exist the same profile
func (p *AggConfigManager) Append(config *AggConfig) error {
	if _, ok := p.configs[config.Profile]; ok {
		return fmt.Errorf("profile %s exists already", config.Profile)
	}

	if config.Active && config.Profile != p.activeProfile {
		if ac, ok := p.configs[p.activeProfile]; ok {
			ac.Active = false
		}
		p.activeProfile = config.Profile
	}
	p.configs[config.Profile] = config
	return p.Save()
}

//UpdateAggConfig  update AggConfig append if not exist
func (p *AggConfigManager) UpdateAggConfig(config *AggConfig) error {
	if _, ok := p.configs[config.Profile]; !ok {
		return p.Append(config)
	}

	if config.Active && config.Profile != p.activeProfile {
		if ac, ok := p.configs[p.activeProfile]; ok {
			ac.Active = false
		}
		p.activeProfile = config.Profile
	}
	return p.Save()
}

//Load AggConfigList from local file  $HOME/.ucloud/config.json+credential.json
func (p *AggConfigManager) Load() error {
	configs, err := p.parseCLIConfigs()
	if err != nil {
		return fmt.Errorf("read config failed: %v", err)
	}
	credentials, err := p.parseCredentials()
	if err != nil {
		return fmt.Errorf("read credential failed: %v", err)
	}

	//key: profile , value: CLIConfig
	configMap := make(map[string]*CLIConfig)
	for _, config := range configs {
		c := config
		configMap[config.Profile] = &c
		if config.Active {
			p.activeProfile = config.Profile
		}
	}
	credMap := make(map[string]*CredentialConfig)
	for _, cred := range credentials {
		c := cred
		credMap[cred.Profile] = &c
	}

	for profile, config := range configMap {
		cred, ok := credMap[profile]
		if !ok {
			LogError("profile: %s don't exist in credential")
			continue
		}

		p.configs[profile] = &AggConfig{
			PrivateKey:    cred.PrivateKey,
			PublicKey:     cred.PublicKey,
			Profile:       config.Profile,
			ProjectID:     config.ProjectID,
			Region:        config.Region,
			Zone:          config.Zone,
			BaseURL:       config.BaseURL,
			Timeout:       config.Timeout,
			Active:        config.Active,
			MaxRetryTimes: config.MaxRetryTimes,
		}
	}

	if p.activeProfile == "" && len(configMap) > 0 {
		return fmt.Errorf("no active config found, run 'ucloud config list' to check")
	}
	if _, ok := credMap[p.activeProfile]; p.activeProfile != "" && !ok {
		return fmt.Errorf("profile %s's credential don't exist, run 'ucloud config list' to check", p.activeProfile)
	}

	return nil
}

//Save configs to local file
func (p *AggConfigManager) Save() error {
	clics := []*CLIConfig{}
	credcs := []*CredentialConfig{}
	for _, aggConfig := range p.configs {
		cliConfig := &CLIConfig{}
		aggConfig.copyToCLIConfig(cliConfig)
		clics = append(clics, cliConfig)

		credConfig := &CredentialConfig{}
		aggConfig.copyToCredentialConfig(credConfig)
		credcs = append(credcs, credConfig)
	}
	aerr := WriteJSONFile(clics, p.configFile.Name())
	berr := WriteJSONFile(credcs, p.credFile.Name())

	if aerr != nil && berr != nil {
		return fmt.Errorf("save cli config failed: %v | save credentail failed: %v", aerr, berr)
	}
	if aerr != nil {
		return fmt.Errorf("save cli config failed: %v", aerr)
	}
	if berr != nil {
		return fmt.Errorf("save cerdentail failed: %v", berr)
	}
	return nil
}

//DeleteByProfile 从AggConfigList和本地文件中删除此配置
func (p *AggConfigManager) DeleteByProfile(profile string) error {
	if _, ok := p.configs[profile]; !ok {
		return fmt.Errorf("profile: %s is not exist", profile)
	}

	ac := p.configs[profile]
	if ac.Active {
		return fmt.Errorf("can't delete active profile")
	}

	delete(p.configs, profile)

	err := p.Save()
	if err != nil {
		return fmt.Errorf("delete profile %s failed: %v", profile, err)
	}
	return nil
}

//GetProfileNameList 获取所有profiles 用于ucloud config --profile 补全
func (p *AggConfigManager) GetProfileNameList() []string {
	profiles := []string{}
	for _, item := range p.configs {
		profiles = append(profiles, item.Profile)
	}
	return profiles
}

//GetAggConfigList get all profile config
func (p *AggConfigManager) GetAggConfigList() []AggConfig {
	configs := []AggConfig{}
	for _, cfg := range p.configs {
		configs = append(configs, *cfg)
	}
	return configs
}

//GetAggConfigByProfile get config of specify profile
func (p *AggConfigManager) GetAggConfigByProfile(profile string) (*AggConfig, bool) {
	if ac, ok := p.configs[profile]; ok {
		return ac, true
	}
	return nil, false
}

//GetActiveAggConfig get active agg config
func (p *AggConfigManager) GetActiveAggConfig() (*AggConfig, error) {
	if ac, ok := p.configs[p.activeProfile]; ok {
		return ac, nil
	}
	return nil, fmt.Errorf("active profile not found. see 'ucloud config list'")
}

//GetActiveAggConfigName get active config name
func (p *AggConfigManager) GetActiveAggConfigName() string {
	if ac, ok := p.configs[p.activeProfile]; ok {
		return ac.Profile
	}
	return ""
}

func (p *AggConfigManager) parseCLIConfigs() ([]CLIConfig, error) {
	var configs []CLIConfig
	rawConfig, err := ioutil.ReadAll(p.configFile)
	if err != nil {
		return nil, err
	}
	if len(rawConfig) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(rawConfig, &configs)
	if err != nil {
		return nil, fmt.Errorf("parse cli config faild: %v", err)
	}
	//特殊处理未配置max_retry_times的情况，v0.1.21之前硬编码重试次数为3
	for idx := range configs {
		if configs[idx].MaxRetryTimes == nil {
			configs[idx].MaxRetryTimes = sdk.Int(DefaultMaxRetryTimes)
		}
	}
	return configs, nil
}

func (p *AggConfigManager) parseCredentials() ([]CredentialConfig, error) {
	var credentials []CredentialConfig
	rawCred, err := ioutil.ReadAll(p.credFile)
	if err != nil {
		return nil, err
	}

	if len(rawCred) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(rawCred, &credentials)
	if err != nil {
		return nil, fmt.Errorf("parse credential failed: %v", err)
	}
	return credentials, nil
}

//ListAggConfig ucloud --config + ucloud config list
func ListAggConfig(json bool) {
	aggConfigs := AggConfigListIns.GetAggConfigList()
	for idx, ac := range aggConfigs {
		aggConfigs[idx].PrivateKey = MosaicString(ac.PrivateKey, 8, 5)
		aggConfigs[idx].PublicKey = MosaicString(ac.PublicKey, 8, 5)
	}
	if json {
		PrintJSON(aggConfigs, os.Stdout)
	} else {
		PrintTableS(aggConfigs)
	}
}

//LoadUserInfo 从~/.ucloud/user.json加载用户信息
func LoadUserInfo() (*uaccount.UserInfo, error) {
	filePath := GetConfigDir() + "/user.json"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("user.json is not exist")
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var user uaccount.UserInfo
	err = json.Unmarshal(content, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//OldConfig 0.1.7以及之前版本的配置struct
type OldConfig struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Region     string `json:"region"`
	Zone       string `json:"zone"`
	ProjectID  string `json:"project_id"`
}

//Load 从本地文件加载配置
func (p *OldConfig) Load() error {
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		p = new(OldConfig)
	} else {
		content, err := ioutil.ReadFile(ConfigFilePath)
		if err != nil {
			return err
		}
		json.Unmarshal(content, p)
	}
	return nil
}

func adaptOldConfig() error {
	oc := &OldConfig{}
	err := oc.Load()
	if err != nil {
		return err
	}
	ac := &AggConfig{
		Profile:       DefaultProfile,
		ProjectID:     oc.ProjectID,
		Region:        oc.Region,
		Zone:          oc.Zone,
		BaseURL:       DefaultBaseURL,
		Timeout:       DefaultTimeoutSec,
		Active:        true,
		PrivateKey:    oc.PrivateKey,
		PublicKey:     oc.PublicKey,
		MaxRetryTimes: sdk.Int(DefaultMaxRetryTimes),
	}
	err = os.Rename(ConfigFilePath, ConfigFilePath+".old")
	if err != nil {
		return err
	}
	return AggConfigListIns.Append(ac)
}

func init() {
	bc, err := GetBizClient(ConfigIns)
	if err != nil {
		HandleError(err)
	}
	BizClient = bc
}

//GetBizClient 初始化BizClient
func GetBizClient(ac *AggConfig) (*Client, error) {
	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", ac.Timeout))
	if err != nil {
		err = fmt.Errorf("parse timeout %ds failed: %v", ac.Timeout, err)
	}
	ClientConfig = &sdk.Config{
		BaseUrl:    ac.BaseURL,
		Timeout:    timeout,
		UserAgent:  fmt.Sprintf("UCloud-CLI/%s", Version),
		LogLevel:   log.FatalLevel,
		Region:     ac.Region,
		ProjectId:  ac.ProjectID,
		MaxRetries: *ac.MaxRetryTimes,
	}

	if Global.PublicKey != "" && Global.PrivateKey != "" {
		AuthCredential = &auth.Credential{
			PublicKey:  Global.PublicKey,
			PrivateKey: Global.PrivateKey,
		}
	} else {
		AuthCredential = &auth.Credential{
			PublicKey:  ac.PublicKey,
			PrivateKey: ac.PrivateKey,
		}
	}

	return NewClient(ClientConfig, AuthCredential), err
}

//InitConfig 初始化配置
func InitConfig() {
	configFile, err := os.OpenFile(ConfigFilePath, os.O_CREATE|os.O_RDONLY, LocalFileMode)
	if err != nil && !os.IsNotExist(err) {
		HandleError(err)
	}
	credFile, err := os.OpenFile(CredentialFilePath, os.O_CREATE|os.O_RDONLY, LocalFileMode)
	if err != nil && !os.IsNotExist(err) {
		HandleError(err)
	}

	AggConfigListIns, err = NewAggConfigManager(configFile, credFile)
	if err != nil {
		LogError(err.Error())
	} else {
		var ins *AggConfig
		if Global.Profile == "" {
			ins, err = AggConfigListIns.GetActiveAggConfig()
			if err != nil && len(AggConfigListIns.GetAggConfigList()) != 0 {
				HandleError(err)
			}
		} else {
			var ok bool
			ins, ok = AggConfigListIns.GetAggConfigByProfile(Global.Profile)
			if !ok {
				LogError("Profile %s does not exist", Global.Profile)
			}
		}

		if ins != nil {
			ConfigIns = ins
		} else {
			ins = ConfigIns
		}

		tmpIns := *ins
		tmpIns.PublicKey = MosaicString(tmpIns.PublicKey, 5, 5)
		tmpIns.PrivateKey = MosaicString(tmpIns.PrivateKey, 5, 5)
		LogInfo(fmt.Sprintf("load active config : %#v", tmpIns))
		bc, err := GetBizClient(ConfigIns)
		if err != nil {
			HandleError(err)
		} else {
			BizClient = bc
		}
	}
}

func init() {
	//配置日志
	err := initLog()
	if err != nil {
		fmt.Println(err)
	}
}
