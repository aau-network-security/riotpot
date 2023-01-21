import { useEffect, useRef, useState } from "react";
import {
  Button,
  ButtonGroup,
  Col,
  Form,
  InputGroup,
  Row,
} from "react-bootstrap";
import { AiOutlineInfoCircle } from "react-icons/ai";
import { BsArrowRepeat, BsCheck, BsX } from "react-icons/bs";
import { FaNetworkWired } from "react-icons/fa";
import { useRecoilCallback, useRecoilState, useRecoilValue } from "recoil";
import { Pop } from "../../components/pop/Pop";
import { Table, TableRow } from "../../components/table/Table";
import { useToast } from "../../components/toast/Toast";
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
import { ErrorToastVariant } from "../../recoil/atoms/toast";
import {
  changeProxyStatus,
  changeProxyPort,
  deleteProxyService,
  fetchProxy,
} from "./InstanceAPI";
import InstanceUtils from "./InstanceUtils";

// Simple helper to create the address string from an instance atom
const instanceAddress = (instance: Instance) =>
  instance.host + ":" + instance.port;

const ProxyServiceRowOptions = ({
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
    deleted
      .then((data) => {
        if ("success" in data) {
          removeService(id);
        }
      })
      .catch((err) => {
        console.log(err);
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
      <small className="info" ref={target}>
        <AiOutlineInfoCircle
          onMouseEnter={() => setShow(true)}
          onMouseLeave={() => setShow(false)}
        />
      </small>
      <Pop target={target} show={show} placement="left">
        <Row>
          <Col className="badges">
            <NetworkBadge {...network} />
            <InteractionBadge {...interaction} />
          </Col>
        </Row>
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
  const { showToast } = useToast();

  const [getProxy, setProxy] = useRecoilState(instanceService(proxyID));
  const [getProxyPort, setProxyPort] = useState(getProxy.port);

  const handler = () => {
    const proxyPort = changeProxyPort(host, proxyID, getProxyPort);
    proxyPort
      .then((data) => {
        if ("error" in data) {
          showToast(data["error"], ErrorToastVariant);
          return;
        }

        if ("port" in data) {
          setProxy(data);
        }
      })
      .catch((error) => {
        showToast(error.message, ErrorToastVariant);
        return;
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
  const { showToast } = useToast();
  const [status, setStatus] = useState(proxy.status);

  const handler = (isRunning: string) => {
    // Change the status of the thing
    if (status !== isRunning) {
      const response = changeProxyStatus(host, proxy.id, isRunning);

      response
        .then((data) => {
          if ("error" in data) {
            showToast(data["error"], ErrorToastVariant);
            return;
          }

          if (["running", "stopped"].includes(data.status))
            setStatus(data.status);
        })
        .catch((err) => {
          showToast(err.message, ErrorToastVariant);
        });
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
  // Get the address string of the instance
  const address = instanceAddress(instance);

  const cells = [
    <InstanceServiceInfo key={0} service={proxy.service} />,
    <InstanceServiceProxy key={1} host={address} proxyID={proxy.id} />,
    <InstanceServiceToggle key={2} host={address} proxy={proxy} />,
    <ProxyServiceRowOptions
      key={3}
      host={address}
      proxyID={proxy.id}
      serviceName={proxy.service.name}
    />,
  ];

  return <TableRow cells={cells} />;
};

const InstanceServicesTable = ({ instance }: { instance: Instance }) => {
  // Get all the proxy services set and create a row for each of them
  const proxyServices = useRecoilValue(instanceProxyServiceSelector);
  // Get the list of proxy service ids
  let ids = useRecoilValue(instanceServiceIDs);

  // If this variable receives a value, the table will not load
  let [error, setErr] = useState(Error);

  // Get the address string
  const address = instanceAddress(instance);

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
    const proxyList = fetchProxy(address);

    // For each of the proxy received add it to the state
    proxyList
      .then((proxies: InstanceProxyService[]) => {
        // Check if there was no response. This may happen
        if (proxies instanceof Error) {
          setErr(proxies);
          return;
        }

        proxies.forEach((x) => {
          addProxyService(x);
        });
      })
      .catch((err) => {
        setErr(err);
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (error.message) {
    return (
      <>
        <h4>Oooops! something went wrong...</h4>
        <p>
          We could not reach the address of the instance. Check the console
          (F12) to see if you can find a clue of the issue. At the very least,
          you should see an empty table here instead of this message. Once you
          have troubleshooted the issue, reload the page.
        </p>
        <small>
          No luck? Perhaps riotpot is not running. Maybe not in the given
          address; check your <a href="/settings">settings</a>!
        </small>
      </>
    );
  }

  // Map the rows into a proxy service
  const rows = proxyServices.map((proxy: InstanceProxyService) => (
    <InstanceServiceRow key={proxy.id} instance={instance} proxy={proxy} />
  ));

  // Send the data
  const data = {
    headers: [`${proxyServices.length} Services`, "", "", ""],
    rows: [],
  };

  return (
    <>
      <InstanceUtils host={address} />
      <Table data={data} rows={rows}></Table>
    </>
  );
};

export default InstanceServicesTable;
