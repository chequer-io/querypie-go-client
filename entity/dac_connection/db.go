package dac_connection

import (
	"errors"
	"gorm.io/gorm"
	"qpc/config"
	"qpc/model"
	"qpc/utils"

	"github.com/sirupsen/logrus"
)

// FetchAllAndForEach /*
// Q: Why is this function defined as a method of SummarizedConnectionV2?
// A: Actually the receiver of this method, `cl`, is not used in the function.
// Although FetchAllAndForEach could be defined as a static function,
// Go does not support static functions. Therefore, it is defined as a method of SummarizedConnectionV2.
// It is implemented as a method of SummarizedConnectionV2 to provide a shorter function name,
// improve readability, and highlight its association with the Entity model.
func (sc *SummarizedConnectionV2) FetchAllAndForEach(
	forEachFunc func(fetched *SummarizedConnectionV2) bool,
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
	forEachFunc func(found *SummarizedConnectionV2) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&SummarizedConnectionV2{}).Count(total)
		},
		func(tx *gorm.DB, items *[]SummarizedConnectionV2) *gorm.DB {
			return tx.Model(&SummarizedConnectionV2{}).Find(items)
		},
		forEachFunc,
	)
}

func (sc *SummarizedConnectionV2) Save() *SummarizedConnectionV2 {
	// NOTE Don’t use Save with Model, it’s an Undefined Behavior.
	// https://gorm.io/docs/update.html#Save
	db := config.LocalDatabase.Save(sc)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	return sc
}

func (c *ConnectionV2) Save() *ConnectionV2 {
	// NOTE Don’t use Save with Model, it’s an Undefined Behavior.
	// https://gorm.io/docs/update.html#Save
	db := config.LocalDatabase.Save(c)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	return c
}

func (c *ConnectionV2) FetchByUuid(uuid string) *ConnectionV2 {
	conn, err := utils.Fetch(
		"/api/external/v2/dac/connections/"+uuid,
		&ConnectionV2{},
	)
	if err != nil {
		logrus.Fatalf("Failed to fetch a resource: %v", err)
	} else if conn.HttpResponse.IsError() {
		logrus.Warnf("Error from API: %s", conn.HttpResponse.Status())
		// When the API returns an error,
		// it is necessary to save an empty object with the Uuid,
		// so that following Save() can save a non-nil object.
		conn = &ConnectionV2{Uuid: uuid}
	}
	return conn
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
	} else {
		logrus.Debugf("AutoMigrate for dac_connection is done")
	}
}
