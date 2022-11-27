package library

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"year-end/model"
)

const noneCOOKIE = ""
const libraryUrl = "http://opac.ccnu.edu.cn/reader/book_hist.php"
const hostUrl = "http://opac.ccnu.edu.cn"

func GetHistoryBooks(ctx context.Context, uname, psd string, client *http.Client) error {
	cookie, err := getLibraryCookie(client)
	if err != nil || cookie == noneCOOKIE {
		return err
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
				return nil
			}
			book := &model.Book{Name: one, Account: uname}
			if err := book.AddBook(tx); err != nil {
				tx.Rollback()
				log.Fatal("事务提交失败:", err)
				return err
			}
		case <-ctx.Done():
			tx.Rollback()
			return err
		}
	}
}
