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
	"time"
)

type Review struct {

	Id int64 `json:"id,omitempty"`

	Content string `json:"content,omitempty"`

	Blogtitle string `json:"blogtitle,omitempty"`

	Blogowner string `json:"blogowner,omitempty"`

	Reviewer string `json:"reviewer,omitempty"`

	Createtime time.Time `json:"createtime,omitempty"`
}
