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
	if sc == nil {
		logrus.Debugf("Did not save it as it is nil")
		return sc
	} else if sc.CreatedByUuid == "" {
		sc.CreatedByUuid = sc.CreatedBy.Uuid
	} else if sc.UpdatedByUuid == "" {
		sc.UpdatedByUuid = sc.UpdatedBy.Uuid
	}
	// NOTE Don’t use Save with Model, it’s an Undefined Behavior.
	// https://gorm.io/docs/update.html#Save
	db := config.LocalDatabase.Save(sc)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	return sc
}

func (c *ConnectionV2) FindAllAndForEach(
	forEachFunc func(found *ConnectionV2) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&ConnectionV2{}).Count(total)
		},
		func(tx *gorm.DB, items *[]ConnectionV2) *gorm.DB {
			return tx.Model(&ConnectionV2{}).
				Preload("UpdatedBy").
				Preload("CreatedBy").
				Preload("Clusters").
				Preload("ConnectionOwners").
				Preload("ConnectionOwners.Role").
				Preload("ConnectionOwners.OwnedBy").
				Find(items)
		},
		forEachFunc,
	)
}

func (c *ConnectionV2) Save() *ConnectionV2 {
	if c == nil {
		logrus.Debugf("Did not save it as it is nil")
		return c
	} else if c.CreatedByUuid == "" {
		c.CreatedByUuid = c.CreatedBy.Uuid
	} else if c.UpdatedByUuid == "" {
		c.UpdatedByUuid = c.UpdatedBy.Uuid
	}
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
		// Early exit to prevent further processing
	}
	conn.Uuid = uuid // Uuid might be missing when 500 Internal Server Error occurs
	return conn
}

func (c *ConnectionV2) SaveAlsoForServerError() *ConnectionV2 {
	if c.HttpResponse.IsSuccess() {
		c.Save()
	} else if utils.IsServerError(c.HttpResponse) {
		// When the API returns a server error,
		// it is necessary to save an empty object with the Uuid,
		// as a workaround to prevent error due to missing entities in local database.
		empty := &ConnectionV2{Uuid: c.Uuid}
		empty.Save()
	} else if utils.IsClientError(c.HttpResponse) {
		// When the API returns a client error such as 404 Not Found,
		// it should not save the object to the local database.
		logrus.Debugf("Did not save it as it is a client error: %s", c.HttpResponse.Status())
	} else {
		logrus.Fatalf("Unexpected error: %s", c.HttpResponse.Status())
	}
	return c
}

func (c *ConnectionV2) FirstByUuid(uuid string) *ConnectionV2 {
	var item ConnectionV2
	return utils.First[ConnectionV2](&item, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&ConnectionV2{}).
			Preload("UpdatedBy").
			Preload("CreatedBy").
			Preload("Clusters").
			Preload("ConnectionOwners").
			Preload("ConnectionOwners.Role").
			Preload("ConnectionOwners.OwnedBy").
			Where("uuid = ?", uuid).
			First(&item)
	})
}

func (c *ConnectionV2) FindByNameOrUuid(query string, connections *[]ConnectionV2) {
	utils.FindMultiple(connections, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&ConnectionV2{}).
			Preload("Clusters").
			Preload("Clusters.Connection").
			// Note: Column names are snake_case in the database.
			Where("name = ? OR uuid = ?", query, query).
			Find(connections)
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
		&model.Modifier{},
		&SummarizedConnectionV2{},

		&model.Modifier{},
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
