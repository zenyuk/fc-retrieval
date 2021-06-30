import { useState, useEffect } from "react";

import { Col, Card, Row, Form, Input, Button, PageHeader } from "antd";
// import { CreateSettings, Client, Offer } from "fc-retrieval-client-js";

const OffersPage = () => {
  // const settings = new CreateSettings();
  // const config = settings.build();

  const [cid, setCid] = useState<string>("");

  // const client = new Client(config);

  // useEffect(() => {
  //   const fetchDataAsync = async () => {
  //     const offers = await client.findOffers();
  //     setOffers(offers);
  //   };
  //   fetchDataAsync();
  // }, [setOffers]);

  const onFinish = (values: any) => {
    console.log("Success:", values);
    setCid(values.cid);
  };

  const onFinishFailed = (errorInfo: any) => {
    console.log("Failed:", errorInfo);
  };

  return (
    <>
      <PageHeader className="page-header" title="Offers" />
      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card>
            <Form
              name="basic"
              labelCol={{ span: 8 }}
              wrapperCol={{ span: 16 }}
              initialValues={{ remember: true }}
              onFinish={onFinish}
              onFinishFailed={onFinishFailed}
            >
              <Form.Item
                label="Cid"
                name="cid"
                rules={[{ required: true, message: "Please input your Cid!" }]}
              >
                <Input />
              </Form.Item>

              <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
                <Button type="primary" htmlType="submit">
                  Submit
                </Button>
              </Form.Item>
            </Form>
          </Card>
          {cid && (
            <Card>
              <h2>Offer found</h2>
              <p>Cid: {cid}</p>
              <p>Price: 42</p>
            </Card>
          )}

          {/* {offers.map((offer) => (
            // test
            <Card>[...]</Card>
          ))} */}
        </Col>
      </Row>
    </>
  );
};

export default OffersPage;
