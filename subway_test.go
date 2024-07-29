package subway

import (
	"github.com/simonks2016/Subway/Filter"
	"github.com/simonks2016/Subway/Relationship"
	"github.com/simonks2016/Subway/Sorter"
	"testing"
)

type Customer struct {
	Name            string                           `json:"name"`
	Sex             int                              `json:"sex"`
	BrandName       string                           `json:"brand_name"`
	Follow          *Relationship.ManyRefs[Customer] `json:"follow"`
	Fans            *Relationship.ManyRefs[Customer] `json:"fans"`
	FavoriteVideos  *Relationship.ManyRefs[Video]    `json:"favorite_videos"`
	FavoriteProgram *Relationship.ManyRefs[Program]  `json:"favorite_program"`
	Videos          *Relationship.ManyRefs[Video]    `json:"videos"`
	Program         *Relationship.ManyRefs[Program]  `json:"program"`
}

type Video struct {
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Id          string                      `json:"id"`
	Creator     *Relationship.Ref[Customer] `json:"creator"`
	Tags        *Relationship.ManyRefs[Tag] `json:"tags"`
	State       *Filter.FilterField[int]    `json:"state"`
	CreateTime  *Sorter.SortField           `json:"create_time"`
	Click       *Sorter.SortField           `json:"click"`
	Viewers     int64                       `json:"viewers"`
}

type Program struct {
	Id          string                      `json:"id"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Creator     *Relationship.Ref[Customer] `json:"creator"`
	Tags        *Relationship.ManyRefs[Tag] `json:"tags"`
}

type Tag struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type TimeLine struct {
	Id         string                      `json:"id"`
	Content    *TimeLineContent            `json:"content"`
	UserId     *Filter.FilterField[string] `json:"user_id"`
	CreateTime *Sorter.SortField           `json:"create_time"`
}

type TimeLineCreator struct {
	CreatorName string `json:"creator_name"`
}
type TimeLineContent struct {
	Text            string           `json:"text"`
	TimeLineCreator *TimeLineCreator `json:"time_line_creator"`
}

func TestNewSubway(t *testing.T) {

	NewRedisConnWithSubway("127.0.0.1:6379", "root", "")

}
