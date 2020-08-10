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

// todo: cssリファクタ
// 色をutilsへ
// typography共通化
// materialuiに寄せる

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

const reducer = (state: Fzf[], action: { type: ActionType; payload: Fzf }) => {
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
  const [filter, setFilter] = useState<string>("");
  const [fzfviews, setFzfviews] = useState<Fuse.FuseResult<string>[]>([]);
  const prURL = `https://github.com/${org}/${repo}/pulls?q=is%3Apr+hash%3A`;
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

  useEffect(() => {
    if (!fzfs[idx]) return;

    const fuse = new Fuse(fzfs[idx].files, {
      isCaseSensitive: true,
      minMatchCharLength: 2,
    })
      .search(filter)
      .map((f) => f);
    setFzfviews(fuse);
  }, [filter]);

  return (
    <div
      style={{
        width: "100%",
        minHeight: "100%",
        backgroundColor: "#002b36",
        padding: "4px",
        ...(fzfs[idx] && fzfs[idx].files.length > 0 &&
          { height: `${fzfs[idx].files.length * 16}px` }),
      }}
    >
      {fzfs.length > 0
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
                href={`${prURL} + ${fzfs[idx].hash}`}
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
                  {fzfs[idx].hash}
                </div>
              </a>
              <div>&nbsp;&nbsp;</div>
              <div
                style={{ color: "#268bd2", fontSize: 10, fontFamily: "Monaco" }}
              >
                {fzfs[idx].commit.author}
              </div>
              <div>&nbsp;&nbsp;</div>
              <div
                style={{
                  color: "#eee8d5",
                  fontSize: 10,
                  fontFamily: "Monaco",
                }}
              >
                {fzfs[idx].commit.message}
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
                    {fzfviews.map((f) => (
                      <div
                        style={{
                          color: "#eee8d5",
                          fontSize: 10,
                          fontFamily: "Monaco",
                        }}
                        key={f.refIndex}
                      >
                        {f.item}
                      </div>
                    ))}
                  </Fragment>
                )
                : (
                  <Fragment>
                    {fzfs[idx].files.map((f, i) => (
                      <div
                        key={i}
                        style={{
                          color: "#93a1a1",
                          fontSize: 10,
                          fontFamily: "Monaco",
                        }}
                      >
                        {f}
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

export default FzfFC;
