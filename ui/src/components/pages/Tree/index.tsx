import React, { FC, Fragment, useEffect, useReducer } from "react";
import { useParams } from "react-router-dom";
import { WS_API_HOST } from "~/utils/constants";

type Tree = {
  i: number;
  hash: string;
  t: string[];
};

const actionTypes = {
  success: "success",
  failed: "failed",
} as const;
type ActionType = typeof actionTypes[keyof typeof actionTypes];

const reducer = (
  state: Tree[],
  action: { type: ActionType; payload: Tree },
) => {
  switch (action.type) {
    case actionTypes.success:
      return [...state, action.payload];
    case actionTypes.failed:
      return [];
  }
};

const TreeFC: FC = () => {
  const { org, repo } = useParams();

  const [trees, dispatch] = useReducer(reducer, []);
  useEffect(() => {
    const ws = new WebSocket(`${WS_API_HOST}/ws?repo=${org}/${repo}`);
    ws.onmessage = ({ data }) => {
      try {
        const packet = JSON.parse(data) as Tree;
        dispatch({ type: actionTypes.success, payload: packet });
      } catch (e) {
        dispatch(
          { type: actionTypes.failed, payload: { i: 0, hash: "", t: [] } },
        );
      }
    };
    return (): void => {
      try {
        ws.close();
      } catch (e) {}
    };
  }, []);

  return (
    <Fragment>
      {trees.map((tree) => (
        <div>{tree.hash}</div>
      ))}
    </Fragment>
  );
};

export default TreeFC;
