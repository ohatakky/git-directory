import React, { FC } from "react";
import { Link } from "react-router-dom";
import Logo from "~/assets/logo_transparent.png";
import Chrome from "~/assets/chrome.png";
import Demo from "~/assets/demo.gif";

const Root: FC = () => {
  return (
    <div
      style={{
        padding: "96px 48px 0",
        display: "flex",
        justifyContent: "space-evenly",
        alignItems: "center",
      }}
    >
      <img src={Demo} style={{ borderRadius: "3px" }} />
      <div>
        <div style={{ display: "flex", alignItems: "center" }}>
          <img src={Logo} style={{ width: "120px" }} />
          <div
            style={{ fontFamily: "Monaco", fontSize: "16", color: "#545474" }}
          >
            Visualizes repository directory history
          </div>
        </div>
        <div
          style={{
            display: "flex",
            alignItems: "center",
            flexDirection: "column",
          }}
        >
          <div
            style={{
              marginTop: "24px",
              display: "flex",
              alignItems: "center",
              flexDirection: "column",
            }}
          >
            <div
              style={{ fontFamily: "Monaco", fontSize: "12", color: "#545474" }}
            >
              Sample Repository
            </div>
            <div style={{ marginTop: "12px" }}></div>
            <Link to="gorilla/websocket">
              <button
                style={{
                  backgroundColor: "transparent",
                  border: "solid 1px #545474",
                  cursor: "pointer",
                  borderRadius: "2px",
                  outline: "none",
                  padding: "4px",
                  fontFamily: "Monaco",
                  fontSize: "12",
                  color: "#545474",
                }}
              >
                Try It
              </button>
            </Link>
          </div>
          <div
            style={{
              marginTop: "24px",
              display: "flex",
              alignItems: "center",
              flexDirection: "column",
            }}
          >
            <div
              style={{ fontFamily: "Monaco", fontSize: "12", color: "#545474" }}
            >
              Also available as chrome extension
            </div>
            <div style={{ marginTop: "12px" }}></div>
            <img src={Chrome} style={{ width: "36px" }} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Root;
