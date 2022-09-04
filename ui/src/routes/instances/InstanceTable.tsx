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
  InstanceProxyService,
  instanceProxySelector,
} from "../../recoil/atoms/instances";
import { Service } from "../../recoil/atoms/services";
import { changeProxyStatus, patchProxy } from "./InstanceAPI";

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

const InstanceServiceInfo = ({ service }: { service: Service }) => {
  return (
    <span>
      {service.name}
      <ServiceInfoHelp {...service} />
    </span>
  );
};

const InstanceServiceProxy = ({
  id,
  proxy,
}: {
  id: number;
  proxy: InstanceProxyService;
}) => {
  const [getProxy, setProxy] = useState(proxy.port);

  const onClick = () => {
    patchProxy(id, proxy);
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
  proxy,
}: {
  id: number;
  proxy: InstanceProxyService;
}) => {
  const [status, setStatus] = useState(proxy.status);

  const handler = (isRunning: string) => {
    setStatus(isRunning);
    changeProxyStatus(id, proxy.id, proxy.status);
  };

  return (
    <ButtonGroup aria-label="toggle-proxy" size="sm">
      <Button
        className="service-running running-true"
        variant="secondary"
        active={status === "running"}
        onClick={() => handler("running")}
      >
        <BsCheck />
      </Button>
      <Button
        className="service-running running-false"
        variant="secondary"
        active={status === "stopped"}
        onClick={() => handler("stopped")}
      >
        <BsX />
      </Button>
    </ButtonGroup>
  );
};

const InstanceServiceRow = ({
  instance,
  proxy,
}: {
  instance: Instance;
  proxy: InstanceProxyService;
}) => {
  const cells = [
    <InstanceServiceInfo service={proxy.service} />,
    <InstanceServiceProxy id={instance.id} proxy={proxy} />,
    <InstanceServiceToggle id={instance.id} proxy={proxy} />,
  ];

  return <Row cells={cells} />;
};

const InstanceServicesTable = ({ instance }: { instance: Instance }) => {
  const proxyList = useRecoilValue(instanceProxySelector(instance.id));
  const rows = proxyList.map((proxy, index) => (
    <InstanceServiceRow key={index} instance={instance} proxy={proxy} />
  ));
  const data = {
    headers: [`${rows.length} Services`, "", ""],
    rows: [],
  };

  return <Table data={data} rows={rows} />;
};

export default InstanceServicesTable;
