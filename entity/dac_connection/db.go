package dac_connection

import (
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
	var item ConnectionV2
	return utils.First[ConnectionV2](&item, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&ConnectionV2{}).
			Preload("Clusters").
			Preload("ConnectionOwners").
			Preload("ConnectionOwners.Role").
			Preload("ConnectionOwners.OwnedBy").
			Where("uuid = ?", uuid).
			First(&item)
	})
}

func (c *Cluster) FindAllAndForEach(
	forEachFunc func(found *Cluster) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&Cluster{}).Count(total)
		},
		func(tx *gorm.DB, items *[]Cluster) *gorm.DB {
			return tx.Model(&Cluster{}).
				Preload("Connection").
				Find(items)
		},
		forEachFunc,
	)
}

func (c *Cluster) FindByHostAndPort(query string, clusters *[]Cluster) {
	utils.FindMultiple(clusters, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&Cluster{}).
			Preload("Connection").
			// Note: Column names are snake_case in the database.
			Where("CONCAT(host, ':', port) = ?", query).
			Find(clusters)

	})
}

func (c *Cluster) FindByCloudIdentifier(
	query string,
	clusters *[]Cluster,
) {
	utils.FindMultiple(clusters, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&Cluster{}).
			Preload("Connection").
			// Note: Column names are snake_case in the database.
			Where("cloud_identifier = ?", query).
			Find(clusters)
	})
}

func (c *Cluster) FirstByUuid(uuid string) *Cluster {
	var item Cluster
	return utils.First(&item, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&Cluster{}).
			Preload("Connection").
			Where("uuid = ?", uuid).
			First(&item)
	})
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
