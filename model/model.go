package model

type Event struct {
	Tag 				string
	Name 				string
	Frontends 			string
	Exercises 			string
	Available 			uint
	Capacity 			uint
	StartedAt			string
	ExpectedFinishTime 	string
	FinishedAt			string
}

type Team struct {
	Id 					string
	EventTag  			string
	Email 				string
	Name 				string
	Password		    string
	CreatedAt		    string
	LastAccess  		string
	SolvedChallenges	string
}