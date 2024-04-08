package subway

import (
	"fmt"
	"github.com/simonks2016/Subway/Filter"
	"github.com/simonks2016/Subway/Relationship"
	"github.com/simonks2016/Subway/Sorter"
	"testing"
	"time"
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
	Uid         *Filter.FilterField[string] `json:"uid"`
	IsPublic    *Filter.FilterField[int]
	CreateTime  *Sorter.SortField `json:"create_time"`
	Click       *Sorter.SortField `json:"click"`
	Viewers     int64             `json:"viewers"`
	Ids         []string          `json:"ids"`
	T1          []int             `json:"t1"`
	T2          []float64         `json:"t2"`
	T3          []bool            `json:"t3"`
	T4          []byte            `json:"t4"`
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

	var v1 = Video{
		Title:       "aaaa",
		Description: "bbb",
		Id:          "a2a2a2a2",
		Creator:     nil,
		Tags:        Relationship.NewManyRefs[Video, Tag]("a2a2a2a2", "Tags"),
		State:       Filter.NewFilter[int]("Video", "State", 1),
		Uid:         Filter.NewFilter[string]("Video", "Uid", "a1"),
		IsPublic:    Filter.NewFilter[int]("Video", "IsPublic", 1),
		CreateTime:  Sorter.NewSorter("Video", "CreateTime", float64(time.Now().Unix())),
		Click:       Sorter.NewSorter("Video", "CreateTime", 1),
		Viewers:     10,
		Ids:         []string{"a", "b", "c"},
		T1:          []int{0, 2, 1},
		T2:          []float64{0.15, 0.23},
		T3:          []bool{false, false},
		T4:          []byte("a"),
	}

	/*
		var v2 = Video{
			Title:       "aaaa",
			Description: "bbb",
			Id:          "a3a3a3a3",
			Creator:     Relationship.NewRef[Customer]("u1"),
			Tags:        Relationship.NewManyRefs[Video, Tag]("a3a3a3a3", "Tags"),
			State:       Filter.NewFilter[int]("Video", "State", 4),
			Uid:         Filter.NewFilter[string]("Video", "Uid", "a1"),
			IsPublic:    Filter.NewFilter[int]("Video", "IsPublic", 1),
			CreateTime:  Sorter.NewSorter("Video", "CreateTime", float64(time.Now().Unix())),
			Click:       Sorter.NewSorter("Video", "CreateTime", 10),
			Viewers:     10,
		}*/

	NewRedisConnWithSubway("127.0.0.1:6379", "root", "")

	exists, err := Exists[Video](v1.Id)
	if err != nil {
		return
	}
	fmt.Println(exists)
	/*
		read, err := Read[Video](v1.Id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(read.Creator.Query())*/
	/*
		err := Update[Video](v1, v1.Id)
		if err != nil {
			return
		}*/

	/*list, err := List[Video, int](&QueryRequest[int]{
		FieldName:   "State",
		FilterFunc:  Filter.Equal[int],
		Condition:   1,
		SortRequest: nil,
	}, 0, 20)
	if err != nil {
		return
	}

	for _, v := range list {
		fmt.Println(v.IsPublic.Set(1, v.Id))
	}

	err := Del[Video](v1.Id, v2.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}*/

	/*
		err := BulkInsert[Video](map[string]*Video{
			v1.Id: &v1,
			v2.Id: &v2,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}*/

}
