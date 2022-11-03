package library

import (
	"context"
	"fmt"
	"log"
	"year-end/model"
)

const noneCOOKIE = ""
const libraryUrl = "http://opac.ccnu.edu.cn/reader/book_hist.php"
const hostUrl = "http://opac.ccnu.edu.cn"

func GetHistoryBooks(ctx context.Context, uname, psd string) {
	cookie, err := getLibraryCookie(uname, psd)
	if err != nil || cookie == noneCOOKIE {
		log.Fatal("cookie获取失败:", err)
	}
	fmt.Println("cookie获取成功")
	data := make(chan string)
	tx := model.D.DB.Begin()
	
	go getHistoryBooks(cookie, data)
	for {
		select {
		case one, ok := <-data:
			if !ok {
				tx.Commit()
				return
			}
			book := &model.LibraryModel{Name: one, Id: uname}
			if err := book.AddBook(tx); err != nil {
				tx.Rollback()
				log.Fatal("事务提交失败:", err)
				return
			}
		case <-ctx.Done():
			tx.Rollback()
			return
		}
	}
}
