import React, { FC } from "react";
import { createBrowserHistory } from "history";
import Router from "~/Router";

const App: FC = () => {
  const history = createBrowserHistory();
  return <Router history={history} />;
};

export default App;
