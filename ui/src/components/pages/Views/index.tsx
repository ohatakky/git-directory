import React, {
  FC,
  Fragment,
  useState,
  useEffect,
  useReducer,
  ChangeEvent,
} from "react";
import { useParams } from "react-router-dom";
import { WS_API_HOST } from "~/utils/constants";
import { LinearProgress } from "@material-ui/core";
import Fuse from "fuse.js";

type Object = {
  name: string;
  is_file: string;
};

type Commit = {
  hash: string;
  message: string;
  author: string;
  objects: Object[];
};

const actionTypes = {
  success: "success",
  failed: "failed",
} as const;
type ActionType = typeof actionTypes[keyof typeof actionTypes];

const reducer = (
  state: Commit[],
  action: { type: ActionType; payload: Commit },
) => {
  switch (action.type) {
    case actionTypes.success:
      return [...state, action.payload];
    case actionTypes.failed:
      return [];
  }
};

const Views: FC = () => {
  const { org, repo } = useParams();
  const [commits, dispatch] = useReducer(reducer, []);
  const [idx, setIdx] = useState(0);
  const [filter, setFilter] = useState<string>("");

  const [fzfs, setFzfs] = useState<Fuse.FuseResult<Object>[]>([]);
  const prURL = `https://github.com/${org}/${repo}/pulls?q=is%3Apr+hash%3A`;

  // webscoket read message
  useEffect(() => {
    const ws = new WebSocket(`${WS_API_HOST}/ws?org=${org}&repo=${repo}`);
    ws.onmessage = ({ data }) => {
      try {
        const packet = JSON.parse(data) as Commit;
        dispatch({ type: actionTypes.success, payload: packet });
      } catch (e) {
        dispatch({
          type: actionTypes.failed,
          payload: { hash: "", objects: [], author: "", message: "" },
        });
      }
    };
    return (): void => {
      try {
        ws.close();
      } catch (e) {}
    };
  }, []);

  // keyboard
  useEffect(() => {
    const onKeyDown = (e) => {
      if (e.keyCode === 37) {
        setIdx((i) => (i <= 0 ? i : i - 1));
      } else if (e.keyCode === 39) {
        setIdx((i) => (i >= commits.length - 1 ? i : i + 1));
      }
    };
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [commits.length]);

  // fzf
  useEffect(() => {
    if (!commits[idx]) return;
    const fuse = new Fuse(commits[idx].objects, {
      isCaseSensitive: true,
      minMatchCharLength: 2,
      keys: ["name"],
    })
      .search(filter)
      .map((f) => f);
    setFzfs(fuse);
  }, [filter]);

  return (
    <div
      style={{
        width: "100%",
        minHeight: "100%",
        backgroundColor: "#002b36",
        padding: "4px",
        ...(commits[idx] && commits[idx].objects.length > 0 &&
          { height: `${commits[idx].objects.length * 16}px` }),
      }}
    >
      {commits.length > 0
        ? (
          <Fragment>
            <div
              style={{
                display: "flex",
                overflow: "hidden",
                whiteSpace: "nowrap",
              }}
            >
              <div
                style={{ color: "#657b83", fontSize: 10, fontFamily: "Monaco" }}
              >
                {`@ ~/git-directory/${org}/${repo}`}
              </div>
              <div>&nbsp;&nbsp;</div>
              <a
                href={`${prURL} + ${commits[idx].hash}`}
                style={{ textDecoration: "none" }}
                target="_blank"
              >
                <div
                  style={{
                    color: "#859900",
                    fontSize: 10,
                    fontFamily: "Monaco",
                  }}
                >
                  {commits[idx].hash}
                </div>
              </a>
              <div>&nbsp;&nbsp;</div>
              <div
                style={{ color: "#268bd2", fontSize: 10, fontFamily: "Monaco" }}
              >
                {commits[idx].author}
              </div>
              <div>&nbsp;&nbsp;</div>
              <div
                style={{
                  color: "#eee8d5",
                  fontSize: 10,
                  fontFamily: "Monaco",
                }}
              >
                {commits[idx].message}
              </div>
            </div>

            <div>
              <div style={{ display: "flex" }}>
                <div style={{ marginRight: "4px", display: "flex" }}>
                  <div
                    style={{
                      color: "#b58900",
                      fontSize: 10,
                      fontFamily: "Monaco",
                    }}
                  >
                    mode:fzf&nbsp;$
                  </div>
                </div>
                <input
                  style={{
                    width: "90%",
                    backgroundColor: "#002b36",
                    padding: 0,
                    border: "none",
                    borderRadius: 0,
                    outline: "none",
                    background: "none",
                    color: "#eee8d5",
                    fontSize: 10,
                    fontFamily: "Monaco",
                  }}
                  onChange={(e: ChangeEvent<HTMLInputElement>): void =>
                    setFilter(e.target.value)}
                />
              </div>
              {filter
                ? (
                  <Fragment>
                    {fzfs.map((f) => (
                      <div
                        style={{
                          color: "#eee8d5",
                          fontSize: 10,
                          fontFamily: "Monaco",
                        }}
                        key={f.refIndex}
                      >
                        {f.item.name}
                      </div>
                    ))}
                  </Fragment>
                )
                : (
                  <Fragment>
                    {commits[idx].objects.map((f, i) => (
                      <div
                        key={i}
                        style={{
                          color: "#93a1a1",
                          fontSize: 10,
                          fontFamily: "Monaco",
                        }}
                      >
                        {f.name}
                      </div>
                    ))}
                  </Fragment>
                )}
            </div>
          </Fragment>
        )
        : (
          <LinearProgress />
        )}
    </div>
  );
};

export default Views;
