import React from "react";
import { Layout, Menu } from "antd";
import { useRouter } from "next/router";

import {
  AppstoreOutlined,
  ApartmentOutlined,
  OrderedListOutlined,
  SyncOutlined,
} from "@ant-design/icons";

import Link from "next/link";

const { Sider } = Layout;

export const Header: React.FC = (): JSX.Element => {
  const router = useRouter();

  const regexp = /^\/[^\/]*/;
  const match = router.pathname.match(regexp);
  const pathname = match?.length ? match[0] : "";

  return (
    <Sider
      className="site-layout-background"
      width={180}
      style={{ fontSize: "24px", backgroundColor: "#ffffff" }}
    >
      <Link href="/">
        <div
          className="header-color"
          style={{
            height: "64px",
            fontSize: "16px",
            color: "#ffffff",
            backgroundColor: "#2f3438",
            paddingLeft: "20px",
            paddingTop: "17px",
          }}
        >
          FC Retrieval
        </div>
      </Link>
      <Menu
        mode="inline"
        defaultSelectedKeys={["/"]}
        selectedKeys={[pathname]}
        style={{ fontSize: "16px" }}
      >
        <Menu.Item
          key="/"
          icon={<AppstoreOutlined style={{ fontSize: "16px" }} />}
        >
          <Link href="/">Home</Link>
        </Menu.Item>

        <Menu.Item
          key="/gateways"
          icon={<ApartmentOutlined style={{ fontSize: "16px" }} />}
        >
          <Link href="/gateways">Gateways</Link>
        </Menu.Item>

        <Menu.Item
          key="/offers"
          icon={<OrderedListOutlined style={{ fontSize: "16px" }} />}
        >
          <Link href="/offers">Offers</Link>
        </Menu.Item>
      </Menu>
    </Sider>
  );
};

export default Header;
