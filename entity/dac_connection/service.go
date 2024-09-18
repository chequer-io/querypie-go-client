package dac_connection

import "fmt"

func PrintHeaderOfDetailedConnection() {
	(&SummarizedConnectionV2{}).PrintHeader()
	(&ConnectionV2{}).PrintHeader("DETAILED  ")
	(&Cluster{}).PrintHeader("CLUSTER  ")
	fmt.Println()
}

func (sc *SummarizedConnectionV2) FirstDetailedConnectionAndPrint() bool {
	conn := (&ConnectionV2{}).FirstByUuid(sc.Uuid)
	conn.Print("DETAILED  ")
	for _, cluster := range conn.Clusters {
		cluster.Print("CLUSTER  ")
	}
	fmt.Println()
	return true // OK to continue fetching
}

func (sc *SummarizedConnectionV2) FetchDetailedConnectionAndPrintAndSave() bool {
	conn := (&ConnectionV2{}).FetchByUuid(sc.Uuid)
	conn.Print("DETAILED  ").SaveAlsoForServerError()
	for _, cluster := range conn.Clusters {
		cluster.Print("CLUSTER  ")
	}
	fmt.Println()
	return true // OK to continue fetching
}
