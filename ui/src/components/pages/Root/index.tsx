import React, { FC } from "react";
import Logo from "~/assets/logo.png"

// todo: ロゴとデモ動画 (URLにリポジトリ名入力させる)
const Root: FC = () => {
  return (
    <div>
      <img src={Logo} />
    </div>
  );
};

export default Root;
