import { useState } from "react";
import { Button, Form, InputGroup, Row, Stack } from "react-bootstrap";
import { useRecoilState } from "recoil";
import { instance } from "../../recoil/atoms/instances";

export const SettingsAPIAddress = () => {
  // Get the current values of the instance address
  const [ins, setIns] = useRecoilState(instance);

  // Create a dummy that can hold the values until the submit button is clicked
  const [dummy, setDummy] = useState(ins);

  // Create a setter for the submit button
  const editAddress = () => {
    setIns({ ...ins, ...dummy });
  };

  // Updates a field on the dummy
  const updateDummyField = (field: string, value: string | number) => {
    setDummy({ ...dummy, [field]: value });
  };

  return (
    <Stack gap={5}>
      <Row>
        <h1>API Address</h1>
        <small>Change the location of the instance API address</small>
      </Row>

      <Row>
        <Form>
          <Form.Label htmlFor="basic-url">Address</Form.Label>
          <InputGroup className="mb-3">
            <Form.Control
              placeholder="host"
              aria-label="host"
              defaultValue={ins.host}
              onChange={(e) => updateDummyField("host", e.target.value)}
            />
            <Form.Control
              type="number"
              placeholder="port"
              aria-label="port"
              defaultValue={ins.port}
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
