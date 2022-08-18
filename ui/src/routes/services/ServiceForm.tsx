import { Form } from "react-bootstrap";
import Select from "react-select";
import { SimpleForm } from "../../components/forms/Form";
import {
  customTheme,
  InteractionOptions,
  NetworkOption,
  NetworkOptions,
} from "../../constants/globals";

const InteractionLevelSelectfield = ({
  option = InteractionOptions[0],
}: {
  option?: any;
}) => {
  return (
    <Form.Group className="mb-3" controlId="formServiceInteraction">
      <Form.Label>Interaction Level</Form.Label>
      <Select
        className="basic-single"
        classNamePrefix="select"
        value={option}
        options={InteractionOptions}
        isSearchable={true}
        name="interactionLevel"
        styles={customTheme}
      />
      <Form.Text className="text-muted">
        Service interaction level. An emulated service should be low; whereas a
        fully flagged service will be high.
      </Form.Text>
    </Form.Group>
  );
};

const NetworkSelectfield = ({
  option = NetworkOptions[0],
}: {
  option?: NetworkOption;
}) => {
  return (
    <Form.Group className="mb-3" controlId="formServiceNetwork">
      <Form.Label>Network</Form.Label>
      <Select
        className="basic-single"
        classNamePrefix="select"
        value={option}
        options={NetworkOptions}
        isSearchable={true}
        name="network"
        styles={customTheme}
      />
      <Form.Text className="text-muted">Network layer protocol</Form.Text>
    </Form.Group>
  );
};

type ServiceFormSerial = {
  name: string;
  network: NetworkOption;
  interaction: any;
  host: string;
  port: string;
};

const ServiceForm = ({
  create = true,
  data,
}: {
  create: boolean;
  data?: ServiceFormSerial;
}) => {
  let fields = [
    {
      name: "name",
      help: "Name of the service",
      fieldattrs: {
        type: "text",
        required: true,
        value: data?.name,
        placeholder: "Type a name for the service",
      },
    },
    {
      name: "host",
      help: "Host in where the service can be reached. E.g.: localhost, 0.0.0.0, etc.",
      fieldattrs: {
        type: "text",
        required: true,
        value: data?.host,
        placeholder: "localhost",
      },
    },
    {
      name: "port",
      help: "Port address of the service. A number between 1 and 65535",
      fieldattrs: {
        type: "number",
        required: true,
        min: 1,
        max: 65535,
        placeholder: 1111,
      },
    },
  ];

  return (
    <SimpleForm page="Services" fields={fields}>
      {create && <NetworkSelectfield option={data && data.network} />}
      <InteractionLevelSelectfield option={data && data.interaction} />
    </SimpleForm>
  );
};

export default ServiceForm;
