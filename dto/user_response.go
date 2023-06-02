package dto

import "github.com/davepokpong/eticket-backend/models"

type UserSignInDto struct {
    Username 	string 			   `json:"username"`
    Password 	string 			   `json:"password"`
}

type UserRegisterDto struct {
	FirstName   string             `json:"firstname"`
	LastName	string			   `json:"lastname"`
	Username	string 			   `json:"username"`
	Email		string			   `json:"email"`
	Phone		string			   `json:"phone"`
	Password	string			   `json:"password"`
	Role		string			   `json:"role"`
	Members		int				   `json:"members"`
	Star 		int				   `json:"star"`
}

type UserDeleteDto struct {
    Username    string 			   `json:"username"`
}

type UserUpdateDto struct {
	FirstName   string             `json:"firstname"`
	LastName	string			   `json:"lastname"`
	Username	string 			   `json:"username"`
	Email		string			   `json:"email"`
	Phone		string			   `json:"phone"`
	Members		int				   `json:"members"`
	Star 		int				   `json:"star"`
	Disable		bool			   `json:"disable"`
}

type UserPasswordUpdateDto struct {
	Username    string			   `json:"username"`
	Password    string			   `json:"password"`
	NewPassword string			   `json:"newPassword"`
}

type PostUserDetail struct {
	To			string     		   `json:"to"`
	Code		string			   `json:"code"`
	Name		string 			   `json:"name"`
	Password	string			   `json:"password"`
	Lang		string			   `json:"lang"`
	Username	string			   `json:"username"`
}

type WeekDayStat struct {
	Sun			int				   `json:"sun"`
	Mon			int				   `json:"mon"`
	Tue			int				   `json:"tue"`
	Wed			int				   `json:"wed"`
	Thu			int				   `json:"thu"`
	Fri			int				   `json:"fri"`
	Sat			int				   `json:"sat"`
}

type StatPerMonth struct {
	Jan			int				   `json:"jan"`
	Feb			int				   `json:"feb"`
	Mar			int				   `json:"mar"`
	Apr			int				   `json:"apr"`
	May			int				   `json:"may"`
	Jun			int				   `json:"jun"`
	Jul			int				   `json:"jul"`
	Aug			int				   `json:"aug"`
	Sep			int				   `json:"sep"`
	Oct			int				   `json:"oct"`
	Nov			int				   `json:"nov"`
	Dec			int				   `json:"dec"`
}

type CountMemberOfUser struct {
	Member		int				   `json:"member"`
	Count 		int				   `json:"count"`
}

type UserActivityListDto struct {
	Activity  	[]models.ActivityList `json:"activity"`
}
