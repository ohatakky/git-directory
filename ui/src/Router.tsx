import React, { FC } from "react";
import {
  Router as ReactRouter,
  Switch,
  Route,
  Redirect,
} from "react-router-dom";
import { History } from "history";
import Root from "~/components/pages/Root";
import Fzf from "~/components/pages/Fzf";

const Router: FC<{ history: History }> = ({ history }) => {
  return (
    <ReactRouter history={history}>
      <Switch>
        <Route path="/" exact>
          <Root />
        </Route>
        <Route path="/:org/:repo" exact>
          <Fzf />
        </Route>
        <Route path="/">
          <Redirect
            to={{
              pathname: "/",
            }}
          />
        </Route>
      </Switch>
    </ReactRouter>
  );
};

export default Router;
