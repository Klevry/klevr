import 'react-perfect-scrollbar/dist/css/styles.css';
import { useRoutes, Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';

import DashboardLayout from 'src/components/DashboardLayout';
import MainLayout from 'src/components/MainLayout';
import TaskList from 'src/pages/TaskList';
import Overview from 'src/pages/Overview';
import ZoneList from 'src/pages/ZoneList';
import Login from 'src/pages/Login';
import Activate from 'src/pages/Activate';
import NotFound from 'src/pages/NotFound';
import Settings from 'src/pages/Settings';
import Credentials from 'src/pages/Credentials';

const App = () => {
  const isLoggedIn = useSelector((store) => store.loginReducer);
  const loggedInView = useRoutes([
    {
      path: 'app',
      element: <DashboardLayout />,
      children: [
        { path: 'tasks', element: <TaskList /> },
        { path: 'overview', element: <Overview /> },
        { path: 'zones', element: <ZoneList /> },
        { path: 'credentials', element: <Credentials /> },
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
        { path: '/', element: <Navigate to="/app/overview" /> },
        { path: '*', element: <Navigate to="/404" /> }
      ]
    }
  ]);
  const loggedOutView = useRoutes([
    {
      path: '/',
      element: <MainLayout />,
      children: [
        { path: 'login', element: <Login /> },
        { path: 'activate', element: <Activate /> },
        { path: '404', element: <NotFound /> },
        { path: '/', element: <Navigate to="/login" /> },
        { path: '/app/tasks', element: <Navigate to="/login" /> },
        { path: '/app/overview', element: <Navigate to="/login" /> },
        { path: '/app/zones', element: <Navigate to="/login" /> },
        { path: '/app/credentials', element: <Navigate to="/login" /> },
        { path: '/app/settings', element: <Navigate to="/login" /> },
        { path: '*', element: <Navigate to="/404" /> }
      ]
    }
  ]);

  return <>{isLoggedIn ? loggedInView : loggedOutView}</>;
};

export default App;
