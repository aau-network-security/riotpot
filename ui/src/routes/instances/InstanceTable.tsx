import { useEffect, useRef, useState } from "react";
import { Button, ButtonGroup, Col, Form, InputGroup } from "react-bootstrap";
import { AiOutlineInfoCircle } from "react-icons/ai";
import { BsArrowRepeat, BsCheck, BsX } from "react-icons/bs";
import { FaNetworkWired } from "react-icons/fa";
import { useRecoilCallback, useRecoilState, useRecoilValue } from "recoil";
import { Pop } from "../../components/pop/Pop";
import Table, { Row } from "../../components/table/Table";
import {
  DeleteDropdownItem,
  InteractionBadge,
  NetworkBadge,
  OptionsDropdown,
} from "../../components/utils/Common";
import {
  getPage,
  InteractionOption,
  NetworkOption,
} from "../../constants/globals";
import {
  Instance,
  InstanceProxyService,
  instanceServiceIDs,
  instanceService,
  instanceProxyServiceSelector,
} from "../../recoil/atoms/instances";
import { Service } from "../../recoil/atoms/services";
import {
  changeProxyStatus,
  changeProxyPort,
  deleteProxyService,
  fetchProxy,
} from "./InstanceAPI";

const ProxyServicceRowOptions = ({
  host,
  proxyID,
  serviceName,
}: {
  host: string;
  proxyID: string;
  serviceName: string;
}) => {
  const removeService = useRecoilCallback(({ set }) => (id: string) => {
    set(instanceServiceIDs, (prev) => prev.filter((x) => x !== id));
  });

  const deleteCallback = (id: string) => {
    const deleted = deleteProxyService(host, id);
    deleted.then((data) => {
      if ("success" in data) {
        removeService(id);
      }
    });
  };

  const page = getPage("Services");
  const note = "The service will be stopped and removed from the instance";

  return (
    <OptionsDropdown>
      {page && (
        <DeleteDropdownItem
          page={page}
          note={note}
          name={serviceName}
          onClick={() => deleteCallback(proxyID)}
        />
      )}
    </OptionsDropdown>
  );
};

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
  host,
  proxyID,
}: {
  host: string;
  proxyID: string;
}) => {
  const [getProxy, setProxy] = useRecoilState(instanceService(proxyID));
  const [getProxyPort, setProxyPort] = useState(getProxy.port);

  const handler = () => {
    const proxyPort = changeProxyPort(host, proxyID, getProxyPort);
    proxyPort.then((data) => {
      if ("port" in data) {
        setProxy(data);
      }
    });
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
          defaultValue={getProxyPort}
          onChange={(e) => setProxyPort(parseInt(e.target.value))}
        />
        <Button
          variant="outline-secondary"
          id="button-addon2"
          onClick={() => handler()}
        >
          <BsArrowRepeat />
        </Button>
      </InputGroup>
    </Col>
  );
};

const InstanceServiceToggle = ({
  host,
  proxy,
}: {
  host: string;
  proxy: InstanceProxyService;
}) => {
  const [status, setStatus] = useState(proxy.status);

  const handler = (isRunning: string) => {
    // Change the status of the thing
    if (status !== isRunning) {
      changeProxyStatus(host, proxy.id, isRunning, setStatus);
    }
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
    <InstanceServiceProxy host={instance.host} proxyID={proxy.id} />,
    <InstanceServiceToggle host={instance.host} proxy={proxy} />,
    <ProxyServicceRowOptions
      host={instance.host}
      proxyID={proxy.id}
      serviceName={proxy.service.name}
    />,
  ];

  return <Row cells={cells} />;
};

const InstanceServicesTable = ({ instance }: { instance: Instance }) => {
  // Get all the proxy services set and create a row for each of them
  const proxyServices = useRecoilValue(instanceProxyServiceSelector);
  // Get the list of proxy service ids
  let ids = useRecoilValue(instanceServiceIDs);

  // Callback to add a service to the list.
  // This is used to track and update the state of the proxies
  const addProxyService = useRecoilCallback(
    ({ set }) =>
      (proxyService: InstanceProxyService) => {
        // Set the new id in the list if it is not there yet
        if (!ids.includes(proxyService.id)) {
          ids = [...ids, proxyService.id];
          set(instanceServiceIDs, ids);
        }

        set(instanceService(proxyService.id), proxyService);
      }
  );

  // Fetch the list of proxy services only once
  useEffect(() => {
    // Populate the list of services
    const proxyList = fetchProxy(instance.host);
    // For each of the proxy received add it to the state
    proxyList.then((proxies: InstanceProxyService[]) => {
      proxies.forEach((x) => {
        addProxyService(x);
      });
    });
  }, []);

  // Map the rows into a proxy service
  const rows = proxyServices.map((proxy: any, index: number) => (
    <InstanceServiceRow key={index} instance={instance} proxy={proxy} />
  ));

  // Send the data
  const data = {
    headers: [`${proxyServices.length} Services`, "", "", ""],
    rows: [],
  };

  return <Table data={data} rows={rows}></Table>;
};

export default InstanceServicesTable;
