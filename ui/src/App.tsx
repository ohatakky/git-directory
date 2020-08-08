import React, { FC, Fragment } from "react";
import { createBrowserHistory } from "history";
import CssBaseline from "@material-ui/core/CssBaseline";
import Router from "~/Router";

const App: FC = () => {
  const history = createBrowserHistory();
  return (
    <Fragment>
      <CssBaseline />
      <Router history={history} />
    </Fragment>
  );
};

export default App;
