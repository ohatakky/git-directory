import React, { FC, Fragment, useState, useEffect, useReducer } from "react";
import { useParams } from "react-router-dom";
import { WS_API_HOST } from "~/utils/constants";
import { LinearProgress } from "@material-ui/core";

type Fzf = {
  hash: string;
  files: string[];
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
  state: Fzf[],
  action: { type: ActionType; payload: Fzf },
) => {
  switch (action.type) {
    case actionTypes.success:
      return [...state, action.payload];
    case actionTypes.failed:
      return [];
  }
};

const FzfFC: FC = () => {
  const { org, repo } = useParams();
  const [fzfs, dispatch] = useReducer(reducer, []);
  const [idx, setIdx] = useState(0);
  useEffect(() => {
    const ws = new WebSocket(`${WS_API_HOST}/ws?repo=${org}/${repo}`);
    ws.onmessage = ({ data }) => {
      try {
        const packet = JSON.parse(data) as Fzf;
        dispatch({ type: actionTypes.success, payload: packet });
      } catch (e) {
        dispatch({
          type: actionTypes.failed,
          payload: { hash: "", files: [], commit: { author: "", message: "" } },
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
        setIdx((i) => (i >= fzfs.length - 1 ? i : i + 1));
      }
    };
    document.addEventListener("keydown", onKeyDown);

    return () => {
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [fzfs.length]);

  // todo: terminal view
  // @ ~/git-directory/gorilla/websocket master ()
  // $ tree
  // author
  // PR

  // todo: loading
  // todo: fuzzy finder表示
  // todo: tree表示

  return (
    <div style={{ width: "100%", height: "100%", backgroundColor: "#002b36" }}>
      {fzfs.length > 0
        ? (
          <div>
            <div style={{ display: "flex" }}>
              <div
                style={{ color: "#657b83", fontSize: 10, fontFamily: "Monaco" }}
              >
                {`@ ~/git-directory/${org}/${repo}`}
              </div>
              <div
                style={{ color: "#859900", fontSize: 10, fontFamily: "Monaco" }}
              >
                {fzfs[idx].hash}
              </div>
              <div
                style={{ color: "#268bd2", fontSize: 10, fontFamily: "Monaco" }}
              >
                {fzfs[idx].commit.author}
              </div>
            </div>

            <div
              style={{ color: "#eee8d5", fontSize: 10, fontFamily: "Monaco" }}
            >
              {fzfs[idx].commit.message}
            </div>

            <div
              style={{ color: "#eee8d5", fontSize: 10, fontFamily: "Monaco" }}
            >
              <div>$ fzf</div>
              {fzfs[idx].files.map((f, i) => (
                <div>{f}</div>
              ))}
            </div>
          </div>
        )
        : (
          <LinearProgress />
        )}
    </div>
  );
};

export default FzfFC;
