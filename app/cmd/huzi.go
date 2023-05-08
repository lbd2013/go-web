package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-ego/gse"
	"github.com/spf13/cobra"
	"goweb/pkg/console"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var CmdHuzi = &cobra.Command{
	Use:   "huzi",
	Short: "huzi",
	Run:   runHuzi,
	Args:  cobra.NoArgs, // 不允许传参
}

type huziJson struct {
	Guild struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		IconUrl string `json:"iconUrl"`
	} `json:"guild"`
	Channel struct {
		Id         string      `json:"id"`
		Type       string      `json:"type"`
		CategoryId string      `json:"categoryId"`
		Category   string      `json:"category"`
		Name       string      `json:"name"`
		Topic      interface{} `json:"topic"`
	} `json:"channel"`
	DateRange struct {
		After  time.Time   `json:"after"`
		Before interface{} `json:"before"`
	} `json:"dateRange"`
	Messages []struct {
		Id                 string      `json:"id"`
		Type               string      `json:"type"`
		Timestamp          time.Time   `json:"timestamp"`
		TimestampEdited    interface{} `json:"timestampEdited"`
		CallEndedTimestamp interface{} `json:"callEndedTimestamp"`
		IsPinned           bool        `json:"isPinned"`
		Content            string      `json:"content"`
		Author             struct {
			Id            string      `json:"id"`
			Name          string      `json:"name"`
			Discriminator string      `json:"discriminator"`
			Nickname      string      `json:"nickname"`
			Color         interface{} `json:"color"`
			IsBot         bool        `json:"isBot"`
			AvatarUrl     string      `json:"avatarUrl"`
		} `json:"author"`
		Attachments []interface{} `json:"attachments"`
		Embeds      []interface{} `json:"embeds"`
		Stickers    []interface{} `json:"stickers"`
		Reactions   []interface{} `json:"reactions"`
		Mentions    []struct {
			Id            string `json:"id"`
			Name          string `json:"name"`
			Discriminator string `json:"discriminator"`
			Nickname      string `json:"nickname"`
			IsBot         bool   `json:"isBot"`
		} `json:"mentions"`
		Reference struct {
			MessageId string `json:"messageId"`
			ChannelId string `json:"channelId"`
			GuildId   string `json:"guildId"`
		} `json:"reference,omitempty"`
	} `json:"messages"`
	MessageCount int `json:"messageCount"`
}

var (
	seg gse.Segmenter
)

func SegmentCutSearchMode(text string) []string {
	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中 ToString 函数的注释。
	// 搜索模式主要用于给搜索引擎提供尽可能多的关键字
	// seg.String, seg.Slice 输出的类型不同
	return seg.Slice(text, true)
}

func initSeg() {
	// Skip log print
	seg.SkipLog = true

	// 加载简体中文词典
	err := seg.LoadDict("zh_s")
	if err != nil {
		fmt.Println(err)
	}

	err = seg.LoadDictEmbed("zh_s")
	if err != nil {
		fmt.Println(err)
	}

	// 加载繁体中文词典
	err = seg.LoadDict("zh_t")
	if err != nil {
		fmt.Println(err)
	}

	// 载入词典
	path, _ := os.Getwd()
	err = seg.LoadDict(path + "/storage/jieba/dictionary.txt")
	if err != nil {
		fmt.Println(err)
	}
}

func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

func runHuzi(cmd *cobra.Command, args []string) {
	console.Log("huzi...")

	// 初始化分词器
	initSeg()

	// 打开json文件,跟main.go文件并列，或者使用绝对路径
	path, _ := os.Getwd()
	jsonFile, err := os.Open(path + "/storage/jieba/huzi.json")

	// 最好要处理以下错误
	if err != nil {
		panic(err)
	}

	// 要记得关闭
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var jsonData huziJson
	json.Unmarshal([]byte(byteValue), &jsonData)

	var wordMap map[string]int /*创建集合 */
	wordMap = make(map[string]int)

	for _, message := range jsonData.Messages {
		rs := SegmentCutSearchMode(message.Content)
		for i := 0; i < len(rs); i++ {
			word := strings.TrimSpace(rs[i])
			word = FilterEmoji(word)
			// 过滤符号 、 中文字符
			if len(word) <= 2 ||
				word == "！" || word == "@" || word == "#" || word == "￥" || word == "%" || word == "……" || word == "&" || word == "*" ||
				word == "（" || word == "）" || word == "——" || word == "-" || word == "=" || word == "【" || word == "】" || word == "；" ||
				word == "‘" || word == "，" || word == "。" || word == "、" || word == "？" || word == "�" || word == "：" || word == "｜" ||
				word == "·" || word == "—" || word == "”" || word == "“" || word == "「" || word == "」" || word == "…" || word == "jb" || word == "Ｏ" || word == "♂" ||
				word == "》" || word == "《" || word == "￣" || word == "～" || word == "へ" {
				continue
			}

			// 过滤特殊单词
			if word == "jpg" || word == "png" || word == "gif" || word == "https" || word == "http" || word == "com" || word == "www" ||
				word == "vip" || word == "qq" || word == "lol" || word == "1u" || word == "100u" || word == "10u" || word == "20u" || word == "200u" || word == "2000u" ||
				word == "gm" || word == "id" || word == "xxx" || word == "what" || word == "about" || word == "12u" || word == "of" ||
				word == "tenor" || word == "view" || word == "dk" || word == "rmb" || word == "one" || word == "everyone" ||
				word == "b5" || word == "b7" || word == "5w" || word == "e5" || word == "e6" || word == "e7" || word == "15f" || word == "18wu" || word == "discord" ||
				word == "youtube" || word == "you" || word == "it" || word == "channels" || word == "e9" || word == "mint" || word == "v神" || word == "20k" ||
				word == "u本位" || word == "xdm" || word == "500w" || word == "80u" || word == "kc0526" || word == "1000w" ||
				word == "notion" || word == "dahuzi" || word == "twitter" || word == "the" {
				continue
			}

			// 过滤纯数字
			if _, err1 := strconv.Atoi(word); err1 == nil {
				continue
			}

			// 过滤人名
			if word == "wx202205021204212x" || word == "0u1274945486b2802139570c26" || word == "liontungalex" ||
				word == "0a1dc75257c3495cbdf9cb89828f1b73" || word == "kc" || word == "220510discordemoji2" {
				continue
			}

			// 过滤单个字
			isChinese := regexp.MustCompile("^[\u4e00-\u9fa5]") //我们要匹配中文的匹配规则
			if len(word) == 3 && isChinese.MatchString(word) {
				continue
			}

			// 过滤中文
			if isChinese.MatchString(word) {
				continue
			}

			_, ok := wordMap[word]
			if ok {
				wordMap[word] = wordMap[word] + 1
			} else {
				wordMap[word] = 1
			}
		}
	}

	keys := make([]string, 0, len(wordMap))

	for key := range wordMap {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return wordMap[keys[i]] < wordMap[keys[j]]
	})

	for _, k := range keys {
		fmt.Printf("%s\t%d\n", k, wordMap[k])
	}

	return
}
