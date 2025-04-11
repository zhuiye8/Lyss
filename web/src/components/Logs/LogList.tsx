import React, { useState, useEffect } from 'react';
import { Table, Tag, Space, Button, Modal, Card, Row, Col, Form, Select, DatePicker, Input, Badge } from 'antd';
import { EyeOutlined, CheckCircleOutlined } from '@ant-design/icons';
import moment from 'moment';
import { getLogs, markErrorAsResolved } from '../../services/logService';
import { LogLevel, LogCategory, LogType, Log, APILog, ErrorLog, ModelCallLog, LogListResponse } from '../../types/log';
import LogDetail from './LogDetail';

const { RangePicker } = DatePicker;
const { Option } = Select;

interface LogListProps {
  type?: LogType;
  initialFilters?: any;
}

const LogList: React.FC<LogListProps> = ({ type = LogType.ALL, initialFilters = {} }) => {
  const [logs, setLogs] = useState<(Log | APILog | ErrorLog | ModelCallLog)[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [filters, setFilters] = useState<any>(initialFilters);
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);
  const [modalVisible, setModalVisible] = useState<boolean>(false);
  const [form] = Form.useForm();

  // 获取日志列表
  const fetchLogs = async () => {
    setLoading(true);
    try {
      const queryParams = {
        page,
        page_size: pageSize,
        type,
        ...filters,
      };

      // 转换日期范围为ISO字符串
      if (filters.dateRange && filters.dateRange.length === 2) {
        queryParams.start_time = filters.dateRange[0].toISOString();
        queryParams.end_time = filters.dateRange[1].toISOString();
        delete queryParams.dateRange;
      }

      const response = await getLogs(queryParams);
      const data = response as LogListResponse;
      setLogs(data.data);
      setTotal(data.meta.total);
    } catch (error) {
      console.error('获取日志失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 首次加载和筛选条件变化时获取日志
  useEffect(() => {
    fetchLogs();
  }, [page, pageSize, type, JSON.stringify(filters)]);

  // 处理筛选条件变化
  const handleFilterChange = (values: any) => {
    const newFilters = { ...values };
    setFilters(newFilters);
    setPage(1); // 重置页码
  };

  // 查看日志详情
  const handleViewLog = (record: Log) => {
    setSelectedLog(record);
    setModalVisible(true);
  };

  // 标记错误为已解决
  const handleResolveError = async (id: string) => {
    try {
      await markErrorAsResolved(id);
      fetchLogs(); // 刷新日志列表
    } catch (error) {
      console.error('标记错误失败:', error);
    }
  };

  // 日志级别对应的标签颜色
  const getLevelColor = (level: LogLevel) => {
    switch (level) {
      case LogLevel.DEBUG:
        return 'blue';
      case LogLevel.INFO:
        return 'green';
      case LogLevel.WARN:
        return 'orange';
      case LogLevel.ERROR:
        return 'red';
      case LogLevel.FATAL:
        return 'purple';
      default:
        return 'default';
    }
  };

  // 生成表格列配置
  const getColumns = () => {
    const baseColumns = [
      {
        title: '级别',
        dataIndex: 'level',
        key: 'level',
        render: (level: LogLevel) => (
          <Tag color={getLevelColor(level)}>
            {level.toUpperCase()}
          </Tag>
        ),
        width: 90,
      },
      {
        title: '类别',
        dataIndex: 'category',
        key: 'category',
        render: (category: LogCategory) => <Tag>{category}</Tag>,
        width: 110,
      },
      {
        title: '消息',
        dataIndex: 'message',
        key: 'message',
        ellipsis: true,
      },
      {
        title: '时间',
        dataIndex: 'created_at',
        key: 'created_at',
        render: (time: string) => moment(time).format('YYYY-MM-DD HH:mm:ss'),
        width: 180,
      },
      {
        title: '操作',
        key: 'action',
        render: (_, record: Log) => (
          <Space size="small">
            <Button 
              type="link" 
              icon={<EyeOutlined />} 
              onClick={() => handleViewLog(record)}
              size="small"
            >
              详情
            </Button>
            {'error_code' in record && (
              <Button
                type="link"
                icon={<CheckCircleOutlined />}
                onClick={() => handleResolveError(record.id)}
                disabled={!!(record as ErrorLog).resolved_at}
                size="small"
              >
                {(record as ErrorLog).resolved_at ? '已解决' : '标记解决'}
              </Button>
            )}
          </Space>
        ),
        width: 160,
      },
    ];

    // 根据日志类型添加特定列
    if (type === LogType.API || type === LogType.ALL) {
      baseColumns.splice(3, 0, {
        title: '状态码',
        dataIndex: 'status_code',
        key: 'status_code',
        render: (code: number) => {
          let color = 'green';
          if (code >= 400) color = 'red';
          else if (code >= 300) color = 'orange';
          return <Badge status={code < 400 ? 'success' : 'error'} text={code} />;
        },
        width: 90,
      });
    }

    if (type === LogType.MODEL_CALL || type === LogType.ALL) {
      baseColumns.splice(3, 0, {
        title: '模型',
        dataIndex: 'model_name',
        key: 'model_name',
        width: 120,
      });
    }

    return baseColumns;
  };

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Form
          form={form}
          layout="horizontal"
          onFinish={handleFilterChange}
          initialValues={initialFilters}
        >
          <Row gutter={16}>
            <Col span={6}>
              <Form.Item name="level" label="日志级别">
                <Select allowClear placeholder="选择日志级别">
                  {Object.values(LogLevel).map(level => (
                    <Option key={level} value={level}>{level.toUpperCase()}</Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="category" label="日志类别">
                <Select allowClear placeholder="选择日志类别">
                  {Object.values(LogCategory).map(category => (
                    <Option key={category} value={category}>{category}</Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="dateRange" label="时间范围">
                <RangePicker showTime />
              </Form.Item>
            </Col>
            {type === LogType.API && (
              <Col span={6}>
                <Form.Item name="path" label="API路径">
                  <Input placeholder="请输入API路径" />
                </Form.Item>
              </Col>
            )}
            {type === LogType.ERROR && (
              <Col span={6}>
                <Form.Item name="error_code" label="错误代码">
                  <Input placeholder="请输入错误代码" />
                </Form.Item>
              </Col>
            )}
            {type === LogType.MODEL_CALL && (
              <Col span={6}>
                <Form.Item name="model_name" label="模型名称">
                  <Input placeholder="请输入模型名称" />
                </Form.Item>
              </Col>
            )}
            <Col span={24} style={{ textAlign: 'right' }}>
              <Button type="primary" htmlType="submit">
                筛选
              </Button>
              <Button 
                style={{ marginLeft: 8 }} 
                onClick={() => {
                  form.resetFields();
                  setFilters({});
                }}
              >
                重置
              </Button>
            </Col>
          </Row>
        </Form>
      </Card>

      <Table
        columns={getColumns()}
        dataSource={logs}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条记录`,
        }}
      />

      <Modal
        title="日志详情"
        visible={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={800}
      >
        {selectedLog && <LogDetail log={selectedLog} />}
      </Modal>
    </div>
  );
};

export default LogList; 