import React, {
  FC,
  Fragment,
  useState,
  useEffect,
  useReducer,
  useCallback,
} from "react";
import { useParams } from "react-router-dom";
import { WS_API_HOST } from "~/utils/constants";

type Tree = {
  hash: string;
  list: string[];
  commit: {
    author: string;
    message: string;
  };
};

const actionTypes = {
  success: "success",
  failed: "failed",
} as const;
type ActionType = typeof actionTypes[keyof typeof actionTypes];

const reducer = (
  state: Tree[],
  action: { type: ActionType; payload: Tree }
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
  const [idx, setIdx] = useState(0);
  useEffect(() => {
    const ws = new WebSocket(`${WS_API_HOST}/ws?repo=${org}/${repo}`);
    ws.onmessage = ({ data }) => {
      try {
        const packet = JSON.parse(data) as Tree;
        dispatch({ type: actionTypes.success, payload: packet });
      } catch (e) {
        dispatch({
          type: actionTypes.failed,
          payload: { hash: "", list: [], commit: { author: "", message: "" } },
        });
      }
    };
    return (): void => {
      try {
        ws.close();
      } catch (e) {}
    };
  }, []);

  useEffect(() => {
    const onKeyDown = (e) => {
      if (e.keyCode === 37) {
        setIdx((i) => (i <= 0 ? i : i - 1));
      } else if (e.keyCode === 39) {
        setIdx((i) => (i >= trees.length - 1 ? i : i + 1));
      }
    };
    document.addEventListener("keydown", onKeyDown);

    return () => {
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [trees.length]);

  // todo: reset css
  // todo: loading
  // todo: author
  // todo: PR
  // todo: tree表示

  return (
    <div style={{ width: "100%", height: "100%", backgroundColor: "#002b36" }}>
      {trees.length > 0 && (
        <Fragment>
          <div style={{ color: "859900", fontSize: 10, fontFamily: "Monaco" }}>
            {trees[idx].hash}
          </div>
          <div style={{ color: "#eee8d5", fontSize: 10, fontFamily: "Monaco" }}>
            <div>tree</div>
            {trees[idx].list.map((d) => (
              <div key={d}>{d}</div>
            ))}
          </div>
        </Fragment>
      )}
    </div>
  );
};

export default TreeFC;
