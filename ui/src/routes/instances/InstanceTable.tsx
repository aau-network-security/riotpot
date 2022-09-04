import { useRef, useState } from "react";
import { Button, ButtonGroup, Col, Form, InputGroup } from "react-bootstrap";
import { AiOutlineInfoCircle } from "react-icons/ai";
import { BsArrowRepeat, BsCheck, BsX } from "react-icons/bs";
import { FaNetworkWired } from "react-icons/fa";
import { useRecoilValue } from "recoil";
import { Pop } from "../../components/pop/Pop";
import Table, { Row } from "../../components/table/Table";
import { InteractionBadge, NetworkBadge } from "../../components/utils/Common";
import { InteractionOption, NetworkOption } from "../../constants/globals";
import {
  Instance,
  InstanceService,
  instanceServicesSelector,
} from "../../recoil/atoms/instances";
import { patchService } from "./InstanceAPI";

const ServiceInfoHelp = ({
  network,
  interaction,
}: {
  network: NetworkOption;
  interaction: InteractionOption;
}) => {
  const [show, setShow] = useState(false);
  const target = useRef(null);

  return (
    <>
      <small
        className="info"
        ref={target}
        onMouseEnter={() => setShow(true)}
        onMouseLeave={() => setShow(false)}
      >
        <AiOutlineInfoCircle />
      </small>
      <Pop target={target} show={show} placement="left">
        <span>
          <NetworkBadge {...network} />
          <InteractionBadge {...interaction} />
        </span>
      </Pop>
    </>
  );
};

const InstanceServiceInfo = ({ service }: { service: InstanceService }) => {
  return (
    <>
      <span>
        {service.name}
        <ServiceInfoHelp {...service} />
      </span>
    </>
  );
};

const InstanceServiceProxy = ({
  id,
  service,
}: {
  id: number;
  service: InstanceService;
}) => {
  const [getProxy, setProxy] = useState(service.proxy);

  const onClick = () => {
    patchService(id, service);
  };

  return (
    <Col xs="8">
      <InputGroup size="sm" className="address">
        <InputGroup.Text>
          <FaNetworkWired />
        </InputGroup.Text>
        <Form.Control
          type="number"
          min={1}
          max={65535}
          defaultValue={getProxy}
          onChange={(e) => setProxy(parseInt(e.target.value))}
        />
        <Button
          variant="outline-secondary"
          id="button-addon2"
          onClick={onClick}
        >
          <BsArrowRepeat />
        </Button>
      </InputGroup>
    </Col>
  );
};

const InstanceServiceToggle = ({
  id,
  service,
}: {
  id: number;
  service: InstanceService;
}) => {
  const [running, setRunning] = useState(service.running);

  const handler = (isRunning: boolean) => {
    setRunning(isRunning);
    patchService(id, { ...service, running: isRunning });
  };

  return (
    <ButtonGroup aria-label="toggle-proxy" size="sm">
      <Button
        className="service-running running-true"
        variant="secondary"
        active={running === true}
        onClick={() => handler(true)}
      >
        <BsCheck />
      </Button>
      <Button
        className="service-running running-false"
        variant="secondary"
        active={running === false}
        onClick={() => handler(false)}
      >
        <BsX />
      </Button>
    </ButtonGroup>
  );
};

const InstanceServiceRow = ({
  instance,
  service,
}: {
  instance: Instance;
  service: InstanceService;
}) => {
  const cells = [
    <InstanceServiceInfo service={service} />,
    <InstanceServiceProxy id={instance.id} service={service} />,
    <InstanceServiceToggle id={instance.id} service={service} />,
  ];

  return <Row cells={cells} />;
};

const InstanceServicesTable = ({ instance }: { instance: Instance }) => {
  const services = useRecoilValue(instanceServicesSelector(instance.id));
  const rows = services.map((service, index) => (
    <InstanceServiceRow key={index} instance={instance} service={service} />
  ));
  const data = {
    headers: [`${rows.length} Services`, "", ""],
    rows: [],
  };

  return <Table data={data} rows={rows} />;
};

export default InstanceServicesTable;
