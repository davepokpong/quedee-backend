package dto

type CreateQueueDto struct {
	Username     string `json:"username"`
	ActivityCode string `json:"activityCode"`
	Status       string `json:"status"`
	Size         int    `json:"size"`
	Star         int    `json:"star"`
}

type QueueDeleteDto struct {
	ID string `json:"_id"`
}

type QueueUpdateDto struct {
	ID           string `json:"_id"`
	Username     string `json:"username"`
	ActivityCode string `json:"activityCode"`
	QueueNo      int    `json:"queueNo"`
	Round        int    `json:"round"`
	Disable      bool   `json:"disable"`
	Status       string `json:"status"`
	Size         int    `json:"size"`
}

type QueueDisableDto struct {
	QueueId      []string `json:"queueId"`
	ActivityCode string   `json:"activityCode"`
}

type QueueCancelDto struct {
	QueueId 	 string   `json:"queueId"`
}

type ActivityPlayedTime struct {
	Code		 string   `json:"code"`
	Name		 []string `json:"name"`
	Count        int      `json:"count"`
}

type CountQueueDto struct {
	All 		 int64 	  `json:"all"`
	One 		 int64	  `json:"one"`   
}

type DataCountStatDto struct {
	CustomerNow	 int64	  `json:"customerNow"`
	CustomerAll	 int64	  `json:"customer"`
	StaffNow	 int64	  `json:"staffNow"`
	StaffAll	 int64	  `json:"staff"`
	ActivityNow	 int64	  `json:"activityNow"`
	ActivityAll	 int64	  `json:"activity"`
}

type ActiveQueueDto struct {
	QueueId 	 string	  `json:"queueId"`
	QueueSize	 int	  `json:"queueSize"`
	Round		 int	  `json:"round"`
}

type CreateQueueSpecificRoundDto struct {
	Username     string `json:"username"`
	ActivityCode string `json:"activityCode"`
	Status       string `json:"status"`
	Size         int    `json:"size"`
	Star         int    `json:"star"`
	Round		 int 	`json:"round"`
}