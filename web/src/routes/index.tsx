import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from '../components/Layout';
import MiddlewarePage from '../pages/Middleware';
import MonitoringPage from '../pages/Monitoring';
import AlertsPage from '../pages/Alerts';
import HostPage from '../pages/Host';

const AppRoutes: React.FC = () => {
  console.log('Routes rendering'); // 添加调试日志
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<Navigate to="/middleware" replace />} />
          <Route path="/middleware" element={<MiddlewarePage />} />
          <Route path="/monitoring" element={<MonitoringPage />} />
          <Route path="/alerts" element={<AlertsPage />} />
          <Route path="/hosts" element={<HostPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default AppRoutes; 