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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetBlogByLabel(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	label := data["labelname"]
	fmt.Println("labelname:", label)
	// TODO query blogs by label
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
