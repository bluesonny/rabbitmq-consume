package handlers

import (
	"log"
	. "mq-consume/models"
	"strconv"
	"strings"
)

func do(msg Message) (sta bool) {
	log.Printf("开始处理业务...........")
	co := Comment{}
	comment, err := co.Get(msg.Id)
	if err != nil {
		log.Printf("错误 %v", err)
	}
	log.Printf("当前数据处理的 %#v", comment)
	if comment.IsEmpty() {
		log.Printf("当前数据处理的---空对象")
		return true
	}

	rc := RedisClient.Get()
	defer rc.Close()
	//redis 计数处理
	if msg.Type == "comment" {
		key := "red_" + strconv.Itoa(comment.ToUserId)
		rc.Do("HINCRBY", key, "comment", 1)
		log.Printf("写入缓存---评论别人 HINCRBY %s comment", key)
	}

	if msg.Type == "digg" {
		key := "red_" + strconv.Itoa(comment.UserId)
		rc.Do("HINCRBY", key, "digg", 1)
		log.Printf("写入缓存---点赞别人HINCRBY %s digg", key)
	}

	if comment.ParentId == 0 {
		log.Printf("当前数据处理的---一级评论")
		return true
	}

	if msg.Type == "comment" || msg.Type == "digg" || msg.Type == "del" {
		if comment.ParentId > 0 {
			cos, err := co.Cos(comment.ParentId)
			if err != nil {
				log.Printf("获取点赞最多评论sql错误--- %v", err)
				return false
			}
			log.Printf("当前数据点赞最多评论  %v", cos)
			var strs []string
			for _, v := range cos {
				id := strconv.Itoa(v.Id)
				strs = append(strs, id)
			}
			str := strings.Join(strs, ",")
			log.Printf("当前数据点赞最多评论 生成字符串  %v", str)
			count, err := co.Count(comment.ParentId)

			if err != nil {
				log.Printf("总数错误 %v", err)
				return false
			}
			log.Printf("当前数据点赞最多评论 总数  %v", count)
			err = comment.Update(str, count)
			if err != nil {
				log.Printf("更新错误 %v", err)
				return false
			}
			return true
		}
	}
	return
}
