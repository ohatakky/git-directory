import React, { FC, Fragment, useEffect, useReducer } from "react";
import { useParams } from "react-router-dom";
import { WS_API_HOST } from "~/utils/constants";

const actionTypes = {
  success: "success",
  reset: "reset",
} as const;
type ActionType = typeof actionTypes[keyof typeof actionTypes];

const reducer = (
  state: string[],
  action: { type: ActionType; payload: string },
) => {
  switch (action.type) {
    case actionTypes.success:
      return [...state, action.payload];
    case actionTypes.reset:
      return [];
    default:
      throw new Error("Unexpected action");
  }
};

const Tree: FC = () => {
  const { org, repo } = useParams();

  const [trees, dispatch] = useReducer(reducer, []);
  useEffect(() => {
    const ws = new WebSocket(`${WS_API_HOST}/ws?repo=${org}/${repo}`);
    ws.onmessage = ({ data }) => {
      try {
        const packet = JSON.parse(data);
        dispatch({ type: actionTypes.success, payload: packet });
      } catch (e) {
        return;
      }
    };
    return (): void => {
      try {
        ws.close();
      } catch (e) {}
    };
  }, []);

  console.log(trees);

  return (
    <Fragment>
      <Fragment></Fragment>
    </Fragment>
  );
};

export default Tree;
