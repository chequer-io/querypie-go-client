package dac_policy

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"qpc/config"
	"qpc/utils"
)

func (r *SensitiveDataRule) FetchAllOfPolicyAndForEach(
	policy *Policy,
	forEachFunc func(fetched *SensitiveDataRule) bool,
) {
	uri := fmt.Sprintf(
		"/api/external/policies/%s/rules",
		policy.Uuid,
	)
	utils.FetchListAndForEach[SensitiveDataRule, []SensitiveDataRule](
		uri,
		[]SensitiveDataRule{},
		forEachFunc,
	)
}

func (r *SensitiveDataRule) FindAllOfPolicyAndForEach(
	policy *Policy,
	forEachFunc func(fetched *SensitiveDataRule) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&SensitiveDataRule{}).
				Where("policy_uuid = ?", policy.Uuid).
				Count(total)
		},
		func(tx *gorm.DB, items *[]SensitiveDataRule) *gorm.DB {
			return tx.Model(&SensitiveDataRule{}).
				Preload("Policy.Connection").
				Preload("Policy.UpdatedBy").
				Preload("Policy.CreatedBy").
				Preload("Policy").
				Where("policy_uuid = ?", policy.Uuid).
				Find(items)
		},
		forEachFunc,
	)
}

func (r *SensitiveDataRule) SaveAndLoad() *SensitiveDataRule {
	db := config.LocalDatabase.Save(r)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	db.
		Preload("Policy.Connection").
		Preload("Policy.UpdatedBy").
		Preload("Policy.CreatedBy").
		Preload("Policy").
		First(r, "uuid = ?", r.Uuid)
	return r
}
