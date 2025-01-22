import React from 'react';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/lib/locale/zh_CN';
import Routes from './routes';
import './App.css';
import 'antd/dist/reset.css';

const App: React.FC = () => {
  console.log('App rendering');
  return (
    <ConfigProvider locale={zhCN}>
      <div className="app">
        <Routes />
      </div>
    </ConfigProvider>
  );
};

export default App; 