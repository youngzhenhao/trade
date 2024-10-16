package custodyBase

var AwardTypes = make(map[string]string)

func init() {
	AwardTypes["tip1"] = "推广大使奖励"
	AwardTypes["tip2"] = "直推奖励"
	AwardTypes["tip3"] = "收藏奖励"
	AwardTypes["tip5"] = "腾会奖励"
	AwardTypes["activity"] = "进群奖励"
}

func GetAwardType(award string) string {
	value, exists := AwardTypes[award]
	if exists {
		return value
	} else {
		return award
	}
}
