package dashboard

import (
	"time"

	"gorm.io/gorm"
)

// Service 提供仪表盘功能
type Service struct {
	db *gorm.DB
}

// NewService 创建仪表盘服务
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// GetStatistics 获取统计数据
func (s *Service) GetStatistics() (StatisticsResponse, error) {
	var agentCount, convCount, userCount int64
	var tokenUsage int

	// 查询智能体数量
	if err := s.db.Table("agents").Count(&agentCount).Error; err != nil {
		return StatisticsResponse{}, err
	}

	// 查询对话数量
	if err := s.db.Table("conversations").Count(&convCount).Error; err != nil {
		return StatisticsResponse{}, err
	}

	// 查询用户数量
	if err := s.db.Table("users").Count(&userCount).Error; err != nil {
		return StatisticsResponse{}, err
	}

	// 查询Token使用量
	if err := s.db.Table("usage_logs").Select("COALESCE(SUM(tokens_used), 0)").Scan(&tokenUsage).Error; err != nil {
		return StatisticsResponse{}, err
	}

	return StatisticsResponse{
		AgentCount:        int(agentCount),
		ConversationCount: int(convCount),
		UserCount:         int(userCount),
		TokenUsage:        tokenUsage,
	}, nil
}

// GetUsageTrend 获取使用趋势数据
func (s *Service) GetUsageTrend(days int) ([]UsageData, error) {
	result := make([]UsageData, 0, days)
	
	// 计算开始日期
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	
	// 按日期分组查询对话数量和Token使用量
	type DailyUsage struct {
		Date          string
		Conversations int
		Tokens        int
	}
	
	var usageData []DailyUsage
	
	// 这里查询数据库中的实际使用数据
	// 此处模拟数据以演示API结构
	// SQL语句应类似: 
	// SELECT DATE(created_at) as date, COUNT(*) as conversations, SUM(tokens_used) as tokens
	// FROM conversations 
	// WHERE created_at BETWEEN ? AND ?
	// GROUP BY DATE(created_at)
	
	// 生成所有日期
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		
		// 查找是否有对应日期的数据
		found := false
		for _, usage := range usageData {
			if usage.Date == dateStr {
				result = append(result, UsageData{
					Date:          dateStr,
					Conversations: usage.Conversations,
					Tokens:        usage.Tokens,
				})
				found = true
				break
			}
		}
		
		// 如果没有数据，添加0值
		if !found {
			result = append(result, UsageData{
				Date:          dateStr,
				Conversations: 0,
				Tokens:        0,
			})
		}
	}
	
	return result, nil
}

// GetTopAgents 获取热门智能体
func (s *Service) GetTopAgents(limit int) ([]TopAgent, error) {
	// 此处应查询数据库获取真实数据
	// 例如:
	// var agents []struct {
	//     ID          uuid.UUID
	//     Name        string
	//     UsageCount  int
	//     SuccessRate float64
	// }
	//
	// SELECT a.id, a.name, COUNT(c.id) as usage_count, AVG(CASE WHEN c.status = 'success' THEN 1 ELSE 0 END) as success_rate
	// FROM agents a
	// LEFT JOIN conversations c ON c.agent_id = a.id
	// GROUP BY a.id, a.name
	// ORDER BY usage_count DESC
	// LIMIT ?
	
	// 生成示例数据(实际应从数据库查询)
	result := []TopAgent{
		{ID: "1", Name: "客服助手", Usage: 87, SuccessRate: 0.95},
		{ID: "2", Name: "数据分析师", Usage: 65, SuccessRate: 0.92},
		{ID: "3", Name: "营销助手", Usage: 53, SuccessRate: 0.89},
		{ID: "4", Name: "产品顾问", Usage: 42, SuccessRate: 0.91},
		{ID: "5", Name: "技术支持", Usage: 38, SuccessRate: 0.93},
	}
	
	// 限制返回数量
	if len(result) > limit {
		result = result[:limit]
	}
	
	return result, nil
}

// GetRecentActivities 获取最近活动
func (s *Service) GetRecentActivities(limit int) ([]RecentActivity, error) {
	// 此处应查询数据库获取真实数据
	// 例如:
	// var activities []struct {
	//     ID        uuid.UUID
	//     Type      string
	//     Content   string
	//     CreatedAt time.Time
	//     UserID    uuid.UUID
	// }
	//
	// SELECT id, type, content, created_at, user_id
	// FROM activities
	// ORDER BY created_at DESC
	// LIMIT ?
	
	// 生成示例数据(实际应从数据库查询)
	result := []RecentActivity{
		{ID: "1", Type: "agent_created", Content: "创建了新智能体 \"客服助手\"", Time: "10分钟前", UserID: "user1"},
		{ID: "2", Type: "conversation", Content: "与 \"数据分析师\" 进行了对话", Time: "30分钟前", UserID: "user2"},
		{ID: "3", Type: "knowledge_base", Content: "更新了 \"产品手册\" 知识库", Time: "1小时前", UserID: "user1"},
		{ID: "4", Type: "agent_updated", Content: "更新了 \"营销助手\" 配置", Time: "2小时前", UserID: "user3"},
		{ID: "5", Type: "conversation", Content: "与 \"技术支持\" 进行了对话", Time: "3小时前", UserID: "user4"},
	}
	
	// 限制返回数量
	if len(result) > limit {
		result = result[:limit]
	}
	
	return result, nil
} 