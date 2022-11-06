/*
   Copyright 2022 CESS scheduler authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package node

import (
	"net/http"
	"unsafe"

	"github.com/CESSProject/cess-oss/configs"
	"github.com/CESSProject/cess-oss/pkg/chain"
	"github.com/CESSProject/cess-oss/pkg/utils"
	"github.com/gin-gonic/gin"
)

// It is used to authorize users
func (n *Node) GetHandle(c *gin.Context) {
	// account
	account := c.Request.Header.Get(configs.Header_Account)
	if account == "" {
		//Uld.Sugar().Infof("[%v] head missing token", c.ClientIP())
		c.JSON(400, "Invalid.Account")
		return
	}

	pkey, err := utils.DecodePublicKeyOfCessAccount(account)
	if err != nil {
		c.JSON(400, "InvalidParameter.Account")
		return
	}

	getName := c.Param("name")

	if VerifyBucketName(getName) {
		bucketInfo, err := n.Chain.GetBucketInfo(pkey, getName)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		filesHash := make([]string, len(bucketInfo.Objects_list))
		for i := 0; i < len(bucketInfo.Objects_list); i++ {
			filesHash[i] = string(bucketInfo.Objects_list[i][:])
		}
		data := struct {
			Num   uint32
			Files []string
		}{
			Num:   uint32(bucketInfo.Objects_num),
			Files: filesHash,
		}
		c.JSON(http.StatusOK, data)
		return
	}
	//
	if len(getName) == int(unsafe.Sizeof(chain.FileHash{})) {
		fmeta, err := n.Chain.GetFileMetaInfo(getName)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		c.JSON(http.StatusOK, fmeta)
	}
	//
	bucketList, err := n.Chain.GetBucketList(pkey)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	bucket := make([]string, len(bucketList))
	for i := 0; i < len(bucketList); i++ {
		bucket[i] = string(bucketList[i][:])
	}
	c.JSON(http.StatusOK, bucket)
}