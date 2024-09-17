package dac_connection

import "fmt"

func PrintHeaderOfDetailedConnection() {
	(&SummarizedConnectionV2{}).PrintHeader()
	(&ConnectionV2{}).PrintHeader("DETAILED  ")
	(&Cluster{}).PrintHeader("CLUSTER  ")
	fmt.Println()
}

func (sc *SummarizedConnectionV2) FindDetailedConnectionAndPrint() bool {
	conn := (&ConnectionV2{}).FindByUuid(sc.Uuid)
	conn.Print("DETAILED  ")
	for _, cluster := range conn.Clusters {
		cluster.Print("CLUSTER  ")
	}
	fmt.Println()
	return true // OK to continue fetching
}

func (sc *SummarizedConnectionV2) FetchDetailedConnectionAndPrintAndSave() bool {
	sc.Print()
	conn := (&ConnectionV2{}).FetchByUuid(sc.Uuid)
	conn.Print("DETAILED  ").Save()
	for _, cluster := range conn.Clusters {
		cluster.Print("CLUSTER  ")
	}
	fmt.Println()
	return true // OK to continue fetching
}
