import { useState, useEffect } from "react";

import { Col, Card, Row, PageHeader, Form, Input, Button } from "antd";
// import { CreateSettings, Client, Gateway } from "fc-retrieval-client-js";

const GatewaysPage = () => {
  // const settings = new CreateSettings();
  // const config = settings.build();

  const [gatewayId, setGatewayId] = useState<string>("");
  // const [gateways, setGateways] = useState<Gateway[]>([]);

  // const client = new Client(config);

  // useEffect(() => {
  //   const fetchDataAsync = async () => {
  //     const gateways = await client.findGateways();
  //     setGateways(gateways);
  //   };
  //   fetchDataAsync();
  // }, [setGateways]);

  const onFinish = (values: any) => {
    console.log("Success:", values.gatewayId);
    setGatewayId(values.gatewayId);
  };

  const onFinishFailed = (errorInfo: any) => {
    console.log("Failed:", errorInfo);
  };

  return (
    <>
      <PageHeader className="page-header" title="Gateways" />
      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card>
            <h2>Active Gateways</h2>
            <Form
              name="basic"
              labelCol={{ span: 8 }}
              wrapperCol={{ span: 16 }}
              initialValues={{ remember: true }}
              onFinish={onFinish}
              onFinishFailed={onFinishFailed}
            >
              <Form.Item
                label="Gateway Id"
                name="gatewayId"
                rules={[
                  { required: true, message: "Please input your Gateway Id!" },
                ]}
              >
                <Input />
              </Form.Item>

              <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
                <Button type="primary" htmlType="submit">
                  Add new active gateway
                </Button>
              </Form.Item>
            </Form>

            {gatewayId ? <>{gatewayId}</> : "No active Gateway"}
          </Card>

          <Card>
            <h2>Registerd Gateways</h2>
            <ul>
              <li>
                0123456789{" "}
                <Button type="primary" htmlType="submit">
                  Add new active gateway
                </Button>
              </li>
              <li>98765465432</li>
            </ul>
          </Card>

          {/* {gateways.map((gateway) => (
            // test
            <Card>
              address: {gateway.address}
              <br />
              networkInfoAdmin: {gateway.networkInfoAdmin}
              <br />
              networkInfoClient: {gateway.networkInfoClient}
              <br />
              networkInfoGateway: {gateway.networkInfoGateway}
              <br />
              networkInfoProvider: {gateway.networkInfoProvider}
              <br />
              nodeId: {gateway.nodeId}
              <br />
              regionCode: {gateway.regionCode}
              <br />
              rootSigningKey: {gateway.rootSigningKey}
              <br />
              sigingKey: {gateway.sigingKey}
              <hr />
            </Card>
          ))} */}
        </Col>
      </Row>
    </>
  );
};

export default GatewaysPage;
