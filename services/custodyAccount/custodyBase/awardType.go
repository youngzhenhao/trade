package custodyBase

var AwardTypes = make(map[string]string)

func init() {
	AwardTypes["tip1"] = "推广大使奖励"
	AwardTypes["tip2"] = "直推奖励"
	AwardTypes["tip3"] = "收藏奖励"
	AwardTypes["tip5"] = "腾讯会议奖励"
	AwardTypes["tip6"] = "补发奖励"
	AwardTypes["activity"] = "进群奖励"
	AwardTypes["tip22"] = "领取资产（直推）"
	AwardTypes["tip21"] = "领取资产（推广）"
	AwardTypes["tip32"] = "购买预售NFT奖励(直推)"
	AwardTypes["tip31"] = "购买预售NFT奖励(推广)"
	AwardTypes["tip41"] = "质押收益"
	AwardTypes["swapLP"] = "LP奖励提现"

}

func GetAwardType(award string) string {
	value, exists := AwardTypes[award]
	if exists {
		return value
	} else {
		return award
	}
}
