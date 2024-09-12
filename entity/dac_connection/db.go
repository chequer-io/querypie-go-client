package dac_connection

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
	"qpc/model"
	"qpc/utils"
)

func (cl *PagedConnectionV2List) Save() {
	for _, sac := range cl.GetList() {
		sac.Save()
	}
}

// FetchAllAndPrintAndSave /*
// Q: Why is this function defined as a method of PagedConnectionV2List?
// A: Actually the receiver of this method, `cl`, is not used in the function.
// Although FetchAllAndPrintAndSave could be defined as a static function,
// Go does not support static functions. Therefore, it is defined as a method of PagedConnectionV2List.
// It is implemented as a method of PagedConnectionV2List to provide a shorter function name,
// improve readability, and highlight its association with the Entity model.
func (cl *PagedConnectionV2List) FetchAllAndPrintAndSave() {
	utils.FetchPrintAndSave(
		"/api/external/v2/dac/connections",
		&PagedConnectionV2List{},
		func(result *PagedConnectionV2List, first bool, last bool) {
			result.Print()
		},
		func(result *PagedConnectionV2List) bool {
			result.Save()
			return !result.Page.HasNext()
		},
	)
}

func (cl *PagedConnectionV2List) FindAllAndPrint() {
	var total, fetched int64 = 0, 0
	result := config.LocalDatabase.Model(&SummarizedConnectionV2{}).Count(&total)
	if result.Error != nil {
		logrus.Fatalf("Failed to count dac connections: %v", result.Error)
	}
	logrus.Debugf("Found %d dac connections", total)

	page := 0
	size := 30 // Set the desired page size

	for {
		list, err := cl.FindAllAsPagedList(page, size, int(total))
		if err != nil {
			logrus.Fatalf("Failed to select data from local database: %v", err)
		}
		logrus.Debugf("Selected %d, page %d, size %d, total %d",
			len(list.List), page, size, total)
		fetched += int64(len(list.List))
		list.Print()

		if !list.Page.HasNext() {
			break
		}
		page++
	}
	logrus.Debugf("Selected %d, whereas total count was %d, difference: %d",
		fetched, total, total-fetched)

}

func (cl *PagedConnectionV2List) FindAllAsPagedList(currentPage, pageSize, totalElements int) (PagedConnectionV2List, error) {
	var list PagedConnectionV2List

	var page model.Page
	page.CurrentPage = currentPage
	page.PageSize = pageSize
	page.TotalElements = totalElements
	page.TotalPages = (totalElements + pageSize - 1) / pageSize

	var connections []SummarizedConnectionV2
	offset := currentPage * pageSize
	result := config.LocalDatabase.
		Offset(offset).
		Limit(pageSize).
		Find(&connections)
	if result.Error != nil {
		return PagedConnectionV2List{}, result.Error
	}
	list.List = connections
	list.Page = page
	return list, nil
}

func (sc *SummarizedConnectionV2) Save() {

	// Attempt to update the user
	result := config.LocalDatabase.Model(&SummarizedConnectionV2{}).Where("uuid = ?", sc.Uuid).Updates(&sc)

	// If no rows were affected, create a new user
	if result.RowsAffected == 0 {
		if err := config.LocalDatabase.Create(&sc).Error; err != nil {
			logrus.Errorf("Failed to save connection %s: %v", sc.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Errorf("Failed to update connection %s: %v", sc.ShortID(), result.Error)
	}
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&SummarizedConnectionV2{},
		&model.Role{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
