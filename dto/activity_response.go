package dto

import (
	"github.com/davepokpong/eticket-backend/models"
)

type CreateActivityDto struct {
	Code     string               `json:"code"`
	Status   string               `json:"status"`
	Name     []string             `json:"name"`
	Size     int                  `json:"size"`
	Duration int                  `json:"duration"`
	Star     int                  `json:"star"`
	Picture  string               `json:"picture"`
	Position []float32            `json:"position"`
	Rating   float32              `json:"rating"`
	Comment  []models.UserComment `json:"comment"`
}

type ActivityDeleteDto struct {
	Code string `json:"code"`
}

type ActivityUpdateDto struct {
	Code          string               `json:"code"`
	Status        string               `json:"status"`
	Name          []string             `json:"name"`
	Size          int                  `json:"size"`
	Duration      int                  `json:"duration"`
	Star          int                  `json:"star"`
	Rating        float32              `json:"rating"`
	Comment       []models.UserComment `json:"comment"`
	CommentNumber int                  `json:"commentNumber"`
	Picture       string               `json:"picture"`
	Position      []float32            `json:"position"`
}

type UpdateActivityComment struct {
	Code    string             `json:"code"`
	Comment models.UserComment `json:"comment"`
}

// type UpdateActivityQueue struct {
// 	Code 	 		string				 `json:"code"`
// 	QueueNo  		int				  	 `json:"queueNo"`
// 	Round	 		int				     `json:"round"`
// 	SumQueue 		int				     `json:"sumQueue"`
// }

type ActivityWaitRoundDto struct {
	Activity  models.Activity `json:"activity"`
	WaitRound int             `json:"waitRound"`
}
