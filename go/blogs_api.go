/*
 * Blog for service computing
 *
 *
 *
 * API version: 1.0.0
 * Contact: 895118352@qq.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetALLBlog(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	pages, err2 := strconv.Atoi(r.URL.Query().Get("pages"))
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(offset, pages)
	// FIXME: get all blogs
	if blogs, err := DBgetAllBlog(offset, pages); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		buf, _ := json.Marshal(blogs)
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetBlogByTitle(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	username := data["username"]
	title := data["title"]
	fmt.Println(username, title)
	// FIXME: get blog by title
	if blog, err := DBgetBolgByBlogTitleAndAuthor(title, username); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		buf, _ := json.Marshal(blog)
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetReviewByBlog(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	username := data["username"]
	title := data["title"]
	fmt.Println(username, title)
	// FIXME: query reviews by blog title (and username ?)
	if reviews, err := DBgetReviewByBlogTitleAndAuthor(title, username); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		buf, _ := json.Marshal(reviews)
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func AddReview(w http.ResponseWriter, r *http.Request) {
	// FIXME:
	username := r.Header.Get("username")
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	data := mux.Vars(r)
	re, _ := ioutil.ReadAll(r.Body)
	content := string(re)
	var review Review
	review.Content = content
	review.Blogtitle = data["title"]
	review.Blogowner = data["username"]
	review.Createtime = time.Now()
	review.Reviewer = username

	if review, err := DBCreateReview(&review); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		buf, _ := json.Marshal(review)
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
