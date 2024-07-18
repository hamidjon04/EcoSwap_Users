package redis

import (
	"context"
	pb "ecoswap/genproto/users"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRedis struct {
	Rdb *redis.Client
}

func NewRedisRepo(rdb *redis.Client) *UserRedis {
	return &UserRedis{
		Rdb: rdb,
	}
}

func (U *UserRedis) AddToBlacklist(ctx context.Context, token string, duration time.Duration) error {
	err := U.Rdb.Set(ctx, token, "blacklisted", duration).Err()
	return err
}

func (U *UserRedis) UpdateEcoPointsInRedis(userID string, ecoPoints int32, reason string) error {
	err := U.Rdb.Set(context.TODO(), fmt.Sprintf("user_id:%s:Points", userID), ecoPoints, 0).Err()
	if err != nil {
		log.Println("Failed to set eco points in Redis:", err)
		return err
	}

	err = U.Rdb.Set(context.TODO(), fmt.Sprintf("user_id:%s:reason", userID), reason, 0).Err()
	if err != nil {
		log.Println("Failed to set reason in Redis:", err)
		return err
	}

	historyKey := fmt.Sprintf("user:%s:eco_points_history", userID)
	historyEntry := pb.HistoryUserPoints{
		Points: ecoPoints,
		Reason: reason,
	}
	historyData, err := json.Marshal(historyEntry)
	if err != nil {
		log.Println("Failed to marshal history entry:", err)
		return err
	}

	err = U.Rdb.RPush(context.TODO(), historyKey, historyData).Err()
	if err != nil {
		log.Println("Failed to push history entry to Redis:", err)
		return err
	}

	return nil
}

func (U *UserRedis) HistoryEcoPointsByUser(req *pb.HistoryReq) (*pb.Histories, error) {
	key := fmt.Sprintf("user:%s:eco_points_history", req.UserId)

	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	totalEntries, err := U.Rdb.LLen(context.TODO(), key).Result()
	if err != nil {
		log.Println("Failed to get total entries from Redis:", err)
		return nil, err
	}

	start := offset
	end := offset + limit - 1

	entries, err := U.Rdb.LRange(context.TODO(), key, int64(start), int64(end)).Result()
	if err != nil {
		log.Println("Failed to get eco points history from Redis:", err)
		return nil, err
	}

	history := []*pb.HistoryUserPoints{}
	for _, entry := range entries {
		histEntry := pb.HistoryUserPoints{}
		err := json.Unmarshal([]byte(entry), &histEntry)
		if err != nil {
			log.Println("Failed to unmarshal eco points history entry:", err)
			continue
		}
		history = append(history, &histEntry)
	}

	response := &pb.Histories{
		Histories: history,
		Total:     int32(totalEntries),
		Page:      req.Offset / req.Limit,
		Limit:     req.Limit,
	}

	return response, nil
}