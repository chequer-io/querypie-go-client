package dac_connection

import (
	"errors"
	"gorm.io/gorm"
	"qpc/config"
	"qpc/model"
	"qpc/utils"

	"github.com/sirupsen/logrus"
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

func (c *ConnectionV2) Save() *ConnectionV2 {
	// Attempt to update the detailed connection
	result := config.LocalDatabase.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(c)
	logrus.Debugf("Updated detailed connection %s: %v", c.ShortID(), result)

	// If no rows were affected, create a new detailed connection
	if result.RowsAffected == 0 {
		logrus.Debugf("RowsAffected == 0")
		if err := config.LocalDatabase.Create(c).Error; err != nil {
			logrus.Errorf("Failed to save detailed connection %s: %v", c.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Errorf("Failed to update detailed connection %s: %v", c.ShortID(), result.Error)
	} else {
		logrus.Debugf("RowsAffected: %d", result.RowsAffected)
	}
	config.LocalDatabase.Save(c)
	logrus.Debugf("Saved detailed connection %s: %v", c.ShortID(), result)
	return c
}

func (c *ConnectionV2) FetchByUuid(uuid string) *ConnectionV2 {
	conn, err := utils.Fetch(
		"/api/external/v2/dac/connections/"+uuid,
		&ConnectionV2{},
	)
	if err != nil {
		logrus.Fatalf("Failed to fetch a resource: %v", err)
	} else {
		logrus.Debugf("Fetched connection: %v", *conn)
	}
	return conn
}

func (c *ConnectionV2) AndSave() *ConnectionV2 {
	if c.HttpResponse == nil {
		logrus.Debugf("No HttpResponse found, save detailed connection")
		c.Save()
	} else if c.HttpResponse.StatusCode() == 200 {
		logrus.Debugf("Got 200 OK, save detailed connection")
		c.Save()
	} else {
		logrus.Debugf("None-200 http status code: %d, skip saving", c.HttpResponse.StatusCode())
	}
	logrus.Debugf("Done saving detailed connection: %v", c.ShortID())
	return c
}

func (c *ConnectionV2) FindByUuid(uuid string) *ConnectionV2 {
	var connection ConnectionV2
	result := config.LocalDatabase.
		Model(&ConnectionV2{}).
		Where("uuid = ?", uuid).
		Preload("Clusters").
		Preload("ConnectionOwners").
		Preload("ConnectionOwners.Role").
		Preload("ConnectionOwners.OwnedBy").
		First(&connection)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logrus.Errorf("Connection not found: %s", uuid)
		} else {
			logrus.Fatalf("Failed to find a connection: %s", uuid)
		}
		return nil
	}
	logrus.Debugf("Found: %v", connection)
	return &connection
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&SummarizedConnectionV2{},
		&model.Role{},
		&OwnedBy{},
		&ConnectionOwner{},
		&Cluster{},
		&ConnectionV2{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
