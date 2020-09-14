package es

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Client      *http.Client //http.Clinet类型，结构体嵌套
	BaseURL     string
	IndexName   string
	KeepTimeDay int
}

// NewClient 函数
func NewClient(ip, port, indexName string, keepTimeDay int) *Client {
	client := &http.Client{}
	return &Client{
		Client:      client,
		IndexName:   indexName,
		BaseURL:     "http://" + ip + ":" + port,
		KeepTimeDay: keepTimeDay,
	}
}

// GetAllIndex 获取所有index
func (c *Client) GetAllIndex() (indexSliceSub [][]string, err error) {
	resp, err := c.Client.Get(c.BaseURL + "/_cat/indices")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	indexSlice := strings.Fields(string(body)) //Fields 用一个或多个连续的空格分隔字符串,返回切片
	for i := 2; i < len(indexSlice); i += 10 {
		if len(indexSlice[i]) < 15 || !strings.Contains(indexSlice[i], "k8s-") {
			continue
		}
		indexTime := indexSlice[i][len(indexSlice[i])-10:]
		indexSliceSub = append(indexSliceSub, []string{indexSlice[i], indexTime}) //组成一个二维数组，包含索引名和日期
	}
	return indexSliceSub, nil
}

//DeleteIndex 删除索引
func (c *Client) DeleteIndex(indexSliceSub [][]string, indexName string, keepTimeDay int) (err error) {

	for _, index := range indexSliceSub {
		if indexName != "all" && !strings.Contains(index[0], indexName) { //indexName不是all，且索引名字不匹配，跳过
			fmt.Printf("%v* 不包含 %v, 跳过删除...\n", indexName, index[0])
			continue
		}
		indexTime, _ := time.Parse("2006-01-02", index[1]) //字符串转换为时间类型
		now := time.Now()
		subM := now.Sub(indexTime)
		if subM.Hours()/24 > float64(keepTimeDay) {
			request, err := http.NewRequest("DELETE", c.BaseURL+"/"+index[0], nil)
			if err != nil {
				continue
			}
			resp, err := c.Client.Do(request)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			fmt.Printf("%v* 【包含】 %v,已存在%.2f天,需要保留%3d天,  开始删除索引%v\n", indexName, index[0], subM.Hours()/24, keepTimeDay, index[0])
		} else {
			fmt.Printf("%v* 【包含】 %v,距离过期还剩%.2f天, 跳过删除...\n", indexName, index[0], float64(keepTimeDay)-subM.Hours()/24)
		}
	}
	return
}
