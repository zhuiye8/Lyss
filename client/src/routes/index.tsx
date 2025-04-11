import { createBrowserRouter } from 'react-router-dom';
import MainLayout from '../layouts/MainLayout';
import Login from '../pages/Login';
import Dashboard from '../pages/Dashboard';
import AgentBuilder from '../pages/AgentBuilder';
import NotFound from '../pages/NotFound';
import KnowledgeBase from '../pages/KnowledgeBase';
import Settings from '../pages/Settings';

const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <Dashboard />,
      },
      {
        path: 'agents',
        element: <AgentBuilder />,
      },
      {
        path: 'knowledge',
        element: <KnowledgeBase />,
      },
      {
        path: 'settings',
        element: <Settings />,
      },
      {
        path: '*',
        element: <NotFound />,
      },
    ],
  },
]);

export default router; 