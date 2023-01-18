import { useState } from "react";
import { Button, Form, InputGroup, Row, Stack } from "react-bootstrap";
import { useRecoilState } from "recoil";
import { backend } from "../../recoil/atoms/settings";

export const SettingsAPIAddress = () => {
  // Get the current values of the backend address
  const [bk, setBk] = useRecoilState(backend);

  // Create a dummy that can hold the values until the submit button is clicked
  const [dummy, setDummy] = useState(bk);

  // Create a setter for the submit button
  const editAddress = () => {
    setBk({ ...bk, ...dummy });
  };

  // Updates a field on the dummy
  const updateDummyField = (field: string, value: string | number) => {
    setDummy({ ...dummy, [field]: value });
  };

  return (
    <Stack gap={5}>
      <Row>
        <h1>API Address</h1>
        <small>Change the location of the backend API address</small>
      </Row>

      <Row>
        <Form>
          <Form.Label htmlFor="basic-url">Address</Form.Label>
          <InputGroup className="mb-3">
            <Form.Control
              placeholder="host"
              aria-label="host"
              defaultValue={bk.host}
              onChange={(e) => updateDummyField("host", e.target.value)}
            />
            <Form.Control
              type="number"
              placeholder="port"
              aria-label="port"
              defaultValue={bk.port}
              onChange={(e) => updateDummyField("port", e.target.value)}
            />
            <Button variant="primary" onClick={editAddress}>
              Save
            </Button>
          </InputGroup>
          <Form.Text className="text-muted">
            Host and port of the address (i.e., localhost:2022)
          </Form.Text>
        </Form>
      </Row>
    </Stack>
  );
};
