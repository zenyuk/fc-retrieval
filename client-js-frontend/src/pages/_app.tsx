import React from "react";
import { Layout } from "antd";
import { NextPage } from "next";
import "antd/dist/antd.css";
import "@src/styles/theme.css";

import PageLayout from "@src/layouts/PageLayout";

export type PageWithLayout = NextPage & {
  Layout: null;
};
interface Props {
  Component: any;
  pageProps: any;
}

const App: React.FC<Props> = (props: Props) => {
  const { Component, pageProps } = props;
  const useLayout = (Component as PageWithLayout).Layout;
  const AppLayout = useLayout || PageLayout;

  return (
    <Layout>
      <AppLayout>
        <Component {...pageProps} />
      </AppLayout>
    </Layout>
  );
};

export default App;
