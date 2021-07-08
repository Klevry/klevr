import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import * as serviceWorker from './serviceWorker';
import App from './App';
import { ThemeProvider } from '@material-ui/core';
import theme from 'src/theme';
import GlobalStyles from 'src/components/GlobalStyles';
import { createStore } from 'redux';
import { Provider } from 'react-redux';
import rootReducer from './components/store/reducers';

const store = createStore(rootReducer);

ReactDOM.render(
  <BrowserRouter>
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <GlobalStyles />
        <App />
      </ThemeProvider>
    </Provider>
  </BrowserRouter>,
  document.getElementById('root')
);

serviceWorker.unregister();
