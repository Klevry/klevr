import 'react-perfect-scrollbar/dist/css/styles.css';
import { useRoutes } from 'react-router-dom';
import routes from 'src/routes';

import {
  BrowserRouter as Router,
  Redirect,
  Route,
  Switch,
  useLocation
} from 'react-router-dom';
import Login from './pages/Login';
import Activate from './pages/Activate';
import NotFound from './pages/NotFound';
import Overview from 'src/pages/Overview';
const App = () => {
  const routing = useRoutes(routes);

  return <>{routing}</>;
};

// const App = () => {
//   return (
//     <Router>
//       <Switch>
//         <Redirect from="/" to="/login" exact />
//         <Route exact path="/login" component={Login} />
//         <Route exact path="/activate" component={Activate} />
//         <Route exact path="/404" component={NotFound} />
//         <Route exact path="/app/overview" component={Overview} />
//         <Redirect from="*" to="/404" exact />
//       </Switch>
//     </Router>
//   );
// };

export default App;
