package dac_policy

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"qpc/config"
	"qpc/model"
	"qpc/utils"
)

func (p *Policy) FetchByUuid(uuid string) *Policy {
	policy, err := utils.Fetch(
		"/api/external/policies/"+uuid,
		&Policy{},
	)
	if err != nil {
		logrus.Fatalf("Failed to fetch a resource: %v", err)
		// Early exit to prevent further processing
	}
	return policy
}

func (p *Policy) FindByConnectionAndNameAndUuid(
	connection string,
	policyName string,
	uuid string,
	policies *[]Policy,
) {
	if len(connection) == 0 {
		connection = "%" // Match all
	} else {
		connection = "%" + connection + "%"
	}
	if len(policyName) == 0 {
		policyName = "%" // Match all
	} else {
		policyName = "%" + policyName + "%"
	}
	if len(uuid) == 0 {
		uuid = "%" // Match all
	} else {
		uuid = "%" + uuid + "%"
	}

	utils.FindMultiple(policies, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&Policy{}).
			Joins("JOIN connection_v2 ON policies.cluster_group_uuid = connection_v2.uuid").
			Where("connection_v2.name LIKE ? AND (policies.name LIKE ? AND policies.uuid LIKE ?)", connection, policyName, uuid).
			Preload("UpdatedBy").
			Preload("CreatedBy").
			Preload("Connection").
			Find(policies)
	})
}

func (p *Policy) FirstByClusterGroupUuidAndPolicyType(uuid string, policyType PolicyType) *Policy {
	policy := &Policy{}
	where := "cluster_group_uuid = ? AND policy_type = ?"
	return utils.First(policy, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&Policy{}).
			// Note: Column names are snake_case in the database.
			Where(where, uuid, policyType).
			Preload("UpdatedBy").
			Preload("CreatedBy").
			Preload("Connection").
			First(policy)
	})
}

func (p *Policy) FetchAllOfPolicyTypeAndForEach(
	policyType PolicyType,
	forEachFunc func(fetched *Policy) bool,
) {
	utils.FetchPagedListAndForEach(
		"/api/external/policies?policyType="+string(policyType),
		&PolicyPagedList{},
		func(page *PolicyPagedList) bool {
			for _, it := range page.List {
				forEachFunc(&it)
			}
			return true
		},
	)
}

func (p *Policy) FindAllOfPolicyTypeAndForEach(
	policyType PolicyType,
	forEachFunc func(found *Policy) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&Policy{}).
				Where("policy_type = ?", policyType).
				Count(total)
		},
		func(tx *gorm.DB, items *[]Policy) *gorm.DB {
			return tx.Model(&Policy{}).
				Preload("UpdatedBy").
				Preload("CreatedBy").
				Preload("Connection").
				Where("policy_type = ?", policyType).
				Find(items)
		},
		forEachFunc,
	)
}

func (p *Policy) Save() *Policy {
	// NOTE Don’t use Save with Model, it’s an Undefined Behavior.
	// https://gorm.io/docs/update.html#Save
	db := config.LocalDatabase.Save(p)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	return p
}

func (p *Policy) SaveAndLoad() *Policy {
	db := config.LocalDatabase.Save(p)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	db.
		Preload("UpdatedBy").
		Preload("CreatedBy").
		Preload("Connection").
		First(p, "uuid = ?", p.Uuid)
	return p
}

func (p *Policy) Delete() *Policy {
	db := config.LocalDatabase.Delete(p)
	logrus.Debugf("Deleted it, RowsAffected: %d", db.RowsAffected)
	return p
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&model.Modifier{},
		&Policy{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
