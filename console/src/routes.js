import { Navigate } from 'react-router-dom';
import DashboardLayout from 'src/components/DashboardLayout';
import MainLayout from 'src/components/MainLayout';
import TaskList from 'src/pages/TaskList';
import Overview from 'src/pages/Overview';
import ZoneList from 'src/pages/ZoneList';
import Login from 'src/pages/Login';
import Activate from 'src/pages/Activate';
import NotFound from 'src/pages/NotFound';
import Settings from 'src/pages/Settings';

const routes = [
  {
    path: 'app',
    element: <DashboardLayout />,
    children: [
      { path: 'tasks', element: <TaskList /> },
      { path: 'overview', element: <Overview /> },
      { path: 'zones', element: <ZoneList /> },
      { path: 'settings', element: <Settings /> },
      { path: '*', element: <Navigate to="/404" /> }
    ]
  },
  {
    path: '/',
    element: <MainLayout />,
    children: [
      { path: 'login', element: <Login /> },
      { path: 'activate', element: <Activate /> },
      { path: '404', element: <NotFound /> },
      { path: '/', element: <Navigate to="/login" /> },
      { path: '*', element: <Navigate to="/404" /> }
    ]
  }
];

export default routes;
