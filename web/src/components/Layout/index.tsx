import React from 'react';
import { Layout, Menu } from 'antd';
import { Link, Outlet, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  AlertOutlined,
  DatabaseOutlined,
  CloudServerOutlined,
} from '@ant-design/icons';

const { Header, Sider, Content } = Layout;

const AppLayout: React.FC = () => {
  console.log('Layout rendering'); // 添加调试日志
  const location = useLocation();

  const menuItems = [
    {
      key: '/middleware',
      icon: <DatabaseOutlined />,
      label: <Link to="/middleware">Middleware</Link>,
    },
    {
      key: '/monitoring',
      icon: <DashboardOutlined />,
      label: <Link to="/monitoring">Monitoring</Link>,
    },
    {
      key: '/alerts',
      icon: <AlertOutlined />,
      label: <Link to="/alerts">Alerts</Link>,
    },
    {
      key: '/hosts',
      icon: <CloudServerOutlined />,
      label: <Link to="/hosts">Host Management</Link>,
    },
  ];

  return (
    <Layout className="app">
      <Header>
        <div className="logo">Middleware Platform</div>
      </Header>
      <Layout>
        <Sider width={200} className="site-layout-background">
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            style={{ height: '100%', borderRight: 0 }}
            items={menuItems}
          />
        </Sider>
        <Layout style={{ padding: '24px' }}>
          <Content className="site-layout-background">
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default AppLayout; 