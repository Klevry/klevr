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
import { persistStore } from 'redux-persist';
import { PersistGate } from 'redux-persist/integration/react';
const store = createStore(rootReducer);
const persistor = persistStore(store);

ReactDOM.render(
  <BrowserRouter>
    <Provider store={store}>
      <PersistGate loading={null} persistor={persistor}>
        <ThemeProvider theme={theme}>
          <GlobalStyles />
          <App />
        </ThemeProvider>
      </PersistGate>
    </Provider>
  </BrowserRouter>,
  document.getElementById('root')
);

serviceWorker.unregister();
