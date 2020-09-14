package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zlingqu/es-clean/es"
)

var (
	ip          string
	port        string
	indexName   string
	keepTimeDay int
)

//NewEsCleanCommand cobra实例化
func NewEsCleanCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "es-clean",
		Short:   "es 索引清理",
		Long:    "es-clean 用于清理es中的索引，以释放存储资源",
		Example: "es-clean --ip 1.1.1.1  --port 9200 --indexName k8s-dev* --keepTimeDay 200",
		Run: func(cmd *cobra.Command, args []string) {
			esClient := es.NewClient(ip, port, indexName, keepTimeDay)
			allIndex, _ := esClient.GetAllIndex()
			esClient.DeleteIndex(allIndex, indexName, keepTimeDay)
		},
	}
	rootCmd.Flags().StringVarP(&ip, "ip", "i", "", "例如：1.1.1.1")
	rootCmd.Flags().StringVarP(&port, "port", "p", "", "端口，例如：9200")
	rootCmd.Flags().StringVarP(&indexName, "indexName", "n", "", "密码")
	rootCmd.Flags().IntVarP(&keepTimeDay, "keepTimeDay", "k", 0, "保留索引的天数，单位是天，比如60")
	return rootCmd
}
