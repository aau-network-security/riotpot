import { useRef, useState } from "react";
import { Button, Col, Form, InputGroup } from "react-bootstrap";
import {
  AiFillCheckCircle,
  AiFillCloseCircle,
  AiOutlineInfoCircle,
} from "react-icons/ai";
import { BsArrowRepeat, BsGlobe } from "react-icons/bs";
import { FaNetworkWired } from "react-icons/fa";
import { Pop } from "../../components/pop/Pop";
import Table from "../../components/table/Table";
import {
  DeleteDropdownItem,
  InteractionBadge,
  OptionsDropdown,
} from "../../components/utils/Common";
import {
  getPage,
  InteractionOption,
  InteractionOptions,
} from "../../constants/globals";

interface InstanceService {
  name: string;
  proxy: Number;
  interaction: InteractionOption;
  running: boolean;
}

interface Instance {
  name: string;
  services: InstanceService[];
  address?: string;
}

const ProxyInfoPop = ({ proxy }: { proxy: Number }) => {
  return (
    <span className="proxy-info">
      <FaNetworkWired />
      {`${proxy}`}
    </span>
  );
};

const InstancePop = ({ service }: { service: InstanceService }) => {
  const cross = <AiFillCloseCircle style={{ fill: "#FF8686" }} />;
  const tick = <AiFillCheckCircle style={{ fill: "#A4E18F" }} />;

  return (
    <Col>
      {service.running ? tick : cross} {service.name}
      <ProxyInfoPop proxy={service.proxy} />
      <InteractionBadge {...service.interaction} />
    </Col>
  );
};

const InstanceRowInfoPop = ({ services }: { services: InstanceService[] }) => {
  const [show, setShow] = useState(false);
  const target = useRef(null);

  return (
    <>
      <small className="info" ref={target} onClick={() => setShow(!show)}>
        <AiOutlineInfoCircle />
        {`${services.length} ${services.length === 1 ? "service" : "services"}`}
      </small>
      <Pop target={target} show={show} placement="right">
        {services.map((service) => {
          return <InstancePop service={service} />;
        })}
      </Pop>
    </>
  );
};

const InstanceRowInfo = ({ name, services }: Instance) => {
  return (
    <>
      <div>{name}</div>
      <InstanceRowInfoPop services={services} />
    </>
  );
};

const InstanceRowAddress = ({ value }: { value?: string }) => {
  return (
    <Col xs="8">
      <InputGroup size="sm" className="address">
        <InputGroup.Text>
          <BsGlobe />
        </InputGroup.Text>
        <Form.Control
          placeholder="Address"
          aria-label="Instance's address"
          aria-describedby="basic-addon2"
          value={value}
        />
        <Button variant="outline-secondary" id="button-addon2">
          <BsArrowRepeat />
        </Button>
      </InputGroup>
    </Col>
  );
};

const InstanceRowOptions = ({ name }: { name: string }) => {
  const page = getPage("Instances");
  const note =
    "Instance services will be stopped and removed from the instance register.";
  return (
    <OptionsDropdown>
      {page && <DeleteDropdownItem page={page} note={note} name={name} />}
    </OptionsDropdown>
  );
};

const InstanceRow = ({ name, address, services }: Instance) => {
  return [
    <InstanceRowInfo name={name} services={services} />,
    <InstanceRowAddress value={address} />,
    <InstanceRowOptions name={name} />,
  ];
};

export const InstancesTable = () => {
  const props: Instance = {
    name: "Lab1",
    address: "127.0.0.11",
    services: [
      {
        name: "CoAP",
        proxy: 5683,
        interaction: InteractionOptions[0],
        running: true,
      },
    ],
  };

  const rows = [InstanceRow(props), InstanceRow(props)];

  const data = {
    headers: [`${rows.length} Instances`, "", ""],
    rows: rows,
  };

  return <Table data={data} />;
};
