import React from "react";
import PropTypes from "prop-types";

import { Layout } from "antd";
import Header from "@src/components/pages/Header";

const { Content } = Layout;

const MinimalLayout: React.FC = ({ children }) => {
  return (
    <Content style={{ padding: "0 50px" }}>
      <Layout className="site-layout-background" style={{ padding: "0 0" }}>
        <Content style={{ padding: "0 0", minHeight: 280 }}>
          <div
            style={{
              background: "#f8f8f8",
              minHeight: "calc(100vh)",
              height: "100%",
            }}
          >
            <Header />
            <div className="container">{children}</div>
          </div>
        </Content>
      </Layout>

      <style jsx>{`
        .container {
          padding: 1.25rem 1.25rem;
        }
      `}</style>
    </Content>
  );
};

MinimalLayout.propTypes = { children: PropTypes.node.isRequired };

export default MinimalLayout;
