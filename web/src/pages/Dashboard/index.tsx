import { PageContainer, StatisticCard } from '@ant-design/pro-components';
import { Col, Row } from 'antd';
import { ApiOutlined, DatabaseOutlined, ClockCircleOutlined, CheckCircleOutlined } from '@ant-design/icons';

const { Statistic } = StatisticCard;

const DashboardPage: React.FC = () => {
  return (
    <PageContainer>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: 'API Endpoints',
              value: 12,
              icon: <ApiOutlined style={{ fontSize: 24, color: '#1890ff' }} />,
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: 'Database Records',
              value: 128,
              icon: <DatabaseOutlined style={{ fontSize: 24, color: '#52c41a' }} />,
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: 'Uptime',
              value: '99.9%',
              icon: <ClockCircleOutlined style={{ fontSize: 24, color: '#faad14' }} />,
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: 'Health Status',
              value: 'OK',
              icon: <CheckCircleOutlined style={{ fontSize: 24, color: '#52c41a' }} />,
            }}
          />
        </Col>
      </Row>

      <StatisticCard
        title="Welcome"
        style={{ marginTop: 16 }}
        statistic={{
          value: 'Project scaffold is ready. Edit this page to customize your dashboard.',
        }}
      />
    </PageContainer>
  );
};

export default DashboardPage;
