package campaign

import "time"

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Descripton       string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	//membuat relasi campaign dengan campaign_images
	CampaignImages []CampaignImage
}

type CampaignImage struct {
	ID 				int
	CampaignID 		int
	FileName 		string
	IsPrimary 		int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}