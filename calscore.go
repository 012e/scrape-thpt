package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/012e/scrape-thptscore/stuparser"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func applyFuncBetween(start int, end int, db *gorm.DB, f func(row *gorm.DB), wg *sync.WaitGroup) {
	defer wg.Done()
	var student stuparser.Student
	for id := start; id < end+1; id++ {
		var totalCount int64 = 3
		log.Infof("finding student with id %d", id)
		db.Raw("select * from students where id = ?", id).Scan(&student).Count(&totalCount)
		if totalCount == 0 {
			panic(fmt.Sprintf("wtf, can't find student id %d", id))
		}
		log.Infof("found %d student with id %d", totalCount, start)
		// db.Where("id=", id).Select("COUNT(*)").Find(&totalCount)
		// f(student)
	}

	// for i := start; i < end+1; i++ {
	// 	log.Infof("Currently applying student %d", i)
	// 	var student stuparser.Student
	// 	var err error
	// 	for tries := int(0); tries < TryNum+1; tries++ {
	// 		student, err = getStudent(i)
	// 		if err != nil {
	// 			log.Warnf("Student with id %d failed %d times", i, tries)
	// 			time.Sleep(1 * time.Second)
	// 		} else {
	// 			db.Create(&student)
	// 			break
	// 		}
	// 	}
	// 	log.Infof("Finished student id %d", i)
	// }
}

func applyScore(row *gorm.DB) {
	row.Update("hoa", nil)
	var name string
	row.Pluck("name", &name)
	log.Info("updated ", name)
}

func calScore(db *gorm.DB) {
	var wg sync.WaitGroup
	wg.Add(ConcurrentConnection)
	var student stuparser.Student
	db.Model(&student).WithContext(context.Background()).Session(&gorm.Session{})
	for _, split := range efficientTaskSplit(StartIndex, EndIndex, ConcurrentConnection) {
		log.Infof("Starting goroutine for task between %d and %d", split[0], split[1])
		go applyFuncBetween(split[0], split[1], db, applyScore, &wg)
	}
	wg.Wait()
	defer f.Close()
}
