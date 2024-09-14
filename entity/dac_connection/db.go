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

// FetchAllAndForEach /*
// Q: Why is this function defined as a method of SummarizedConnectionV2?
// A: Actually the receiver of this method, `cl`, is not used in the function.
// Although FetchAllAndForEach could be defined as a static function,
// Go does not support static functions. Therefore, it is defined as a method of SummarizedConnectionV2.
// It is implemented as a method of SummarizedConnectionV2 to provide a shorter function name,
// improve readability, and highlight its association with the Entity model.
func (sc *SummarizedConnectionV2) FetchAllAndForEach(
	forEachFunc func(sc *SummarizedConnectionV2) bool,
) {
	utils.FetchPagedListAndForEach(
		"/api/external/v2/dac/connections",
		&PagedConnectionV2List{},
		func(page *PagedConnectionV2List) bool {
			for _, sc := range page.List {
				forEachFunc(&sc)
			}
			return true
		},
	)
}

func (sc *SummarizedConnectionV2) FindAllAndForEach(
	forEachFunc func(sc *SummarizedConnectionV2) bool,
) {
	var connections []SummarizedConnectionV2
	var total, fetched int64 = 0, 0
	result := config.LocalDatabase.Model(&SummarizedConnectionV2{}).Count(&total)
	if result.Error != nil {
		logrus.Fatalf("Failed to count dac connections: %v", result.Error)
	}
	logrus.Debugf("Found %d dac connections", total)

	result = config.LocalDatabase.Find(&connections)
	if result.Error != nil {
		logrus.Fatalf("Failed to select data from local database: %v", result.Error)
		return
	}
	fetched = int64(len(connections))
	for _, sc := range connections {
		forEachFunc(&sc)
	}
	if fetched != total {
		logrus.Errorf("Selected %d, whereas total count was %d, difference: %d",
			fetched, total, total-fetched)
	}
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
