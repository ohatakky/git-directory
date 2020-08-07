import React, { FC } from "react";
import {
  Router as ReactRouter,
  Switch,
  Route,
  Redirect,
} from "react-router-dom";
import { History } from "history";
import Root from "~/components/pages/Root";
import Tree from "~/components/pages/Tree";

const Router: FC<{ history: History }> = ({ history }) => {
  return (
    <ReactRouter history={history}>
      <Switch>
        <Route path="/root" exact>
          <Root />
        </Route>
        <Route path="/:org/:repo" exact>
          <Tree />
        </Route>
        <Route path="/">
          <Redirect
            to={{
              pathname: "/root",
            }}
          />
        </Route>
      </Switch>
    </ReactRouter>
  );
};

export default Router;
